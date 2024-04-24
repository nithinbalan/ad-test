package repositories

import (
	"jethur/grc/risk/models"

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
	if result.Error != nil {
		return nil, result.Error
	}

	return riskMatrix, nil
}
