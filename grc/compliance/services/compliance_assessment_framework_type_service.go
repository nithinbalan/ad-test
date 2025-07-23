package services

import (
	"errors"
	"jethur/grc/compliance/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentFrameworkTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentFrameworkType, error)
	GetByID(id uint) (*models.ComplianceAssessmentFrameworkType, error)
	GetByType(typeValue string) (*models.ComplianceAssessmentFrameworkType, error)
}

type ComplianceAssessmentFrameworkTypeService struct {
	repo repositories.ComplianceAssessmentFrameworkTypeRepositoryInterface
}

func NewComplianceAssessmentFrameworkTypeService(repo repositories.ComplianceAssessmentFrameworkTypeRepositoryInterface) ComplianceAssessmentFrameworkTypeServiceInterface {
	return &ComplianceAssessmentFrameworkTypeService{
		repo: repo,
	}
}

func (caft *ComplianceAssessmentFrameworkTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentFrameworkType, error) {
	return caft.repo.GetAll(sort, sortBy, keyword, status)
}

func (caft *ComplianceAssessmentFrameworkTypeService) GetByID(id uint) (*models.ComplianceAssessmentFrameworkType, error) {
	frameworkType, err := caft.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return frameworkType, nil
}

func (caft *ComplianceAssessmentFrameworkTypeService) GetByType(typeValue string) (*models.ComplianceAssessmentFrameworkType, error) {
	frameworkTypeResult, err := caft.repo.GetByType(typeValue)
	if err != nil {
		return nil, err
	}
	return frameworkTypeResult, nil
}
