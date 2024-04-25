package services

import (
	"jethur/grc/organization/models"
	"jethur/grc/organization/repositories"
)

type ControlImplementationStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlImplementationStatus, error)
	GetByID(id uint) (*models.ControlImplementationStatus, error)
	GetByType(status string) (*models.ControlImplementationStatus, error) // Renamed method
}

// Service and implementations
type ControlImplementationStatusService struct {
	repo repositories.ControlImplementationStatusRepositoryInterface
}

func NewControlImplementationStatusService(repo repositories.ControlImplementationStatusRepositoryInterface) ControlImplementationStatusServiceInterface {
	return &ControlImplementationStatusService{
		repo: repo,
	}
}

func (ciss *ControlImplementationStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlImplementationStatus, error) {
	statusResult, err := ciss.repo.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}

func (ciss ControlImplementationStatusService) GetByID(id uint) (*models.ControlImplementationStatus, error) {
	return ciss.repo.GetByID(id)
}

func (ciss *ControlImplementationStatusService) GetByType(status string) (*models.ControlImplementationStatus, error) {
	statusResult, err := ciss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
