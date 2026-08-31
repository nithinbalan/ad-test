package services

import (
	"jethur/grc/models"
	"jethur/grc/project/repositories"
)

type ProjectPhaseServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectPhase, error)
	GetByID(id uint) (*models.ProjectPhase, error)
	GetByType(status string) (*models.ProjectPhase, error)
}

// Service and implementations
type ProjectPhaseService struct {
	repo repositories.ProjectPhaseRepositoryInterface
}

func NewProjectPhaseService(repo repositories.ProjectPhaseRepositoryInterface) ProjectPhaseServiceInterface {
	return &ProjectPhaseService{
		repo: repo,
	}
}

func (pps *ProjectPhaseService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectPhase, error) {
	return pps.repo.GetAll(sort, sortBy, keyword, status)
}

func (pps *ProjectPhaseService) GetByID(id uint) (*models.ProjectPhase, error) {
	return pps.repo.GetByID(id)
}

func (pps *ProjectPhaseService) GetByType(status string) (*models.ProjectPhase, error) {
	return pps.repo.GetByType(status)
}