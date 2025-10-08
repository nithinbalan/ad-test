package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
)

type ComplianceDetailDashboardServiceInterface interface {
	GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDetailSummaryDTO, error)
	GetComplianceClauseSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceClauseSummaryDTO, error)
	GetTop5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error)
	GetBottom5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error)
	GetComplianceDomainWiseOverview(params *complianceDTOs.ComplianceDashboardQueryParamsDTO, complianceRegisterID uint) (*complianceDTOs.ComplianceDomainWiseOverviewDTO, error)
	GetPolicyByOrganizationUnit(complianceRegisterID uint) ([]complianceDTOs.CompliancePolicyByOrganizationUnitDTO, error)
	GetControlEffectivenessSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlEffectivenessDTO, error)
	GetControlImplementationSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlImplementationSummaryDTO, error)
	GetControlSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlSummaryDTO, error)
	GetControlEfficiencySummary(complianceRegisterID uint) ([]complianceDTOs.ComplianceSummaryItemDTO, error)
	GetControlCoverageByDomain(complianceRegisterID uint) ([]complianceDTOs.ComplianceControlCoverageByDomainDTO, error)
	GetCompliancePolicySummary(complianceRegisterID uint) (*complianceDTOs.CompliancePolicySummaryDTO, error)
}

type ComplianceDetailDashboardService struct {
	repo repositories.ComplianceDetailDashboardRepositoryInterface
}

func NewComplianceDetailDashboardService(repo repositories.ComplianceDetailDashboardRepositoryInterface) *ComplianceDetailDashboardService {
	return &ComplianceDetailDashboardService{
		repo: repo,
	}
}

func (cdds *ComplianceDetailDashboardService) GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDetailSummaryDTO, error) {
	return cdds.repo.GetComplianceDetailSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetComplianceClauseSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceClauseSummaryDTO, error) {
	return cdds.repo.GetComplianceClauseSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetTop5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error) {
	return cdds.repo.GetTop5ClausePerformance(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetBottom5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error) {
	return cdds.repo.GetBottom5ClausePerformance(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetComplianceDomainWiseOverview(params *complianceDTOs.ComplianceDashboardQueryParamsDTO, complianceRegisterID uint) (*complianceDTOs.ComplianceDomainWiseOverviewDTO, error) {
	return cdds.repo.GetComplianceDomainWiseOverview(params, complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetPolicyByOrganizationUnit(complianceRegisterID uint) ([]complianceDTOs.CompliancePolicyByOrganizationUnitDTO, error) {
	return cdds.repo.GetPolicyByOrganizationUnit(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetControlEffectivenessSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlEffectivenessDTO, error) {
	return cdds.repo.GetControlEffectivenessSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetControlImplementationSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlImplementationSummaryDTO, error) {
	return cdds.repo.GetControlImplementationSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetControlSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlSummaryDTO, error) {
	return cdds.repo.GetControlSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetControlEfficiencySummary(complianceRegisterID uint) ([]complianceDTOs.ComplianceSummaryItemDTO, error) {
	return cdds.repo.GetControlEfficiencySummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetControlCoverageByDomain(complianceRegisterID uint) ([]complianceDTOs.ComplianceControlCoverageByDomainDTO, error) {
	return cdds.repo.GetControlCoverageByDomain(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetCompliancePolicySummary(complianceRegisterID uint) (*complianceDTOs.CompliancePolicySummaryDTO, error) {
	return cdds.repo.GetCompliancePolicySummary(complianceRegisterID)
}
