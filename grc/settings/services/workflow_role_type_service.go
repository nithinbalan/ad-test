package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/settings/repositories" // Ensure this path is correct and the repository handles WorkflowRoleType
	"jethur/utility"

	"gorm.io/gorm"
)

type WorkflowRoleTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowRoleType, error)
	GetByID(id uint) (*models.WorkflowRoleType, error)
	GetByType(roleType string) (*models.WorkflowRoleType, error)
}

// Service and implementations
type WorkflowRoleTypeService struct {
	repo repositories.WorkflowRoleTypeRepositoryInterface
}

func NewWorkflowRoleTypeService(repo repositories.WorkflowRoleTypeRepositoryInterface) WorkflowRoleTypeServiceInterface {
	return &WorkflowRoleTypeService{
		repo: repo,
	}
}

func (wrt *WorkflowRoleTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowRoleType, error) {
	return wrt.repo.GetAll(sort, sortBy, keyword, status)
}

func (wrt *WorkflowRoleTypeService) GetByID(id uint) (*models.WorkflowRoleType, error) {
	workflowRoleType, err := wrt.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return workflowRoleType, nil
}

func (wrt *WorkflowRoleTypeService) GetByType(roleType string) (*models.WorkflowRoleType, error) {
	workflowRoleType, err := wrt.repo.GetByType(roleType)
	if err != nil {
		return nil, err
	}
	return workflowRoleType, nil
}
