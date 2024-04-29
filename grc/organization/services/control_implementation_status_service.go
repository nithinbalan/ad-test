package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
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
	return ciss.repo.GetAll(sort, sortBy, keyword, status)
}

func (ciss *ControlImplementationStatusService) GetByID(id uint) (*models.ControlImplementationStatus, error) {
	controlImplementationStatus, err := ciss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return controlImplementationStatus, nil
}

func (ciss *ControlImplementationStatusService) GetByType(status string) (*models.ControlImplementationStatus, error) {
	controlImplementationStatus, err := ciss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return controlImplementationStatus, nil
}
