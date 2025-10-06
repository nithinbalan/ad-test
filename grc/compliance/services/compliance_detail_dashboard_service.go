package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
)

type ComplianceDetailDashboardServiceInterface interface {
	GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardDetailSummaryDTO, error)
	GetComplianceClauseSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClauseSummaryDTO, error)
	GetTop5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClausePerformanceDTO, error)
	GetBottom5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClausePerformanceDTO, error)
}

type ComplianceDetailDashboardService struct {
	repo repositories.ComplianceDetailDashboardRepositoryInterface
}

func NewComplianceDetailDashboardService(repo repositories.ComplianceDetailDashboardRepositoryInterface) *ComplianceDetailDashboardService {
	return &ComplianceDetailDashboardService{
		repo: repo,
	}
}

func (cdds *ComplianceDetailDashboardService) GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardDetailSummaryDTO, error) {
	return cdds.repo.GetComplianceDetailSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetComplianceClauseSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClauseSummaryDTO, error) {
	return cdds.repo.GetComplianceClauseSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetTop5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClausePerformanceDTO, error) {
	return cdds.repo.GetTop5ClausePerformance(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetBottom5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClausePerformanceDTO, error) {
	return cdds.repo.GetBottom5ClausePerformance(complianceRegisterID)
}
