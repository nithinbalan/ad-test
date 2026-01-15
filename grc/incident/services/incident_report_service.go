package services

import (
	"errors"
	"fmt"
	"io"
	"jethur/config"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepo "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	settingSvc "jethur/grc/settings/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentReportServiceInterface interface {
	ExportIncidentReport(id uint, currentUserID uint) (string, error)
	GetFullReportDetailsByID(id uint) (*incidentDTOs.IncidentReportDetailDTO, error)
	GetAll(params *incidentDTOs.IncidentReportListQueryParamsDTO) (*incidentDTOs.IncidentReportListResponseDTO, error)
	Create(currentUserID uint, dto incidentDTOs.IncidentReportRequestDTO) (uint, error)
	GetIncidentReportDetails(id uint) (*incidentDTOs.IncidentReportDetailsDTO, error)
	DeleteIncidentReport(incidentReportID uint, currentUserID uint) error
	GetSectionByID(incidentReportID uint, sectionID uint, tx ...utility.DatabaseTransaction) (*models.IncidentReportSection, error)
	UpdateSection(incidentReportID uint, sectionID uint, currentUserID uint, formDTO *incidentDTOs.IncidentReportSectionFormDTO) error
	DeleteIncidentReportSection(incidentReportID uint, sectionID uint, currentUserID uint) error
	GetIncidentReportExportZip(incidentReportID uint, currentUserID uint) (string, error)
	GetDocumentPath(id uint, token string) (string, error)
}

// Service and implementations
type IncidentReportService struct {
	repo                         incidentRepo.IncidentReportRepositoryInterface
	settigSvc                    settingSvc.ReportTemplateSettingServiceInterface
	incidentRootCauseAnalysisSvc IncidentRootCauseAnalysisServiceInterface
	incidentActionPlanSvc        IncidentActionPlanServiceInterface
	activityLogSvc               logServices.ActivityLogServiceInterface
}

func NewIncidentReportService(
	repo incidentRepo.IncidentReportRepositoryInterface,
	settigSvc settingSvc.ReportTemplateSettingServiceInterface,
	incidentRootCauseAnalysisSvc IncidentRootCauseAnalysisServiceInterface,
	incidentActionPlanSvc IncidentActionPlanServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) IncidentReportServiceInterface {
	return &IncidentReportService{
		repo:                         repo,
		settigSvc:                    settigSvc,
		incidentRootCauseAnalysisSvc: incidentRootCauseAnalysisSvc,
		incidentActionPlanSvc:        incidentActionPlanSvc,
		activityLogSvc:               activityLogSvc,
	}
}

func (irs *IncidentReportService) GetAll(params *incidentDTOs.IncidentReportListQueryParamsDTO) (*incidentDTOs.IncidentReportListResponseDTO, error) {
	return irs.repo.GetAll(params)
}

func (irs IncidentReportService) GetFullReportDetailsByID(id uint) (*incidentDTOs.IncidentReportDetailDTO, error) {
	// Initialize your DTO
	var incidentDTO incidentDTOs.IncidentReportDetailDTO
	incident, err := irs.repo.GetFullReportDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if _, err := irs.settigSvc.GetImagePath("logo"); err != nil {
		return nil, utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := irs.settigSvc.GetImagePath("cover-page-bg"); err != nil {
		return nil, utility.NewForbiddenError("cover_image_missing")
	}

	err = mapstructure.Decode(incident, &incidentDTO)
	if err != nil {
		return nil, err
	}

	kpis, err := irs.repo.GetAllIncidentKPI(id)
	if err != nil {
		return nil, err
	}
	incidentDTO.KPIs = kpis

	strategicObjectives, err := irs.repo.GetAllIncidentStrategicObjectives(id)
	if err != nil {
		return nil, err
	}
	incidentDTO.StrategicObjectives = strategicObjectives

	actionPlans, err := irs.repo.GetAllIncidentActionPlan(id)
	if err != nil {
		return nil, err
	}
	incidentDTO.ActionPlans = actionPlans

	assets, err := irs.repo.GetAllIncidentAssets(id)
	if err != nil {
		return nil, err
	}
	incidentDTO.Assets = assets

	associatedDocuments, err := irs.repo.GetAllIncidentAssociatedDocuments(id)
	if err != nil {
		return nil, err
	}
	incidentDTO.AssociatedDocuments = associatedDocuments

	controls, err := irs.repo.GetAllIncidentControls(id)
	if err != nil {
		return nil, err
	}
	incidentDTO.Controls = controls

	risks, err := irs.repo.GetAllIncidentRisks(id)
	if err != nil {
		return nil, err
	}
	incidentDTO.Risks = risks

	process, err := irs.repo.GetAllIncidentProcesses(id)
	if err != nil {
		return nil, err
	}
	incidentDTO.Processes = process

	complianceRequirements, err := irs.repo.GetAllIncidentComplianceRequirements(id)
	if err != nil {
		return nil, err
	}
	incidentDTO.ComplianceRequirements = complianceRequirements

	if incident.IncidentClosure != nil && incident.IncidentClosure.CompletionDate != nil {
		dateStr := incident.IncidentClosure.CompletionDate.Format("2006-01-02")
		incidentDTO.IncidentClosure.CompletionDateString = &dateStr
	}

	if incident.Investigation != nil && incident.Investigation.InitiatedAt != nil {
		dateStr := incident.Investigation.InitiatedAt.Format("2006-01-02")
		incidentDTO.Investigation.InitiatedAtString = &dateStr
	}

	incidentDTO.IncidentAtFormatted = incident.IncidentAt.Format("01 02, 2006, 03:04 PM")
	incidentDTO.ReportedAtFormatted = incident.ReportedAt.Format("01 02, 2006, 03:04 PM")

	incidentDTO.Timelines = []*incidentDTOs.IncidentTimelineListDTO{}
	for _, t := range incident.Timelines {
		incidentDTO.Timelines = append(incidentDTO.Timelines, &incidentDTOs.IncidentTimelineListDTO{
			ID:                          t.ID,
			IncidentTimelineTypeTitleEn: t.IncidentTimelineType.TitleEn,
			TimelineStatusTitleEn:       t.TimelineStatus.TitleEn,
			PerformedByFirstNameEn:      &t.PerformedBy.FirstNameEn,
			PerformedByLastNameEn:       &t.PerformedBy.LastNameEn,
			PerformedAt:                 t.PerformedAt,
		})
	}

	if len(incident.IncidentRootCauseAnalyses) > 0 {
		incidentDTO.RootCauseAnalyses = make([]*incidentDTOs.IncidentRootCauseAnalysisListDTO, 0, len(incident.IncidentRootCauseAnalyses))
		for _, rca := range incident.IncidentRootCauseAnalyses {
			incidentDTO.RootCauseAnalyses = append(incidentDTO.RootCauseAnalyses, &incidentDTOs.IncidentRootCauseAnalysisListDTO{
				Why:      rca.Why,
				Question: rca.Question,
				Answer:   rca.Answer,
			})
		}
	}

	incidentDTO.AppUrl = os.Getenv("APP_URL")
	incidentDTO.CreatedAt = incident.CreatedAt
	incidentDTO.UpdatedAt = incident.UpdatedAt

	return &incidentDTO, nil
}

func (irs IncidentReportService) ExportIncidentReport(id uint, currentUserID uint) (string, error) {
	incidentDTO, err := irs.GetFullReportDetailsByID(id)
	if err != nil {
		return "", err
	}

	incidentDTO.GeneratedAt = time.Now().Format("2006-01-02 03:04 PM")

	// Main content template
	templatePaths := []string{"incident/templates/pdf/incident_report.html"}
	htmlContent, err := utility.GeneratePDFContent(incidentDTO, templatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	// Header template - process with data to resolve template variables
	headerTemplatePaths := []string{"incident/templates/pdf/incident_report_header.html"}
	headerContent, err := utility.GeneratePDFContent(incidentDTO, headerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	// Footer template - process with data to resolve template variables
	footerTemplatePaths := []string{"incident/templates/pdf/incident_report_footer.html"}
	footerContent, err := utility.GeneratePDFContent(incidentDTO, footerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}
	fcoverTemplatePaths := []string{"incident/templates/pdf/incident_report_cover.html"}
	coverContent, err := utility.GeneratePDFContent(incidentDTO, fcoverTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	// Convert HTML → PDF with header and footer
	pdfFileName := fmt.Sprintf("IncidentReport_%d.pdf", id)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}

func (irs *IncidentReportService) Create(
	currentUserID uint,
	dto incidentDTOs.IncidentReportRequestDTO,
) (uint, error) {

	tx := irs.repo.BeginTransaction()

	template, err := irs.settigSvc.GetFirst()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if template.CoverPageBGImagePath == nil || template.LogoImagePath == nil {
		tx.Rollback()
		return 0, utility.NewForbiddenError("Logo and Background Image not added.")
	}

	dates, err := irs.parseIncidentReportDates(dto)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	report := &models.IncidentReport{
		Title:       dto.Title,
		StartDate:   dates.StartDate,
		EndDate:     dates.EndDate,
		CreatedBYID: currentUserID,
	}

	reportID, err := irs.repo.Create(report, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sectionTypes := []string{
		"executive-summary", "scope", "objectives", "overview",
		"details", "investigation", "incident-closure", "post-incident-review",
	}

	for i, s := range sectionTypes {
		en, ar := formatTitle(s)

		sec := models.IncidentReportSection{
			IncidentReportID: reportID,
			Type:             s,
			TitleEn:          en,
			TitleAr:          ar,
			OrderNum:         i + 1,
		}

		if _, err := irs.repo.CreateSection(&sec, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var incidents []incidentDTOs.IncidentReportIncidentDTO

	switch dto.Type {

	case "date-range":
		if dates.StartDate == nil || dates.EndDate == nil {
			tx.Rollback()
			return 0, utility.NewValidationError(
				errors.New("start_date and end_date are required for date-range reports"),
				"required",
				"start_date",
			)
		}

		incidents, err = irs.repo.GetReportIncidentsByDateRange(
			*dates.StartDate,
			*dates.EndDate,
			tx,
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

	case "individual":
		if dto.IncidentIDs == nil || len(*dto.IncidentIDs) == 0 {
			tx.Rollback()
			return 0, utility.NewValidationError(
				errors.New("incident_ids are required for individual reports"),
				"required",
				"incident_ids",
			)
		}

		incidents, err = irs.repo.GetIncidentsByIDs(*dto.IncidentIDs, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

	default:
		tx.Rollback()
		return 0, utility.NewValidationError(
			errors.New("invalid report type"),
			"invalid",
			"type",
		)
	}

	for _, inc := range incidents {
		reportIncident := &models.IncidentReportIncident{
			IncidentReportID:           reportID,
			IncidentRegisterID:         inc.ID,
			ReferenceNumber:            &inc.ReferenceNumber,
			IncidentTitle:              inc.IncidentTitle,
			IncidentCategoryTitleEn:    &inc.IncidentCategoryTitleEn,
			IncidentCategoryTitleAr:    &inc.IncidentCategoryTitleAr,
			IncidentSubCategoryTitleEn: inc.IncidentSubCategoryTitleEn,
			IncidentSubCategoryTitleAr: inc.IncidentSubCategoryTitleAr,
			SeverityLevelTitleEn:       &inc.SeverityTitleEn,
			SeverityLevelTitleAr:       &inc.SeverityTitleAr,
			SeverityLevelColorCode:     inc.SeverityColorCode,
			IncidentAt:                 &inc.RealTime,
			ReportedAt:                 &inc.DetectedTime,
			IncidentLocation:           inc.IncidentLocation,
			IncidentStatusTitleEn:      &inc.IncidentStatusTitleEn,
			IncidentStatusTitleAr:      &inc.IncidentStatusTitleAr,
			ResponsibleUserNameEn:      inc.ResponsibleUserNameEn,
			ResponsibleUserNameAr:      inc.ResponsibleUserNameAr,
			ReportedUserNameEn:         inc.ReportedUserNameEn,
			ReportedUserNameAr:         inc.ReportedUserNameAr,
			ClosedDate:                 inc.ResolvedDate,
		}

		if err := irs.repo.CreateReportIncident(reportIncident, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		documents, err := irs.repo.GetDocumentsByIncidentID(inc.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		var reportDocs []*models.IncidentReportDocument

		for _, document := range documents {
			if document.Name == nil || *document.Name == "" {
				continue // Skip documents without a name
			}

			storagePath := config.GetEnv("STORAGE_PATH", "")

			oldFilePath := filepath.Join(
				storagePath,
				"incident",
				"incident-register",
				strconv.FormatUint(uint64(inc.ID), 10),
				"documents",
				*document.Name,
			)

			if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
				tx.Rollback()
				return 0, fmt.Errorf("source file does not exist: %s", oldFilePath)
			}

			newStoragePath := filepath.Join(
				storagePath,
				"incident",
				"incident-report",
				strconv.FormatUint(uint64(reportID), 10),
				"documents",
			)

			newFilePath := filepath.Join(newStoragePath, *document.Name)

			newURL := utility.GenerateDocumentURL(
				"incident",
				"incident-report",
				strconv.FormatUint(uint64(reportID), 10),
				"documents",
			)

			if err := os.MkdirAll(newStoragePath, os.ModePerm); err != nil {
				tx.Rollback()
				return 0, err
			}

			input, err := os.Open(oldFilePath)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			defer input.Close()

			output, err := os.Create(newFilePath)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			defer output.Close()

			if _, err := io.Copy(output, input); err != nil {
				tx.Rollback()
				return 0, err
			}

			docName := *document.Name
			docSize := document.Size
			docToken := document.Token
			docOriginalName := document.OriginalName
			docType := document.Type
			reportDocs = append(reportDocs, &models.IncidentReportDocument{
				IncidentDocumentID:       document.ID,
				IncidentReportIncidentID: reportIncident.ID,
				IncidentInvestigationID:  nil, // Can be set if document is linked to investigation
				URL:                      &newURL,
				Name:                     &docName,
				Size:                     &docSize,
				Token:                    &docToken,
				OriginalName:             &docOriginalName,
				Type:                     &docType,
				CreatedByID:              &currentUserID,
				UpdatedByID:              nil,
			})
		}

		if len(reportDocs) > 0 {
			if err := irs.repo.CreateIncidentDocuments(reportDocs, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}

		analyses, err := irs.incidentRootCauseAnalysisSvc.GetAllByIncidentID(inc.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		var rcas []*models.IncidentReportRootCauseAnalysis

		for _, a := range analyses {
			var answer *string
			if a.Answer != "" {
				answer = &a.Answer
			}
			rcas = append(rcas, &models.IncidentReportRootCauseAnalysis{
				IncidentReportIncidentID: reportIncident.ID,
				IncidentReportID:         reportID,
				Why:                      a.Why,
				Question:                 a.Question,
				Answer:                   answer,
			})
		}

		if len(rcas) > 0 {
			if err := irs.repo.CreateIncidentRootCauseAnalyses(rcas, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}

		plans, err := irs.incidentActionPlanSvc.GetAllByIncidentID(inc.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		var reportPlans []*models.IncidentReportActionPlan

		for _, p := range plans {
			actionPlanTitle := p.Title
			var actionPlanDescription *string
			if p.Description != nil && *p.Description != "" {
				actionPlanDescription = p.Description
			}

			responsibleUserNameEn := ""
			responsibleUserNameAr := ""
			if p.ResponsibleUser != nil {
				responsibleUserNameEn = fmt.Sprintf("%s %s", p.ResponsibleUser.FirstNameEn, p.ResponsibleUser.LastNameEn)
				responsibleUserNameAr = fmt.Sprintf("%s %s", p.ResponsibleUser.FirstNameAr, p.ResponsibleUser.LastNameAr)
			}

			severityLevelTitleEn := ""
			severityLevelTitleAr := ""
			if p.SeverityLevel != nil {
				severityLevelTitleEn = p.SeverityLevel.TitleEn
				severityLevelTitleAr = p.SeverityLevel.TitleAr
			}

			actionPlanStatusTitleEn := ""
			actionPlanStatusTitleAr := ""
			actionPlanStatusColorCode := ""
			if p.ActionPlanStatus != nil {
				actionPlanStatusTitleEn = p.ActionPlanStatus.TitleEn
				actionPlanStatusTitleAr = p.ActionPlanStatus.TitleAr
				actionPlanStatusColorCode = p.ActionPlanStatus.ColorCode
			}

			// Ensure target date is never nil when inserting into incident_report_action_plans
			var targetDate *time.Time
			if p.TargetDate != nil {
				targetDate = p.TargetDate
			} else {
				// Fallback: use today's date to satisfy NOT NULL constraint in report table
				now := time.Now().UTC()
				targetDate = &now
			}

			reportPlans = append(reportPlans, &models.IncidentReportActionPlan{
				IncidentReportIncidentID: reportIncident.ID,
				IncidentReportID:         reportID,

				ActionPlanTitle:       actionPlanTitle,
				ActionPlanDescription: actionPlanDescription,

				TargetDate: targetDate,

				ResponsibleUserNameEn: responsibleUserNameEn,
				ResponsibleUserNameAr: responsibleUserNameAr,

				SeverityLevelTitleEn: &severityLevelTitleEn,
				SeverityLevelTitleAr: &severityLevelTitleAr,

				ActionPlanStatusTitleEn:   actionPlanStatusTitleEn,
				ActionPlanStatusTitleAr:   actionPlanStatusTitleAr,
				ActionPlanStatusColorCode: actionPlanStatusColorCode,
			})
		}

		if len(reportPlans) > 0 {
			if err := irs.repo.CreateIncidentReportActionPlans(reportPlans, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()
	return reportID, nil
}

func (irs *IncidentReportService) parseIncidentReportDates(dto incidentDTOs.IncidentReportRequestDTO) (*incidentDTOs.IncidentReportParsedDates, error) {
	parsed := &incidentDTOs.IncidentReportParsedDates{}

	var err error

	if dto.StartDate != nil {
		parsed.StartDate, err = utility.ParseDate(*dto.StartDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
		}
	}

	if dto.EndDate != nil {
		parsed.EndDate, err = utility.ParseDate(*dto.EndDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end date")
		}
	}

	return parsed, nil
}

func formatTitle(sectionType string) (string, string) {
	switch sectionType {
	case "executive-summary":
		return "Executive Summary", "الملخص التنفيذي"
	case "scope":
		return "Scope", "النطاق"
	case "objectives":
		return "Objectives", "الأهداف"
	case "overview":
		return "Overview", "نظرة عامة"
	case "details":
		return "Details", "التفاصيل"
	case "investigation":
		return "Investigation", "التحقيق"
	case "incident-closure":
		return "Incident Closure", "إغلاق الحادث"
	case "post-incident-review":
		return "Post-Incident Review", "مراجعة ما بعد الحادث"
	case "rca":
		return "Root Cause Analysis", "تحليل السبب الجذري"
	case "communication-log":
		return "Communication Log", "سجل التواصل"
	case "action-plan":
		return "Action Plan", "خطة العمل"
	default:
		return sectionType, sectionType // Default to sectionType if not found
	}
}

func (irs *IncidentReportService) GetIncidentReportDetails(id uint) (*incidentDTOs.IncidentReportDetailsDTO, error) {
	report, err := irs.repo.GetIncidentReportDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	response := &incidentDTOs.IncidentReportDetailsDTO{
		ID:        report.ID,
		Title:     report.Title,
		CreatedAt: report.CreatedAt,
	}

	if report.StartDate != nil {
		startDateStr := report.StartDate.Format("2006-01-02")
		response.StartDate = &startDateStr
	}

	if report.EndDate != nil {
		endDateStr := report.EndDate.Format("2006-01-02")
		response.EndDate = &endDateStr
	}

	if report.CreatedByUser != nil {
		response.CreatedBy = organizationDTOs.UserMinimalDTO{
			ID:               int(report.CreatedByUser.ID),
			FirstNameEN:      report.CreatedByUser.FirstNameEn,
			LastNameEN:       report.CreatedByUser.LastNameEn,
			FirstNameAR:      report.CreatedByUser.FirstNameAr,
			LastNameAR:       report.CreatedByUser.LastNameAr,
			Initial:          report.CreatedByUser.Initial,
			Color:            report.CreatedByUser.Color,
			Email:            report.CreatedByUser.Email,
			ProfileImagePath: report.CreatedByUser.ProfileImagePath,
		}
	}

	response.Sections = make([]*incidentDTOs.IncidentReportSectionDTO, 0, len(report.Sections))
	for _, section := range report.Sections {
		response.Sections = append(response.Sections, &incidentDTOs.IncidentReportSectionDTO{
			ID:            section.ID,
			Type:          section.Type,
			TitleEn:       section.TitleEn,
			TitleAr:       section.TitleAr,
			DescriptionEn: section.DescriptionEn,
			DescriptionAr: section.DescriptionAr,
			OrderNum:      section.OrderNum,
		})
	}

	response.Incidents = make([]*incidentDTOs.IncidentReportIncidentDetailDTO, 0, len(report.Incidents))
	for _, incident := range report.Incidents {
		categoryTitleEn := ""
		categoryTitleAr := ""
		if incident.IncidentCategoryTitleEn != nil {
			categoryTitleEn = *incident.IncidentCategoryTitleEn
		}
		if incident.IncidentCategoryTitleAr != nil {
			categoryTitleAr = *incident.IncidentCategoryTitleAr
		}
		severityTitleEn := ""
		severityTitleAr := ""
		if incident.SeverityLevelTitleEn != nil {
			severityTitleEn = *incident.SeverityLevelTitleEn
		}
		if incident.SeverityLevelTitleAr != nil {
			severityTitleAr = *incident.SeverityLevelTitleAr
		}
		statusTitleEn := ""
		statusTitleAr := ""
		if incident.IncidentStatusTitleEn != nil {
			statusTitleEn = *incident.IncidentStatusTitleEn
		}
		if incident.IncidentStatusTitleAr != nil {
			statusTitleAr = *incident.IncidentStatusTitleAr
		}
		incidentDTO := &incidentDTOs.IncidentReportIncidentDetailDTO{
			ID:                         incident.ID,
			ReferenceNumber:            incident.ReferenceNumber,
			IncidentCategoryTitleEn:    categoryTitleEn,
			IncidentCategoryTitleAr:    categoryTitleAr,
			IncidentSubCategoryTitleEn: incident.IncidentSubCategoryTitleEn,
			IncidentSubCategoryTitleAr: incident.IncidentSubCategoryTitleAr,
			SeverityTitleEn:            severityTitleEn,
			SeverityTitleAr:            severityTitleAr,
			IncidentAt:                 incident.IncidentAt,
			ReportedAt:                 incident.ReportedAt,
			IncidentStatusTitleEn:      statusTitleEn,
			IncidentStatusTitleAr:      statusTitleAr,
		}

		incidentDTO.Documents = make([]*incidentDTOs.IncidentReportDocumentDTO, 0, len(incident.Documents))
		for _, doc := range incident.Documents {
			docURL := ""
			if doc.URL != nil {
				docURL = *doc.URL
			}
			docName := ""
			if doc.Name != nil {
				docName = *doc.Name
			}
			docSize := 0
			if doc.Size != nil {
				docSize = *doc.Size
			}
			docToken := ""
			if doc.Token != nil {
				docToken = *doc.Token
			}
			docOriginalName := ""
			if doc.OriginalName != nil {
				docOriginalName = *doc.OriginalName
			}
			docType := ""
			if doc.Type != nil {
				docType = *doc.Type
			}
			incidentDTO.Documents = append(incidentDTO.Documents, &incidentDTOs.IncidentReportDocumentDTO{
				ID:                       doc.ID,
				IncidentDocumentID:       doc.IncidentDocumentID,
				IncidentReportIncidentID: doc.IncidentReportIncidentID,
				URL:                      docURL,
				Name:                     &docName,
				Size:                     docSize,
				Token:                    docToken,
				OriginalName:             docOriginalName,
				Type:                     docType,
			})
		}

		incidentDTO.RootCauseAnalyses = make([]*incidentDTOs.IncidentReportRootCauseAnalysisDTO, 0, len(incident.RootCauseAnalyses))
		for _, rca := range incident.RootCauseAnalyses {
			answer := ""
			if rca.Answer != nil {
				answer = *rca.Answer
			}
			rcaDTO := &incidentDTOs.IncidentReportRootCauseAnalysisDTO{
				ID:       rca.ID,
				Why:      rca.Why,
				Question: rca.Question,
				Answer:   answer,
			}
			incidentDTO.RootCauseAnalyses = append(incidentDTO.RootCauseAnalyses, rcaDTO)
		}

		incidentDTO.ActionPlans = make([]*incidentDTOs.IncidentReportActionPlanDTO, 0, len(incident.ActionPlans))
		for _, plan := range incident.ActionPlans {
			planDTO := &incidentDTOs.IncidentReportActionPlanDTO{
				ID:          plan.ID,
				Title:       plan.ActionPlanTitle,
				Description: plan.ActionPlanDescription,
				TargetDate:  plan.TargetDate,
				StatusTitle: plan.ActionPlanStatusTitleEn,
			}
			incidentDTO.ActionPlans = append(incidentDTO.ActionPlans, planDTO)
		}

		// Map Investigation, Communication Logs, Closure, and PIR from IncidentRegister
		if incident.IncidentRegister != nil {
			// Map Investigation
			if incident.IncidentRegister.Investigation != nil {
				var investigationDTO incidentDTOs.IncidentInvestigationDetailDTO
				if err := mapstructure.Decode(incident.IncidentRegister.Investigation, &investigationDTO); err == nil {
					if incident.IncidentRegister.Investigation.InitiatedAt != nil {
						initiatedAtStr := incident.IncidentRegister.Investigation.InitiatedAt.Format("2006-01-02")
						investigationDTO.InitiatedAtString = &initiatedAtStr
					}
					incidentDTO.Investigation = &investigationDTO
				}
			}

			// Map Communication Logs
			if len(incident.IncidentRegister.CommunicationLogs) > 0 {
				incidentDTO.CommunicationLogs = make([]*incidentDTOs.IncidentCommunicationLogDetailDTO, 0, len(incident.IncidentRegister.CommunicationLogs))
				for _, commLog := range incident.IncidentRegister.CommunicationLogs {
					var commLogDTO incidentDTOs.IncidentCommunicationLogDetailDTO
					if err := mapstructure.Decode(commLog, &commLogDTO); err == nil {
						incidentDTO.CommunicationLogs = append(incidentDTO.CommunicationLogs, &commLogDTO)
					}
				}
			}

			// Map Incident Closure
			if incident.IncidentRegister.IncidentClosure != nil {
				var closureDTO incidentDTOs.IncidentClosureDetailDTO
				if err := mapstructure.Decode(incident.IncidentRegister.IncidentClosure, &closureDTO); err == nil {
					if incident.IncidentRegister.IncidentClosure.CompletionDate != nil {
						completionDateStr := incident.IncidentRegister.IncidentClosure.CompletionDate.Format("2006-01-02")
						closureDTO.CompletionDateString = &completionDateStr
					}
					incidentDTO.IncidentClosure = &closureDTO
				}
			}

			// Map Post Incident Review (PIR)
			if incident.IncidentRegister.IncidentPostReview != nil {
				var pirDTO incidentDTOs.IncidentPostReviewDetailDTO
				if err := mapstructure.Decode(incident.IncidentRegister.IncidentPostReview, &pirDTO); err == nil {
					if incident.IncidentRegister.IncidentPostReview.PostReviewDate != nil {
						postReviewDateStr := incident.IncidentRegister.IncidentPostReview.PostReviewDate.Format("2006-01-02")
						pirDTO.PostReviewDateString = &postReviewDateStr
					}
					incidentDTO.IncidentPostReview = &pirDTO
				}
			}
		}

		response.Incidents = append(response.Incidents, incidentDTO)
	}

	return response, nil
}

func (irs *IncidentReportService) DeleteIncidentReport(incidentReportID uint, currentUserID uint) error {
	tx := irs.repo.BeginTransaction()

	existingReport, err := irs.repo.GetIncidentReportByID(incidentReportID, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, incidentReportID)
		}
		return err
	}

	if err := irs.repo.DeleteIncidentReportDocuments(incidentReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := irs.repo.DeleteIncidentReportRootCauseAnalyses(incidentReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := irs.repo.DeleteIncidentReportActionPlans(incidentReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := irs.repo.DeleteIncidentReportIncidents(incidentReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := irs.repo.DeleteIncidentReportSections(incidentReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := irs.repo.DeleteIncidentReport(incidentReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	go func() {
		err := irs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.IncidentReport{}.TableName(),
			incidentReportID,
			existingReport,
			currentUserID,
			"incident", "incident_report", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	tx.Commit()
	return nil
}

func (irs *IncidentReportService) GetSectionByID(incidentReportID uint, sectionID uint, tx ...utility.DatabaseTransaction) (*models.IncidentReportSection, error) {
	section, err := irs.repo.GetSectionByID(sectionID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, sectionID)
		}
		return nil, err
	}

	// Validate that section belongs to the report (security improvement)
	if section.IncidentReportID != incidentReportID {
		return nil, utility.NewResourceNotFoundError(
			errors.New("section does not belong to this report"),
			sectionID,
		)
	}

	return section, nil
}

func (irs *IncidentReportService) UpdateSection(incidentReportID uint, sectionID uint, currentUserID uint, formDTO *incidentDTOs.IncidentReportSectionFormDTO) error {
	tx := irs.repo.BeginTransaction()

	existingSection, err := irs.GetSectionByID(incidentReportID, sectionID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Capture old values before update
	oldValues := map[string]interface{}{
		"title_en":       existingSection.TitleEn,
		"title_ar":       existingSection.TitleAr,
		"description_en": existingSection.DescriptionEn,
		"description_ar": existingSection.DescriptionAr,
	}

	// Update section details
	existingSection.TitleEn = formDTO.TitleEn
	existingSection.TitleAr = formDTO.TitleAr
	if formDTO.DescriptionEn != nil {
		existingSection.DescriptionEn = *formDTO.DescriptionEn
	}
	if formDTO.DescriptionAr != nil {
		existingSection.DescriptionAr = *formDTO.DescriptionAr
	}

	err = irs.repo.UpdateSection(existingSection, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Capture new values after update
	newValues := map[string]interface{}{
		"title_en":       existingSection.TitleEn,
		"title_ar":       existingSection.TitleAr,
		"description_en": existingSection.DescriptionEn,
		"description_ar": existingSection.DescriptionAr,
	}

	// Identify changes
	changedData := utility.CompareFieldChanges(oldValues, newValues)

	// Log activity only if there are changes
	if len(changedData) > 0 {
		go func() {
			err := irs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.IncidentReportSection{}.TableName(),
				sectionID,
				changedData,
				currentUserID,
				"incident", "incident_report", "", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}

	return nil
}

func (irs *IncidentReportService) DeleteIncidentReportSection(incidentReportID uint, sectionID uint, currentUserID uint) error {
	tx := irs.repo.BeginTransaction()

	section, err := irs.GetSectionByID(incidentReportID, sectionID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the section
	if err := irs.repo.DeleteIncidentReportSection(sectionID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Reorder remaining sections
	if err := irs.repo.ReorderIncidentReportSections(incidentReportID, section.OrderNum, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		err := irs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.IncidentReportSection{}.TableName(),
			sectionID,
			section,
			currentUserID,
			"incident", "incident_report", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()
	return nil
}

func (irs *IncidentReportService) GetIncidentReportExportZip(incidentReportID uint, currentUserID uint) (string, error) {
	// TODO: PDF template is not ready yet
	// When template is ready, uncomment and implement the following:
	// 1. Fetch incident report detail: incidentReport, err := irs.GetIncidentReportDetails(incidentReportID)
	// 2. Build section data with incidents (similar to risk report pattern)
	// 3. Generate PDF content using template
	// 4. Convert HTML to PDF
	// 5. Collect all documents from all incidents
	// 6. Create ZIP file with PDF + documents

	// For now, return error indicating template is not ready
	return "", fmt.Errorf("PDF template for incident report export is not ready yet")

	// Placeholder code structure (to be implemented when template is ready):
	/*
		var sectionData []incidentDTOs.IncidentReportSectionWithDataDTO

		for _, section := range incidentReport.Sections {
			sectionItem := incidentDTOs.IncidentReportSectionWithDataDTO{
				ID:            section.ID,
				Type:          section.Type,
				TitleEn:       section.TitleEn,
				TitleAr:       section.TitleAr,
				DescriptionEn: utility.StripPTags(section.DescriptionEn),
				DescriptionAr: utility.StripPTags(section.DescriptionAr),
			}

			// Map incidents to sections based on section type
			for _, incident := range incidentReport.Incidents {
				// Build incident item for PDF
				// sectionItem.Incidents = append(sectionItem.Incidents, incidentItem)
			}

			sectionData = append(sectionData, sectionItem)
		}

		parsedCreatedAt, err := time.Parse(time.RFC3339, incidentReport.CreatedAt.Format(time.RFC3339))
		if err != nil {
			return "", fmt.Errorf("invalid CreatedAt format: %w", err)
		}

		// Final DTO to render
		reportDetail := incidentDTOs.IncidentReportPDFDetailDTO{
			ID:                       incidentReport.ID,
			Title:                    incidentReport.Title,
			StartDate:                incidentReport.StartDate,
			EndDate:                  incidentReport.EndDate,
			CreatedUserFirstNameEn:   incidentReport.CreatedBy.FirstNameEN,
			CreatedUserFirstNameAr:   incidentReport.CreatedBy.FirstNameAR,
			CreatedUserLastNameEn:    incidentReport.CreatedBy.LastNameEN,
			CreatedUserLastNameAr:     incidentReport.CreatedBy.LastNameAR,
			CreatedAt:                parsedCreatedAt.Format("2006-01-02"),
			Sections:                 sectionData,
		}

		// Generate PDF content
		templatePaths := []string{"incident/templates/pdf/incident_report_summary.html"}
		htmlContent, err := utility.GeneratePDFContent(reportDetail, templatePaths, "", "")
		if err != nil {
			return "", fmt.Errorf("failed to generate HTML content: %w", err)
		}

		// Convert to PDF
		pdfFileName := fmt.Sprintf("IncidentReport_%d.pdf", incidentReportID)
		if err := utility.ConvertHTMLContentToPDF(htmlContent, pdfFileName); err != nil {
			return "", fmt.Errorf("failed to generate PDF: %w", err)
		}

		filesToZip := []string{pdfFileName}
		// Collect all documents from all incidents
		for _, incident := range incidentReport.Incidents {
			for _, doc := range incident.Documents {
				path, err := irs.GetDocumentPath(incidentReportID, doc.Token)
				if err != nil {
					fmt.Printf("Error getting document path for ID %d: %v\n", doc.ID, err)
					continue
				}
				filesToZip = append(filesToZip, path)
			}
		}

		// Create ZIP archive
		zipDir := filepath.Join("storage", "zip")
		zipFilePath := filepath.Join(zipDir, fmt.Sprintf("incident_report_%d.zip", incidentReportID))

		// Ensure the directory exists
		if err := os.MkdirAll(zipDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %w", zipDir, err)
		}

		err = utility.CreateZipFile(zipFilePath, filesToZip)
		if err != nil {
			return "", fmt.Errorf("failed to create ZIP file: %w", err)
		}

		go func() {
			time.Sleep(2 * time.Second)
			// Delete the temporary file
			if err := os.Remove(pdfFileName); err != nil {
				utility.LogError("Failed to delete temp file", map[string]interface{}{"file": pdfFileName, "error": err.Error()})
			}
		}()

		return zipFilePath, nil
	*/
}

func (irs *IncidentReportService) GetDocumentPath(id uint, token string) (string, error) {
	// Get full details of the incident report by ID
	existingIncidentReport, err := irs.GetIncidentReportDetails(id)
	if err != nil {
		return "", err
	}

	// Check if there are any associated documents in any incident
	hasDocuments := false
	for _, incident := range existingIncidentReport.Incidents {
		if len(incident.Documents) > 0 {
			hasDocuments = true
			break
		}
	}

	if !hasDocuments {
		return "", fmt.Errorf("no documents found for the incident report")
	}

	// Decrypt the token to get the file name
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	// Define the storage path
	storagePath := filepath.Join("incident", "incident-report", strconv.FormatUint(uint64(id), 10), "documents")

	// Get the base storage path from environment variables
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists at the specified path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
