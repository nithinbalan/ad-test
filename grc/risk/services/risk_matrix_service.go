package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskMatrixServiceInterface interface {
	GetByType(status string) (*models.RiskMatrix, error)
	GetDetailsByID(id uint) (*models.RiskMatrix, error)
	Update(id uint, riskMatrix *models.RiskMatrix, tx utility.DatabaseTransaction) error
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, status)
		}
		return nil, err
	}

	return result, nil
}
func (rms *RiskMatrixService) GetDetailsByID(id uint) (*models.RiskMatrix, error) {
	// Fetch the risk matrix by its ID
	result, err := rms.repo.GetDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return result, nil

}

func (rms *RiskMatrixService) Update(id uint, riskMatrix *models.RiskMatrix, tx utility.DatabaseTransaction) error {
	existingRiskMatrix, err := rms.repo.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utility.NewResourceNotFoundError(err, id)
	}

	existingRiskMatrix.CalculationType = riskMatrix.CalculationType

	err = rms.repo.Update(existingRiskMatrix, tx)
	if err != nil {
		return err
	}

	return nil
}
