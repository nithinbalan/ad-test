package services

import (
	"jethur/grc/organization/models"
	"jethur/grc/organization/repositories"
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
	statusResult, err := cess.repo.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}

func (cess ControlEfficiencyStatusService) GetByID(id uint) (*models.ControlEfficiencyStatus, error) {
	return cess.repo.GetByID(id)
}

func (cess *ControlEfficiencyStatusService) GetByType(status string) (*models.ControlEfficiencyStatus, error) {
	statusResult, err := cess.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
