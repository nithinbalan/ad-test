package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateRiskReview(auditReview *models.RiskWorkflowReview, tx ...utility.DatabaseTransaction) error

	UpdateRisk(risk *models.RiskRegister, tx ...utility.DatabaseTransaction) error
}

type RiskReviewRepository struct {
	db *gorm.DB
}

func NewRiskReviewRepository(db *gorm.DB) *RiskReviewRepository {
	return &RiskReviewRepository{db: db}
}

// BeginTransaction starts a new transaction.
func (repo *RiskReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskReviewRepository) CreateRiskReview(auditReview *models.RiskWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditReview).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskReviewRepository) UpdateRisk(section *models.RiskRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(section).Error
}
