package services

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
)

type GRCDashboardServiceInterface interface {
	GetGRCComplianceSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceSummaryDTO, error)
	GetGRCRiskSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskSummaryDTO, error)
}

type GRCDashboardService struct {
	repo repositories.GRCDashboardRepositoryInteface
}

func NewGRCDashboardService(repo repositories.GRCDashboardRepositoryInteface) GRCDashboardServiceInterface {
	return &GRCDashboardService{repo: repo}
}

// Compliance Summary
func (gds *GRCDashboardService) GetGRCComplianceSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceSummaryDTO, error) {
	return gds.repo.GetGRCComplianceSummary(params)
}

// Risk Summary
func (gds *GRCDashboardService) GetGRCRiskSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskSummaryDTO, error) {
	return gds.repo.GetGRCRiskSummary(params)
}
