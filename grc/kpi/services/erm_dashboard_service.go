package services

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
)

type ERMDashboardServiceInterface interface {
	GetERMRiskOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskOverviewDTO, error)
	GetERMRiskDistribution(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskDistributionDTO, error)
	GetERMOverdueRisks(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMOverdueRisksDTO, error)
	GetERMRisksAboveAppetite(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRisksAboveAppetiteDTO, error)
	GetERMLeadershipAttentionMetrics(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMLeadershipAttentionDTO, error)
	GetERMControlEfficiencySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlEfficiencySummaryDTO, error)
	GetERMControlEfficiencyBreakdown(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlEfficiencyBreakdownDTO, error)
	GetERMControlsOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlsOverviewDTO, error)
	GetERMRiskAppetiteOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAppetiteOverviewDTO, error)
	GetERMRiskAppetiteSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAppetiteSummarytDTO, error)
}

type ERMDashboardService struct {
	repo repositories.ERMDashboardRepositoryInterface
}

func NewERMDashboardService(repo repositories.ERMDashboardRepositoryInterface) ERMDashboardServiceInterface {
	return &ERMDashboardService{repo: repo}
}

func (eds *ERMDashboardService) GetERMRiskOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskOverviewDTO, error) {
	return eds.repo.GetERMRiskOverview(params)
}

func (eds *ERMDashboardService) GetERMRiskDistribution(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskDistributionDTO, error) {
	return eds.repo.GetERMRiskDistribution(params)
}

func (eds *ERMDashboardService) GetERMOverdueRisks(parms *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMOverdueRisksDTO, error) {
	return eds.repo.GetERMOverdueRisks(parms)
}

func (eds *ERMDashboardService) GetERMRisksAboveAppetite(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRisksAboveAppetiteDTO, error) {
	return eds.repo.GetERMRisksAboveAppetite(params)
}

func (eds *ERMDashboardService) GetERMLeadershipAttentionMetrics(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMLeadershipAttentionDTO, error) {
	return eds.repo.GetERMLeadershipAttentionMetrics(params)
}

func (eds *ERMDashboardService) GetERMControlEfficiencySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlEfficiencySummaryDTO, error) {
	return eds.repo.GetERMControlEfficiencySummary(params)
}

func (eds *ERMDashboardService) GetERMControlEfficiencyBreakdown(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlEfficiencyBreakdownDTO, error) {
	return eds.repo.GetERMControlEfficiencyBreakdown(params)
}

func (eds *ERMDashboardService) GetERMControlsOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlsOverviewDTO, error) {
	return eds.repo.GetERMControlsOverview(params)
}

func (eds *ERMDashboardService) GetERMRiskAppetiteOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAppetiteOverviewDTO, error) {
	return eds.repo.GetERMRiskAppetiteOverview(params)
}

func (eds *ERMDashboardService) GetERMRiskAppetiteSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAppetiteSummarytDTO, error) {
	return eds.repo.GetERMRiskAppetiteSummary(params)
}
