package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskResponseReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateRiskResponseReview(review *models.RiskResponseWorkflowReview, tx ...utility.DatabaseTransaction) error
	UpdateRiskResponse(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error
}

type RiskResponseReviewRepository struct {
	db *gorm.DB
}

func NewRiskResponseReviewRepository(db *gorm.DB) *RiskResponseReviewRepository {
	return &RiskResponseReviewRepository{db: db}
}

func (repo *RiskResponseReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskResponseReviewRepository) CreateRiskResponseReview(review *models.RiskResponseWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(review).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskResponseReviewRepository) UpdateRiskResponse(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(riskResponse).Error
}
