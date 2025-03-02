package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type ApplicabilityStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ApplicabilityStatus, error)
	GetByID(id uint) (*models.ApplicabilityStatus, error)
	GetByType(status string) (*models.ApplicabilityStatus, error)
}

// Service and implementations
type ApplicabilityStatusService struct {
	repo repositories.ApplicabilityStatusRepositoryInterface
}

func NewApplicabilityStatusService(repo repositories.ApplicabilityStatusRepositoryInterface) ApplicabilityStatusServiceInterface {
	return &ApplicabilityStatusService{
		repo: repo,
	}
}

func (ass *ApplicabilityStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ApplicabilityStatus, error) {
	return ass.repo.GetAll(sort, sortBy, keyword, status)
}

func (ass ApplicabilityStatusService) GetByID(id uint) (*models.ApplicabilityStatus, error) {
	return ass.repo.GetByID(id)
}

func (ass *ApplicabilityStatusService) GetByType(status string) (*models.ApplicabilityStatus, error) {
	// Fetch the applicability status by its status
	statusResult, err := ass.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
