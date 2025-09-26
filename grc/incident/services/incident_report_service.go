package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepo "jethur/grc/incident/repositories"
	"jethur/utility"
	"time"
)

type IncidentReportServiceInterface interface {
	ExportIncidentReportReport(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (string, error)
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

func (imr IncidentReportService) ExportIncidentReportReport(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (string, error) {

	incident, err := imr.incidentRegisterSvc.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", utility.NewResourceNotFoundError(err, id)
		}
		return "", err
	}

	timelines, err := imr.repo.GetAllIncidentTimeline(id)
	if err != nil {
		return "", err
	}

	strategicObjectives, err := imr.repo.GetAllIncidentStrategicObjectives(id)
	if err != nil {
		return "", err
	}

	kpis, err := imr.repo.GetAllIncidentKPI(id)
	if err != nil {
		return "", err
	}

	controls, err := imr.repo.GetAllControls(id)
	if err != nil {
		return "", err
	}

	assets, err := imr.repo.GetAllIncidentAssets(id)
	if err != nil {
		return "", err
	}

	investigation, err := imr.investigationSvc.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	rca, err := imr.repo.GetAllIncidentRca(id)
	if err != nil {
		return "", err
	}

	communicationLog, err := imr.repo.GetAllCommunicationLog(id)
	if err != nil {
		return "", err
	}

	Documents, err := imr.repo.GetAllIncidentDocuments(id)
	if err != nil {
		return "", err
	}

	// Risks, err := imr.repo.GetAllIncidentRisks(id)
	// if err != nil {
	// 	return "", err
	// }
	Processes, err := imr.repo.GetAllIncidentProcesses(id)
	if err != nil {
		return "", err
	}
	
	Compliance, err := imr.repo.GetAllIncidentComplianceRequirements(id)
	if err != nil {
		return "", err
	}


	actionPlan, err := imr.repo.GetAllIncidentActionPlan(id)
	if err != nil {
		return "", err
	}

	closure, err := imr.incidentClosureSvc.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	postReview, err := imr.incidentPostReviewSvc.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}
	pdfDTO := incidentDTOs.IncidentReportPDFDetailDTO{
		Incident:            *incident,
		Timelines:           timelines.Items,
		StrategicObjectives: strategicObjectives,
		KPIs:                kpis,
		Controls:            controls,
		Compliance:          Compliance,
		Processes:           Processes,
		Risks:               incident.Risks,
		Documents:           Documents,
		Assets:              assets,
		Investigation:       investigation,
		RootCauses:          rca,
		CommunicationLogs:   communicationLog,
		ActionPlans:         actionPlan,
		Closure:             closure,
		PostReview:          postReview,
		GeneratedAt:         time.Now().Format("2006-01-02  03:04: PM"),
		IncidentAt: incident.IncidentAt.Format("2006-01-02  03:04: PM"),
		ReportedAt: incident.ReportedAt.Format("2006-01-02  03:04: PM"),
	}


	templatePaths := []string{"incident/templates/pdf/incident_report.html"}
	htmlContent, err := utility.GeneratePDFContent(pdfDTO, templatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	// 4. Convert HTML → PDF
	pdfFileName := fmt.Sprintf("IncidentReport_%d.pdf", id)
	if err := utility.ConvertHTMLContentToPDF(htmlContent, pdfFileName); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil

}
