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
	GetHighPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceClausePerformanceDTO, error)
	GetLowPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceClausePerformanceDTO, error)
	GetComplianceControlSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlSummaryDTO, error)
	GetComplianceEfficiencySummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceEfficiencySummaryDTO, error)
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

func (cds *ComplianceDashboardService) GetHighPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceClausePerformanceDTO, error) {
	return cds.repo.GetHighPerformanceComplianceClauses(params)
}

func (cds *ComplianceDashboardService) GetLowPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceClausePerformanceDTO, error) {
	return cds.repo.GetLowPerformanceComplianceClauses(params)
}

func (cds *ComplianceDashboardService) GetComplianceControlSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlSummaryDTO, error) {
	return cds.repo.GetComplianceControlSummary(params)
}

func (cds *ComplianceDashboardService) GetComplianceEfficiencySummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceEfficiencySummaryDTO, error) {
	return cds.repo.GetComplianceEfficiencySummary(params)
}
