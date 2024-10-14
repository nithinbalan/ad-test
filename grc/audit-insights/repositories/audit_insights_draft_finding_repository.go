package repositories

import (
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditInsightsDraftFindingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
    
}

type AuditInsightsDraftFindingRepository struct {
	db *gorm.DB
}

func NewAuditInsightsDraftFindingRepository(db *gorm.DB) *AuditInsightsDraftFindingRepository {
	return &AuditInsightsDraftFindingRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditInsightsDraftFindingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}