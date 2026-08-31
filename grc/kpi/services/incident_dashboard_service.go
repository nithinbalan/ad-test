package services

import (
	incidentDTOs "jethur/grc/kpi/dtos"
	incidentRepositories "jethur/grc/kpi/repositories"
)

type IncidentDashboardServiceInterface interface {
	GetIncidentStatusSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentStatusSummaryDTO, error)
	GetIncidentSeverityBreakdown(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentSeverityBreakdownDTO, error)
	GetIncidentAcknowledgementSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentAcknowledgementSummaryDTO, error)
	GetIncidentSLASummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentSLASummaryDTO, error)
	GetIncidentCategorySummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentCategoryCountDTO, error)
	GetIncidentTypeSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentTypeCountDTO, error)
	GetIncidentTrend(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentTrendDTO, error)
	GetIncidentActionPlanSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentActionPlanSummaryDTO, error)
	GetOverdueActionPlans(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.OverdueActionPlanDTO, error)
}

type IncidentDashboardService struct {
	repo incidentRepositories.IncidentDashboardRepositoryInterface
}

func NewIncidentDashboardService(repo incidentRepositories.IncidentDashboardRepositoryInterface) IncidentDashboardServiceInterface {
	return &IncidentDashboardService{repo: repo}
}

func (ids *IncidentDashboardService) GetIncidentStatusSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentStatusSummaryDTO, error) {
	return ids.repo.GetIncidentStatusSummary(params)
}

func (ids *IncidentDashboardService) GetIncidentSeverityBreakdown(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentSeverityBreakdownDTO, error) {
	return ids.repo.GetIncidentSeverityBreakdown(params)
}

func (ids *IncidentDashboardService) GetIncidentAcknowledgementSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentAcknowledgementSummaryDTO, error) {
	return ids.repo.GetIncidentAcknowledgementSummary(params)
}

func (ids *IncidentDashboardService) GetIncidentSLASummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentSLASummaryDTO, error) {
	return ids.repo.GetIncidentSLASummary(params)
}

func (ids *IncidentDashboardService) GetIncidentCategorySummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentCategoryCountDTO, error) {
	return ids.repo.GetIncidentCategorySummary(params)
}

func (ids *IncidentDashboardService) GetIncidentTypeSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentTypeCountDTO, error) {
	return ids.repo.GetIncidentTypeSummary(params)
}

func (ids *IncidentDashboardService) GetIncidentTrend(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentTrendDTO, error) {
	return ids.repo.GetIncidentTrend(params)
}

func (ids *IncidentDashboardService) GetIncidentActionPlanSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentActionPlanSummaryDTO, error) {
	return ids.repo.GetIncidentActionPlanSummary(params)
}

func (ids *IncidentDashboardService) GetOverdueActionPlans(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.OverdueActionPlanDTO, error) {
	return ids.repo.GetOverdueActionPlans(params)
}
