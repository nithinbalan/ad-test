package services

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepo "jethur/grc/incident/repositories"
	"jethur/utility"
	"time"
)

type IncidentReportServiceInterface interface {
	ExportIncidentReport(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (string, error)
	GetFullReportDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentReportDetailDTO, error)
}

// Service and implementations
type IncidentReportService struct {
	incidentRegisterSvc   IncidentRegisterServiceInterface
	repo                  incidentRepo.IncidentReportRepositoryInterface
	investigationSvc      IncidentInvestigationServiceInterface
	incidentClosureSvc    IncidentClosureServiceInterface
	incidentPostReviewSvc IncidentPostReviewServiceInterface
}

func NewIncidentReportService(
	repo incidentRepo.IncidentReportRepositoryInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
	investigationSvc IncidentInvestigationServiceInterface,
	incidentClosureSvc IncidentClosureServiceInterface,
	incidentPostReviewSvc IncidentPostReviewServiceInterface) IncidentReportServiceInterface {
	return &IncidentReportService{
		repo:                  repo,
		incidentRegisterSvc:   incidentRegisterSvc,
		investigationSvc:      investigationSvc,
		incidentClosureSvc:    incidentClosureSvc,
		incidentPostReviewSvc: incidentPostReviewSvc,
	}
}

func (irs IncidentReportService) GetFullReportDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentReportDetailDTO, error) {
	// Initialize your DTO
	var incidentDTO incidentDTOs.IncidentReportDetailDTO
	incident, err := irs.repo.GetFullReportDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
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

	incidentDTO.IncidentAtFormatted = incident.IncidentAt.Format("Jan 02, 2006, 03:04 PM")
	incidentDTO.ReportedAtFormatted = incident.ReportedAt.Format("Jan 02, 2006, 03:04 PM")

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

// func (irs IncidentReportService) ExportIncidentReport(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (string, error) {

// 	incidentDTO, err := irs.GetFullReportDetailsByID(id, tx...)
// 	if err != nil {
// 		return "", err
// 	}

// 	incidentDTO.GeneratedAt = time.Now().Format("2006-01-02 03:04 PM")

// 	templatePaths := []string{"incident/templates/pdf/incident_report.html"}
// 	htmlContent, err := utility.GeneratePDFContent(incidentDTO, templatePaths, "", "")
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate HTML content: %w", err)
// 	}

// 	// 4. Convert HTML → PDF
// 	pdfFileName := fmt.Sprintf("IncidentReport_%d.pdf", id)
// 	if err := utility.ConvertHTMLContentToPDF(htmlContent, pdfFileName); err != nil {
// 		return "", fmt.Errorf("failed to generate PDF: %w", err)
// 	}

// 	return pdfFileName, nil
// }

func (irs IncidentReportService) ExportIncidentReport(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (string, error) {
	incidentDTO, err := irs.GetFullReportDetailsByID(id, tx...)
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

	// Convert HTML → PDF with header and footer
	pdfFileName := fmt.Sprintf("IncidentReport_%d.pdf", id)
	if err := utility.ConvertHTMLContentToPDFWithHeaderAndFooter(htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}