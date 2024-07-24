package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/settings/repositories" // Ensure this path is correct and the repository handles WorkflowLevelType
	"jethur/utility"

	"gorm.io/gorm"
)

type WorkflowLevelTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowLevelType, error)
	GetByID(id uint) (*models.WorkflowLevelType, error)
	GetByType(levelType string) (*models.WorkflowLevelType, error)
}

// Service and implementations
type WorkflowLevelTypeService struct {
	repo repositories.WorkflowLevelTypeRepositoryInterface
}

func NewWorkflowLevelTypeService(repo repositories.WorkflowLevelTypeRepositoryInterface) WorkflowLevelTypeServiceInterface {
	return &WorkflowLevelTypeService{
		repo: repo,
	}
}

func (wlt *WorkflowLevelTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowLevelType, error) {
	return wlt.repo.GetAll(sort, sortBy, keyword, status)
}

func (wlt *WorkflowLevelTypeService) GetByID(id uint) (*models.WorkflowLevelType, error) {
	workflowLevelType, err := wlt.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return workflowLevelType, nil
}

func (wlt *WorkflowLevelTypeService) GetByType(levelType string) (*models.WorkflowLevelType, error) {
	workflowLevelType, err := wlt.repo.GetByType(levelType)
	if err != nil {
		return nil, err
	}
	return workflowLevelType, nil
}
