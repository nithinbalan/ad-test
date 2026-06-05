package services

import (
	"jethur/grc/models"
	"jethur/grc/settings/repositories"
)

type WorkflowCustomRoleServiceInterface interface {
	GetAll() ([]*models.WorkflowCustomRole, error)
	GetByType(roleType string) (*models.WorkflowCustomRole, error)
}

type WorkflowCustomRoleService struct {
	repo repositories.WorkflowCustomRoleRepositoryInterface
}

func NewWorkflowCustomRoleService(
	repo repositories.WorkflowCustomRoleRepositoryInterface,
) WorkflowCustomRoleServiceInterface {
	return &WorkflowCustomRoleService{
		repo: repo,
	}
}

func (s *WorkflowCustomRoleService) GetAll() ([]*models.WorkflowCustomRole, error) {
	return s.repo.GetAll()
}

func (s *WorkflowCustomRoleService) GetByType(roleType string) (*models.WorkflowCustomRole, error) {
	result, err := s.repo.GetByType(roleType)
	if err != nil {
		return nil, err
	}
	return result, nil
}
