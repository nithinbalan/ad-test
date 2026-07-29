package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentReportServiceInterface interface {
	GetAll(params *dtos.DocumentReportListQueryParamsDTO) (*dtos.DocumentReportListResponseDTO, error)
	GetDocumentReportDetails(reportID uint) (*dtos.DocumentReportDetailDTO, error)
	CreateDocumentReport(currentUserID uint, form *dtos.DocumentReportFormDTO) (uint, error)
	DeleteDocumentReport(reportID uint, currentUserID uint) error
}

type DocumentReportActivityLogger interface {
	LogActivity(action models.Action, tableName string, recordID uint, changedData interface{}, userID uint, moduleIdentifier string, subModuleIdentifier string, section string, tx *gorm.DB) error
}

type DocumentReportService struct {
	repo           repositories.DocumentReportRepositoryInterface
	activityLogSvc DocumentReportActivityLogger
}

func NewDocumentReportService(
	repo repositories.DocumentReportRepositoryInterface,
	activityLogSvc DocumentReportActivityLogger,
) DocumentReportServiceInterface {
	return &DocumentReportService{repo: repo, activityLogSvc: activityLogSvc}
}

func (service *DocumentReportService) GetAll(params *dtos.DocumentReportListQueryParamsDTO) (*dtos.DocumentReportListResponseDTO, error) {
	return service.repo.GetAll(params)
}

func (service *DocumentReportService) GetDocumentReportDetails(reportID uint) (*dtos.DocumentReportDetailDTO, error) {
	report, err := service.repo.GetDocumentReportDetails(reportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, reportID)
		}
		return nil, err
	}
	sections, err := service.repo.GetDocumentReportSections(reportID)
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
		document.PublishedAt = reportDocument.PublishedAt
		document.Clauses = []dtos.DocumentReportClauseDTO{}
		document.Controls = []dtos.DocumentReportControlDTO{}
		document.AssociatedItems = emptyAssociatedItems()

		if err := json.Unmarshal([]byte(reportDocument.ClausesJSON), &document.Clauses); err != nil {
			return nil, fmt.Errorf("failed to decode clauses snapshot for document %d: %w", reportDocument.SourceDocumentID, err)
		}
		if err := json.Unmarshal([]byte(reportDocument.ControlsJSON), &document.Controls); err != nil {
			return nil, fmt.Errorf("failed to decode controls snapshot for document %d: %w", reportDocument.SourceDocumentID, err)
		}
		if err := json.Unmarshal([]byte(reportDocument.AssociatedItemsJSON), &document.AssociatedItems); err != nil {
			return nil, fmt.Errorf("failed to decode associated items snapshot for document %d: %w", reportDocument.SourceDocumentID, err)
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

func (service *DocumentReportService) CreateDocumentReport(currentUserID uint, form *dtos.DocumentReportFormDTO) (uint, error) {
	tx := service.repo.BeginTransaction()
	sourceDocuments, err := service.repo.GetPublishedDocumentsForReporting(form.DocumentCategoryID, form.DocumentTypeID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if len(sourceDocuments) == 0 {
		tx.Rollback()
		return 0, utility.NewValidationError(errors.New("no published documents match the selected category and type"), "not_found", "documents")
	}

	reportID, err := service.repo.Create(&models.DocumentReport{
		Title: form.Title, DocumentCategoryID: form.DocumentCategoryID,
		DocumentTypeID: form.DocumentTypeID, CreatedByID: currentUserID,
	}, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for index, sectionType := range []string{"documents", "clauses-controls", "associated-items"} {
		titleEn, titleAr := documentReportSectionTitle(sectionType)
		if err := service.repo.CreateSection(&models.DocumentReportSection{
			DocumentReportID: reportID, Type: sectionType, TitleEn: titleEn, TitleAr: titleAr, OrderNum: index + 1,
		}, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	snapshots := make([]*models.DocumentReportDocument, 0, len(sourceDocuments))
	for _, source := range sourceDocuments {
		clauses, err := service.repo.GetDocumentClausesForReporting(source.SourceDocumentID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		controls, err := service.repo.GetDocumentControlsForReporting(source.SourceDocumentID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		associated, err := service.repo.GetDocumentAssociatedItemsForReporting(source.SourceDocumentID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		clausesJSON, err := json.Marshal(clauses)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		controlsJSON, err := json.Marshal(controls)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		associatedJSON, err := json.Marshal(associated)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		snapshots = append(snapshots, &models.DocumentReportDocument{
			DocumentReportID: reportID, SourceDocumentID: source.SourceDocumentID,
			ReferenceNumber: source.ReferenceNumber, TitleEn: source.TitleEn, TitleAr: source.TitleAr,
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
			ClausesJSON: string(clausesJSON), ControlsJSON: string(controlsJSON), AssociatedItemsJSON: string(associatedJSON),
			CreatedByID: currentUserID,
		})
	}
	if err := service.repo.CreateDocumentSnapshots(snapshots, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return reportID, nil
}

func (service *DocumentReportService) DeleteDocumentReport(reportID uint, currentUserID uint) error {
	tx := service.repo.BeginTransaction()
	existingReport, err := service.repo.GetDocumentReportByID(reportID, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, reportID)
		}
		return err
	}
	if err := service.repo.DeleteDocumentReport(reportID, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	if service.activityLogSvc != nil {
		go func() {
			if err := service.activityLogSvc.LogActivity(
				models.DeleteAction,
				models.DocumentReport{}.TableName(),
				reportID,
				existingReport,
				currentUserID,
				"", "", "", nil,
			); err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func documentReportSectionTitle(sectionType string) (string, string) {
	switch sectionType {
	case "documents":
		return "Documents", "المستندات"
	case "clauses-controls":
		return "Clauses and Controls", "البنود والضوابط"
	case "associated-items":
		return "Associated Items", "العناصر المرتبطة"
	default:
		return sectionType, sectionType
	}
}

func emptyAssociatedItems() dtos.DocumentReportAssociatedItemsDTO {
	return dtos.DocumentReportAssociatedItemsDTO{
		Risks: []dtos.DocumentReportAssociatedItemDTO{}, Incidents: []dtos.DocumentReportAssociatedItemDTO{},
		Assets: []dtos.DocumentReportAssociatedItemDTO{}, StrategicObjectives: []dtos.DocumentReportAssociatedItemDTO{},
	}
}
