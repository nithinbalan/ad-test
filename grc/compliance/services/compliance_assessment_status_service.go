package services

import (
	"errors"
	"jethur/grc/compliance/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentStatus, error)
	GetByID(id uint) (*models.ComplianceAssessmentStatus, error)
	GetByType(status string) (*models.ComplianceAssessmentStatus, error)
}

type ComplianceAssessmentStatusService struct {
	repo repositories.ComplianceAssessmentStatusRepositoryInterface
}

func NewComplianceAssessmentStatusService(repo repositories.ComplianceAssessmentStatusRepositoryInterface) ComplianceAssessmentStatusServiceInterface {
	return &ComplianceAssessmentStatusService{
		repo: repo,
	}
}

func (cass *ComplianceAssessmentStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentStatus, error) {
	return cass.repo.GetAll(sort, sortBy, keyword, status)
}

func (cass *ComplianceAssessmentStatusService) GetByID(id uint) (*models.ComplianceAssessmentStatus, error) {
	status, err := cass.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return status, nil
}

func (cass *ComplianceAssessmentStatusService) GetByType(status string) (*models.ComplianceAssessmentStatus, error) {
	statusResult, err := cass.repo.GetByType(status)
	if err != nil {
		return nil, err
	}
	return statusResult, nil
}
