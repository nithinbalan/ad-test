package services

import (
	"errors"
	"jethur/grc/incident/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentStatus, error)
	GetByID(id uint) (*models.IncidentStatus, error)
	GetByType(status string) (*models.IncidentStatus, error)
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
	incidentStatus, err := iss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentStatus, nil
}

func (iss *IncidentStatusService) GetByType(status string) (*models.IncidentStatus, error) {
	statusResult, err := iss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
