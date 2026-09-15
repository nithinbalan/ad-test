package services

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"jethur/config"
	"jethur/grc/models"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	settingSvc "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
)

type ProjectSecurityGateReportServiceInterface interface {
	GetFullReportDetailsByID(id uint) (*dtos.ProjectSecurityGateReportDTO, error)
	GetProjectSecurityGateReportExportZip(id uint, currentUserID uint, language string) (string, error)

	GetExceptionReportDetails() (*dtos.ProjectSecurityGateExceptionReportDTO, error)
	ExportProjectSecurityGateExceptionReport(currentUserID uint, language string) (string, error)
}

type ProjectSecurityGateReportService struct {
	repo         repositories.ProjectSecurityGateReportRepositoryInterface
	checklistSvc ProjectSecurityGateChecklistServiceInterface
	settingSvc   settingSvc.ReportTemplateSettingServiceInterface
}

func NewProjectSecurityGateReportService(
	repo repositories.ProjectSecurityGateReportRepositoryInterface,
	checklistSvc ProjectSecurityGateChecklistServiceInterface,
	settingSvc settingSvc.ReportTemplateSettingServiceInterface,
) ProjectSecurityGateReportServiceInterface {
	return &ProjectSecurityGateReportService{
		repo:         repo,
		checklistSvc: checklistSvc,
		settingSvc:   settingSvc,
	}
}

