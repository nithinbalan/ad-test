package services

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type IncidentReportServiceInterface interface {
	GetIncidentByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByCategoryDTO, error)
	GetIncidentByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByStatusDTO, error)
	GetIncidentBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentBySeverityLevelDTO, error)
	GetIncidentByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByResponsibleUnitDTO, error)
	GetIncidentByReportedUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByReportedUserDTO, error)
	GetIncidentByReportedYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByReportedYearDTO, error)

	GetIncidentByAsset(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByAssetDTO, error)
	GetIncidentBySLA(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentBySLADTO, error)

	GetIncidentInvestigationSummary(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentInvestigationSummaryDTO, error)
	GetIncidentRCA(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentRCADTO, error)
	GetIncidentClosure(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentClosureDTO, error)
	GetOpenIncidentsWithRisk(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.OpenIncidentWithRiskDTO, error)

	GetIncidentPendingActionPlan(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentPendingActionPlanDTO, error)
	GetIncidentActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByStatusDTO, error)
	GetIncidentActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByResponsibleUserDTO, error)
	GetIncidentActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByResponsibleUnitDTO, error)
}

type IncidentReportService struct {
	repo repositories.IncidentReportRepositoryInterface
}

func NewIncidentReportService(repo repositories.IncidentReportRepositoryInterface) IncidentReportServiceInterface {
	return &IncidentReportService{repo: repo}
}

func (irs *IncidentReportService) GetIncidentByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByCategoryDTO, error) {
	return irs.repo.GetIncidentByCategory(params)
}

func (irs *IncidentReportService) GetIncidentByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByStatusDTO, error) {
	return irs.repo.GetIncidentByStatus(params)
}

func (irs *IncidentReportService) GetIncidentBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentBySeverityLevelDTO, error) {
	return irs.repo.GetIncidentBySeverityLevel(params)
}

func (irs *IncidentReportService) GetIncidentByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByResponsibleUnitDTO, error) {
	return irs.repo.GetIncidentByResponsibleUnit(params)
}

func (irs *IncidentReportService) GetIncidentByReportedUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByReportedUserDTO, error) {
	return irs.repo.GetIncidentByReportedUser(params)
}

func (irs *IncidentReportService) GetIncidentByReportedYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByReportedYearDTO, error) {
	return irs.repo.GetIncidentByReportedYear(params)
}

func (irs *IncidentReportService) GetIncidentByAsset(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByAssetDTO, error) {
	return irs.repo.GetIncidentByAsset(params)
}

func (irs *IncidentReportService) GetIncidentBySLA(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentBySLADTO, error) {
	return irs.repo.GetIncidentBySLA(params)
}

func (irs *IncidentReportService) GetIncidentInvestigationSummary(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentInvestigationSummaryDTO, error) {
	return irs.repo.GetIncidentInvestigationSummary(params)
}

func (irs *IncidentReportService) GetIncidentRCA(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentRCADTO, error) {
	return irs.repo.GetIncidentRCA(params)
}

func (irs *IncidentReportService) GetIncidentClosure(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentClosureDTO, error) {
	return irs.repo.GetIncidentClosure(params)
}

func (irs *IncidentReportService) GetOpenIncidentsWithRisk(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.OpenIncidentWithRiskDTO, error) {
	return irs.repo.GetOpenIncidentsWithRisk(params)
}

func (irs *IncidentReportService) GetIncidentPendingActionPlan(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentPendingActionPlanDTO, error) {
	return irs.repo.GetIncidentPendingActionPlan(params)
}

func (irs *IncidentReportService) GetIncidentActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByStatusDTO, error) {
	return irs.repo.GetIncidentActionPlanByStatus(params)
}

func (irs *IncidentReportService) GetIncidentActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByResponsibleUserDTO, error) {
	return irs.repo.GetIncidentActionPlanByResponsibleUser(params)
}

func (irs *IncidentReportService) GetIncidentActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByResponsibleUnitDTO, error) {
	return irs.repo.GetIncidentActionPlanByResponsibleUnit(params)
}
