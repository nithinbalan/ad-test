package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type DocumentAcknowledgeUserRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	Create(documentAcknowledgeUser *models.DocumentAcknowledgeUser, tx ...utility.DatabaseTransaction) (*models.DocumentAcknowledgeUser, error)
	Delete(documentID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
}

// Implementation
type DocumentAcknowledgeUserRepository struct {
	db *gorm.DB
}

func NewDocumentAcknowledgeUserRepository(db *gorm.DB) *DocumentAcknowledgeUserRepository {
	return &DocumentAcknowledgeUserRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *DocumentAcknowledgeUserRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentAcknowledgeUserRepository) Create(documentAcknowledgeUser *models.DocumentAcknowledgeUser, tx ...utility.DatabaseTransaction) (*models.DocumentAcknowledgeUser, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	result := db.Create(&documentAcknowledgeUser)

	return documentAcknowledgeUser, result.Error
}

func (repo *DocumentAcknowledgeUserRepository) Delete(documentID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	condition := map[string]interface{}{
		"document_id": documentID,
		"user_id":     currentUserID,
	}

	// Find the record based on the condition
	var documentAcknowledgeUser models.DocumentAcknowledgeUser
	if err := db.Where(condition).First(&documentAcknowledgeUser).Error; err != nil {
		return err
	}

	// Delete the record
	result := db.Delete(&documentAcknowledgeUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}


