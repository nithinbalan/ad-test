package services

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepo "jethur/grc/incident/repositories"
	organizationDTOs "jethur/grc/organization/dtos"
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

	for _, so := range incident.StrategicObjectives {
		if so.StrategicObjective != nil {
			incidentDTO.StrategicObjectives = append(incidentDTO.StrategicObjectives, &organizationDTOs.StrategicObjectiveDTO{
				ID:      so.StrategicObjective.ID,
				TitleEn: so.StrategicObjective.TitleEn,
				TitleAr: so.StrategicObjective.TitleAr,
				Status:  so.StrategicObjective.Status,
			})
		}
	}

	for _, so := range incident.StrategicObjectives {
		if so.StrategicObjective != nil {
			incidentDTO.StrategicObjectives = append(incidentDTO.StrategicObjectives, &organizationDTOs.StrategicObjectiveDTO{
				ID:      so.StrategicObjective.ID,
				TitleEn: so.StrategicObjective.TitleEn,
				TitleAr: so.StrategicObjective.TitleAr,
				Status:  so.StrategicObjective.Status,
			})
		}
	}
	incidentDTO.KPIs = []*commonDTOs.AssociatedKPIRegisterListDTO{}
	for _, kpi := range incident.KPIs {
		if kpi.KPIRegister != nil {
			incidentDTO.KPIs = append(incidentDTO.KPIs, &commonDTOs.AssociatedKPIRegisterListDTO{
				ID:                 int(kpi.ID),
				KPIRegisterID:      kpi.KPIRegister.ID,
				Title:              kpi.KPIRegister.Title,
				KPICategoryID:      kpi.KPIRegister.KPICategoryID,
				KPICategoryTitleEn: kpi.KPIRegister.KPICategory.TitleEn,
				KPICategoryTitleAr: kpi.KPIRegister.KPICategory.TitleAr,
				KPIStatusID:        kpi.KPIRegister.KPIStatusID,
				KPIStatusTitleEn:   kpi.KPIRegister.KPIStatus.TitleEn,
				KPIStatusTitleAr:   kpi.KPIRegister.KPIStatus.TitleAr,
				KPIStatusColorCode: kpi.KPIRegister.KPIStatus.ColorCode,
				KPIStatusType:      kpi.KPIRegister.KPIStatus.Type,
			})
		}
	}

	incidentDTO.CreatedAt = incident.CreatedAt
	incidentDTO.UpdatedAt = incident.UpdatedAt

	return &incidentDTO, nil
}

func (irs IncidentReportService) ExportIncidentReport(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (string, error) {

	incidentDTO, err := irs.GetFullReportDetailsByID(id, tx...)
	if err != nil {
		return "", err
	}
	incidentDTO.GeneratedAt = time.Now().Format("2006-01-02 03:04 PM")

	templatePaths := []string{"incident/templates/pdf/incident_report.html"}
	htmlContent, err := utility.GeneratePDFContent(incidentDTO, templatePaths, "", "")
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
