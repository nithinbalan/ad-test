package services

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
)

type DocumentDashboardServiceInterface interface {
	GetDocumentSummary(params *kpiDTOs.DocumentDashboardQueryParamsDTO) (*kpiDTOs.DocumentDashboardSummaryDTO, error)
	GetDocumentStatusByType(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentStatusByTypeDTO, error)
	GetDocumentsByConfidentialityClassification(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentConfidentialityClassificationDTO, error)
	GetDocumentsByResponsibleUnit(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentResponsibleUnitDTO, error)
	GetDocumentReviewCountTrend(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentReviewCountTrendDTO, error)
}

type DocumentDashboardService struct {
	repo repositories.DocumentDashboardRepositoryInterface
}

func NewDocumentDashboardService(repo repositories.DocumentDashboardRepositoryInterface) DocumentDashboardServiceInterface {
	return &DocumentDashboardService{repo: repo}
}

func (service *DocumentDashboardService) GetDocumentSummary(params *kpiDTOs.DocumentDashboardQueryParamsDTO) (*kpiDTOs.DocumentDashboardSummaryDTO, error) {
	return service.repo.GetDocumentSummary(params)
}

func (service *DocumentDashboardService) GetDocumentStatusByType(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentStatusByTypeDTO, error) {
	return service.repo.GetDocumentStatusByType(params)
}

func (service *DocumentDashboardService) GetDocumentsByConfidentialityClassification(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentConfidentialityClassificationDTO, error) {
	return service.repo.GetDocumentsByConfidentialityClassification(params)
}

func (service *DocumentDashboardService) GetDocumentsByResponsibleUnit(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentResponsibleUnitDTO, error) {
	return service.repo.GetDocumentsByResponsibleUnit(params)
}

func (service *DocumentDashboardService) GetDocumentReviewCountTrend(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentReviewCountTrendDTO, error) {
	return service.repo.GetDocumentReviewCountTrend(params)
}
