package repositories

import (
	"errors"

	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskInherentEvaluationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetFullDetailsByRiskID(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskEvaluation, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskEvaluation, error)
	Create(riskEvaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskEvaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AssetRiskInherentEvaluationRepository struct {
	db *gorm.DB
}

func NewAssetRiskInherentEvaluationRepository(db *gorm.DB) AssetRiskInherentEvaluationRepositoryInterface {
	return &AssetRiskInherentEvaluationRepository{
		db: db,
	}
}

func (repo *AssetRiskInherentEvaluationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetRiskInherentEvaluationRepository) GetFullDetailsByRiskID(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskEvaluation, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskEvaluation models.RiskEvaluation
	result := db.
		Preload("RiskMatrixLikelihood").
		Preload("RiskMatrixImpact").
		Preload("RiskMatrixRiskRating").
		Preload("RiskMatrixRiskRatingValue").
		Where("risk_register_id = ? AND type = ?", riskID, "inherent").
		First(&riskEvaluation)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &riskEvaluation, nil
}

func (repo *AssetRiskInherentEvaluationRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskEvaluation, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskEvaluation models.RiskEvaluation
	if err := db.First(&riskEvaluation, id).Error; err != nil {
		return nil, err
	}

	return &riskEvaluation, nil
}

func (repo *AssetRiskInherentEvaluationRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.RiskEvaluation{}, id).Error
}

func (repo *AssetRiskInherentEvaluationRepository) Create(riskEvaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(riskEvaluation).Error; err != nil {
		return 0, err
	}

	return riskEvaluation.ID, nil
}

func (repo *AssetRiskInherentEvaluationRepository) Update(riskEvaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(riskEvaluation).Error
}
