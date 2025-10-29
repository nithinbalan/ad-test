package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	kpiDTOs "jethur/grc/kpi/dtos"

	"jethur/grc/kpi/repositories"
)

type ComplianceDashboardServiceInterface interface {
	GetComplianceOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceOverviewDTO, error)
	GetLeadershipAttentionSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceLeadershipAttentionDTO, error)
	GetComplianceBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceBreakdownDTO, error)
	GetHighPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePerformanceRadarDTO, error)
	GetLeastPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePerformanceRadarDTO, error)
	GetHighPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceClausePerformanceDTO, error)
	GetLeastPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceClausePerformanceDTO, error)
	GetPolicyOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePolicyOverviewDTO, error)
	GetPolicyStatusBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePolicyStatusBreakdownDTO, error)
	GetPolicyByOrganizationUnit(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePolicyByOrganizationUnitDTO, error)
	GetControlSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlSummaryDTO, error)
	GetComplianceEfficiencySummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceEfficiencySummaryDTO, error)
	GetControlOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlOverviewDTO, error)
	GetOrganizationUnitWiseAssessmentCoverage(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceAssessmentControlOrganizationUnit, error)
	GetActionPlanOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceActionPlanOverviewDTO, error)
	GetActionPlanStatusSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceActionPlanStatusSummaryDTO, error)
	GetActionPlansByOrganizationUnit(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceActionPlanByOrganizationUnitDTO, error)

	GetTopCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]complianceDTOs.ComplianceRegisterListDTO, error)
}

type ComplianceDashboardService struct {
	repo repositories.ComplianceDashboardRepositoryInterface
}

func NewComplianceDashboardService(repo repositories.ComplianceDashboardRepositoryInterface) ComplianceDashboardServiceInterface {
	return &ComplianceDashboardService{repo: repo}
}

func (cds *ComplianceDashboardService) GetComplianceOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceOverviewDTO, error) {
	return cds.repo.GetComplianceOverview(params)
}

func (cds *ComplianceDashboardService) GetLeadershipAttentionSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceLeadershipAttentionDTO, error) {
	return cds.repo.GetLeadershipAttentionSummary(params)
}

func (cds *ComplianceDashboardService) GetComplianceBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceBreakdownDTO, error) {
	return cds.repo.GetComplianceBreakdown(params)
}

func (cds *ComplianceDashboardService) GetHighPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePerformanceRadarDTO, error) {
	return cds.repo.GetHighPerformanceCompliances(params)
}

func (cds *ComplianceDashboardService) GetLeastPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePerformanceRadarDTO, error) {
	return cds.repo.GetLeastPerformanceCompliances(params)
}

func (cds *ComplianceDashboardService) GetHighPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceClausePerformanceDTO, error) {
	return cds.repo.GetHighPerformanceComplianceClauses(params)
}

func (cds *ComplianceDashboardService) GetLeastPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceClausePerformanceDTO, error) {
	return cds.repo.GetLeastPerformanceComplianceClauses(params)
}

func (cds *ComplianceDashboardService) GetPolicyOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePolicyOverviewDTO, error) {
	return cds.repo.GetPolicyOverview(params)
}

func (cds *ComplianceDashboardService) GetPolicyStatusBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePolicyStatusBreakdownDTO, error) {
	return cds.repo.GetPolicyStatusBreakdown(params)
}

func (cds *ComplianceDashboardService) GetPolicyByOrganizationUnit(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePolicyByOrganizationUnitDTO, error) {
	return cds.repo.GetPolicyByOrganizationUnit(params)
}

func (cds *ComplianceDashboardService) GetControlSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlSummaryDTO, error) {
	return cds.repo.GetControlSummary(params)
}

func (cds *ComplianceDashboardService) GetComplianceEfficiencySummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceEfficiencySummaryDTO, error) {
	return cds.repo.GetComplianceEfficiencySummary(params)
}

func (cds *ComplianceDashboardService) GetControlOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlOverviewDTO, error) {
	return cds.repo.GetControlOverview(params)
}

func (cds *ComplianceDashboardService) GetOrganizationUnitWiseAssessmentCoverage(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceAssessmentControlOrganizationUnit, error) {
	return cds.repo.GetOrganizationUnitWiseAssessmentCoverage(params)
}

func (cds *ComplianceDashboardService) GetActionPlanOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceActionPlanOverviewDTO, error) {
	return cds.repo.GetActionPlanOverview(params)
}

func (cds *ComplianceDashboardService) GetActionPlanStatusSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceActionPlanStatusSummaryDTO, error) {
	return cds.repo.GetActionPlanStatusSummary(params)
}

func (cds *ComplianceDashboardService) GetActionPlansByOrganizationUnit(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceActionPlanByOrganizationUnitDTO, error) {
	return cds.repo.GetActionPlansByOrganizationUnit(params)
}

func (cds *ComplianceDashboardService) GetTopCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]complianceDTOs.ComplianceRegisterListDTO, error) {
	return cds.repo.GetTopCompliances(params)
}
