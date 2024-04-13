package services

import (
	"jethur/grc/masters/models"
	"jethur/grc/masters/repositories"
)

type LifecycleStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LifecycleStatus, error)
	GetByID(id uint) (*models.LifecycleStatus, error)
	GetByType(status string) (*models.LifecycleStatus, error) // Renamed method

}

// Service and implementations
type LifecycleStatusService struct {
	repo repositories.LifecycleStatusRepositoryInterface
}

func NewLifecycleStatusService(repo repositories.LifecycleStatusRepositoryInterface) LifecycleStatusServiceInterface {
	return &LifecycleStatusService{
		repo: repo,
	}
}

func (lss *LifecycleStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LifecycleStatus, error) {
	return lss.repo.GetAll(sort, sortBy, keyword, status)
}

func (lss LifecycleStatusService) GetByID(id uint) (*models.LifecycleStatus, error) {
	return lss.repo.GetByID(id)
}

func (lss *LifecycleStatusService) GetByType(status string) (*models.LifecycleStatus, error) {
	// Fetch the lifecycle status by its status
	statusResult, err := lss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
