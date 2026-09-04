package services

import (
	"jethur/grc/models"
	"jethur/grc/project/repositories"
	"jethur/utility"
)

type ProjectSecurityGateStageServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateStage, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateStage, error)
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

func (psgss *ProjectSecurityGateStageService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateStage, error) {
	return psgss.repo.GetByID(id, tx...)
}

func (psgss *ProjectSecurityGateStageService) GetByType(stageType string) (*models.ProjectSecurityGateStage, error) {
	return psgss.repo.GetByType(stageType)
}
