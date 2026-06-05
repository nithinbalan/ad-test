package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceRegisterReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateComplianceRegisterReview(review *models.ComplianceRegisterWorkflowReview, tx ...utility.DatabaseTransaction) error
	CreateComplianceRegisterWorkflowHistory(historyReview *models.ComplianceRegisterWorkflowHistory, tx ...utility.DatabaseTransaction) error
	UpdateComplianceRegister(complianceRegister *models.ComplianceRegister, tx ...utility.DatabaseTransaction) error
}

type ComplianceRegisterReviewRepository struct {
	db *gorm.DB
}

func NewComplianceRegisterReviewRepository(db *gorm.DB) *ComplianceRegisterReviewRepository {
	return &ComplianceRegisterReviewRepository{db: db}
}

func (repo *ComplianceRegisterReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceRegisterReviewRepository) CreateComplianceRegisterReview(review *models.ComplianceRegisterWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(review).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRegisterReviewRepository) CreateComplianceRegisterWorkflowHistory(historyReview *models.ComplianceRegisterWorkflowHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(historyReview).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRegisterReviewRepository) UpdateComplianceRegister(complianceRegister *models.ComplianceRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(complianceRegister).Error
}
