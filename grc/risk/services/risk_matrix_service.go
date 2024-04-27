package services

import (
	"jethur/grc/risk/models"
	riskModels "jethur/grc/risk/models"
	"jethur/grc/risk/repositories"
)

type RiskMatrixServiceInterface interface {
	GetByType(status string) (*models.RiskMatrix, error)
	Update(id uint, riskMatrix *riskModels.RiskMatrix) error
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

func (rms *RiskMatrixService) GetByType(status string) (*models.RiskMatrix, error) {
	// Fetch the risk matrix by its type
	result, err := rms.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return result, nil
}

func (rms *RiskMatrixService) Update(id uint, riskMatrix *riskModels.RiskMatrix) error {
	existingRiskMatrix, err := rms.repo.GetByID(id)
	if err != nil {
		return err
	}

	existingRiskMatrix.CalculationType = riskMatrix.CalculationType

	err = rms.repo.Update(existingRiskMatrix)
	if err != nil {
		return err
	}

	return nil
}
