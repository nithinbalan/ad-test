package services

import (
	"errors"
	"fmt"
	"html/template"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepo "jethur/grc/incident/repositories"
	"jethur/grc/models"
	settingSvc "jethur/grc/settings/services"
	"jethur/utility"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentReportServiceInterface interface {
	ExportIncidentReport(id uint, currentUserID uint) (string, error)
	GetFullReportDetailsByID(id uint) (*incidentDTOs.IncidentReportDetailDTO, error)
	CreateIncidentReport(request incidentDTOs.IncidentReportRequestDTO, currentUserID uint) (string, error)
}

// Service and implementations
type IncidentReportService struct {
	repo      incidentRepo.IncidentReportRepositoryInterface
	settigSvc settingSvc.ReportTemplateSettingServiceInterface
}

func NewIncidentReportService(
	repo incidentRepo.IncidentReportRepositoryInterface,
	settigSvc settingSvc.ReportTemplateSettingServiceInterface,

) IncidentReportServiceInterface {
	return &IncidentReportService{
		repo:      repo,
		settigSvc: settigSvc,
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

func (irs IncidentReportService) CreateIncidentReport(request incidentDTOs.IncidentReportRequestDTO, currentUserID uint) (string, error) {
	tx := irs.repo.BeginTransaction()

	template, err := irs.settigSvc.GetFirst()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if template.CoverPageBGImagePath == nil || template.LogoImagePath == nil {
		tx.Rollback()
		return "", utility.NewForbiddenError("Logo and Background Image not added.")
	}

	var startDate, endDate *time.Time

	if request.Type == "date-range" {
		dates, err := irs.parseIncidentReportDates(request)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		startDate = dates.StartDate
		endDate = dates.EndDate
	}

	incidentReport := &models.IncidentReport{
		Title:       request.Title,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedBYID: currentUserID,
	}

	reportID, err := irs.repo.Create(incidentReport, tx)

	if err != nil {
		tx.Rollback()
		return "", err
	}

	var incidents []*models.IncidentRegister

	switch request.Type {

	case "individual":
		if request.IncidentIDs == nil || len(*request.IncidentIDs) == 0 {
			tx.Rollback()
			return "", utility.NewValidationError(
				fmt.Errorf("incident_ids cannot be empty for individual report"),
				"incident_ids_empty",
				"incident_ids",
			)
		}

		incidents, err = irs.repo.GetIncidentsByIDs(*request.IncidentIDs, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}

	case "date-range":
		incidents, err = irs.repo.GetIncidentsByDateRange(*startDate, *endDate, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}

	default:
		tx.Rollback()
		return "", utility.NewValidationError(
			fmt.Errorf("invalid type"),
			"invalid_type",
			"type",
		)
	}

	for _, incident := range incidents {
	incidentReportIncident := &models.IncidentReportIncident{
		
	}
	}


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
