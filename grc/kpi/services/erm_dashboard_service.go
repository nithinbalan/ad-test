package services

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
)

type ERMDashboardServiceInterface interface {
	GetERMRiskOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskOverviewDTO, error)
	GetERMRiskDistribution(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskDistributionDTO, error)
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
