package repositories

import (
	"errors"
	"jethur/grc/risk/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RiskMatrixRepositoryInterface interface {
	GetByType(status string) (*models.RiskMatrix, error)
}

// Implementation
type RiskMatrixRepository struct {
	db *gorm.DB
}

func NewRiskMatrixRepository(db *gorm.DB) *RiskMatrixRepository {
	return &RiskMatrixRepository{
		db: db,
	}
}

func (repo *RiskMatrixRepository) GetByType(status string) (*models.RiskMatrix, error) {
	var riskMatrix *models.RiskMatrix

	result := repo.db.Where("type = ?", status).First(&riskMatrix)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Risk matrix with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskMatrix, nil
}
