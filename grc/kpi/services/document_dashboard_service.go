package services

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
)

type DocumentDashboardServiceInterface interface {
	GetDocumentSummary(params *kpiDTOs.DocumentDashboardQueryParamsDTO) (*kpiDTOs.DocumentDashboardSummaryDTO, error)
	GetDocumentStatusByType(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentStatusByTypeDTO, error)
	GetDocumentOnTimeReviewTrend(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentOnTimeReviewTrendDTO, error)
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

func (service *DocumentDashboardService) GetDocumentOnTimeReviewTrend(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentOnTimeReviewTrendDTO, error) {
	return service.repo.GetDocumentOnTimeReviewTrend(params)
}
