package services

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
)

type DocumentDashboardServiceInterface interface {
	GetDocumentSummary() (*kpiDTOs.DocumentDashboardSummaryDTO, error)
	GetDocumentStatusByType() ([]kpiDTOs.DocumentStatusByTypeDTO, error)
	GetDocumentOnTimeReviewTrend() ([]kpiDTOs.DocumentOnTimeReviewTrendDTO, error)
}

type DocumentDashboardService struct {
	repo repositories.DocumentDashboardRepositoryInterface
}

func NewDocumentDashboardService(repo repositories.DocumentDashboardRepositoryInterface) DocumentDashboardServiceInterface {
	return &DocumentDashboardService{repo: repo}
}

func (service *DocumentDashboardService) GetDocumentSummary() (*kpiDTOs.DocumentDashboardSummaryDTO, error) {
	return service.repo.GetDocumentSummary()
}

func (service *DocumentDashboardService) GetDocumentStatusByType() ([]kpiDTOs.DocumentStatusByTypeDTO, error) {
	return service.repo.GetDocumentStatusByType()
}

func (service *DocumentDashboardService) GetDocumentOnTimeReviewTrend() ([]kpiDTOs.DocumentOnTimeReviewTrendDTO, error) {
	return service.repo.GetDocumentOnTimeReviewTrend()
}
