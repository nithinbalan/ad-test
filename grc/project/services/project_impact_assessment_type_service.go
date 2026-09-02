package services

import (
	"jethur/grc/models"
	"jethur/grc/project/repositories"
)

type ProjectImpactAssessmentTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectImpactAssessmentType, error)
	GetByID(id uint) (*models.ProjectImpactAssessmentType, error)
	GetByType(status string) (*models.ProjectImpactAssessmentType, error)
}

type ProjectImpactAssessmentTypeService struct {
	repo repositories.ProjectImpactAssessmentTypeRepositoryInterface
}

func NewProjectImpactAssessmentTypeService(repo repositories.ProjectImpactAssessmentTypeRepositoryInterface) ProjectImpactAssessmentTypeServiceInterface {
	return &ProjectImpactAssessmentTypeService{
		repo: repo,
	}
}

func (piats *ProjectImpactAssessmentTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectImpactAssessmentType, error) {
	return piats.repo.GetAll(sort, sortBy, keyword, status)
}

func (piats *ProjectImpactAssessmentTypeService) GetByID(id uint) (*models.ProjectImpactAssessmentType, error) {
	return piats.repo.GetByID(id)
}

func (piats *ProjectImpactAssessmentTypeService) GetByType(status string) (*models.ProjectImpactAssessmentType, error) {
	return piats.repo.GetByType(status)
}
