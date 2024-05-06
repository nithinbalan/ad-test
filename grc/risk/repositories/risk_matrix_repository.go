package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RiskMatrixRepositoryInterface interface {
	GetByID(id uint) (*models.RiskMatrix, error)
	GetByType(status string) (*models.RiskMatrix, error)
	GetDetailsByID(id uint) (*models.RiskMatrix, error)
	Update(matrix *models.RiskMatrix, tx ...utility.DatabaseTransaction) error
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

func (repo *RiskMatrixRepository) GetByID(id uint) (*models.RiskMatrix, error) {
	var riskMatrix *models.RiskMatrix

	result := repo.db.First(&riskMatrix, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskMatrix, nil
}

func (repo *RiskMatrixRepository) GetByType(status string) (*models.RiskMatrix, error) {
	var riskMatrix *models.RiskMatrix

	result := repo.db.Where("type = ?", status).First(&riskMatrix)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskMatrix, nil
}

func (repo *RiskMatrixRepository) GetDetailsByID(id uint) (*models.RiskMatrix, error) {
	var riskMatrix *models.RiskMatrix

	result := repo.db.Preload("RiskMatrixImpact").
		Preload("RiskMatrixLikelihood").
		Preload("RiskMatrixRiskRating.RiskMatrixRiskRatingValue").
		First(&riskMatrix, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskMatrix, nil
}

func (repo *RiskMatrixRepository) Update(riskMatrix *models.RiskMatrix, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskMatrix).Error; err != nil {
		return err
	}

	return nil
}
