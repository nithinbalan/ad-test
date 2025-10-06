package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
)

type ComplianceDetailDashboardServiceInterface interface {
	GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardDetailSummaryDTO, error)
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
