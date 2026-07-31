package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateRiskReview(riskReview *models.RiskWorkflowReview, tx ...utility.DatabaseTransaction) error
	CreateRiskWorkflowHistory(history *models.AssetRiskWorkflowHistory, tx ...utility.DatabaseTransaction) error
	UpdateRisk(risk *models.RiskRegister, tx ...utility.DatabaseTransaction) error
	UpdateVulnerabilityStatusByID(vulnerabilityID uint, vulnerabilityStatusID uint, tx ...utility.DatabaseTransaction) error
}

type AssetRiskReviewRepository struct {
	db *gorm.DB
}

func NewAssetRiskReviewRepository(db *gorm.DB) AssetRiskReviewRepositoryInterface {
	return &AssetRiskReviewRepository{db: db}
}

func (repo *AssetRiskReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetRiskReviewRepository) CreateRiskReview(riskReview *models.RiskWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskReview).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRiskReviewRepository) CreateRiskWorkflowHistory(history *models.AssetRiskWorkflowHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(history).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRiskReviewRepository) UpdateRisk(risk *models.RiskRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(risk).Error
}

func (repo *AssetRiskReviewRepository) UpdateVulnerabilityStatusByID(vulnerabilityID uint, vulnerabilityStatusID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	updates := map[string]interface{}{
		"vulnerability_status_id": vulnerabilityStatusID,
	}

	return db.
		Model(&models.Vulnerability{}).
		Where("id = ?", vulnerabilityID).
		Updates(updates).Error
}
