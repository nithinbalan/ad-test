package services

import (
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
)

type ComplianceStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceStatus, error)
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
