package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditInsightsDraftFindingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateConversation(conversation *models.AuditDraftFindingConversation, tx ...utility.DatabaseTransaction) (uint, error)
    
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

func (repo *AuditInsightsDraftFindingRepository) CreateConversation(conversation *models.AuditDraftFindingConversation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(conversation).Error; err != nil {
		return 0, err
	}

	return conversation.ID, nil
}