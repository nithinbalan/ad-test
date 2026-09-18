package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"jethur/config"
	logSvc "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	riskDTOs "jethur/grc/risk/dtos"
	settingSvc "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
)

type ProjectSecurityGateReportServiceInterface interface {
	GetFullReportDetailsByID(id uint) (*dtos.ProjectSecurityGateReportDTO, error)
	GetProjectSecurityGateReportExportZip(id uint, currentUserID uint, language string) (string, error)

	CreateSecurityGateExceptionReport(currentUserID uint, form *dtos.ProjectSecurityGateExceptionReportFormDTO) (uint, error)
	GetAllSecurityGateExceptionReports(params *dtos.ProjectSecurityGateReportListQueryParamsDTO) (*dtos.ProjectSecurityGateReportListResponseDTO, error)
	GetSecurityGateExceptionReportDetails(id uint) (*dtos.ProjectSecurityGateExceptionReportDetailsDTO, error)
	ExportSecurityGateExceptionReport(id uint, currentUserID uint, language string) (string, error)
	DeleteSecurityGateExceptionReport(id uint, currentUserID uint) error
	GetSecurityGateExceptionReportSection(reportID uint, sectionID uint) (*dtos.ProjectSecurityGateReportSectionDTO, error)
	UpdateSecurityGateExceptionReportSection(reportID uint, sectionID uint, form *dtos.ProjectSecurityGateReportSectionFormDTO, currentUserID uint) error
	DeleteSecurityGateExceptionReportSection(reportID uint, sectionID uint, currentUserID uint) error

	CreateSecurityGateReport(currentUserID uint, form *dtos.ProjectSecurityGateReportFormDTO) (uint, error)
	GetAllSecurityGateReports(params *dtos.ProjectSecurityGateReportListQueryParamsDTO) (*dtos.ProjectSecurityGateReportListResponseDTO, error)
	GetSecurityGateReportDetails(id uint) (*dtos.ProjectSecurityGateReportDetailsDTO, error)
	ExportSecurityGateReport(id uint, currentUserID uint, language string) (string, error)
	DeleteSecurityGateReport(id uint, currentUserID uint) error
	GetSecurityGateReportSection(reportID uint, sectionID uint) (*dtos.ProjectSecurityGateReportSectionDTO, error)
	UpdateSecurityGateReportSection(reportID uint, sectionID uint, form *dtos.ProjectSecurityGateReportSectionFormDTO, currentUserID uint) error
	DeleteSecurityGateReportSection(reportID uint, sectionID uint, currentUserID uint) error
}

type ProjectSecurityGateReportService struct {
	repo           repositories.ProjectSecurityGateReportRepositoryInterface
	checklistSvc   ProjectSecurityGateChecklistServiceInterface
	settingSvc     settingSvc.ReportTemplateSettingServiceInterface
	activityLogSvc logSvc.ActivityLogServiceInterface
}

func NewProjectSecurityGateReportService(
	repo repositories.ProjectSecurityGateReportRepositoryInterface,
	checklistSvc ProjectSecurityGateChecklistServiceInterface,
	settingSvc settingSvc.ReportTemplateSettingServiceInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
) ProjectSecurityGateReportServiceInterface {
	return &ProjectSecurityGateReportService{
		repo:           repo,
		checklistSvc:   checklistSvc,
		settingSvc:     settingSvc,
		activityLogSvc: activityLogSvc,
	}
}

func (psgrs *ProjectSecurityGateReportService) GetFullReportDetailsByID(id uint) (*dtos.ProjectSecurityGateReportDTO, error) {
	return psgrs.buildGateReportDTO(id, nil)
}

