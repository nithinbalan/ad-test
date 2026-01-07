package services

import (
	"errors"
	"fmt"
	"io"
	"jethur/config"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepo "jethur/grc/incident/repositories"
	"jethur/grc/models"
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
	Create(currentUserID uint, dto incidentDTOs.IncidentReportRequestDTO) (uint, error)
}

// Service and implementations
type IncidentReportService struct {
	repo                         incidentRepo.IncidentReportRepositoryInterface
	settigSvc                    settingSvc.ReportTemplateSettingServiceInterface
	incidentRootCauseAnalysisSvc IncidentRootCauseAnalysisServiceInterface
	incidentActionPlanSvc        IncidentActionPlanServiceInterface
}

func NewIncidentReportService(
	repo incidentRepo.IncidentReportRepositoryInterface,
	settigSvc settingSvc.ReportTemplateSettingServiceInterface,
	incidentRootCauseAnalysisSvc IncidentRootCauseAnalysisServiceInterface,
	incidentActionPlanSvc IncidentActionPlanServiceInterface,

) IncidentReportServiceInterface {
	return &IncidentReportService{
		repo:                         repo,
		settigSvc:                    settigSvc,
		incidentRootCauseAnalysisSvc: incidentRootCauseAnalysisSvc,
		incidentActionPlanSvc:        incidentActionPlanSvc,
	}
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
		"details", "containment", "recovery", "post-incident",
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
			IncidentCategoryTitleEn:    inc.IncidentCategoryTitleEn,
			IncidentCategoryTitleAr:    inc.IncidentCategoryTitleAr,
			IncidentSubCategoryTitleEn: inc.IncidentSubCategoryTitleEn,
			IncidentSubCategoryTitleAr: inc.IncidentSubCategoryTitleAr,
			SeverityTitleEn:            inc.SeverityTitleEn,
			SeverityTitleAr:            inc.SeverityTitleAr,
			IncidentAt:                 inc.RealTime,
			ReportedAt:                 inc.DetectedTime,
			IncidentStatusTitleEn:      inc.IncidentStatusTitleEn,
			IncidentStatusTitleAr:      inc.IncidentStatusTitleAr,
			ContainmentStrategy:        inc.ContainmentStrategy,
			NextSteps:                  inc.NextSteps,
			ResolvedDate:               inc.ResolvedDate,
			RecoveryMeasures:           inc.RecoveryMeasures,
			AssignedUserInitial:        inc.AssignedUserInitial,
			AssignedUserFirstNameEn:    inc.AssignedUserFirstNameEn,
			AssignedUserLastNameEn:     inc.AssignedUserLastNameEn,
			AssignedUserFirstNameAr:    inc.AssignedUserFirstNameAr,
			AssignedUserLastNameAr:     inc.AssignedUserLastNameAr,
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

			reportDocs = append(reportDocs, &models.IncidentReportDocument{
				IncidentDocumentID:       document.ID,
				IncidentReportIncidentID: reportIncident.ID,
				URL:                      newURL,
				Name:                     *document.Name,
				Size:                     document.Size,
				Token:                    document.Token,
				OriginalName:             document.OriginalName,
				Type:                     document.Type,
				CreatedByID:              currentUserID,
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
			rcas = append(rcas, &models.IncidentReportRootCauseAnalysis{
				IncidentReportIncidentID: reportIncident.ID,
				IncidentReportID:         reportID,
				Why:                      a.Why,
				Question:                 a.Question,
				Answer:                   a.Answer,
				CreatedByID:              a.CreatedByID,
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
			reportPlans = append(reportPlans, &models.IncidentReportActionPlan{
				IncidentReportIncidentID: reportIncident.ID,
				IncidentReportID:         reportID,

				Title:       p.Title,
				Description: p.Description,
				TargetDate:  p.TargetDate,

				StatusTitle: p.ActionPlanStatus.TitleEn,

				CreatedByID: &p.CreatedBYID,
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

	parsed.StartDate, err = utility.ParseDate(*dto.StartDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
	}

	parsed.EndDate, err = utility.ParseDate(*dto.EndDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end date")
	}

	return parsed, nil
}

func formatTitle(sectionType string) (string, string) {
	switch sectionType {
	case "executive-summary":
		return "Executive Summary", "الملخص التنفيذي"
	case "context":
		return "Risk Governance & Context", "حوكمة المخاطر والسياق"
	case "identification":
		return "Risk Identification", "تحديد المخاطر"
	case "assessment":
		return "Risk Assessment Details", "تفاصيل تقييم المخاطر"
	case "mitigation":
		return "Risk Response & Mitigation", "الاستجابة للمخاطر والتخفيف منها"
	case "kri":
		return "Key Risk Indicators (KRIs) & Performance", "مؤشرات المخاطر الرئيسية والأداء"
	case "communication":
		return "Risk Monitoring, Reporting & Communication", "مراقبة المخاطر والإبلاغ والتواصل"
	case "recommendation":
		return "Recommendations & Action Plan", "التوصيات وخطة العمل"
	case "documents":
		return "Appendices & Supporting Documents", "الملاحق والمستندات الداعمة"
	default:
		return sectionType, sectionType // Default to sectionType if not found
	}
}
