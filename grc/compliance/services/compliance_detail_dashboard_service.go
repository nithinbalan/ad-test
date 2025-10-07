package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
)

type ComplianceDetailDashboardServiceInterface interface {
	GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDetailSummaryDTO, error)
	GetComplianceClauseSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceClauseSummaryDTO, error)
	GetTop5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceClausePerformanceDTO, error)
	GetBottom5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceClausePerformanceDTO, error)
	GetComplianceDomainWiseOverview(params *complianceDTOs.ComplianceDashboardQueryParamsDTO, complianceRegisterID uint) (*complianceDTOs.ComplianceDomainWiseOverviewDTO, error)
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

func (cdds *ComplianceDetailDashboardService) GetTop5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceClausePerformanceDTO, error) {
	return cdds.repo.GetTop5ClausePerformance(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetBottom5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceClausePerformanceDTO, error) {
	return cdds.repo.GetBottom5ClausePerformance(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetComplianceDomainWiseOverview(params *complianceDTOs.ComplianceDashboardQueryParamsDTO, complianceRegisterID uint) (*complianceDTOs.ComplianceDomainWiseOverviewDTO, error) {
	return cdds.repo.GetComplianceDomainWiseOverview(params, complianceRegisterID)
}
