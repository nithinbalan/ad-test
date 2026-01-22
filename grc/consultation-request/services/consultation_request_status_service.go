package services

import (
	"jethur/grc/consultation-request/repositories"
	"jethur/grc/models"
)

// Interface
type ConsultationRequestStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ConsultationRequestStatus, error)
	GetByID(id uint) (*models.ConsultationRequestStatus, error)
	GetByType(status string) (*models.ConsultationRequestStatus, error)
}

// Service implementation
type ConsultationRequestStatusService struct {
	repo repositories.ConsultationRequestStatusRepositoryInterface
}

func NewConsultationRequestStatusService(
	repo repositories.ConsultationRequestStatusRepositoryInterface,
) ConsultationRequestStatusServiceInterface {
	return &ConsultationRequestStatusService{
		repo: repo,
	}
}

func (crs *ConsultationRequestStatusService) GetAll(
	sort string, sortBy string, keyword string, status string,
) ([]*models.ConsultationRequestStatus, error) {
	return crs.repo.GetAll(sort, sortBy, keyword, status)
}

func (crs *ConsultationRequestStatusService) GetByID(id uint) (*models.ConsultationRequestStatus, error) {
	return crs.repo.GetByID(id)
}

func (crs *ConsultationRequestStatusService) GetByType(status string) (*models.ConsultationRequestStatus, error) {
	statusResult, err := crs.repo.GetByType(status)
	if err != nil {
		return nil, err
	}
	return statusResult, nil
}
