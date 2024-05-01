package services

import (
	"jethur/grc/incident/repositories"
	"jethur/grc/models"
)

type IncidentStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentStatus, error)
	GetByID(id uint) (*models.IncidentStatus, error)
}

// Service and implementations
type IncidentStatusService struct {
	repo repositories.IncidentStatusRepositoryInterface
}

func NewIncidentStatusService(repo repositories.IncidentStatusRepositoryInterface) IncidentStatusServiceInterface {
	return &IncidentStatusService{
		repo: repo,
	}
}

func (iss *IncidentStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentStatus, error) {
	return iss.repo.GetAll(sort, sortBy, keyword, status)
}

func (iss IncidentStatusService) GetByID(id uint) (*models.IncidentStatus, error) {
	return iss.repo.GetByID(id)
}
