package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type ActionPlanStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ActionPlanStatus, error)
	GetByID(id uint) (*models.ActionPlanStatus, error)
	GetByType(status string) (*models.ActionPlanStatus, error)
}

// Service and implementations
type ActionPlanStatusService struct {
	repo repositories.ActionPlanStatusRepositoryInterface
}

func NewActionPlanStatusService(repo repositories.ActionPlanStatusRepositoryInterface) ActionPlanStatusServiceInterface {
	return &ActionPlanStatusService{
		repo: repo,
	}
}

func (apss *ActionPlanStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ActionPlanStatus, error) {
	return apss.repo.GetAll(sort, sortBy, keyword, status)
}

func (apss ActionPlanStatusService) GetByID(id uint) (*models.ActionPlanStatus, error) {
	return apss.repo.GetByID(id)
}

func (apss *ActionPlanStatusService) GetByType(status string) (*models.ActionPlanStatus, error) {
	// Fetch the action plan status by its status
	statusResult, err := apss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
