package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditPlanCommencementLetterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	Create(auditPlanCommencementLetter *models.AuditPlanCommencementLetter, tx ...utility.DatabaseTransaction) (uint, error)
}

type AuditPlanCommencementLetterRepository struct {
	db *gorm.DB
}

func NewAuditPlanCommencementLetterRepository(db *gorm.DB) *AuditPlanCommencementLetterRepository {
	return &AuditPlanCommencementLetterRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditPlanCommencementLetterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditPlanCommencementLetterRepository) Create(auditPlanCommencementLetter *models.AuditPlanCommencementLetter, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditPlanCommencementLetter).Error; err != nil {
		return 0, err
	}

	return auditPlanCommencementLetter.ID, nil
}
