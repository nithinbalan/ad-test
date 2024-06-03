package services

import (
	"jethur/grc/compliance/repositories"
	"jethur/grc/models"
)

type ComplianceReviewStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceReviewStatus, error)
	GetByID(id uint) (*models.ComplianceReviewStatus, error)
	GetByType(status string) (*models.ComplianceReviewStatus, error)
}

// Service and implementations
type ComplianceReviewStatusService struct {
	repo repositories.ComplianceReviewStatusRepositoryInterface
}

func NewComplianceReviewStatusService(repo repositories.ComplianceReviewStatusRepositoryInterface) ComplianceReviewStatusServiceInterface {
	return &ComplianceReviewStatusService{
		repo: repo,
	}
}

func (crss *ComplianceReviewStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceReviewStatus, error) {
	return crss.repo.GetAll(sort, sortBy, keyword, status)
}

func (crss ComplianceReviewStatusService) GetByID(id uint) (*models.ComplianceReviewStatus, error) {
	return crss.repo.GetByID(id)
}

func (crss *ComplianceReviewStatusService) GetByType(status string) (*models.ComplianceReviewStatus, error) {
	statusResult, err := crss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
