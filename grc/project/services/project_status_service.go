package services

import (
	"jethur/grc/models"
	"jethur/grc/project/repositories"
)

type ProjectStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectStatus, error)
	GetByID(id uint) (*models.ProjectStatus, error)
	GetByType(status string) (*models.ProjectStatus, error)
}

// Service and implementations
type ProjectStatusService struct {
	repo repositories.ProjectStatusRepositoryInterface
}

func NewProjectStatusService(repo repositories.ProjectStatusRepositoryInterface) ProjectStatusServiceInterface {
	return &ProjectStatusService{
		repo: repo,
	}
}

func (pss *ProjectStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectStatus, error) {
	return pss.repo.GetAll(sort, sortBy, keyword, status)
}

func (pss *ProjectStatusService) GetByID(id uint) (*models.ProjectStatus, error) {
	return pss.repo.GetByID(id)
}

func (pss *ProjectStatusService) GetByType(status string) (*models.ProjectStatus, error) {
	return pss.repo.GetByType(status)
}
