package services

import (
	"jethur/grc/compliance/repositories"
	"jethur/grc/models"
)

type ComplianceStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceStatus, error)
	GetByID(id uint) (*models.ComplianceStatus, error)
}

// Service and implementations
type ComplianceStatusService struct {
	repo repositories.ComplianceStatusRepositoryInterface
}

func NewComplianceStatusService(repo repositories.ComplianceStatusRepositoryInterface) ComplianceStatusServiceInterface {
	return &ComplianceStatusService{
		repo: repo,
	}
}

func (css *ComplianceStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceStatus, error) {
	return css.repo.GetAll(sort, sortBy, keyword, status)
}

func (css ComplianceStatusService) GetByID(id uint) (*models.ComplianceStatus, error) {
	return css.repo.GetByID(id)
}
