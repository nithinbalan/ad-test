package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditUniverseRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	Create(auditableItems []*models.AuditableItem, tx ...utility.DatabaseTransaction) error
	UpdateMarkedAsAuditableItem(ids []uint, entityType interface{}, tx ...utility.DatabaseTransaction) error
}

// implementation
type AuditUniverseRepository struct {
	db *gorm.DB
}

func NewAuditUniverseRepository(db *gorm.DB) *AuditUniverseRepository {
	return &AuditUniverseRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditUniverseRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditUniverseRepository) Create(auditableItems []*models.AuditableItem, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all auditableItems in a single transaction
	if err := db.Create(&auditableItems).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AuditUniverseRepository) UpdateMarkedAsAuditableItem(ids []uint, entityType interface{}, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	updates := map[string]interface{}{"marked_as_auditable_item": true}

	if err := db.Model(entityType).Where("id IN (?)", ids).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}