func (psgrs *ProjectSecurityGateReportService) GetFullReportDetailsByID(id uint) (*dtos.ProjectSecurityGateReportDTO, error) {
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
		review := buildExceptionReviewItem(history, len(reviewsByChecklistID[history.ChecklistID])+1)
		reviewsByChecklistID[history.ChecklistID] = append(reviewsByChecklistID[history.ChecklistID], review)
	}

	checklists, err := psgrs.repo.GetAllChecklists(id)
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
			checklistItem.ResponseStatusTitleEn = checklist.ResponseStatus.TitleEn
			checklistItem.ResponseStatusTitleAr = checklist.ResponseStatus.TitleAr
			checklistItem.ResponseStatusColorCode = checklist.ResponseStatus.ColorCode
		}
		if checklist.Comment != nil {
			checklistItem.Comment = *checklist.Comment
		}
		if checklist.UpdatedByUser != nil {
			checklistItem.RespondedByEn = strings.TrimSpace(checklist.UpdatedByUser.FirstNameEn + " " + checklist.UpdatedByUser.LastNameEn)
			checklistItem.RespondedByAr = strings.TrimSpace(checklist.UpdatedByUser.FirstNameAr + " " + checklist.UpdatedByUser.LastNameAr)
		}
		if checklistItem.RespondedByAr == "" {
			checklistItem.RespondedByAr = checklistItem.RespondedByEn
		}
		if checklist.UpdatedAt != nil {
			checklistItem.RespondedAtFormatted = checklist.UpdatedAt.Format("01/02/2006")
		}
		report.Checklists = append(report.Checklists, checklistItem)

		if checklist.IsException {
			exception := &dtos.ProjectSecurityGateReportExceptionDTO{
				ChecklistID: checklist.ID,
				TitleEn:     checklistTitleEn,
				TitleAr:     checklistTitleAr,
				IsCritical:  checklistItem.IsCritical,
				Reviews:     reviewsByChecklistID[checklist.ID],
			}

			if checklist.ExceptionJustification != nil {
				exception.Justification = *checklist.ExceptionJustification
			}
			if checklist.WorkflowStatus != nil {
				exception.WorkflowStatusTitleEn = checklist.WorkflowStatus.TitleEn
				exception.WorkflowStatusTitleAr = checklist.WorkflowStatus.TitleAr
				exception.WorkflowStatusColorCode = checklist.WorkflowStatus.ColorCode
			}
			if checklist.SubmittedByUser != nil {
				exception.SubmittedByEn = strings.TrimSpace(checklist.SubmittedByUser.FirstNameEn + " " + checklist.SubmittedByUser.LastNameEn)
				exception.SubmittedByAr = strings.TrimSpace(checklist.SubmittedByUser.FirstNameAr + " " + checklist.SubmittedByUser.LastNameAr)
			}
			if checklist.SubmittedAt != nil {
				exception.SubmittedAtFormatted = checklist.SubmittedAt.Format("01/02/2006")
			}
			if checklist.ApprovedByUser != nil {
				exception.ApprovedByEn = strings.TrimSpace(checklist.ApprovedByUser.FirstNameEn + " " + checklist.ApprovedByUser.LastNameEn)
				exception.ApprovedByAr = strings.TrimSpace(checklist.ApprovedByUser.FirstNameAr + " " + checklist.ApprovedByUser.LastNameAr)
			}
			if exception.SubmittedByAr == "" {
				exception.SubmittedByAr = exception.SubmittedByEn
			}
			if exception.ApprovedByAr == "" {
				exception.ApprovedByAr = exception.ApprovedByEn
			}
			if checklist.ApprovedAt != nil {
				exception.ApprovedAtFormatted = checklist.ApprovedAt.Format("01/02/2006")
			}

			report.Exceptions = append(report.Exceptions, exception)
		}

		for _, mapping := range checklist.Controls {
			if mapping.Control == nil {
				continue
			}

			control := &dtos.ProjectSecurityGateReportControlDTO{
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
				risk.RiskStatusTitleEn = mapping.RiskRegister.RiskStatus.TitleEn
				risk.RiskStatusTitleAr = mapping.RiskRegister.RiskStatus.TitleAr
			}
			if mapping.RiskRegister.RiskCategory != nil {
				risk.RiskCategoryTitleEn = mapping.RiskRegister.RiskCategory.TitleEn
				risk.RiskCategoryTitleAr = mapping.RiskRegister.RiskCategory.TitleAr
			}
			if mapping.RiskRegister.ResponsibleUnit != nil {
				risk.ResponsibleUnitTitleEn = mapping.RiskRegister.ResponsibleUnit.TitleEn
				risk.ResponsibleUnitTitleAr = mapping.RiskRegister.ResponsibleUnit.TitleAr
			}
			if mapping.RiskRegister.ResponsibleUser != nil {
				risk.ResponsibleUserEn = strings.TrimSpace(mapping.RiskRegister.ResponsibleUser.FirstNameEn + " " + mapping.RiskRegister.ResponsibleUser.LastNameEn)
				risk.ResponsibleUserAr = strings.TrimSpace(mapping.RiskRegister.ResponsibleUser.FirstNameAr + " " + mapping.RiskRegister.ResponsibleUser.LastNameAr)
				if risk.ResponsibleUserAr == "" {
					risk.ResponsibleUserAr = risk.ResponsibleUserEn
				}
			}
			if mapping.RiskRegister.RegistrationDate != nil {
				risk.RegistrationDateFormatted = mapping.RiskRegister.RegistrationDate.Format("01/02/2006")
			}
			report.Risks = append(report.Risks, risk)
		}

	}

	return &report, nil
}

