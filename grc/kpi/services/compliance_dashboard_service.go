package services

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
)

type ComplianceDashboardServiceInterface interface {
	GetComplianceOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceOverviewDTO, error)
	GetComplianceBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceBreakdownDTO, error)
	GetHighPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePerformanceRadarDTO, error)
	GetLeastPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePerformanceRadarDTO, error)
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

func (cds *ComplianceDashboardService) GetComplianceBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceBreakdownDTO, error) {
	return cds.repo.GetComplianceBreakdown(params)
}

func (cds *ComplianceDashboardService) GetHighPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePerformanceRadarDTO, error) {
	return cds.repo.GetHighPerformanceCompliances(params)
}

func (cds *ComplianceDashboardService) GetLeastPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePerformanceRadarDTO, error) {
	return cds.repo.GetLeastPerformanceCompliances(params)
}
