package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type WorkflowStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowStatus, error)
	GetByID(id uint) (*models.WorkflowStatus, error)
	GetByType(status string) (*models.WorkflowStatus, error)
}

// Service and implementations
type WorkflowStatusService struct {
	repo repositories.WorkflowStatusRepositoryInterface
}

func NewWorkflowStatusService(repo repositories.WorkflowStatusRepositoryInterface) WorkflowStatusServiceInterface {
	return &WorkflowStatusService{
		repo: repo,
	}
}

func (wss *WorkflowStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowStatus, error) {
	return wss.repo.GetAll(sort, sortBy, keyword, status)
}

func (wss WorkflowStatusService) GetByID(id uint) (*models.WorkflowStatus, error) {
	return wss.repo.GetByID(id)
}

func (wss *WorkflowStatusService) GetByType(status string) (*models.WorkflowStatus, error) {
	// Fetch the action plan status by its status
	statusResult, err := wss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
