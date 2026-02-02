package services

import (
	"errors"
	"jethur/grc/compliance/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentImpactTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentImpactType, error)
	GetByID(id uint) (*models.ComplianceAssessmentImpactType, error)
	GetByType(status string) (*models.ComplianceAssessmentImpactType, error)
}

type ComplianceAssessmentImpactTypeService struct {
	repo repositories.ComplianceAssessmentImpactTypeRepositoryInterface
}

func NewComplianceAssessmentImpactTypeService(repo repositories.ComplianceAssessmentImpactTypeRepositoryInterface) ComplianceAssessmentImpactTypeServiceInterface {
	return &ComplianceAssessmentImpactTypeService{
		repo: repo,
	}
}

func (service *ComplianceAssessmentImpactTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentImpactType, error) {
	return service.repo.GetAll(sort, sortBy, keyword, status)
}

func (service *ComplianceAssessmentImpactTypeService) GetByID(id uint) (*models.ComplianceAssessmentImpactType, error) {
	impactType, err := service.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return impactType, nil
}

func (service *ComplianceAssessmentImpactTypeService) GetByType(status string) (*models.ComplianceAssessmentImpactType, error) {
	return service.repo.GetByType(status)
}
