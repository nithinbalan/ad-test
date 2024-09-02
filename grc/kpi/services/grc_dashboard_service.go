package services

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
)

type GRCDashboardServiceInterface interface {
	GetGRCComplianceSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceSummaryDTO, error)
	GetGRCRiskSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskSummaryDTO, error)
	GetGRCAuditSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCAuditSummaryDTO, error)
	GetGRCDocumentSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCDocumentSummaryDTO, error)
	GetGRCOrganizationSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCOrganizationSummaryDTO, error)

	GetGRCComplianceOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceOverviewDTO, error)
	GetGRCRiskOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskOverviewDTO, error)

	GetGRCTopCompliances(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopComplianceRegistersDTO, error)
	GetGRCTopRisks(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopRiskRegistersDTO, error)
	GetGRCTopPublishedDocuments(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopDocumentsDTO, error)
}

type GRCDashboardService struct {
	repo repositories.GRCDashboardRepositoryInteface
}

func NewGRCDashboardService(repo repositories.GRCDashboardRepositoryInteface) GRCDashboardServiceInterface {
	return &GRCDashboardService{repo: repo}
}

func (gds *GRCDashboardService) GetGRCComplianceSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceSummaryDTO, error) {
	return gds.repo.GetGRCComplianceSummary(params)
}

func (gds *GRCDashboardService) GetGRCRiskSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskSummaryDTO, error) {
	return gds.repo.GetGRCRiskSummary(params)
}

func (gds *GRCDashboardService) GetGRCAuditSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCAuditSummaryDTO, error) {
	return gds.repo.GetGRCAuditSummary(params)
}

func (gds *GRCDashboardService) GetGRCDocumentSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCDocumentSummaryDTO, error) {
	return gds.repo.GetGRCDocumentSummary(params)
}

func (gds *GRCDashboardService) GetGRCOrganizationSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCOrganizationSummaryDTO, error) {
	return gds.repo.GetGRCOrganizationSummary(params)
}

func (gds *GRCDashboardService) GetGRCComplianceOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceOverviewDTO, error) {
	return gds.repo.GetGRCComplianceOverview(params)
}

func (gds *GRCDashboardService) GetGRCRiskOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskOverviewDTO, error) {
	return gds.repo.GetGRCRiskOverview(params)
}

func (gds *GRCDashboardService) GetGRCTopCompliances(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopComplianceRegistersDTO, error) {
	return gds.repo.GetTopComplianceRegisters(params)
}

func (gds *GRCDashboardService) GetGRCTopRisks(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopRiskRegistersDTO, error) {
	return gds.repo.GetTopRiskRegisters(params)
}

func (gds *GRCDashboardService) GetGRCTopPublishedDocuments(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopDocumentsDTO, error) {
	return gds.repo.GetTopPublishedDocuments(params)
}
