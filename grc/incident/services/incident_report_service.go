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

	vulnerabilities, err := irs.repo.GetAllIncidentVulnerabilities(id)
	if err != nil {
		return nil, err
	}
	incidentDTO.Vulnerabilities = vulnerabilities

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

func (irs *IncidentReportService) Create(currentUserID uint, dto incidentDTOs.IncidentReportRequestDTO) (uint, error) {
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

		if dates.StartDate.After(*dates.EndDate) {
			tx.Rollback()
			return 0, utility.NewValidationError(
				errors.New("start_date must be less than or equal to end_date"),
				"invalid",
				"start_date",
			)
		}

		incidents, err = irs.repo.GetReportIncidentsByDateRange(
			*dates.StartDate,
			*dates.EndDate,
			tx,
		)

		if len(incidents) == 0 {
			tx.Rollback()
			return 0, utility.NewValidationError(
				errors.New("no incidents found for the specified criteria"),
				"empty_result",
				"date_range",
			)
		}

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

	var incidentDocumentsMap = make(map[uint][]incidentDTOs.IncidentReportDocumentDTO)
	var incidentEvidencesMap = make(map[uint][]incidentDTOs.IncidentReportEvidenceDTO)

	for _, inc := range incidents {
		reportIncident := &models.IncidentReportIncident{
			IncidentReportID:                         reportID,
			IncidentRegisterID:                       inc.ID,
			ReferenceNumber:                          &inc.ReferenceNumber,
			IncidentTitle:                            inc.IncidentTitle,
			IncidentCategoryTitleEn:                  &inc.IncidentCategoryTitleEn,
			IncidentCategoryTitleAr:                  &inc.IncidentCategoryTitleAr,
			IncidentSubCategoryTitleEn:               &inc.IncidentSubCategoryTitleEn,
			IncidentSubCategoryTitleAr:               &inc.IncidentSubCategoryTitleAr,
			IncidentTypeTitleEn:                      (*string)(&inc.IncidentTypeTitleEn),
			IncidentTypeTitleAr:                      (*string)(&inc.IncidentTypeTitleAr),
			SeverityLevelTitleEn:                     &inc.SeverityTitleEn,
			SeverityLevelTitleAr:                     &inc.SeverityTitleAr,
			SeverityLevelColorCode:                   inc.SeverityColorCode,
			IncidentAt:                               &inc.RealTime,
			ReportedAt:                               &inc.DetectedTime,
			IncidentLocation:                         inc.IncidentLocation,
			IncidentStatusTitleEn:                    &inc.IncidentStatusTitleEn,
			IncidentStatusTitleAr:                    &inc.IncidentStatusTitleAr,
			IncidentStatusColorCode:                  (*string)(&inc.IncidentStatusColorCode),
			ResponsibleUserNameEn:                    inc.ResponsibleUserNameEn,
			ResponsibleUserNameAr:                    inc.ResponsibleUserNameAr,
			ReportedUserNameEn:                       inc.ReportedUserNameEn,
			ReportedUserNameAr:                       inc.ReportedUserNameAr,
			ResponsibleUnitTitleEn:                   (*string)(&inc.ResponsibleUnitTitleEn),
			ResponsibleUnitTitleAr:                   (*string)(&inc.ResponsibleUnitTitleAr),
			ResponsibleUnitColorCode:                 (*string)(&inc.ResponsibleUnitColorCode),
			BusinessOwnerNameEn:                      (*string)(&inc.BusinessOwnerNameEn),
			BusinessOwnerNameAr:                      (*string)(&inc.BusinessOwnerNameAr),
			ResponseTeamLeadNameEn:                   (*string)(&inc.ResponseTeamLeadNameEn),
			ResponseTeamLeadNameAr:                   (*string)(&inc.ResponseTeamLeadNameAr),
			ResponseTeamMembersEn:                    (*string)(&inc.ResponseTeamMembersEn),
			ResponseTeamMembersAr:                    (*string)(&inc.ResponseTeamMembersAr),
			IncidentScopeTitleEn:                     (*string)(&inc.IncidentScopeTitleEn),
			IncidentScopeTitleAr:                     (*string)(&inc.IncidentScopeTitleAr),
			IncidentSensitivityClassificationTitleEn: (*string)(&inc.IncidentSensitivityClassificationTitleEn),
			IncidentSensitivityClassificationTitleAr: (*string)(&inc.IncidentSensitivityClassificationTitleAr),
			ConfidentialitySeverityLevelTitleEn:      (*string)(&inc.ConfidentialitySeverityLevelTitleEn),
			ConfidentialitySeverityLevelTitleAr:      (*string)(&inc.ConfidentialitySeverityLevelTitleAr),
			ConfidentialitySeverityLevelColorCode:    (*string)(&inc.ConfidentialitySeverityLevelColorCode),
			IntegritySeverityLevelTitleEn:            (*string)(&inc.IntegritySeverityLevelTitleEn),
			IntegritySeverityLevelTitleAr:            (*string)(&inc.IntegritySeverityLevelTitleAr),
			IntegritySeverityLevelColorCode:          (*string)(&inc.IntegritySeverityLevelColorCode),
			AvailabilitySeverityLevelTitleEn:         (*string)(&inc.AvailabilitySeverityLevelTitleEn),
			AvailabilitySeverityLevelTitleAr:         (*string)(&inc.AvailabilitySeverityLevelTitleAr),
			AvailabilitySeverityLevelColorCode:       (*string)(&inc.AvailabilitySeverityLevelColorCode),
			Description:                              (*string)(&inc.Description),
			QuickAction:                              (*string)(&inc.QuickAction),
			IsAcknowledge:                            inc.IsAcknowledge,
			AcknowledgedAt:                           inc.AcknowledgedAt,
			AcknowledgedByNameEn:                     (*string)(&inc.AcknowledgedByNameEn),
			AcknowledgedByNameAr:                     (*string)(&inc.AcknowledgedByNameAr),
			IsSLATimeAchieved:                        inc.IsSLATimeAchieved,
			IsSLATimeExceeded:                        inc.IsSLATimeExceeded,
			ClosureResolutionSummary:                 (*string)(&inc.ClosureResolutionSummary),
			ClosureNotes:                             (*string)(&inc.ClosureNotes),
			ClosureSubmittedAt:                       inc.ClosureSubmittedAt,
			ClosureSubmittedByNameEn:                 (*string)(&inc.ClosureSubmittedByNameEn),
			ClosureSubmittedByNameAr:                 (*string)(&inc.ClosureSubmittedByNameAr),
			ClosureApprovedAt:                        inc.ClosureApprovedAt,
			ClosureApprovedByNameEn:                  (*string)(&inc.ClosureApprovedByNameEn),
			ClosureApprovedByNameAr:                  (*string)(&inc.ClosureApprovedByNameAr),
			ClosureWorkflowStatusTitleEn:             (*string)(&inc.ClosureWorkflowStatusTitleEn),
			ClosureWorkflowStatusTitleAr:             (*string)(&inc.ClosureWorkflowStatusTitleAr),
			ClosureWorkflowStatusColorCode:           (*string)(&inc.ClosureWorkflowStatusColorCode),
			ClosedDate:                               inc.ResolvedDate,
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

		incidentDocumentsMap[inc.ID] = documents

		var reportDocs []*models.IncidentReportDocument

		storagePath := config.GetEnv("STORAGE_PATH", "")

		for _, document := range documents {
			if document.Name == nil || *document.Name == "" {
				continue // Skip documents without a name
			}

			oldFilePath := filepath.Join(
				storagePath,
				"incident",
				"incident-register",
				strconv.FormatUint(uint64(inc.ID), 10),
				"documents",
				*document.Name,
			)

			if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
				fmt.Printf("Warning: Document file not found, skipping: %s\n", oldFilePath)
				continue // Skip missing files instead of failing
			}

			newURL := utility.GenerateDocumentURL(
				"incident",
				"incident-report",
				strconv.FormatUint(uint64(reportID), 10),
				"documents",
			)

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
				Size:                     docSize,
				Token:                    docToken,
				OriginalName:             docOriginalName,
				Type:                     docType,
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

		evidences, err := irs.repo.GetEvidencesByIncidentID(inc.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		incidentEvidencesMap[inc.ID] = evidences

		var reportEvidences []*models.IncidentReportEvidence

		for _, evidence := range evidences {
			if evidence.Name == nil || *evidence.Name == "" {
				continue // Skip evidences without a name
			}

			oldEvidenceFilePath := filepath.Join(
				storagePath,
				"incident",
				"incident-register",
				strconv.FormatUint(uint64(inc.ID), 10),
				"evidences",
				*evidence.Name,
			)

			if _, err := os.Stat(oldEvidenceFilePath); os.IsNotExist(err) {
				fmt.Printf("Warning: Evidence file not found, skipping: %s\n", oldEvidenceFilePath)
				continue // Skip missing files instead of failing
			}

			newEvidenceURL := utility.GenerateDocumentURL(
				"incident",
				"incident-report",
				strconv.FormatUint(uint64(reportID), 10),
				"evidences",
			)

			evidenceName := *evidence.Name
			evidenceSize := evidence.Size
			evidenceToken := evidence.Token
			evidenceOriginalName := evidence.OriginalName
			evidenceType := evidence.Type
			reportEvidences = append(reportEvidences, &models.IncidentReportEvidence{
				IncidentEvidenceID:       evidence.ID,
				IncidentReportIncidentID: reportIncident.ID,
				URL:                      &newEvidenceURL,
				Name:                     &evidenceName,
				Size:                     evidenceSize,
				Token:                    evidenceToken,
				OriginalName:             evidenceOriginalName,
				Type:                     evidenceType,
				CreatedByID:              &currentUserID,
				UpdatedByID:              nil,
			})
		}

		if len(reportEvidences) > 0 {
			if err := irs.repo.CreateIncidentEvidences(reportEvidences, tx); err != nil {
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

			// Responsible User Details
			//var responsibleUserID *uint
			responsibleUserNameEn := ""
			responsibleUserNameAr := ""
			//var responsibleUserInitial *string
			//var responsibleUserColor *string
			//var responsibleUserEmail *string
			//var responsibleUserProfileImagePath *string

			if p.ResponsibleUser != nil {
				//responsibleUserID = &p.ResponsibleUser.ID
				responsibleUserNameEn = fmt.Sprintf("%s %s", p.ResponsibleUser.FirstNameEn, p.ResponsibleUser.LastNameEn)
				responsibleUserNameAr = fmt.Sprintf("%s %s", p.ResponsibleUser.FirstNameAr, p.ResponsibleUser.LastNameAr)
				//responsibleUserInitial = &p.ResponsibleUser.Initial
				//responsibleUserColor = &p.ResponsibleUser.Color
				//responsibleUserEmail = &p.ResponsibleUser.Email
				//responsibleUserProfileImagePath = p.ResponsibleUser.ProfileImagePath
			}

			// Severity Level Details
			//var severityLevelID *uint
			var severityLevelTitleEn *string
			var severityLevelTitleAr *string
			//var severityLevelColorCode *string
			//var severityLevelType *string

			if p.SeverityLevel != nil {
				//severityLevelID = &p.SeverityLevel.ID
				severityLevelTitleEn = &p.SeverityLevel.TitleEn
				severityLevelTitleAr = &p.SeverityLevel.TitleAr
				//severityLevelColorCode = &p.SeverityLevel.ColorCode
				//severityLevelType = &p.SeverityLevel.Type
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
				TargetDate:            targetDate,

				// Responsible User
				//ResponsibleUserID:               responsibleUserID,
				ResponsibleUserNameEn: responsibleUserNameEn,
				ResponsibleUserNameAr: responsibleUserNameAr,
				// ResponsibleUserInitial:          responsibleUserInitial,
				// ResponsibleUserColor:            responsibleUserColor,
				// ResponsibleUserEmail:            responsibleUserEmail,
				// ResponsibleUserProfileImagePath: responsibleUserProfileImagePath,

				// Severity Level
				//SeverityLevelID:        severityLevelID,
				SeverityLevelTitleEn: severityLevelTitleEn,
				SeverityLevelTitleAr: severityLevelTitleAr,
				// SeverityLevelColorCode: severityLevelColorCode,
				// SeverityLevelType:      severityLevelType,

				// Action Plan Progress
				//ActionPlanPercentage: p.Percentage,

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

	// Copy files after transaction commits to avoid orphaned files on rollback
	storagePath := config.GetEnv("STORAGE_PATH", "")
	newStoragePath := filepath.Join(
		storagePath,
		"incident",
		"incident-report",
		strconv.FormatUint(uint64(reportID), 10),
		"documents",
	)

	if err := os.MkdirAll(newStoragePath, os.ModePerm); err != nil {
		_ = irs.DeleteIncidentReport(reportID, currentUserID)
		return 0, fmt.Errorf("failed to create storage directory: %w", err)
	}

	for _, inc := range incidents {
		documents := incidentDocumentsMap[inc.ID]

		for _, document := range documents {
			if document.Name == nil || *document.Name == "" {
				continue
			}

			oldFilePath := filepath.Join(
				storagePath,
				"incident",
				"incident-register",
				strconv.FormatUint(uint64(inc.ID), 10),
				"documents",
				*document.Name,
			)

			newFilePath := filepath.Join(newStoragePath, *document.Name)

			if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
				fmt.Printf("Warning: Source document not found during copy, skipping: %s\n", oldFilePath)
				continue
			}

			input, err := os.Open(oldFilePath)
			if err != nil {
				fmt.Printf("Warning: Failed to open source document, skipping: %s - %v\n", oldFilePath, err)
				continue
			}

			output, err := os.Create(newFilePath)
			if err != nil {
				input.Close()
				_ = irs.DeleteIncidentReport(reportID, currentUserID)
				return 0, fmt.Errorf("failed to create destination file %s: %w", newFilePath, err)
			}

			_, err = io.Copy(output, input)
			output.Close()
			input.Close()

			if err != nil {
				_ = irs.DeleteIncidentReport(reportID, currentUserID)
				return 0, fmt.Errorf("failed to copy file from %s to %s: %w", oldFilePath, newFilePath, err)
			}
		}
	}

	newEvidenceStoragePath := filepath.Join(
		storagePath,
		"incident",
		"incident-report",
		strconv.FormatUint(uint64(reportID), 10),
		"evidences",
	)

	if err := os.MkdirAll(newEvidenceStoragePath, os.ModePerm); err != nil {
		_ = irs.DeleteIncidentReport(reportID, currentUserID)
		return 0, fmt.Errorf("failed to create evidence storage directory: %w", err)
	}

	for _, inc := range incidents {
		evidences := incidentEvidencesMap[inc.ID]

		for _, evidence := range evidences {
			if evidence.Name == nil || *evidence.Name == "" {
				continue
			}

			oldEvidenceFilePath := filepath.Join(
				storagePath,
				"incident",
				"incident-register",
				strconv.FormatUint(uint64(inc.ID), 10),
				"evidences",
				*evidence.Name,
			)

			newEvidenceFilePath := filepath.Join(newEvidenceStoragePath, *evidence.Name)

			if _, err := os.Stat(oldEvidenceFilePath); os.IsNotExist(err) {
				fmt.Printf("Warning: Source evidence not found during copy, skipping: %s\n", oldEvidenceFilePath)
				continue
			}

			input, err := os.Open(oldEvidenceFilePath)
			if err != nil {
				fmt.Printf("Warning: Failed to open source evidence, skipping: %s - %v\n", oldEvidenceFilePath, err)
				continue
			}

			output, err := os.Create(newEvidenceFilePath)
			if err != nil {
				input.Close()
				_ = irs.DeleteIncidentReport(reportID, currentUserID)
				return 0, fmt.Errorf("failed to create destination evidence file %s: %w", newEvidenceFilePath, err)
			}

			_, err = io.Copy(output, input)
			output.Close()
			input.Close()

			if err != nil {
				_ = irs.DeleteIncidentReport(reportID, currentUserID)
				return 0, fmt.Errorf("failed to copy evidence file from %s to %s: %w", oldEvidenceFilePath, newEvidenceFilePath, err)
			}
		}
	}

	go func() {
		report, err := irs.repo.GetIncidentReportByID(reportID)
		if err != nil {
			utility.LogError("Failed to get incident report", map[string]interface{}{"error": err.Error()})
			return
		}

		err = irs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentReport{}.TableName(),
			reportID,
			report,
			currentUserID,
			"incident", "incident_report", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

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

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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

	// Clear preloaded action plans before decode - they will be fetched live via GetAllIncidentActionPlan
	for i := range report.Incidents {
		report.Incidents[i].ActionPlans = nil
	}

	// Use single mapstructure.Decode() for the entire structure - DONE!
	var response incidentDTOs.IncidentReportDetailsDTO
	if err := mapstructure.Decode(report, &response); err != nil {
		return nil, err
	}

	// Flatten IncidentRegister fields into Incident level for API response
	for i, incident := range response.Incidents {
		if incident.IncidentRegister != nil {
			incident.Investigation = incident.IncidentRegister.Investigation
			incident.CommunicationLogs = incident.IncidentRegister.CommunicationLogs
			incident.IncidentClosure = incident.IncidentRegister.IncidentClosure
			incident.IncidentPostReview = incident.IncidentRegister.IncidentPostReview
			incident.IncidentRegister = nil // Hide nested structure from API response
		}

		// Map ResponseTeamLead and ResponseTeamMembers from preloaded IncidentRegister
		if i < len(report.Incidents) && report.Incidents[i].IncidentRegister != nil {
			reg := report.Incidents[i].IncidentRegister
			if reg.ResponseTeamLead != nil {
				incident.ResponseTeamLead = &organizationDTOs.UserMinimalDTO{
					ID:               int(reg.ResponseTeamLead.ID),
					FirstNameEN:      reg.ResponseTeamLead.FirstNameEn,
					LastNameEN:       reg.ResponseTeamLead.LastNameEn,
					FirstNameAR:      reg.ResponseTeamLead.FirstNameAr,
					LastNameAR:       reg.ResponseTeamLead.LastNameAr,
					Email:            reg.ResponseTeamLead.Email,
					Initial:          reg.ResponseTeamLead.Initial,
					Color:            reg.ResponseTeamLead.Color,
					ProfileImagePath: reg.ResponseTeamLead.ProfileImagePath,
				}
			}
			if len(reg.ResponseTeamMembers) > 0 {
				for _, m := range reg.ResponseTeamMembers {
					if m.User != nil {
						member := &incidentDTOs.IncidentRegisterResponseTeamMemberDTO{
							ID:                 m.ID,
							IncidentRegisterID: m.IncidentRegisterID,
							UserID:             m.UserID,
							User: organizationDTOs.UserMinimalDTO{
								ID:               int(m.User.ID),
								FirstNameEN:      m.User.FirstNameEn,
								LastNameEN:       m.User.LastNameEn,
								FirstNameAR:      m.User.FirstNameAr,
								LastNameAR:       m.User.LastNameAr,
								Email:            m.User.Email,
								Initial:          m.User.Initial,
								Color:            m.User.Color,
								ProfileImagePath: m.User.ProfileImagePath,
							},
						}
						incident.ResponseTeamMembers = append(incident.ResponseTeamMembers, member)
					}
				}
			}
		}

		// Build timeline from the raw preloaded model (flat DTO requires manual mapping)
		if i < len(report.Incidents) && report.Incidents[i].IncidentRegister != nil {
			for _, t := range report.Incidents[i].IncidentRegister.Timelines {
				timelineItem := &incidentDTOs.IncidentTimelineListDTO{
					ID:                          t.ID,
					IncidentTimelineTypeTitleEn: t.IncidentTimelineType.TitleEn,
					IncidentTimelineTypeTitleAr: t.IncidentTimelineType.TitleAr,
					TimelineStatusTitleEn:       t.TimelineStatus.TitleEn,
					TimelineStatusTitleAr:       t.TimelineStatus.TitleAr,
					PerformedByFirstNameEn:      &t.PerformedBy.FirstNameEn,
					PerformedByLastNameEn:       &t.PerformedBy.LastNameEn,
					PerformedAt:                 t.PerformedAt,
				}
				incident.Timelines = append(incident.Timelines, timelineItem)
			}
		}

		// Fetch live associated items using the underlying incident register ID
		if incident.IncidentRegisterID != 0 {
			if actionPlans, err := irs.repo.GetAllIncidentActionPlan(incident.IncidentRegisterID); err == nil {
				incident.ActionPlans = actionPlans
			}
			if kpis, err := irs.repo.GetAllIncidentKPI(incident.IncidentRegisterID); err == nil {
				incident.KPIs = kpis
			}
			if strategicObjectives, err := irs.repo.GetAllIncidentStrategicObjectives(incident.IncidentRegisterID); err == nil {
				incident.StrategicObjectives = strategicObjectives
			}
			if assets, err := irs.repo.GetAllIncidentAssets(incident.IncidentRegisterID); err == nil {
				incident.Assets = assets
			}
			if associatedDocuments, err := irs.repo.GetAllIncidentAssociatedDocuments(incident.IncidentRegisterID); err == nil {
				incident.AssociatedDocuments = associatedDocuments
			}
			if controls, err := irs.repo.GetAllIncidentControls(incident.IncidentRegisterID); err == nil {
				incident.ControlsAssociated = controls
			}
			if risks, err := irs.repo.GetAllIncidentRisks(incident.IncidentRegisterID); err == nil {
				incident.Risks = risks
			}
			if processes, err := irs.repo.GetAllIncidentProcesses(incident.IncidentRegisterID); err == nil {
				incident.Processes = processes
			}
			if complianceRequirements, err := irs.repo.GetAllIncidentComplianceRequirements(incident.IncidentRegisterID); err == nil {
				incident.ComplianceRequirements = complianceRequirements
			}
			if vulnerabilities, err := irs.repo.GetAllIncidentVulnerabilities(incident.IncidentRegisterID); err == nil {
				incident.Vulnerabilities = vulnerabilities
			}
		}
	}

	return &response, nil
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

	// Fetch documents before deleting from DB
	var documentsToDelete []*models.IncidentReportDocument
	db := tx.DB()
	if err := db.Where("incident_report_incident_id IN (SELECT id FROM incident_report_incidents WHERE incident_report_id = ?)", incidentReportID).Find(&documentsToDelete).Error; err != nil {
		tx.Rollback()
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

	// Delete files from storage
	if len(documentsToDelete) > 0 {
		storagePath := filepath.Join("incident", "incident-report", strconv.FormatUint(uint64(incidentReportID), 10), "documents")
		for _, doc := range documentsToDelete {
			if doc.Name != nil && *doc.Name != "" {
				if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
					utility.LogError("Failed to delete file", map[string]interface{}{"file": *doc.Name, "error": err.Error()})
				}
			}
		}
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
	incidentReport, err := irs.GetIncidentReportDetails(incidentReportID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch incident report details: %w", err)
	}

	var incidents []*incidentDTOs.IncidentReportIncidentPDFDTO

	for _, incident := range incidentReport.Incidents {
		var incidentAtFormatted, reportedAtFormatted string
		if incident.IncidentAt != nil {
			incidentAtFormatted = incident.IncidentAt.Format("Jan _2, 2006, 03:04 PM")
		}
		if incident.ReportedAt != nil {
			reportedAtFormatted = incident.ReportedAt.Format("Jan _2, 2006, 03:04 PM")
		}

		var acknowledgedAtString, closedDateString, closureSubmittedAtString, closureApprovedAtString string
		if incident.AcknowledgedAt != nil {
			acknowledgedAtString = incident.AcknowledgedAt.Format("Jan _2, 2006, 03:04 PM")
		}
		if incident.ClosedDate != nil {
			closedDateString = incident.ClosedDate.Format("Jan _2, 2006, 03:04 PM")
		}
		if incident.ClosureSubmittedAt != nil {
			closureSubmittedAtString = incident.ClosureSubmittedAt.Format("Jan _2, 2006, 03:04 PM")
		}
		if incident.ClosureApprovedAt != nil {
			closureApprovedAtString = incident.ClosureApprovedAt.Format("Jan _2, 2006, 03:04 PM")
		}

		incidentItem := &incidentDTOs.IncidentReportIncidentPDFDTO{
			ID:                                       incident.ID,
			ReferenceNumber:                          incident.ReferenceNumber,
			Title:                                    incident.IncidentTitle,
			IncidentCategoryTitleEn:                  incident.IncidentCategoryTitleEn,
			IncidentCategoryTitleAr:                  incident.IncidentCategoryTitleAr,
			IncidentSubCategoryTitleEn:               incident.IncidentSubCategoryTitleEn,
			IncidentSubCategoryTitleAr:               incident.IncidentSubCategoryTitleAr,
			IncidentTypeTitleEn:                      incident.IncidentTypeTitleEn,
			IncidentTypeTitleAr:                      incident.IncidentTypeTitleAr,
			SeverityLevelTitleEn:                     incident.SeverityLevelTitleEn,
			SeverityLevelTitleAr:                     incident.SeverityLevelTitleAr,
			SeverityLevelColorCode:                   incident.SeverityLevelColorCode,
			IncidentAtFormatted:                      incidentAtFormatted,
			ReportedAtFormatted:                      reportedAtFormatted,
			IncidentLocation:                         incident.IncidentLocation,
			IncidentStatusTitleEn:                    incident.IncidentStatusTitleEn,
			IncidentStatusTitleAr:                    incident.IncidentStatusTitleAr,
			IncidentStatusColorCode:                  incident.IncidentStatusColorCode,
			ResponsibleUserNameEn:                    incident.ResponsibleUserNameEn,
			ResponsibleUserNameAr:                    incident.ResponsibleUserNameAr,
			ReportedUserNameEn:                       incident.ReportedUserNameEn,
			ReportedUserNameAr:                       incident.ReportedUserNameAr,
			ResponsibleUnitTitleEn:                   incident.ResponsibleUnitTitleEn,
			ResponsibleUnitTitleAr:                   incident.ResponsibleUnitTitleAr,
			ResponsibleUnitColorCode:                 incident.ResponsibleUnitColorCode,
			BusinessOwnerNameEn:                      safeStr(incident.BusinessOwnerNameEn),
			BusinessOwnerNameAr:                      safeStr(incident.BusinessOwnerNameAr),
			ResponseTeamLeadNameEn:                   safeStr(incident.ResponseTeamLeadNameEn),
			ResponseTeamLeadNameAr:                   safeStr(incident.ResponseTeamLeadNameAr),
			ResponseTeamLead:                         incident.ResponseTeamLead,
			ResponseTeamMembersEn:                    safeStr(incident.ResponseTeamMembersEn),
			ResponseTeamMembersAr:                    safeStr(incident.ResponseTeamMembersAr),
			ResponseTeamMembers:                      incident.ResponseTeamMembers,
			IncidentScopeTitleEn:                     safeStr(incident.IncidentScopeTitleEn),
			IncidentScopeTitleAr:                     safeStr(incident.IncidentScopeTitleAr),
			IncidentSensitivityClassificationTitleEn: safeStr(incident.IncidentSensitivityClassificationTitleEn),
			IncidentSensitivityClassificationTitleAr: safeStr(incident.IncidentSensitivityClassificationTitleAr),
			ConfidentialitySeverityLevelTitleEn:      safeStr(incident.ConfidentialitySeverityLevelTitleEn),
			ConfidentialitySeverityLevelTitleAr:      safeStr(incident.ConfidentialitySeverityLevelTitleAr),
			ConfidentialitySeverityLevelColorCode:    safeStr(incident.ConfidentialitySeverityLevelColorCode),
			IntegritySeverityLevelTitleEn:            safeStr(incident.IntegritySeverityLevelTitleEn),
			IntegritySeverityLevelTitleAr:            safeStr(incident.IntegritySeverityLevelTitleAr),
			IntegritySeverityLevelColorCode:          safeStr(incident.IntegritySeverityLevelColorCode),
			AvailabilitySeverityLevelTitleEn:         safeStr(incident.AvailabilitySeverityLevelTitleEn),
			AvailabilitySeverityLevelTitleAr:         safeStr(incident.AvailabilitySeverityLevelTitleAr),
			AvailabilitySeverityLevelColorCode:       safeStr(incident.AvailabilitySeverityLevelColorCode),
			Description:                              safeStr(incident.Description),
			QuickAction:                              safeStr(incident.QuickAction),
			IsAcknowledge:                            incident.IsAcknowledge,
			AcknowledgedAt:                           incident.AcknowledgedAt,
			AcknowledgedAtString:                     acknowledgedAtString,
			AcknowledgedByNameEn:                     safeStr(incident.AcknowledgedByNameEn),
			AcknowledgedByNameAr:                     safeStr(incident.AcknowledgedByNameAr),
			IsSLATimeAchieved:                        incident.IsSLATimeAchieved,
			IsSLATimeExceeded:                        incident.IsSLATimeExceeded,
			ClosedDateString:                         closedDateString,
			ClosureResolutionSummary:                 safeStr(incident.ClosureResolutionSummary),
			ClosureNotes:                             safeStr(incident.ClosureNotes),
			ClosureSubmittedAtString:                 closureSubmittedAtString,
			ClosureSubmittedByNameEn:                 safeStr(incident.ClosureSubmittedByNameEn),
			ClosureSubmittedByNameAr:                 safeStr(incident.ClosureSubmittedByNameAr),
			ClosureApprovedAtString:                  closureApprovedAtString,
			ClosureApprovedByNameEn:                  safeStr(incident.ClosureApprovedByNameEn),
			ClosureApprovedByNameAr:                  safeStr(incident.ClosureApprovedByNameAr),
			ClosureWorkflowStatusTitleEn:             safeStr(incident.ClosureWorkflowStatusTitleEn),
			ClosureWorkflowStatusTitleAr:             safeStr(incident.ClosureWorkflowStatusTitleAr),
			ClosureWorkflowStatusColorCode:           safeStr(incident.ClosureWorkflowStatusColorCode),
			Documents:                                incident.Documents,
			Evidences:                                incident.Evidences,
			ActionPlans:                              incident.ActionPlans,
			CommunicationLogs:                        incident.CommunicationLogs,
			Investigation:                            incident.Investigation,
			RootCauseAnalyses:                        incident.RootCauseAnalyses,
			IncidentClosure:                          incident.IncidentClosure,
			IncidentPostReview:                       incident.IncidentPostReview,
			Timelines:                                incident.Timelines,
			KPIs:                                     incident.KPIs,
			StrategicObjectives:                      incident.StrategicObjectives,
			ComplianceRequirements:                   incident.ComplianceRequirements,
			Processes:                                incident.Processes,
			Risks:                                    incident.Risks,
			Controls:                                 incident.ControlsAssociated,
			AssociatedDocuments:                      incident.AssociatedDocuments,
			Assets:                                   incident.Assets,
			Vulnerabilities:                          incident.Vulnerabilities,
		}

		incidents = append(incidents, incidentItem)
	}

	var executiveSummary, scope, objective string
	if len(incidentReport.Sections) > 0 {
		for _, section := range incidentReport.Sections {
			switch section.Type {
			case "executive-summary":
				executiveSummary = section.DescriptionEn
			case "scope":
				scope = section.DescriptionEn
			case "objectives":
				objective = section.DescriptionEn
			}
		}
	}

	var startDateStr, endDateStr string
	if incidentReport.StartDate != nil {
		startDateStr = incidentReport.StartDate.Format("Jan _2, 2006, 03:04 PM")
	}
	if incidentReport.EndDate != nil {
		endDateStr = incidentReport.EndDate.Format("Jan _2, 2006")
	}

	parsedCreatedAt := incidentReport.CreatedAt.Format("2006-01-02")
	reportedAt := incidentReport.CreatedAt.Format("02,Jan 2006")
	generatedAt := time.Now().Format("Jan _2, 2006, 03:04 PM")

	var createdUserFirstNameEn, createdUserFirstNameAr, createdUserLastNameEn, createdUserLastNameAr string
	if incidentReport.CreatedByUser != nil {
		createdUserFirstNameEn = incidentReport.CreatedByUser.FirstNameEn
		createdUserFirstNameAr = incidentReport.CreatedByUser.FirstNameAr
		createdUserLastNameEn = incidentReport.CreatedByUser.LastNameEn
		createdUserLastNameAr = incidentReport.CreatedByUser.LastNameAr
	}

	reportDetail := incidentDTOs.IncidentReportPDFDetailDTO{
		ID:                     incidentReport.ID,
		Title:                  incidentReport.Title,
		StartDate:              startDateStr,
		EndDate:                endDateStr,
		ReportedAt:             reportedAt,
		CreatedUserFirstNameEn: createdUserFirstNameEn,
		CreatedUserFirstNameAr: createdUserFirstNameAr,
		CreatedUserLastNameEn:  createdUserLastNameEn,
		CreatedUserLastNameAr:  createdUserLastNameAr,
		CreatedAt:              parsedCreatedAt,
		GeneratedAt:            generatedAt,
		ExecutiveSummary:       executiveSummary,
		Scope:                  scope,
		Objective:              objective,
		Incidents:              incidents,
	}

	bodyTemplatePaths := []string{"incident/templates/pdf/incident_multiple_report_body.html"}
	htmlContent, err := utility.GeneratePDFContent(reportDetail, bodyTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	headerTemplatePaths := []string{"incident/templates/pdf/incident_multiple_report_header.html"}
	headerContent, err := utility.GeneratePDFContent(reportDetail, headerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	footerTemplatePaths := []string{"incident/templates/pdf/incident_multiple_report_footer.html"}
	footerContent, err := utility.GeneratePDFContent(reportDetail, footerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	coverTemplatePaths := []string{"incident/templates/pdf/incident_multiple_report_cover.html"}
	coverContent, err := utility.GeneratePDFContent(reportDetail, coverTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}

	pdfFileName := fmt.Sprintf("IncidentReport_%d.pdf", incidentReportID)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	filesToZip := []string{pdfFileName}

	// Optimize: Direct file path resolution without DB calls
	storagePath := config.GetEnv("STORAGE_PATH", "")
	baseDocPath := filepath.Join(storagePath, "incident", "incident-report",
		strconv.FormatUint(uint64(incidentReportID), 10), "documents")
	baseEvidencePath := filepath.Join(storagePath, "incident", "incident-report",
		strconv.FormatUint(uint64(incidentReportID), 10), "evidences")

	for _, incident := range incidentReport.Incidents {
		for _, doc := range incident.Documents {
			if doc.Token != nil {
				decryptedFilename, err := utility.DecryptData(*doc.Token)
				if err != nil {
					fmt.Printf("Error decrypting token for doc ID %d: %v\n", doc.ID, err)
					continue
				}
				fullPath := filepath.Join(baseDocPath, string(decryptedFilename))
				if _, err := os.Stat(fullPath); err == nil {
					filesToZip = append(filesToZip, fullPath)
				}
			}
		}

		for _, evidence := range incident.Evidences {
			if evidence.Token != nil {
				decryptedFilename, err := utility.DecryptData(*evidence.Token)
				if err != nil {
					fmt.Printf("Error decrypting token for evidence ID %d: %v\n", evidence.ID, err)
					continue
				}
				fullPath := filepath.Join(baseEvidencePath, string(decryptedFilename))
				if _, err := os.Stat(fullPath); err == nil {
					filesToZip = append(filesToZip, fullPath)
				}
			}
		}
	}

	zipDir := filepath.Join("storage", "zip")
	zipFilePath := filepath.Join(zipDir, fmt.Sprintf("incident_report_%d.zip", incidentReportID))

	if err := os.MkdirAll(zipDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", zipDir, err)
	}

	err = utility.CreateZipFile(zipFilePath, filesToZip)
	if err != nil {
		return "", fmt.Errorf("failed to create ZIP file: %w", err)
	}

	go func() {
		time.Sleep(2 * time.Second)
		if err := os.Remove(pdfFileName); err != nil {
			utility.LogError("Failed to delete temp file", map[string]interface{}{"file": pdfFileName, "error": err.Error()})
		}
	}()

	return zipFilePath, nil
}
