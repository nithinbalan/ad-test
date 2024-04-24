package services

import (
	"jethur/grc/organization/models"
	"jethur/grc/organization/repositories"
)

type ControlTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlType, error)
}

// Service and implementations
type ControlTypeService struct {
	repo repositories.ControlTypeRepositoryInterface
}

func NewControlTypeService(repo repositories.ControlTypeRepositoryInterface) ControlTypeServiceInterface {
	return &ControlTypeService{
		repo: repo,
	}
}

func (cts *ControlTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlType, error) {
	return cts.repo.GetAll(sort, sortBy, keyword, status)
}
