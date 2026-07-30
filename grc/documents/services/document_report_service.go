package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentReportServiceInterface interface {
	GetAll(params *dtos.DocumentReportListQueryParamsDTO) (*dtos.DocumentReportListResponseDTO, error)
	GetDocumentReportDetails(reportID uint) (*dtos.DocumentReportDetailDTO, error)
	Create(currentUserID uint, formDTO *dtos.DocumentReportFormDTO) (uint, error)
	DeleteDocumentReport(reportID uint, currentUserID uint) error
	GetDocumentReportExportZip(reportID uint, currentUserID uint) (string, error)
}

type DocumentReportService struct {
	repo           repositories.DocumentReportRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	settingSvc     settingServices.ReportTemplateSettingServiceInterface
}

func NewDocumentReportService(
	repo repositories.DocumentReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	settingSvc settingServices.ReportTemplateSettingServiceInterface,
) DocumentReportServiceInterface {
	return &DocumentReportService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		settingSvc:     settingSvc,
	}
}

func (drs *DocumentReportService) GetAll(params *dtos.DocumentReportListQueryParamsDTO) (*dtos.DocumentReportListResponseDTO, error) {
	return drs.repo.GetAll(params)
}

func (drs *DocumentReportService) GetDocumentReportDetails(reportID uint) (*dtos.DocumentReportDetailDTO, error) {
	report, err := drs.repo.GetDocumentReportDetails(reportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, reportID)
		}
		return nil, err
	}
	sections, err := drs.repo.GetDocumentReportSections(reportID)
	if err != nil {
		return nil, err
	}

	documents := make([]*dtos.DocumentReportDocumentDTO, 0, len(report.Documents))
	for _, reportDocument := range report.Documents {
		var document dtos.DocumentReportDocumentDTO
		if err := mapstructure.Decode(reportDocument, &document); err != nil {
			return nil, fmt.Errorf("failed to map document report snapshot %d: %w", reportDocument.ID, err)
		}

		document.DocumentCreatedAt = reportDocument.DocumentCreatedAt
		document.NextReviewDate = reportDocument.NextReviewDate
		if reportDocument.PublishedAt != nil {
			publishedAt := reportDocument.PublishedAt.Format("2006-01-02")
			document.PublishedAt = &publishedAt
		}
		document.Clauses = []dtos.DocumentReportClauseDTO{}

		if err := json.Unmarshal([]byte(reportDocument.ClausesJSON), &document.Clauses); err != nil {
			return nil, fmt.Errorf("failed to decode clauses snapshot for document %s: %w", reportDocument.ReferenceNumber, err)
		}
		documents = append(documents, &document)
	}

	var createdBy organizationDTOs.UserMinimalDTO
	if user := report.CreatedByUser; user != nil {
		createdBy = organizationDTOs.UserMinimalDTO{
			ID: int(user.ID), FirstNameEN: user.FirstNameEn, FirstNameAR: user.FirstNameAr,
			LastNameEN: user.LastNameEn, LastNameAR: user.LastNameAr, Email: user.Email,
			Initial: user.Initial, Color: user.Color, ProfileImagePath: user.ProfileImagePath,
			MobileNumber: user.MobileNumber,
		}
	}
	return &dtos.DocumentReportDetailDTO{
		ID: report.ID, Title: report.Title, DocumentCategoryID: report.DocumentCategoryID,
		DocumentTypeID: report.DocumentTypeID, CreatedAt: report.CreatedAt, CreatedBy: createdBy,
		Sections: sections, Documents: documents,
	}, nil
}

