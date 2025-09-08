package services

import (
	"jethur/grc/masters/dtos"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type WorkflowContextServiceInterface interface {
	GetAll(params dtos.WorkflowContextListQueryParamsDTO) ([]*models.WorkflowContext, error)
	GetByID(id uint) (*models.WorkflowContext, error)
	GetByType(status string) (*models.WorkflowContext, error)
}

// Service and implementations
type WorkflowContextService struct {
	repo repositories.WorkflowContextRepositoryInterface
}

func NewWorkflowContextService(repo repositories.WorkflowContextRepositoryInterface) WorkflowContextServiceInterface {
	return &WorkflowContextService{
		repo: repo,
	}
}

func (wss *WorkflowContextService) GetAll(params dtos.WorkflowContextListQueryParamsDTO) ([]*models.WorkflowContext, error) {
	return wss.repo.GetAll(params)
}

func (wss WorkflowContextService) GetByID(id uint) (*models.WorkflowContext, error) {
	return wss.repo.GetByID(id)
}

func (wss *WorkflowContextService) GetByType(status string) (*models.WorkflowContext, error) {
	// Fetch the action plan status by its status
	statusResult, err := wss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
