package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskMatrixRiskRatingValueRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskMatrixID uint) ([]*riskDTOs.RiskMatrixRiskRatingValueListResponseDTO, error)
	Create(riskMatrixRiskRatingValue []*models.RiskMatrixRiskRatingValue, riskMatrixRiskRatingIds []uint, tx ...utility.DatabaseTransaction) error
}

type RiskMatrixRiskRatingValueRepository struct {
	db *gorm.DB
}

func NewRiskMatrixRiskRatingValueRepository(db *gorm.DB) *RiskMatrixRiskRatingValueRepository {
	return &RiskMatrixRiskRatingValueRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *RiskMatrixRiskRatingValueRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskMatrixRiskRatingValueRepository) GetAll(riskMatrixID uint) ([]*riskDTOs.RiskMatrixRiskRatingValueListResponseDTO, error) {
	var riskRatingValues []*riskDTOs.RiskMatrixRiskRatingValueListResponseDTO

	query := repo.db.Table("risk_matrix_risk_rating_values").
		Select(`
		risk_matrix_risk_rating_values.*`).
		Joins(`
		    LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_matrix_risk_rating_values.risk_matrix_risk_rating_id
		`).
		Where("risk_matrix_risk_ratings.risk_matrix_id = ?", riskMatrixID)

	result := query.Find(&riskRatingValues)

	return riskRatingValues, result.Error

}

func (repo *RiskMatrixRiskRatingValueRepository) Create(riskMatrixRiskRatingValue []*models.RiskMatrixRiskRatingValue, riskMatrixRiskRatingIds []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Where("risk_matrix_risk_rating_values.risk_matrix_risk_rating_id IN (?)", riskMatrixRiskRatingIds).Delete(&models.RiskMatrixRiskRatingValue{})
	if result.Error != nil {
		return result.Error
	}

	// Save all values in a single transaction
	if err := db.Create(&riskMatrixRiskRatingValue).Error; err != nil {
		return err
	}
	return nil
}