func (psgrs *ProjectSecurityGateReportService) GetProjectSecurityGateReportExportZip(id uint, currentUserID uint, language string) (string, error) {
	reportDTO, err := psgrs.GetFullReportDetailsByID(id)
	if err != nil {
		return "", err
	}

	reportDTO.GeneratedAt = time.Now().Format("2006-01-02 03:04 PM")

	bodyTemplate := "project/templates/pdf/project_security_gate_report.html"
	coverTemplate := "project/templates/pdf/project_security_gate_report_cover.html"
	footerTemplate := "project/templates/pdf/project_security_gate_report_footer.html"
	if language == "ar" {
		bodyTemplate = "project/templates/pdf/project_security_gate_report_ar.html"
		coverTemplate = "project/templates/pdf/project_security_gate_report_cover_ar.html"
		footerTemplate = "project/templates/pdf/project_security_gate_report_footer_ar.html"
	}

	htmlContent, err := utility.GeneratePDFContent(reportDTO, []string{bodyTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	headerTemplatePaths := []string{"project/templates/pdf/project_security_gate_report_header.html"}
	headerContent, err := utility.GeneratePDFContent(reportDTO, headerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	footerContent, err := utility.GeneratePDFContent(reportDTO, []string{footerTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	coverContent, err := utility.GeneratePDFContent(reportDTO, []string{coverTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}

	pdfFileName := fmt.Sprintf("ProjectSecurityGateReport_%s_%d.pdf", language, id)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	entries := []utility.ZipEntry{{SourcePath: pdfFileName, ZipPath: pdfFileName}}
	usedZipPaths := map[string]int{}
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")

	responseDocuments, err := psgrs.repo.GetAllResponseDocuments(id)
	if err != nil {
		return "", err
	}
	for _, document := range responseDocuments {
		if document.Name == "" {
			continue
		}

		fullPath := filepath.Join(baseStoragePath, "project", "security-gates",
			strconv.FormatUint(uint64(id), 10),
			"checklists", strconv.FormatUint(uint64(document.ChecklistID), 10),
			"documents", document.Name)
		if _, err := os.Stat(fullPath); err != nil {
			continue
		}

		checklistFolder := "Unassigned Checklist"
		if document.Checklist != nil && document.Checklist.Requirement != nil {
			checklistFolder = document.Checklist.Requirement.TitleEn
		}

		fileName := document.OriginalName
		if fileName == "" {
			fileName = document.Name
		}

		zipPath := path.Join("Response Documents", utility.SanitizeExportFileName(checklistFolder), utility.SanitizeExportFileName(fileName))
		usedZipPaths[zipPath]++
		if count := usedZipPaths[zipPath]; count > 1 {
			extension := path.Ext(zipPath)
			zipPath = fmt.Sprintf("%s (%d)%s", strings.TrimSuffix(zipPath, extension), count, extension)
		}

		entries = append(entries, utility.ZipEntry{SourcePath: fullPath, ZipPath: zipPath})
	}

	commentDocuments, err := psgrs.repo.GetAllWorkflowCommentDocuments(id)
	if err != nil {
		return "", err
	}
	for _, document := range commentDocuments {
		if document.Name == "" || document.Comment == nil {
			continue
		}

		fullPath := filepath.Join(baseStoragePath, "project", "security-gates",
			strconv.FormatUint(uint64(id), 10),
			"checklists", strconv.FormatUint(uint64(document.Comment.ChecklistID), 10),
			"comments", strconv.FormatUint(uint64(document.ProjectSecurityGateChecklistCommentID), 10),
			"documents", document.Name)
		if _, err := os.Stat(fullPath); err != nil {
			continue
		}

		checklistFolder := "Unassigned Checklist"
		if document.Comment.Checklist != nil && document.Comment.Checklist.Requirement != nil {
			checklistFolder = document.Comment.Checklist.Requirement.TitleEn
		}

		fileName := document.OriginalName
		if fileName == "" {
			fileName = document.Name
		}

		zipPath := path.Join("Workflow Documents", utility.SanitizeExportFileName(checklistFolder), utility.SanitizeExportFileName(fileName))
		usedZipPaths[zipPath]++
		if count := usedZipPaths[zipPath]; count > 1 {
			extension := path.Ext(zipPath)
			zipPath = fmt.Sprintf("%s (%d)%s", strings.TrimSuffix(zipPath, extension), count, extension)
		}

		entries = append(entries, utility.ZipEntry{SourcePath: fullPath, ZipPath: zipPath})
	}

	zipDir := filepath.Join("storage", "zip")
	if err := os.MkdirAll(zipDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", zipDir, err)
	}

	zipFilePath := filepath.Join(zipDir, fmt.Sprintf("project_security_gate_report_%s_%d.zip", language, id))
	if err := utility.CreateZipFileWithFolders(zipFilePath, entries); err != nil {
		return "", fmt.Errorf("failed to create ZIP file: %w", err)
	}

	go func() {
		time.Sleep(2 * time.Second)
		_ = os.Remove(pdfFileName)
	}()

	return zipFilePath, nil
}

func (psgrs *ProjectSecurityGateReportService) CreateSecurityGateExceptionReport(currentUserID uint, form *dtos.ProjectSecurityGateExceptionReportFormDTO) (uint, error) {
	securityGates, err := psgrs.getSelectedGates(form.SecurityGateIDs)
	if err != nil {
		return 0, err
	}

	reportData, err := psgrs.buildExceptionReportData(currentUserID, form.Title, securityGates)
	if err != nil {
		return 0, err
	}

	payload, err := json.Marshal(reportData)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal report data: %w", err)
	}

	report := &models.ProjectSecurityGateExceptionReport{
		Title:           form.Title,
		SecurityGateIDs: psgrs.joinIDList(form.SecurityGateIDs),
		ReportData:      string(payload),
		CreatedBYID:     currentUserID,
		CreatedAt:       time.Now(),
	}

	tx := psgrs.repo.BeginTransaction()

	reportID, err := psgrs.repo.CreateSecurityGateExceptionReport(report, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for index, sectionType := range projectSecurityGateReportSectionTypes {
		titleEn, titleAr := formatSecurityGateReportSectionTitle(sectionType)
		section := &models.ProjectSecurityGateExceptionReportSection{
			ReportID: reportID,
			Type:     sectionType,
			TitleEn:  titleEn,
			TitleAr:  titleAr,
			OrderNum: index + 1,
		}

		if err := psgrs.repo.CreateSecurityGateExceptionReportSection(section, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := psgrs.activityLogSvc.LogActivity(models.CreateAction, models.ProjectSecurityGateExceptionReport{}.TableName(), reportID, map[string]interface{}{
			"title":             report.Title,
			"security_gate_ids": report.SecurityGateIDs,
		}, currentUserID, "project", "project_report", "security_gate_exception_report", nil)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return reportID, nil
}

func (psgrs *ProjectSecurityGateReportService) GetAllSecurityGateExceptionReports(params *dtos.ProjectSecurityGateReportListQueryParamsDTO) (*dtos.ProjectSecurityGateReportListResponseDTO, error) {
	return psgrs.repo.GetAllSecurityGateExceptionReports(params)
}

func (psgrs *ProjectSecurityGateReportService) GetSecurityGateExceptionReportDetails(id uint) (*dtos.ProjectSecurityGateExceptionReportDetailsDTO, error) {
	report, reportData, err := psgrs.getExceptionReportWithData(id)
	if err != nil {
		return nil, err
	}

	return &dtos.ProjectSecurityGateExceptionReportDetailsDTO{
		ID:            report.ID,
		CreatedAt:     report.CreatedAt,
		CreatedByUser: psgrs.mapReportUser(report.CreatedByUser),
		ProjectSecurityGateExceptionReportDataDTO: reportData,
	}, nil
}

func (psgrs *ProjectSecurityGateReportService) ExportSecurityGateExceptionReport(id uint, currentUserID uint, language string) (string, error) {
	_, reportData, err := psgrs.getExceptionReportWithData(id)
	if err != nil {
		return "", err
	}

	if err := psgrs.validateReportTemplateImages(); err != nil {
		return "", err
	}

	// the payload is frozen, only the export moment is fresh
	reportData.ExportedAt = time.Now().Format("2006-01-02 03:04 PM")

	bodyTemplate := "project/templates/pdf/project_security_gate_exception_report.html"
	coverTemplate := "project/templates/pdf/project_security_gate_exception_report_cover.html"
	footerTemplate := "project/templates/pdf/project_security_gate_report_footer.html"
	if language == "ar" {
		bodyTemplate = "project/templates/pdf/project_security_gate_exception_report_ar.html"
		coverTemplate = "project/templates/pdf/project_security_gate_exception_report_cover_ar.html"
		footerTemplate = "project/templates/pdf/project_security_gate_report_footer_ar.html"
	}

	return psgrs.generateReportPDF(reportData, bodyTemplate, coverTemplate, footerTemplate,
		fmt.Sprintf("ProjectSecurityGateExceptionReport_%s_%d.pdf", language, id))
}

func (psgrs *ProjectSecurityGateReportService) DeleteSecurityGateExceptionReport(id uint, currentUserID uint) error {
	report, err := psgrs.repo.GetSecurityGateExceptionReportByID(id)
	if err != nil {
		return err
	}

	if err := psgrs.repo.DeleteSecurityGateExceptionReport(id); err != nil {
		return err
	}

	go func() {
		err := psgrs.activityLogSvc.LogActivity(models.DeleteAction, models.ProjectSecurityGateExceptionReport{}.TableName(), id, map[string]interface{}{
			"title":             report.Title,
			"security_gate_ids": report.SecurityGateIDs,
		}, currentUserID, "project", "project_report", "security_gate_exception_report", nil)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return nil
}

func (psgrs *ProjectSecurityGateReportService) CreateSecurityGateReport(currentUserID uint, form *dtos.ProjectSecurityGateReportFormDTO) (uint, error) {
	securityGates, err := psgrs.getSelectedGates(form.SecurityGateIDs)
	if err != nil {
		return 0, err
	}

	reportData, err := psgrs.buildReportData(currentUserID, form.Title, securityGates, form.ResponseStatusIDs)
	if err != nil {
		return 0, err
	}

	payload, err := json.Marshal(reportData)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal report data: %w", err)
	}

	report := &models.ProjectSecurityGateReport{
		Title:             form.Title,
		SecurityGateIDs:   psgrs.joinIDList(form.SecurityGateIDs),
		ResponseStatusIDs: psgrs.joinIDList(form.ResponseStatusIDs),
		ReportData:        string(payload),
		CreatedBYID:       currentUserID,
		CreatedAt:         time.Now(),
	}

	tx := psgrs.repo.BeginTransaction()

	reportID, err := psgrs.repo.CreateSecurityGateReport(report, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for index, sectionType := range projectSecurityGateReportSectionTypes {
		titleEn, titleAr := formatSecurityGateReportSectionTitle(sectionType)
		section := &models.ProjectSecurityGateReportSection{
			ReportID: reportID,
			Type:     sectionType,
			TitleEn:  titleEn,
			TitleAr:  titleAr,
			OrderNum: index + 1,
		}

		if err := psgrs.repo.CreateSecurityGateReportSection(section, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := psgrs.activityLogSvc.LogActivity(models.CreateAction, models.ProjectSecurityGateReport{}.TableName(), reportID, map[string]interface{}{
			"title":               report.Title,
			"security_gate_ids":   report.SecurityGateIDs,
			"response_status_ids": report.ResponseStatusIDs,
		}, currentUserID, "project", "project_report", "security_gate_report", nil)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return reportID, nil
}

func (psgrs *ProjectSecurityGateReportService) GetAllSecurityGateReports(params *dtos.ProjectSecurityGateReportListQueryParamsDTO) (*dtos.ProjectSecurityGateReportListResponseDTO, error) {
	return psgrs.repo.GetAllSecurityGateReports(params)
}

func (psgrs *ProjectSecurityGateReportService) GetSecurityGateReportDetails(id uint) (*dtos.ProjectSecurityGateReportDetailsDTO, error) {
	report, reportData, err := psgrs.getReportWithData(id)
	if err != nil {
		return nil, err
	}

	details := &dtos.ProjectSecurityGateReportDetailsDTO{
		ID:                               report.ID,
		CreatedAt:                        report.CreatedAt,
		CreatedByUser:                    psgrs.mapReportUser(report.CreatedByUser),
		ResponseStatuses:                 []*dtos.ProjectSecurityGateChecklistResponseStatusDTO{},
		ProjectSecurityGateReportDataDTO: reportData,
	}

	responseStatusIDs := psgrs.parseIDList(report.ResponseStatusIDs)
	if len(responseStatusIDs) > 0 {
		responseStatuses, err := psgrs.repo.GetResponseStatusesByIDs(responseStatusIDs)
		if err != nil {
			return nil, err
		}

		for _, responseStatus := range responseStatuses {
			var responseStatusDTO dtos.ProjectSecurityGateChecklistResponseStatusDTO
			if err := mapstructure.Decode(responseStatus, &responseStatusDTO); err != nil {
				return nil, err
			}
			details.ResponseStatuses = append(details.ResponseStatuses, &responseStatusDTO)
		}
	}

	return details, nil
}

func (psgrs *ProjectSecurityGateReportService) ExportSecurityGateReport(id uint, currentUserID uint, language string) (string, error) {
	_, reportData, err := psgrs.getReportWithData(id)
	if err != nil {
		return "", err
	}

	if err := psgrs.validateReportTemplateImages(); err != nil {
		return "", err
	}

	// the payload is frozen, only the export moment is fresh
	reportData.ExportedAt = time.Now().Format("2006-01-02 03:04 PM")

	bodyTemplate := "project/templates/pdf/project_security_gate_multi_report.html"
	coverTemplate := "project/templates/pdf/project_security_gate_multi_report_cover.html"
	footerTemplate := "project/templates/pdf/project_security_gate_report_footer.html"
	if language == "ar" {
		bodyTemplate = "project/templates/pdf/project_security_gate_multi_report_ar.html"
		coverTemplate = "project/templates/pdf/project_security_gate_multi_report_cover_ar.html"
		footerTemplate = "project/templates/pdf/project_security_gate_report_footer_ar.html"
	}

	return psgrs.generateReportPDF(reportData, bodyTemplate, coverTemplate, footerTemplate,
		fmt.Sprintf("ProjectSecurityGateReport_%s_%d.pdf", language, id))
}

func (psgrs *ProjectSecurityGateReportService) DeleteSecurityGateReport(id uint, currentUserID uint) error {
	report, err := psgrs.repo.GetSecurityGateReportByID(id)
	if err != nil {
		return err
	}

	if err := psgrs.repo.DeleteSecurityGateReport(id); err != nil {
		return err
	}

	go func() {
		err := psgrs.activityLogSvc.LogActivity(models.DeleteAction, models.ProjectSecurityGateReport{}.TableName(), id, map[string]interface{}{
			"title":               report.Title,
			"security_gate_ids":   report.SecurityGateIDs,
			"response_status_ids": report.ResponseStatusIDs,
		}, currentUserID, "project", "project_report", "security_gate_report", nil)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return nil
}

func (psgrs *ProjectSecurityGateReportService) buildGateReportDTO(id uint, responseStatusIDs []uint) (*dtos.ProjectSecurityGateReportDTO, error) {
	securityGate, err := psgrs.repo.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	if _, err := psgrs.settingSvc.GetImagePath("logo"); err != nil {
		return nil, utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := psgrs.settingSvc.GetImagePath("cover-page-bg"); err != nil {
		return nil, utility.NewForbiddenError("cover_image_missing")
	}

	var report dtos.ProjectSecurityGateReportDTO
	if err := mapstructure.Decode(securityGate, &report); err != nil {
		return nil, err
	}

	report.CreatedAtFormatted = securityGate.CreatedAt.Format("01/02/2006")
	if securityGate.UpdatedAt != nil {
		report.UpdatedAtFormatted = securityGate.UpdatedAt.Format("01/02/2006")
	}

	report.Classifications = make([]*dtos.ProjectSecurityGateClassificationDTO, 0, len(securityGate.Classifications))
	for _, mapping := range securityGate.Classifications {
		if mapping.Classification == nil {
			continue
		}

		var classificationDTO dtos.ProjectSecurityGateClassificationDTO
		if err := mapstructure.Decode(mapping.Classification, &classificationDTO); err != nil {
			return nil, err
		}
		report.Classifications = append(report.Classifications, &classificationDTO)
	}

	summary, err := psgrs.checklistSvc.GetChecklistSummary(id)
	if err != nil {
		return nil, err
	}
	report.Summary = summary

	histories, err := psgrs.repo.GetAllExceptionWorkflowHistories(id)
	if err != nil {
		return nil, err
	}

	reviewsByChecklistID := make(map[uint][]*dtos.ProjectSecurityGateReportReviewDTO)
	for _, history := range histories {
		historyCreatedAt := history.CreatedAt
		if reviewsByChecklistID[history.ChecklistID] == nil {
			reviewsByChecklistID[history.ChecklistID] = []*dtos.ProjectSecurityGateReportReviewDTO{}
		}

		review := psgrs.buildExceptionReviewItem(history.Level, &historyCreatedAt, history.WorkflowStatus, history.CreatedByUser)
		reviewsByChecklistID[history.ChecklistID] = append(reviewsByChecklistID[history.ChecklistID], review)
	}

	var checklists []*models.ProjectSecurityGateChecklist
	if len(responseStatusIDs) > 0 {
		checklists, err = psgrs.repo.GetAllChecklistsByResponseStatusIDs(id, responseStatusIDs)
	} else {
		checklists, err = psgrs.repo.GetAllChecklists(id)
	}
	if err != nil {
		return nil, err
	}

	report.Checklists = make([]*dtos.ProjectSecurityGateReportChecklistDTO, 0, len(checklists))
	report.Exceptions = make([]*dtos.ProjectSecurityGateReportExceptionDTO, 0)
	report.Controls = make([]*dtos.ProjectSecurityGateReportControlDTO, 0)
	report.Risks = make([]*dtos.ProjectSecurityGateReportRiskDTO, 0)

	for _, checklist := range checklists {
		checklistTitleEn := ""
		checklistTitleAr := ""
		checklistItem := &dtos.ProjectSecurityGateReportChecklistDTO{
			ChecklistID: checklist.ID,
			IsException: checklist.IsException,
		}

		if checklist.Requirement != nil {
			checklistTitleEn = checklist.Requirement.TitleEn
			checklistTitleAr = checklist.Requirement.TitleAr
			checklistItem.TitleEn = checklist.Requirement.TitleEn
			checklistItem.TitleAr = checklist.Requirement.TitleAr
			checklistItem.IsCritical = checklist.Requirement.IsCritical
		}
		if checklist.ResponseStatus != nil {
			checklistItem.ResponseStatus = &dtos.ProjectSecurityGateChecklistResponseStatusDTO{
				ID:        int(checklist.ResponseStatus.ID),
				TitleEn:   checklist.ResponseStatus.TitleEn,
				TitleAr:   checklist.ResponseStatus.TitleAr,
				Type:      checklist.ResponseStatus.Type,
				ColorCode: checklist.ResponseStatus.ColorCode,
				Status:    checklist.ResponseStatus.Status,
			}
		}
		if checklist.Comment != nil {
			checklistItem.Comment = *checklist.Comment
		}

		if checklist.ResponseStatusID != nil || checklist.IsException {
			if checklist.UpdatedByUser != nil {
				checklistItem.RespondedBy = psgrs.mapReportUser(checklist.UpdatedByUser)
			} else if checklist.CreatedByUser != nil {
				checklistItem.RespondedBy = psgrs.mapReportUser(checklist.CreatedByUser)
			}

			if checklist.UpdatedAt != nil {
				checklistItem.RespondedAtFormatted = checklist.UpdatedAt.Format("01/02/2006")
			} else if checklist.CreatedAt != nil {
				checklistItem.RespondedAtFormatted = checklist.CreatedAt.Format("01/02/2006")
			}
		}
		report.Checklists = append(report.Checklists, checklistItem)

		if checklist.IsException {
			reviews := reviewsByChecklistID[checklist.ID]
			if reviews == nil {
				reviews = []*dtos.ProjectSecurityGateReportReviewDTO{}
			}

			exception := &dtos.ProjectSecurityGateReportExceptionDTO{
				ChecklistID: checklist.ID,
				TitleEn:     checklistTitleEn,
				TitleAr:     checklistTitleAr,
				IsCritical:  checklistItem.IsCritical,
				Reviews:     reviews,
			}

			if checklist.ExceptionJustification != nil {
				exception.Justification = *checklist.ExceptionJustification
			}
			if checklist.WorkflowStatus != nil {
				exception.WorkflowStatusTitleEn = checklist.WorkflowStatus.TitleEn
				exception.WorkflowStatusTitleAr = checklist.WorkflowStatus.TitleAr
				exception.WorkflowStatusColorCode = checklist.WorkflowStatus.ColorCode
				exception.WorkflowStatusType = checklist.WorkflowStatus.Type
				textColor, borderColor, bgColor := utility.ResolveStatusBadgeColors(checklist.WorkflowStatus.ColorCode, checklist.WorkflowStatus.Type, checklist.WorkflowStatus.TitleEn)
				exception.WorkflowStatus = &dtos.WorkflowStatusDTO{
					TitleEn:         checklist.WorkflowStatus.TitleEn,
					TitleAr:         checklist.WorkflowStatus.TitleAr,
					Type:            checklist.WorkflowStatus.Type,
					ColorCode:       checklist.WorkflowStatus.ColorCode,
					TextColor:       textColor,
					BorderColor:     borderColor,
					BackgroundColor: bgColor,
				}
			}
			if checklist.SubmittedByUser != nil {
				exception.SubmittedByEn = strings.TrimSpace(checklist.SubmittedByUser.FirstNameEn + " " + checklist.SubmittedByUser.LastNameEn)
				exception.SubmittedByAr = strings.TrimSpace(checklist.SubmittedByUser.FirstNameAr + " " + checklist.SubmittedByUser.LastNameAr)
				exception.SubmittedBy = psgrs.mapReportUser(checklist.SubmittedByUser)
			}
			if checklist.SubmittedAt != nil {
				exception.SubmittedAtFormatted = checklist.SubmittedAt.Format("01/02/2006")
				exception.SubmittedAt = checklist.SubmittedAt
			}
			if checklist.ApprovedByUser != nil {
				exception.ApprovedByEn = strings.TrimSpace(checklist.ApprovedByUser.FirstNameEn + " " + checklist.ApprovedByUser.LastNameEn)
				exception.ApprovedByAr = strings.TrimSpace(checklist.ApprovedByUser.FirstNameAr + " " + checklist.ApprovedByUser.LastNameAr)
				exception.ApprovedBy = psgrs.mapReportUser(checklist.ApprovedByUser)
			}
			if exception.SubmittedByAr == "" {
				exception.SubmittedByAr = exception.SubmittedByEn
			}
			if exception.ApprovedByAr == "" {
				exception.ApprovedByAr = exception.ApprovedByEn
			}
			if checklist.ApprovedAt != nil {
				exception.ApprovedAtFormatted = checklist.ApprovedAt.Format("01/02/2006")
				exception.ApprovedAt = checklist.ApprovedAt
			}

			report.Exceptions = append(report.Exceptions, exception)
		}

		for _, mapping := range checklist.Controls {
			if mapping.Control == nil {
				continue
			}

			control := &dtos.ProjectSecurityGateReportControlDTO{
				ControlID:        mapping.Control.ID,
				ChecklistTitleEn: checklistTitleEn,
				ChecklistTitleAr: checklistTitleAr,
				ReferenceNumber:  mapping.Control.ReferenceNumber,
				TitleEn:          mapping.Control.TitleEn,
				TitleAr:          mapping.Control.TitleAr,
			}
			if mapping.Control.Code != nil {
				control.Code = *mapping.Control.Code
			}
			if mapping.Control.Description != nil {
				control.Description = utility.StripHTMLTags(*mapping.Control.Description)
			}
			report.Controls = append(report.Controls, control)
		}

		for _, mapping := range checklist.Risks {
			if mapping.RiskRegister == nil {
				continue
			}

			risk := &dtos.ProjectSecurityGateReportRiskDTO{
				RiskID:           mapping.RiskRegister.ID,
				ChecklistTitleEn: checklistTitleEn,
				ChecklistTitleAr: checklistTitleAr,
				Title:            mapping.RiskRegister.Title,
			}
			if mapping.RiskRegister.ReferenceNumber != nil {
				risk.ReferenceNumber = *mapping.RiskRegister.ReferenceNumber
			}
			if mapping.RiskRegister.RiskCode != nil {
				risk.RiskCode = *mapping.RiskRegister.RiskCode
			}
			if mapping.RiskRegister.RiskStatus != nil {
				risk.RiskStatus = &riskDTOs.RiskStatusDTO{
					ID:            int(mapping.RiskRegister.RiskStatus.ID),
					TitleEn:       mapping.RiskRegister.RiskStatus.TitleEn,
					TitleAr:       mapping.RiskRegister.RiskStatus.TitleAr,
					DescriptionEn: mapping.RiskRegister.RiskStatus.DescriptionEn,
					DescriptionAr: mapping.RiskRegister.RiskStatus.DescriptionAr,
					ColorCode:     mapping.RiskRegister.RiskStatus.ColorCode,
					Type:          mapping.RiskRegister.RiskStatus.Type,
					Status:        mapping.RiskRegister.RiskStatus.Status,
				}
			}
			if mapping.RiskRegister.RiskCategory != nil {
				risk.RiskCategory = &riskDTOs.RiskCategoryDTO{
					ID:      int(mapping.RiskRegister.RiskCategory.ID),
					TitleEn: mapping.RiskRegister.RiskCategory.TitleEn,
					TitleAr: mapping.RiskRegister.RiskCategory.TitleAr,
					Status:  mapping.RiskRegister.RiskCategory.Status,
				}
			}
			if mapping.RiskRegister.Division != nil {
				risk.Division = &organizationDTOs.OrganizationUnitDTO{
					ID:      mapping.RiskRegister.Division.ID,
					TitleEn: mapping.RiskRegister.Division.TitleEn,
					TitleAr: mapping.RiskRegister.Division.TitleAr,
				}
			}
			if mapping.RiskRegister.ResponsibleUnit != nil {
				risk.ResponsibleUnit = &organizationDTOs.OrganizationUnitDTO{
					ID:      mapping.RiskRegister.ResponsibleUnit.ID,
					TitleEn: mapping.RiskRegister.ResponsibleUnit.TitleEn,
					TitleAr: mapping.RiskRegister.ResponsibleUnit.TitleAr,
				}
			}
			risk.ResponsibleUser = psgrs.mapReportUser(mapping.RiskRegister.ResponsibleUser)
			if mapping.RiskRegister.RegistrationDate != nil {
				risk.RegistrationDateFormatted = mapping.RiskRegister.RegistrationDate.Format("01/02/2006")
			}
			report.Risks = append(report.Risks, risk)
		}

	}

	return &report, nil
}

func (psgrs *ProjectSecurityGateReportService) buildReportData(currentUserID uint, title string, securityGates []*models.ProjectSecurityGate, responseStatusIDs []uint) (*dtos.ProjectSecurityGateReportDataDTO, error) {
	generatedByEn, generatedByAr, err := psgrs.getGeneratedByNames(currentUserID)
	if err != nil {
		return nil, err
	}

	reportData := &dtos.ProjectSecurityGateReportDataDTO{
		Title:         title,
		GeneratedByEn: generatedByEn,
		GeneratedByAr: generatedByAr,
		GeneratedAt:   time.Now().Format("2006-01-02 03:04 PM"),
		Projects:      make([]*dtos.ProjectSecurityGateReportProjectDTO, 0, len(securityGates)),
		SecurityGates: make([]*dtos.ProjectSecurityGateReportDTO, 0, len(securityGates)),
	}

	for _, securityGate := range securityGates {
		gateDTO, err := psgrs.buildGateReportDTO(securityGate.ID, responseStatusIDs)
		if err != nil {
			return nil, err
		}

		reportData.Projects = append(reportData.Projects, psgrs.buildReportProject(securityGate.ID, gateDTO))
		reportData.SecurityGates = append(reportData.SecurityGates, gateDTO)
	}

	return reportData, nil
}

func (psgrs *ProjectSecurityGateReportService) buildExceptionReportData(currentUserID uint, title string, securityGates []*models.ProjectSecurityGate) (*dtos.ProjectSecurityGateExceptionReportDataDTO, error) {
	generatedByEn, generatedByAr, err := psgrs.getGeneratedByNames(currentUserID)
	if err != nil {
		return nil, err
	}

	reportData := &dtos.ProjectSecurityGateExceptionReportDataDTO{
		Title:         title,
		GeneratedByEn: generatedByEn,
		GeneratedByAr: generatedByAr,
		GeneratedAt:   time.Now().Format("2006-01-02 03:04 PM"),
		Projects:      make([]*dtos.ProjectSecurityGateReportProjectDTO, 0, len(securityGates)),
		SecurityGates: make([]*dtos.ProjectSecurityGateExceptionReportDTO, 0, len(securityGates)),
		TotalCount:    0,
		ApprovedCount: 0,
		PendingCount:  0,
		CriticalCount: 0,
	}

	for _, securityGate := range securityGates {
		gateDTO, err := psgrs.buildGateReportDTO(securityGate.ID, nil)
		if err != nil {
			return nil, err
		}

		project := psgrs.buildReportProject(securityGate.ID, gateDTO)
		reportData.Projects = append(reportData.Projects, project)

		exceptionGateDTO := &dtos.ProjectSecurityGateExceptionReportDTO{
			ProjectTitleEn:     gateDTO.ProjectTitleEn,
			ProjectTitleAr:     gateDTO.ProjectTitleAr,
			Stage:              gateDTO.Stage,
			TemplateVersion:    gateDTO.TemplateVersion,
			BusinessOwner:      gateDTO.BusinessOwner,
			ReadinessStatus:    gateDTO.ReadinessStatus,
			Notes:              gateDTO.Notes,
			Classifications:    gateDTO.Classifications,
			CreatedByUser:      gateDTO.CreatedByUser,
			CreatedAtFormatted: gateDTO.CreatedAtFormatted,
			UpdatedByUser:      gateDTO.UpdatedByUser,
			UpdatedAtFormatted: gateDTO.UpdatedAtFormatted,
			Exceptions:         make([]*dtos.ProjectSecurityGateExceptionItemDTO, 0, len(gateDTO.Exceptions)),
		}

		for _, exception := range gateDTO.Exceptions {
			item := &dtos.ProjectSecurityGateExceptionItemDTO{
				ChecklistID:    exception.ChecklistID,
				TitleEn:        exception.TitleEn,
				TitleAr:        exception.TitleAr,
				IsCritical:     exception.IsCritical,
				Justification:  exception.Justification,
				WorkflowStatus: exception.WorkflowStatus,
				SubmittedBy:    exception.SubmittedBy,
				SubmittedAt:    exception.SubmittedAt,
				ApprovedBy:     exception.ApprovedBy,
				ApprovedAt:     exception.ApprovedAt,
				Reviews:        exception.Reviews,
			}

			exceptionGateDTO.Exceptions = append(exceptionGateDTO.Exceptions, item)

			if item.ApprovedBy != nil {
				reportData.ApprovedCount++
			} else {
				reportData.PendingCount++
			}
			if item.IsCritical {
				reportData.CriticalCount++
			}
			reportData.TotalCount++
		}

		reportData.SecurityGates = append(reportData.SecurityGates, exceptionGateDTO)
	}

	return reportData, nil
}

func (psgrs *ProjectSecurityGateReportService) getReportWithData(id uint) (*models.ProjectSecurityGateReport, *dtos.ProjectSecurityGateReportDataDTO, error) {
	report, err := psgrs.repo.GetSecurityGateReportByID(id)
	if err != nil {
		return nil, nil, err
	}

	payload, err := psgrs.applyRenamedReportKeys([]byte(report.ReportData))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to migrate report data of report %d: %w", id, err)
	}

	reportData := &dtos.ProjectSecurityGateReportDataDTO{}
	if err := json.Unmarshal(payload, reportData); err != nil {
		return nil, nil, fmt.Errorf("failed to decode report data of report %d: %w", id, err)
	}

	// sections stay editable, so they are read live instead of coming from the payload
	sections, err := psgrs.repo.GetSecurityGateReportSections(id)
	if err != nil {
		return nil, nil, err
	}

	reportData.Sections = make([]*dtos.ProjectSecurityGateReportSectionDTO, 0, len(sections))
	for _, section := range sections {
		reportData.Sections = append(reportData.Sections,
			mapSecurityGateReportSection(section.ID, section.Type, section.TitleEn, section.TitleAr, section.DescriptionEn, section.DescriptionAr, section.OrderNum))
	}

	return report, reportData, nil
}

func (psgrs *ProjectSecurityGateReportService) getExceptionReportWithData(id uint) (*models.ProjectSecurityGateExceptionReport, *dtos.ProjectSecurityGateExceptionReportDataDTO, error) {
	report, err := psgrs.repo.GetSecurityGateExceptionReportByID(id)
	if err != nil {
		return nil, nil, err
	}

	payload, err := psgrs.applyRenamedReportKeys([]byte(report.ReportData))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to migrate report data of report %d: %w", id, err)
	}

	reportData := &dtos.ProjectSecurityGateExceptionReportDataDTO{}
	if err := json.Unmarshal(payload, reportData); err != nil {
		return nil, nil, fmt.Errorf("failed to decode report data of report %d: %w", id, err)
	}

	// sections stay editable, so they are read live instead of coming from the payload
	sections, err := psgrs.repo.GetSecurityGateExceptionReportSections(id)
	if err != nil {
		return nil, nil, err
	}

	reportData.Sections = make([]*dtos.ProjectSecurityGateReportSectionDTO, 0, len(sections))
	for _, section := range sections {
		reportData.Sections = append(reportData.Sections,
			mapSecurityGateReportSection(section.ID, section.Type, section.TitleEn, section.TitleAr, section.DescriptionEn, section.DescriptionAr, section.OrderNum))
	}

	return report, reportData, nil
}

func (psgrs *ProjectSecurityGateReportService) buildReportProject(securityGateID uint, gateDTO *dtos.ProjectSecurityGateReportDTO) *dtos.ProjectSecurityGateReportProjectDTO {
	project := &dtos.ProjectSecurityGateReportProjectDTO{
		ID:             securityGateID,
		ProjectTitleEn: gateDTO.ProjectTitleEn,
		ProjectTitleAr: gateDTO.ProjectTitleAr,
	}

	if gateDTO.Stage != nil {
		project.StageTitleEn = gateDTO.Stage.TitleEn
		project.StageTitleAr = gateDTO.Stage.TitleAr
	}

	return project
}

func (psgrs *ProjectSecurityGateReportService) getGeneratedByNames(currentUserID uint) (string, string, error) {
	user, err := psgrs.repo.GetUserByID(currentUserID)
	if err != nil {
		return "", "", err
	}

	nameEn := strings.TrimSpace(user.FirstNameEn + " " + user.LastNameEn)
	nameAr := strings.TrimSpace(user.FirstNameAr + " " + user.LastNameAr)
	if nameAr == "" {
		nameAr = nameEn
	}

	return nameEn, nameAr, nil
}

var reportSnapshotRenamedKeys = map[string]string{}

func (psgrs *ProjectSecurityGateReportService) applyRenamedReportKeys(payload []byte) ([]byte, error) {
	if len(reportSnapshotRenamedKeys) == 0 {
		return payload, nil
	}

	var tree interface{}
	if err := json.Unmarshal(payload, &tree); err != nil {
		return nil, err
	}

	return json.Marshal(psgrs.renameReportKeys(tree))
}

func (psgrs *ProjectSecurityGateReportService) renameReportKeys(node interface{}) interface{} {
	switch value := node.(type) {
	case map[string]interface{}:
		renamed := make(map[string]interface{}, len(value))
		for key, item := range value {
			if currentKey, exists := reportSnapshotRenamedKeys[key]; exists {
				key = currentKey
			}
			renamed[key] = psgrs.renameReportKeys(item)
		}

		return renamed
	case []interface{}:
		for index, item := range value {
			value[index] = psgrs.renameReportKeys(item)
		}

		return value
	default:
		return node
	}
}

func (psgrs *ProjectSecurityGateReportService) validateReportTemplateImages() error {
	if _, err := psgrs.settingSvc.GetImagePath("logo"); err != nil {
		return utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := psgrs.settingSvc.GetImagePath("cover-page-bg"); err != nil {
		return utility.NewForbiddenError("cover_image_missing")
	}

	return nil
}

func (psgrs *ProjectSecurityGateReportService) generateReportPDF(reportDTO interface{}, bodyTemplate string, coverTemplate string, footerTemplate string, pdfFileName string) (string, error) {
	htmlContent, err := utility.GeneratePDFContent(reportDTO, []string{bodyTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	headerContent, err := utility.GeneratePDFContent(reportDTO, []string{"project/templates/pdf/project_security_gate_report_header.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	footerContent, err := utility.GeneratePDFContent(reportDTO, []string{footerTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	coverContent, err := utility.GeneratePDFContent(reportDTO, []string{coverTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}

	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}

func (psgrs *ProjectSecurityGateReportService) getSelectedGates(securityGateIDs []uint) ([]*models.ProjectSecurityGate, error) {
	securityGates, err := psgrs.repo.GetGatesByIDs(securityGateIDs)
	if err != nil {
		return nil, err
	}

	if len(securityGates) != len(securityGateIDs) {
		return nil, utility.NewValidationError(errors.New("one or more selected security gates do not exist"), "invalid", "security_gate_ids")
	}

	return securityGates, nil
}

func (psgrs *ProjectSecurityGateReportService) mapReportUser(user *models.User) *organizationDTOs.UserMinimalDTO {
	if user == nil {
		return nil
	}

	var userDTO organizationDTOs.UserMinimalDTO
	if err := mapstructure.Decode(user, &userDTO); err != nil {
		return nil
	}

	return &userDTO
}

func (psgrs *ProjectSecurityGateReportService) joinIDList(ids []uint) string {
	values := make([]string, 0, len(ids))
	for _, id := range ids {
		values = append(values, strconv.FormatUint(uint64(id), 10))
	}

	return strings.Join(values, ",")
}

func (psgrs *ProjectSecurityGateReportService) parseIDList(value string) []uint {
	ids := []uint{}
	for _, part := range strings.Split(value, ",") {
		id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 32)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id))
	}

	return ids
}

func (psgrs *ProjectSecurityGateReportService) buildExceptionReviewItem(level int, actionAt *time.Time, workflowStatus *models.WorkflowStatus, actionByUser *models.User) *dtos.ProjectSecurityGateReportReviewDTO {
	review := &dtos.ProjectSecurityGateReportReviewDTO{
		Level: level,
	}

	if actionAt != nil {
		review.CreatedAtString = actionAt.Format("Jan 2, 2006 03:04 PM")
	}

	statusTitleEn := ""
	statusTitleAr := ""
	statusType := ""
	if workflowStatus != nil {
		statusTitleEn = workflowStatus.TitleEn
		statusTitleAr = workflowStatus.TitleAr
		statusType = workflowStatus.Type

		textColor, borderColor, backgroundColor :=
			utility.ResolveStatusBadgeColors(workflowStatus.ColorCode, workflowStatus.Type, workflowStatus.TitleEn)
		review.Status = &dtos.WorkflowStatusDTO{
			ID:              workflowStatus.ID,
			TitleEn:         workflowStatus.TitleEn,
			TitleAr:         workflowStatus.TitleAr,
			Type:            workflowStatus.Type,
			ColorCode:       workflowStatus.ColorCode,
			Status:          workflowStatus.Status,
			TextColor:       textColor,
			BorderColor:     borderColor,
			BackgroundColor: backgroundColor,
		}
	}

	review.ActionDateLabelEn = "Action at"
	review.ActionByLabelEn = "Action By"
	review.ActionDateLabelAr = "الإجراء في"
	review.ActionByLabelAr = "الإجراء بواسطة"
	lowerStatus := strings.ToLower(statusTitleEn)
	if lowerStatus == "approved" || statusType == "approve" {
		review.ActionDateLabelEn = "Approved at"
		review.ActionByLabelEn = "Approved By"
		review.ActionDateLabelAr = "تمت الموافقة في"
		review.ActionByLabelAr = "تمت الموافقة بواسطة"
	} else if lowerStatus == "submitted" || lowerStatus == "resubmitted" || statusType == "submit" {
		review.ActionDateLabelEn = "Submitted at"
		review.ActionByLabelEn = "Submitted By"
		review.ActionDateLabelAr = "تم الإرسال في"
		review.ActionByLabelAr = "تم الإرسال بواسطة"
	} else if lowerStatus == "rejected" || statusType == "reject" {
		review.ActionDateLabelEn = "Rejected at"
		review.ActionByLabelEn = "Rejected By"
		review.ActionDateLabelAr = "تم الرفض في"
		review.ActionByLabelAr = "تم الرفض بواسطة"
	} else if lowerStatus == "reverted" || statusType == "revert" {
		review.ActionDateLabelEn = "Reverted at"
		review.ActionByLabelEn = "Reverted By"
		review.ActionDateLabelAr = "تمت الإعادة في"
		review.ActionByLabelAr = "تمت الإعادة بواسطة"
	} else if statusTitleEn != "" {
		review.ActionDateLabelEn = fmt.Sprintf("%s at", statusTitleEn)
		review.ActionByLabelEn = fmt.Sprintf("%s By", statusTitleEn)
		if statusTitleAr != "" {
			review.ActionDateLabelAr = fmt.Sprintf("%s في", statusTitleAr)
			review.ActionByLabelAr = fmt.Sprintf("%s بواسطة", statusTitleAr)
		}
	}

	review.CreatedBy = psgrs.mapReportUser(actionByUser)

	return review
}

// Sections every security gate report starts with, in render order
var projectSecurityGateReportSectionTypes = []string{"executive-summary", "scope", "objectives"}

func formatSecurityGateReportSectionTitle(sectionType string) (string, string) {
	switch sectionType {
	case "executive-summary":
		return "Executive Summary", "الملخص التنفيذي"
	case "scope":
		return "Scope", "النطاق"
	case "objectives":
		return "Objectives", "الأهداف"
	default:
		return sectionType, sectionType
	}
}

func (psgrs *ProjectSecurityGateReportService) GetSecurityGateReportSection(reportID uint, sectionID uint) (*dtos.ProjectSecurityGateReportSectionDTO, error) {
	section, err := psgrs.repo.GetSecurityGateReportSectionByID(sectionID)
	if err != nil {
		return nil, err
	}

	if section.ReportID != reportID {
		return nil, utility.NewResourceNotFoundError(errors.New("section does not belong to report"), sectionID)
	}

	return mapSecurityGateReportSection(section.ID, section.Type, section.TitleEn, section.TitleAr, section.DescriptionEn, section.DescriptionAr, section.OrderNum), nil
}

func (psgrs *ProjectSecurityGateReportService) UpdateSecurityGateReportSection(reportID uint, sectionID uint, form *dtos.ProjectSecurityGateReportSectionFormDTO, currentUserID uint) error {
	section, err := psgrs.repo.GetSecurityGateReportSectionByID(sectionID)
	if err != nil {
		return err
	}

	if section.ReportID != reportID {
		return utility.NewResourceNotFoundError(errors.New("section does not belong to report"), sectionID)
	}

	previous := *section

	section.TitleEn = form.TitleEn
	section.TitleAr = form.TitleAr
	if form.DescriptionEn != nil {
		section.DescriptionEn = *form.DescriptionEn
	}
	if form.DescriptionAr != nil {
		section.DescriptionAr = *form.DescriptionAr
	}

	if err := psgrs.repo.UpdateSecurityGateReportSection(section); err != nil {
		return err
	}

	go func() {
		fieldsToCompare := []string{"TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr"}
		changes, err := utility.CompareFields(&previous, section, fieldsToCompare)
		if err != nil {
			utility.LogError("Failed to compare section fields", map[string]interface{}{"error": err.Error()})
			return
		}

		changedDataJSON, err := utility.ChangesToJSON(changes)
		if err != nil {
			utility.LogError("Failed to build section changes", map[string]interface{}{"error": err.Error()})
			return
		}

		err = psgrs.activityLogSvc.LogActivity(models.UpdateAction, models.ProjectSecurityGateReportSection{}.TableName(), sectionID,
			changedDataJSON, currentUserID, "project", "project_report", "security_gate_report_section", nil)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return nil
}

func (psgrs *ProjectSecurityGateReportService) DeleteSecurityGateReportSection(reportID uint, sectionID uint, currentUserID uint) (err error) {
	section, err := psgrs.repo.GetSecurityGateReportSectionByID(sectionID)
	if err != nil {
		return err
	}

	if section.ReportID != reportID {
		return utility.NewResourceNotFoundError(errors.New("section does not belong to report"), sectionID)
	}

	tx := psgrs.repo.BeginTransaction()
	defer func() {
		if recovered := recover(); recovered != nil {
			tx.Rollback()
			panic(recovered)
		}
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = psgrs.repo.DeleteSecurityGateReportSection(sectionID, tx); err != nil {
		return err
	}

	// the remaining sections close the gap
	if err = psgrs.repo.ReorderSecurityGateReportSections(reportID, section.OrderNum, tx); err != nil {
		return err
	}

	tx.Commit()

	go func() {
		logErr := psgrs.activityLogSvc.LogActivity(models.DeleteAction, models.ProjectSecurityGateReportSection{}.TableName(), sectionID, map[string]interface{}{
			"report_id": section.ReportID,
			"type":      section.Type,
			"title_en":  section.TitleEn,
			"title_ar":  section.TitleAr,
		}, currentUserID, "project", "project_report", "security_gate_report_section", nil)
		if logErr != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": logErr.Error()})
		}
	}()

	return nil
}

func (psgrs *ProjectSecurityGateReportService) GetSecurityGateExceptionReportSection(reportID uint, sectionID uint) (*dtos.ProjectSecurityGateReportSectionDTO, error) {
	section, err := psgrs.repo.GetSecurityGateExceptionReportSectionByID(sectionID)
	if err != nil {
		return nil, err
	}

	if section.ReportID != reportID {
		return nil, utility.NewResourceNotFoundError(errors.New("section does not belong to report"), sectionID)
	}

	return mapSecurityGateReportSection(section.ID, section.Type, section.TitleEn, section.TitleAr, section.DescriptionEn, section.DescriptionAr, section.OrderNum), nil
}

func (psgrs *ProjectSecurityGateReportService) UpdateSecurityGateExceptionReportSection(reportID uint, sectionID uint, form *dtos.ProjectSecurityGateReportSectionFormDTO, currentUserID uint) error {
	section, err := psgrs.repo.GetSecurityGateExceptionReportSectionByID(sectionID)
	if err != nil {
		return err
	}

	if section.ReportID != reportID {
		return utility.NewResourceNotFoundError(errors.New("section does not belong to report"), sectionID)
	}

	previous := *section

	section.TitleEn = form.TitleEn
	section.TitleAr = form.TitleAr
	if form.DescriptionEn != nil {
		section.DescriptionEn = *form.DescriptionEn
	}
	if form.DescriptionAr != nil {
		section.DescriptionAr = *form.DescriptionAr
	}

	if err := psgrs.repo.UpdateSecurityGateExceptionReportSection(section); err != nil {
		return err
	}

	go func() {
		fieldsToCompare := []string{"TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr"}
		changes, err := utility.CompareFields(&previous, section, fieldsToCompare)
		if err != nil {
			utility.LogError("Failed to compare section fields", map[string]interface{}{"error": err.Error()})
			return
		}

		changedDataJSON, err := utility.ChangesToJSON(changes)
		if err != nil {
			utility.LogError("Failed to build section changes", map[string]interface{}{"error": err.Error()})
			return
		}

		err = psgrs.activityLogSvc.LogActivity(models.UpdateAction, models.ProjectSecurityGateExceptionReportSection{}.TableName(), sectionID,
			changedDataJSON, currentUserID, "project", "project_report", "security_gate_exception_report_section", nil)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return nil
}

func (psgrs *ProjectSecurityGateReportService) DeleteSecurityGateExceptionReportSection(reportID uint, sectionID uint, currentUserID uint) (err error) {
	section, err := psgrs.repo.GetSecurityGateExceptionReportSectionByID(sectionID)
	if err != nil {
		return err
	}

	if section.ReportID != reportID {
		return utility.NewResourceNotFoundError(errors.New("section does not belong to report"), sectionID)
	}

	tx := psgrs.repo.BeginTransaction()
	defer func() {
		if recovered := recover(); recovered != nil {
			tx.Rollback()
			panic(recovered)
		}
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = psgrs.repo.DeleteSecurityGateExceptionReportSection(sectionID, tx); err != nil {
		return err
	}

	// the remaining sections close the gap
	if err = psgrs.repo.ReorderSecurityGateExceptionReportSections(reportID, section.OrderNum, tx); err != nil {
		return err
	}

	tx.Commit()

	go func() {
		logErr := psgrs.activityLogSvc.LogActivity(models.DeleteAction, models.ProjectSecurityGateExceptionReportSection{}.TableName(), sectionID, map[string]interface{}{
			"report_id": section.ReportID,
			"type":      section.Type,
			"title_en":  section.TitleEn,
			"title_ar":  section.TitleAr,
		}, currentUserID, "project", "project_report", "security_gate_exception_report_section", nil)
		if logErr != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": logErr.Error()})
		}
	}()

	return nil
}

func mapSecurityGateReportSection(id uint, sectionType string, titleEn string, titleAr string, descriptionEn string, descriptionAr string, orderNum int) *dtos.ProjectSecurityGateReportSectionDTO {
	return &dtos.ProjectSecurityGateReportSectionDTO{
		ID:            id,
		Type:          sectionType,
		TitleEn:       titleEn,
		TitleAr:       titleAr,
		DescriptionEn: descriptionEn,
		DescriptionAr: descriptionAr,
		OrderNum:      orderNum,
	}
}
