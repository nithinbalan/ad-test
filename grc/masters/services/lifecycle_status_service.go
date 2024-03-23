package services

import (
	"jethur/grc/masters/models"
	"jethur/grc/masters/repositories"
)

type LifecycleStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LifecycleStatus, error)
	GetByID(id uint) (*models.LifecycleStatus, error)
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
