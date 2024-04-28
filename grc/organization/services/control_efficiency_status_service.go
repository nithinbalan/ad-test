package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ControlEfficiencyStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlEfficiencyStatus, error)
	GetByID(id uint) (*models.ControlEfficiencyStatus, error)
	GetByType(status string) (*models.ControlEfficiencyStatus, error) // Renamed method
}

// Service and implementations
type ControlEfficiencyStatusService struct {
	repo repositories.ControlEfficiencyStatusRepositoryInterface
}

func NewControlEfficiencyStatusService(repo repositories.ControlEfficiencyStatusRepositoryInterface) ControlEfficiencyStatusServiceInterface {
	return &ControlEfficiencyStatusService{
		repo: repo,
	}
}

func (cess *ControlEfficiencyStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlEfficiencyStatus, error) {
	return cess.repo.GetAll(sort, sortBy, keyword, status)
}

func (cess ControlEfficiencyStatusService) GetByID(id uint) (*models.ControlEfficiencyStatus, error) {
	controlEfficiencyStatus, err := cess.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return controlEfficiencyStatus, nil
}

func (cess *ControlEfficiencyStatusService) GetByType(status string) (*models.ControlEfficiencyStatus, error) {
	controlEfficiencyStatus, err := cess.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return controlEfficiencyStatus, nil
}
