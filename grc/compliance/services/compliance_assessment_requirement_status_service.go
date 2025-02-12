package services

import (
	"errors"
	"jethur/grc/compliance/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentRequirementStatus, error)
	GetByID(id uint) (*models.ComplianceAssessmentRequirementStatus, error)
	GetByType(status string) (*models.ComplianceAssessmentRequirementStatus, error)
}

type ComplianceAssessmentRequirementStatusService struct {
	repo repositories.ComplianceAssessmentRequirementStatusRepositoryInterface
}

func NewComplianceAssessmentRequirementStatusService(repo repositories.ComplianceAssessmentRequirementStatusRepositoryInterface) ComplianceAssessmentRequirementStatusServiceInterface {
	return &ComplianceAssessmentRequirementStatusService{
		repo: repo,
	}
}

func (cars *ComplianceAssessmentRequirementStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentRequirementStatus, error) {
	return cars.repo.GetAll(sort, sortBy, keyword, status)
}

func (cars *ComplianceAssessmentRequirementStatusService) GetByID(id uint) (*models.ComplianceAssessmentRequirementStatus, error) {
	status, err := cars.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return status, nil
}

func (cars *ComplianceAssessmentRequirementStatusService) GetByType(status string) (*models.ComplianceAssessmentRequirementStatus, error) {
	statusResult, err := cars.repo.GetByType(status)
	if err != nil {
		return nil, err
	}
	return statusResult, nil
}
