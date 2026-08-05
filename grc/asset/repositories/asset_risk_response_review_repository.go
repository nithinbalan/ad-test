package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskResponseReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateRiskResponseReview(review *models.RiskResponseWorkflowReview, tx ...utility.DatabaseTransaction) error
	CreateRiskResponseWorkflowHistory(history *models.RiskResponseWorkflowHistory, tx ...utility.DatabaseTransaction) error
	UpdateRiskResponse(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error
}

type AssetRiskResponseReviewRepository struct {
	db *gorm.DB
}

func NewAssetRiskResponseReviewRepository(db *gorm.DB) AssetRiskResponseReviewRepositoryInterface {
	return &AssetRiskResponseReviewRepository{db: db}
}

func (repo *AssetRiskResponseReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetRiskResponseReviewRepository) CreateRiskResponseReview(review *models.RiskResponseWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(review).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRiskResponseReviewRepository) CreateRiskResponseWorkflowHistory(history *models.RiskResponseWorkflowHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(history).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRiskResponseReviewRepository) UpdateRiskResponse(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(riskResponse).Error
}
