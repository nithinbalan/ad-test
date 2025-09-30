package services

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
	riskDTOs "jethur/grc/risk/dtos"
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
	GetERMRiskAgingStatus(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAgingStatusDTO, error)
	GetERMRiskTreatmentPlanProgress(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskTreatmentPlanProgressDTO, error)
	GetERMStrategicObjectivesRiskOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMStrategicObjectiveRiskOverviewDTO, error)
	GetERMStrategicObjectivesRiskSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMStrategicObjectivesRiskSummaryDTO, error)
	GetERMRiskSummaryByDepartment(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskDepartmentSummaryDTO, error)

	// ERM Risk Details
	GetERMRiskCategorySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskCategorySummaryDTO, error)
	GetERMRiskSourceSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskSourceSummaryDTO, error)
	GetERMRiskStatusSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskStatusSummaryDTO, error)
	GetERMRiskDepartmentSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskSummaryByDepartmentDTO, error)
	GetERMRisksIdentifiedSummary() (*kpiDTOs.ERMRisksIdentifiedSummaryDTO, error)
	GetERMRiskResponseStrategySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskResponseStrategySummaryDTO, error)
	GetERMRiskControlEfficiencySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskControlEfficiencySummaryDTO, error)
	GetERMRiskMitigationPlanStatusSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskMitigationPlanStatusSummaryDTO, error)
	GetERMOpenRisksSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMOpenRisksSummaryDTO, error)

	GetTopRiskRegisters(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]*riskDTOs.RiskRegisterListDTO, error)
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

func (eds *ERMDashboardService) GetERMRiskAgingStatus(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAgingStatusDTO, error) {
	return eds.repo.GetERMRiskAgingStatus(params)
}

func (eds *ERMDashboardService) GetERMRiskTreatmentPlanProgress(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskTreatmentPlanProgressDTO, error) {
	return eds.repo.GetERMRiskTreatmentPlanProgress(params)
}

func (eds *ERMDashboardService) GetERMStrategicObjectivesRiskOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMStrategicObjectiveRiskOverviewDTO, error) {
	return eds.repo.GetERMStrategicObjectivesRiskOverview(params)
}

func (eds *ERMDashboardService) GetERMStrategicObjectivesRiskSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMStrategicObjectivesRiskSummaryDTO, error) {
	return eds.repo.GetERMStrategicObjectivesRiskSummary(params)
}

func (eds *ERMDashboardService) GetERMRiskSummaryByDepartment(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskDepartmentSummaryDTO, error) {
	return eds.repo.GetERMRiskSummaryByDepartment(params)
}

func (eds *ERMDashboardService) GetERMRiskCategorySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskCategorySummaryDTO, error) {
	return eds.repo.GetERMRiskCategorySummary(params)
}

func (eds *ERMDashboardService) GetERMRiskSourceSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskSourceSummaryDTO, error) {
	return eds.repo.GetERMRiskSourceSummary(params)
}

func (eds *ERMDashboardService) GetERMRiskStatusSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskStatusSummaryDTO, error) {
	return eds.repo.GetERMRiskStatusSummary(params)
}

func (eds *ERMDashboardService) GetERMRiskDepartmentSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskSummaryByDepartmentDTO, error) {
	return eds.repo.GetERMRiskDepartmentSummary(params)
}

func (eds *ERMDashboardService) GetERMRisksIdentifiedSummary() (*kpiDTOs.ERMRisksIdentifiedSummaryDTO, error) {
	return eds.repo.GetERMRisksIdentifiedSummary()
}

func (eds *ERMDashboardService) GetERMRiskResponseStrategySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskResponseStrategySummaryDTO, error) {
	return eds.repo.GetERMRiskResponseStrategySummary(params)
}

func (eds *ERMDashboardService) GetERMRiskControlEfficiencySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskControlEfficiencySummaryDTO, error) {
	return eds.repo.GetERMRiskControlEfficiencySummary(params)
}

func (eds *ERMDashboardService) GetERMRiskMitigationPlanStatusSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskMitigationPlanStatusSummaryDTO, error) {
	return eds.repo.GetERMRiskMitigationPlanStatusSummary(params)
}

func (eds *ERMDashboardService) GetERMOpenRisksSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMOpenRisksSummaryDTO, error) {
	return eds.repo.GetERMOpenRisksSummary(params)
}

func (eds *ERMDashboardService) GetTopRiskRegisters(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]*riskDTOs.RiskRegisterListDTO, error) {
	return eds.repo.GetTopRiskRegisters(params)
}
