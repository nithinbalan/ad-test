package services

import (
	"jethur/grc/risk/models"
	"jethur/grc/risk/repositories"
)

type RiskMatrixServiceInterface interface {
	GetByType(status string) (*models.RiskMatrix, error)
}

// Service and implementations
type RiskMatrixService struct {
	repo repositories.RiskMatrixRepositoryInterface
}

func NewRiskMatrixService(repo repositories.RiskMatrixRepositoryInterface) RiskMatrixServiceInterface {
	return &RiskMatrixService{
		repo: repo,
	}
}

func (prss *RiskMatrixService) GetByType(status string) (*models.RiskMatrix, error) {
	// Fetch the risk matrix by its type
	result, err := prss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return result, nil
}
