package services

import (
	"jethur/grc/risk/models"
	"jethur/grc/risk/repositories"
)

type PotentialRiskStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PotentialRiskStatus, error)
	GetByID(id uint) (*models.PotentialRiskStatus, error)
	GetByStatus(status string) (*models.PotentialRiskStatus, error)
}

// Service and implementations
type PotentialRiskStatusService struct {
	repo repositories.PotentialRiskStatusRepositoryInterface
}

func NewPotentialRiskStatusService(repo repositories.PotentialRiskStatusRepositoryInterface) PotentialRiskStatusServiceInterface {
	return &PotentialRiskStatusService{
		repo: repo,
	}
}

func (prss *PotentialRiskStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PotentialRiskStatus, error) {
	return prss.repo.GetAll(sort, sortBy, keyword, status)
}

func (prss PotentialRiskStatusService) GetByID(id uint) (*models.PotentialRiskStatus, error) {
	return prss.repo.GetByID(id)
}

func (prss *PotentialRiskStatusService) GetByStatus(status string) (*models.PotentialRiskStatus, error) {
	// Fetch the potential risk status by its status
	statusResult, err := prss.repo.GetByStatus(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
