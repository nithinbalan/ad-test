package services

import (
	"jethur/grc/logs/services"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type IntegrationCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IntegrationCategory, error)
}

// Service and implementations
type IntegrationCategoryService struct {
	repo           repositories.IntegrationCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewIntegrationCategoryService(repo repositories.IntegrationCategoryRepositoryInterface) IntegrationCategoryServiceInterface {
	return &IntegrationCategoryService{
		repo: repo,
	}
}

func (ics *IntegrationCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IntegrationCategory, error) {
	return ics.repo.GetAll(sort, sortBy, keyword, status)
}
