package services

import (
	"jethur/grc/models"
	"jethur/grc/project/repositories"
)

type ProjectSecurityGateClassificationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateClassification, error)
	GetByID(id uint) (*models.ProjectSecurityGateClassification, error)
	GetByType(classificationType string) (*models.ProjectSecurityGateClassification, error)
}

type ProjectSecurityGateClassificationService struct {
	repo repositories.ProjectSecurityGateClassificationRepositoryInterface
}

func NewProjectSecurityGateClassificationService(repo repositories.ProjectSecurityGateClassificationRepositoryInterface) ProjectSecurityGateClassificationServiceInterface {
	return &ProjectSecurityGateClassificationService{
		repo: repo,
	}
}

func (psgcs *ProjectSecurityGateClassificationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateClassification, error) {
	return psgcs.repo.GetAll(sort, sortBy, keyword, status)
}

func (psgcs *ProjectSecurityGateClassificationService) GetByID(id uint) (*models.ProjectSecurityGateClassification, error) {
	return psgcs.repo.GetByID(id)
}

func (psgcs *ProjectSecurityGateClassificationService) GetByType(classificationType string) (*models.ProjectSecurityGateClassification, error) {
	return psgcs.repo.GetByType(classificationType)
}
