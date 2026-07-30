package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	kpiDTOs "jethur/grc/kpi/dtos"

	"jethur/grc/kpi/repositories"
	organizationServices "jethur/grc/organization/services"
)

type ComplianceDashboardServiceInterface interface {
	GetComplianceOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceOverviewDTO, error)
	GetLeadershipAttentionSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceLeadershipAttentionDTO, error)
	GetComplianceBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceBreakdownDTO, error)
	GetHighPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.CompliancePerformanceRadarDTO, error)
	GetLeastPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.CompliancePerformanceRadarDTO, error)
	GetHighPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.ComplianceClausePerformanceDTO, error)
	GetLeastPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.ComplianceClausePerformanceDTO, error)
	GetPolicyOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePolicyOverviewDTO, error)
	GetPolicyStatusBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePolicyStatusBreakdownDTO, error)
	GetPolicyByOrganizationUnit(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePolicyByOrganizationUnitDTO, error)
	GetControlSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlSummaryDTO, error)
	GetComplianceEfficiencySummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.ComplianceEfficiencySummaryDTO, error)
	GetControlOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlOverviewDTO, error)
	GetOrganizationUnitWiseAssessmentCoverage(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.ComplianceAssessmentControlOrganizationUnit, error)
	GetActionPlanOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceActionPlanOverviewDTO, error)
	GetActionPlanStatusSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceActionPlanStatusSummaryDTO, error)
	GetActionPlansByOrganizationUnit(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.ComplianceActionPlanByOrganizationUnitDTO, error)

	GetTopCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]complianceDTOs.ComplianceRegisterListDTO, error)
	GetComplianceRatingDistribution(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceRatingDistributionDTO, error)
	GetComplianceRegisterClassificationSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceRegisterClassificationSummaryDTO, error)
}

type ComplianceDashboardService struct {
	repo    repositories.ComplianceDashboardRepositoryInterface
	userSvc organizationServices.UserServiceInterface
}

func NewComplianceDashboardService(repo repositories.ComplianceDashboardRepositoryInterface, userSvc organizationServices.UserServiceInterface) ComplianceDashboardServiceInterface {
	return &ComplianceDashboardService{repo: repo, userSvc: userSvc}
}

func (cds *ComplianceDashboardService) GetComplianceOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceOverviewDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetComplianceOverview(params)
}

func (cds *ComplianceDashboardService) GetLeadershipAttentionSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceLeadershipAttentionDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetLeadershipAttentionSummary(params)
}

func (cds *ComplianceDashboardService) GetComplianceBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceBreakdownDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetComplianceBreakdown(params)
}

func (cds *ComplianceDashboardService) GetHighPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.CompliancePerformanceRadarDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetHighPerformanceCompliances(params)
}

func (cds *ComplianceDashboardService) GetLeastPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.CompliancePerformanceRadarDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetLeastPerformanceCompliances(params)
}

func (cds *ComplianceDashboardService) GetHighPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.ComplianceClausePerformanceDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetHighPerformanceComplianceClauses(params)
}

func (cds *ComplianceDashboardService) GetLeastPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.ComplianceClausePerformanceDTO, error) {
	cds.setCyberFilter(params, currentUserID)
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

func (cds *ComplianceDashboardService) GetComplianceEfficiencySummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.ComplianceEfficiencySummaryDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetComplianceEfficiencySummary(params)
}

func (cds *ComplianceDashboardService) GetControlOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlOverviewDTO, error) {
	return cds.repo.GetControlOverview(params)
}

func (cds *ComplianceDashboardService) GetOrganizationUnitWiseAssessmentCoverage(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.ComplianceAssessmentControlOrganizationUnit, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetOrganizationUnitWiseAssessmentCoverage(params)
}

func (cds *ComplianceDashboardService) GetActionPlanOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceActionPlanOverviewDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetActionPlanOverview(params)
}

func (cds *ComplianceDashboardService) GetActionPlanStatusSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceActionPlanStatusSummaryDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetActionPlanStatusSummary(params)
}

func (cds *ComplianceDashboardService) GetActionPlansByOrganizationUnit(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]kpiDTOs.ComplianceActionPlanByOrganizationUnitDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetActionPlansByOrganizationUnit(params)
}

func (cds *ComplianceDashboardService) GetTopCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) ([]complianceDTOs.ComplianceRegisterListDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetTopCompliances(params)
}

func (cds *ComplianceDashboardService) GetComplianceRatingDistribution(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceRatingDistributionDTO, error) {
	cds.setCyberFilter(params, currentUserID)

	registers, err := cds.repo.GetComplianceRegisterRatingDistribution(params)
	if err != nil {
		return nil, err
	}

	assessments, err := cds.repo.GetComplianceAssessmentRatingDistribution(params)
	if err != nil {
		return nil, err
	}

	return &kpiDTOs.ComplianceRatingDistributionDTO{
		Registers:   registers,
		Assessments: assessments,
	}, nil
}

func (cds *ComplianceDashboardService) GetComplianceRegisterClassificationSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.ComplianceRegisterClassificationSummaryDTO, error) {
	cds.setCyberFilter(params, currentUserID)
	return cds.repo.GetComplianceRegisterClassificationSummary(params)
}

func (cds *ComplianceDashboardService) setCyberFilter(params *kpiDTOs.ComplianceDashboardQueryParamsDTO, currentUserID uint) {
	if currentUserID == 0 {
		return
	}
	isCyberUser, err := cds.userSvc.IsCyberUser(currentUserID)
	if err != nil {
		return
	}
	params.IsCyberUser = isCyberUser
}
