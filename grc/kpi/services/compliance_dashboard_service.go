package services

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
)

type ComplianceDashboardServiceInterface interface {
	GetComplianceOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceOverviewDTO, error)
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