func (drs *DocumentReportService) Create(currentUserID uint, formDTO *dtos.DocumentReportFormDTO) (uint, error) {
	tx := drs.repo.BeginTransaction()
	sourceDocuments, err := drs.repo.GetPublishedDocumentsForReporting(formDTO.DocumentCategoryID, formDTO.DocumentTypeID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if len(sourceDocuments) == 0 {
		tx.Rollback()
		return 0, utility.NewValidationError(errors.New("no published documents match the selected category and type"), "not_found", "documents")
	}

	reportID, err := drs.repo.Create(&models.DocumentReport{
		Title: formDTO.Title, DocumentCategoryID: formDTO.DocumentCategoryID,
		DocumentTypeID: formDTO.DocumentTypeID, CreatedBYID: currentUserID,
	}, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for index, sectionType := range []string{"documents", "clauses-controls"} {
		titleEn, titleAr := documentReportSectionTitle(sectionType)
		if err := drs.repo.CreateSection(&models.DocumentReportSection{
			DocumentReportID: reportID, Type: sectionType, TitleEn: titleEn, TitleAr: titleAr, OrderNum: index + 1,
		}, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	snapshots := make([]*models.DocumentReportDocument, 0, len(sourceDocuments))
	for _, source := range sourceDocuments {
		clauses, err := drs.repo.GetDocumentClausesForReporting(source.SourceDocumentID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		controls, err := drs.repo.GetDocumentControlsForReporting(source.SourceDocumentID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		controlsByClauseID := make(map[uint][]dtos.DocumentReportControlDTO)
		for _, control := range controls {
			if control.ClauseID != nil {
				controlsByClauseID[*control.ClauseID] = append(controlsByClauseID[*control.ClauseID], control)
			}
		}
		for index := range clauses {
			clauses[index].Controls = []dtos.DocumentReportControlDTO{}
			clauses[index].Controls = append(clauses[index].Controls, controlsByClauseID[clauses[index].ID]...)
		}

		clausesJSON, err := json.Marshal(clauses)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		snapshots = append(snapshots, &models.DocumentReportDocument{
			DocumentReportID: reportID,
			ReferenceNumber:  source.ReferenceNumber, TitleEn: source.TitleEn, TitleAr: source.TitleAr,
			DescriptionEn: source.DescriptionEn, DescriptionAr: source.DescriptionAr,
			DocumentCategoryTitleEn: source.DocumentCategoryTitleEn, DocumentCategoryTitleAr: source.DocumentCategoryTitleAr,
			DocumentTypeTitleEn: source.DocumentTypeTitleEn, DocumentTypeTitleAr: source.DocumentTypeTitleAr,
			DocumentStatusTitleEn: source.DocumentStatusTitleEn, DocumentStatusTitleAr: source.DocumentStatusTitleAr,
			DocumentStatusColorCode: source.DocumentStatusColorCode, DocumentStatusType: source.DocumentStatusType,
			DocumentAccessTypeTitleEn: source.DocumentAccessTypeTitleEn, DocumentAccessTypeTitleAr: source.DocumentAccessTypeTitleAr,
			DocumentAccessTypeColorCode:          source.DocumentAccessTypeColorCode,
			ConfidentialityClassificationTitleEn: source.ConfidentialityClassificationTitleEn,
			ConfidentialityClassificationTitleAr: source.ConfidentialityClassificationTitleAr,
			ConfidentialityClassificationType:    source.ConfidentialityClassificationType,
			ReviewFrequencyTitleEn:               source.ReviewFrequencyTitleEn, ReviewFrequencyTitleAr: source.ReviewFrequencyTitleAr,
			ReviewFrequencyType:    source.ReviewFrequencyType,
			ResponsibleUnitTitleEn: source.ResponsibleUnitTitleEn, ResponsibleUnitTitleAr: source.ResponsibleUnitTitleAr,
			ResponsibleUserNameEn: source.ResponsibleUserNameEn, ResponsibleUserNameAr: source.ResponsibleUserNameAr,
			Version: source.Version, FileName: source.FileName, FileType: source.FileType, FileSize: source.FileSize,
			DocumentCreatedAt: source.DocumentCreatedAt, NextReviewDate: source.NextReviewDate, PublishedAt: source.PublishedAt,
			ClausesJSON: string(clausesJSON),
			CreatedBYID: currentUserID,
		})
	}
	if err := drs.repo.CreateDocumentSnapshots(snapshots, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return reportID, nil
}

func (drs *DocumentReportService) DeleteDocumentReport(reportID uint, currentUserID uint) error {
	tx := drs.repo.BeginTransaction()
	existingReport, err := drs.repo.GetDocumentReportByID(reportID, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, reportID)
		}
		return err
	}
	if err := drs.repo.DeleteDocumentReport(reportID, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	if drs.activityLogSvc != nil {
		go func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					utility.LogError("Panic recovered in document report activity log goroutine", map[string]interface{}{
						"recover": recovered,
						"stack":   string(debug.Stack()),
					})
				}
			}()

			if err := drs.activityLogSvc.LogActivity(
				models.DeleteAction,
				models.DocumentReport{}.TableName(),
				reportID,
				existingReport,
				currentUserID,
				"document", "document_report", "", nil,
			); err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func (drs *DocumentReportService) GetDocumentReportExportZip(reportID uint, currentUserID uint) (string, error) {
	report, err := drs.GetDocumentReportDetails(reportID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch document report details: %w", err)
	}

	templateSetting, err := drs.settingSvc.GetFirst()
	if err != nil {
		return "", err
	}
	if templateSetting.CoverPageBGImagePath == nil || templateSetting.LogoImagePath == nil {
		return "", utility.NewForbiddenError("Logo and Background Image not added.")
	}
	if _, err := drs.settingSvc.GetImagePath("logo"); err != nil {
		return "", utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := drs.settingSvc.GetImagePath("cover-page-bg"); err != nil {
		return "", utility.NewForbiddenError("cover_image_missing")
	}

	reportDetail := dtos.DocumentReportPDFDetailDTO{
		ID:        report.ID,
		Title:     report.Title,
		CreatedAt: report.CreatedAt.Format("2006-01-02"),
		CreatedBy: report.CreatedBy,
		Documents: report.Documents,
	}

	bodyContent, err := utility.GeneratePDFContent(
		reportDetail,
		[]string{"document/templates/pdf/document_report.html"},
		"",
		"",
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate body HTML content: %w", err)
	}
	headerContent, err := utility.GeneratePDFContent(
		reportDetail,
		[]string{"document/templates/pdf/document_report_header.html"},
		"",
		"",
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}
	footerContent, err := utility.GeneratePDFContent(
		reportDetail,
		[]string{"document/templates/pdf/document_report_footer.html"},
		"",
		"",
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}
	coverContent, err := utility.GeneratePDFContent(
		reportDetail,
		[]string{"document/templates/pdf/document_report_cover.html"},
		"",
		"",
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}

	pdfFileName := fmt.Sprintf("DocumentReport_%d.pdf", reportID)
	if err := utility.GeneratePDFWithCoverHeaderFooter(
		coverContent,
		bodyContent,
		pdfFileName,
		headerContent,
		footerContent,
	); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	zipDir := filepath.Join("storage", "zip")
	zipFilePath := filepath.Join(zipDir, fmt.Sprintf("document_report_%d.zip", reportID))
	if err := os.MkdirAll(zipDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", zipDir, err)
	}
	if err := utility.CreateZipFile(zipFilePath, []string{pdfFileName}); err != nil {
		return "", fmt.Errorf("failed to create ZIP file: %w", err)
	}

	go func() {
		time.Sleep(2 * time.Second)
		if err := os.Remove(pdfFileName); err != nil {
			log.Printf("failed to delete temporary document report PDF %s: %v", pdfFileName, err)
		}
	}()

	return zipFilePath, nil
}

func documentReportSectionTitle(sectionType string) (string, string) {
	switch sectionType {
	case "documents":
		return "Documents", "المستندات"
	case "clauses-controls":
		return "Clauses and Controls", "البنود والضوابط"
	default:
		return sectionType, sectionType
	}
}