func (psgrs *ProjectSecurityGateReportService) GetProjectSecurityGateReportExportZip(id uint, currentUserID uint, language string) (string, error) {
	reportDTO, err := psgrs.GetFullReportDetailsByID(id)
	if err != nil {
		return "", err
	}

	reportDTO.GeneratedAt = time.Now().Format("2006-01-02 03:04 PM")

	// Only body, cover and footer have Arabic variants
	bodyTemplate := "project/templates/pdf/project_security_gate_report.html"
	coverTemplate := "project/templates/pdf/project_security_gate_report_cover.html"
	footerTemplate := "project/templates/pdf/project_security_gate_report_footer.html"
	if language == "ar" {
		bodyTemplate = "project/templates/pdf/project_security_gate_report_ar.html"
		coverTemplate = "project/templates/pdf/project_security_gate_report_cover_ar.html"
		footerTemplate = "project/templates/pdf/project_security_gate_report_footer_ar.html"
	}

	// Main content template
	htmlContent, err := utility.GeneratePDFContent(reportDTO, []string{bodyTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	// Header template
	headerTemplatePaths := []string{"project/templates/pdf/project_security_gate_report_header.html"}
	headerContent, err := utility.GeneratePDFContent(reportDTO, headerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	// Footer template
	footerContent, err := utility.GeneratePDFContent(reportDTO, []string{footerTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	// Cover template
	coverContent, err := utility.GeneratePDFContent(reportDTO, []string{coverTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}

	// Convert HTML -> PDF with cover, header, footer
	pdfFileName := fmt.Sprintf("ProjectSecurityGateReport_%s_%d.pdf", language, id)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Report at the root, documents in folders
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

	// CreateZipFile would flatten the folders
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

// buildExceptionReviewItem maps a workflow history row to a timeline entry
func buildExceptionReviewItem(history *models.ProjectSecurityGateChecklistExceptionWorkflowHistory, stepNumber int) *dtos.ProjectSecurityGateReportReviewDTO {
	review := &dtos.ProjectSecurityGateReportReviewDTO{
		StepNumber:      stepNumber,
		Level:           history.Level,
		CreatedAtString: history.CreatedAt.Format("Jan 2, 2006 03:04 PM"),
	}

	if history.WorkflowStatus != nil {
		review.StatusTitleEn = history.WorkflowStatus.TitleEn
		review.StatusTitleAr = history.WorkflowStatus.TitleAr
		review.StatusType = history.WorkflowStatus.Type
		review.StatusTextColor, review.StatusBorderColor, review.StatusBackgroundColor =
			utility.ResolveStatusBadgeColors(history.WorkflowStatus.ColorCode, history.WorkflowStatus.Type, history.WorkflowStatus.TitleEn)
	}

	review.ActionDateLabelEn = "Action at"
	review.ActionByLabelEn = "Action By"
	review.ActionDateLabelAr = "الإجراء في"
	review.ActionByLabelAr = "الإجراء بواسطة"
	lowerStatus := strings.ToLower(review.StatusTitleEn)
	if lowerStatus == "approved" || review.StatusType == "approve" {
		review.ActionDateLabelEn = "Approved at"
		review.ActionByLabelEn = "Approved By"
		review.ActionDateLabelAr = "تمت الموافقة في"
		review.ActionByLabelAr = "تمت الموافقة بواسطة"
	} else if lowerStatus == "submitted" || lowerStatus == "resubmitted" || review.StatusType == "submit" {
		review.ActionDateLabelEn = "Submitted at"
		review.ActionByLabelEn = "Submitted By"
		review.ActionDateLabelAr = "تم الإرسال في"
		review.ActionByLabelAr = "تم الإرسال بواسطة"
	} else if lowerStatus == "rejected" || review.StatusType == "reject" {
		review.ActionDateLabelEn = "Rejected at"
		review.ActionByLabelEn = "Rejected By"
		review.ActionDateLabelAr = "تم الرفض في"
		review.ActionByLabelAr = "تم الرفض بواسطة"
	} else if lowerStatus == "reverted" || review.StatusType == "revert" {
		review.ActionDateLabelEn = "Reverted at"
		review.ActionByLabelEn = "Reverted By"
		review.ActionDateLabelAr = "تمت الإعادة في"
		review.ActionByLabelAr = "تمت الإعادة بواسطة"
	} else if review.StatusTitleEn != "" {
		review.ActionDateLabelEn = fmt.Sprintf("%s at", review.StatusTitleEn)
		review.ActionByLabelEn = fmt.Sprintf("%s By", review.StatusTitleEn)
		if review.StatusTitleAr != "" {
			review.ActionDateLabelAr = fmt.Sprintf("%s في", review.StatusTitleAr)
			review.ActionByLabelAr = fmt.Sprintf("%s بواسطة", review.StatusTitleAr)
		}
	}

	if history.CreatedByUser != nil {
		review.UserFullNameEn = strings.TrimSpace(history.CreatedByUser.FirstNameEn + " " + history.CreatedByUser.LastNameEn)
		review.UserFullNameAr = strings.TrimSpace(history.CreatedByUser.FirstNameAr + " " + history.CreatedByUser.LastNameAr)
		if review.UserFullNameAr == "" {
			review.UserFullNameAr = review.UserFullNameEn
		}

		review.UserInitial = history.CreatedByUser.Initial
		if review.UserInitial == "" {
			firstName := strings.TrimSpace(history.CreatedByUser.FirstNameEn)
			lastName := strings.TrimSpace(history.CreatedByUser.LastNameEn)
			if len(firstName) > 0 {
				review.UserInitial += string([]rune(firstName)[0:1])
			}
			if len(lastName) > 0 {
				review.UserInitial += string([]rune(lastName)[0:1])
			}
			review.UserInitial = strings.ToUpper(review.UserInitial)
		}

		// users.initial is Latin, build an Arabic one
		firstNameAr := strings.TrimSpace(history.CreatedByUser.FirstNameAr)
		lastNameAr := strings.TrimSpace(history.CreatedByUser.LastNameAr)
		if len(firstNameAr) > 0 {
			review.UserInitialAr += string([]rune(firstNameAr)[0:1])
		}
		if len(lastNameAr) > 0 {
			review.UserInitialAr += string([]rune(lastNameAr)[0:1])
		}
		if review.UserInitialAr == "" {
			review.UserInitialAr = review.UserInitial
		}
	}

	return review
}

func (psgrs *ProjectSecurityGateReportService) GetExceptionReportDetails() (*dtos.ProjectSecurityGateExceptionReportDTO, error) {
	if _, err := psgrs.settingSvc.GetImagePath("logo"); err != nil {
		return nil, utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := psgrs.settingSvc.GetImagePath("cover-page-bg"); err != nil {
		return nil, utility.NewForbiddenError("cover_image_missing")
	}

	histories, err := psgrs.repo.GetAllExceptionWorkflowHistoriesForAllGates()
	if err != nil {
		return nil, err
	}

	reviewsByChecklistID := make(map[uint][]*dtos.ProjectSecurityGateReportReviewDTO)
	for _, history := range histories {
		review := buildExceptionReviewItem(history, len(reviewsByChecklistID[history.ChecklistID])+1)
		reviewsByChecklistID[history.ChecklistID] = append(reviewsByChecklistID[history.ChecklistID], review)
	}

	checklists, err := psgrs.repo.GetAllExceptionChecklists()
	if err != nil {
		return nil, err
	}

	report := &dtos.ProjectSecurityGateExceptionReportDTO{
		Exceptions: make([]*dtos.ProjectSecurityGateExceptionReportItemDTO, 0, len(checklists)),
	}

	for _, checklist := range checklists {
		item := &dtos.ProjectSecurityGateExceptionReportItemDTO{
			Reviews: reviewsByChecklistID[checklist.ID],
		}

		if checklist.SecurityGate != nil {
			item.ProjectTitleEn = checklist.SecurityGate.ProjectTitleEn
			item.ProjectTitleAr = checklist.SecurityGate.ProjectTitleAr

			if checklist.SecurityGate.Stage != nil {
				item.StageTitleEn = checklist.SecurityGate.Stage.TitleEn
				item.StageTitleAr = checklist.SecurityGate.Stage.TitleAr
			}
			if checklist.SecurityGate.ReadinessStatus != nil {
				item.ReadinessStatusTitleEn = checklist.SecurityGate.ReadinessStatus.TitleEn
				item.ReadinessStatusTitleAr = checklist.SecurityGate.ReadinessStatus.TitleAr
			}
			if checklist.SecurityGate.BusinessOwner != nil {
				item.BusinessOwnerEn = strings.TrimSpace(checklist.SecurityGate.BusinessOwner.FirstNameEn + " " + checklist.SecurityGate.BusinessOwner.LastNameEn)
				item.BusinessOwnerAr = strings.TrimSpace(checklist.SecurityGate.BusinessOwner.FirstNameAr + " " + checklist.SecurityGate.BusinessOwner.LastNameAr)
				if item.BusinessOwnerAr == "" {
					item.BusinessOwnerAr = item.BusinessOwnerEn
				}
			}
		}

		if checklist.Requirement != nil {
			item.TitleEn = checklist.Requirement.TitleEn
			item.TitleAr = checklist.Requirement.TitleAr
			item.IsCritical = checklist.Requirement.IsCritical
		}
		if checklist.ExceptionJustification != nil {
			item.Justification = *checklist.ExceptionJustification
		}
		if checklist.WorkflowStatus != nil {
			item.WorkflowStatusTitleEn = checklist.WorkflowStatus.TitleEn
			item.WorkflowStatusTitleAr = checklist.WorkflowStatus.TitleAr
			item.WorkflowStatusTextColor, item.WorkflowStatusBorderColor, item.WorkflowStatusBackgroundColor =
				utility.ResolveStatusBadgeColors(checklist.WorkflowStatus.ColorCode, checklist.WorkflowStatus.Type, checklist.WorkflowStatus.TitleEn)
		}
		if checklist.SubmittedByUser != nil {
			item.SubmittedByEn = strings.TrimSpace(checklist.SubmittedByUser.FirstNameEn + " " + checklist.SubmittedByUser.LastNameEn)
			item.SubmittedByAr = strings.TrimSpace(checklist.SubmittedByUser.FirstNameAr + " " + checklist.SubmittedByUser.LastNameAr)
			if item.SubmittedByAr == "" {
				item.SubmittedByAr = item.SubmittedByEn
			}
		}
		if checklist.SubmittedAt != nil {
			item.SubmittedAtFormatted = checklist.SubmittedAt.Format("01/02/2006")
		}
		if checklist.ApprovedByUser != nil {
			item.ApprovedByEn = strings.TrimSpace(checklist.ApprovedByUser.FirstNameEn + " " + checklist.ApprovedByUser.LastNameEn)
			item.ApprovedByAr = strings.TrimSpace(checklist.ApprovedByUser.FirstNameAr + " " + checklist.ApprovedByUser.LastNameAr)
			if item.ApprovedByAr == "" {
				item.ApprovedByAr = item.ApprovedByEn
			}
		}
		if checklist.ApprovedAt != nil {
			item.ApprovedAtFormatted = checklist.ApprovedAt.Format("01/02/2006")
		}

		if checklist.ApprovedBYID != nil {
			report.ApprovedCount++
		} else {
			report.PendingCount++
		}
		if item.IsCritical {
			report.CriticalCount++
		}

		report.Exceptions = append(report.Exceptions, item)
	}

	report.TotalCount = len(report.Exceptions)

	return report, nil
}

func (psgrs *ProjectSecurityGateReportService) ExportProjectSecurityGateExceptionReport(currentUserID uint, language string) (string, error) {
	reportDTO, err := psgrs.GetExceptionReportDetails()
	if err != nil {
		return "", err
	}

	reportDTO.GeneratedAt = time.Now().Format("2006-01-02 03:04 PM")

	bodyTemplate := "project/templates/pdf/project_security_gate_exception_report.html"
	coverTemplate := "project/templates/pdf/project_security_gate_exception_report_cover.html"
	footerTemplate := "project/templates/pdf/project_security_gate_report_footer.html"
	if language == "ar" {
		bodyTemplate = "project/templates/pdf/project_security_gate_exception_report_ar.html"
		coverTemplate = "project/templates/pdf/project_security_gate_exception_report_cover_ar.html"
		footerTemplate = "project/templates/pdf/project_security_gate_report_footer_ar.html"
	}

	// Main content template
	htmlContent, err := utility.GeneratePDFContent(reportDTO, []string{bodyTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	// Header template
	headerTemplatePaths := []string{"project/templates/pdf/project_security_gate_report_header.html"}
	headerContent, err := utility.GeneratePDFContent(reportDTO, headerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	// Footer template
	footerContent, err := utility.GeneratePDFContent(reportDTO, []string{footerTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	// Cover template
	coverContent, err := utility.GeneratePDFContent(reportDTO, []string{coverTemplate}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}

	// No id: the report spans every gate
	pdfFileName := fmt.Sprintf("ProjectSecurityGateExceptionReport_%s.pdf", language)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}
