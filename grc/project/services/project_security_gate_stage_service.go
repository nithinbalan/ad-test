package services

import (
	"jethur/grc/models"
	"jethur/grc/project/repositories"
)

type ProjectSecurityGateStageServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateStage, error)
	GetByID(id uint) (*models.ProjectSecurityGateStage, error)
	GetByType(stageType string) (*models.ProjectSecurityGateStage, error)
}

type ProjectSecurityGateStageService struct {
	repo repositories.ProjectSecurityGateStageRepositoryInterface
}

func NewProjectSecurityGateStageService(repo repositories.ProjectSecurityGateStageRepositoryInterface) ProjectSecurityGateStageServiceInterface {
	return &ProjectSecurityGateStageService{
		repo: repo,
	}
}

func (psgss *ProjectSecurityGateStageService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateStage, error) {
	return psgss.repo.GetAll(sort, sortBy, keyword, status)
}

func (psgss *ProjectSecurityGateStageService) GetByID(id uint) (*models.ProjectSecurityGateStage, error) {
	return psgss.repo.GetByID(id)
}

func (psgss *ProjectSecurityGateStageService) GetByType(stageType string) (*models.ProjectSecurityGateStage, error) {
	return psgss.repo.GetByType(stageType)
}
