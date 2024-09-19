package repositories

import (
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type DocumentAcknowledgeUserRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(documentID uint) (*documentDTOs.DocumentAcknowledgementListResponseDTO, error)
	Create(documentAcknowledgeUser *models.DocumentAcknowledgeUser, tx ...utility.DatabaseTransaction) (*models.DocumentAcknowledgeUser, error)
	Exists(documentID uint, userID uint, tx ...utility.DatabaseTransaction) (bool, error)
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

func (repo *DocumentAcknowledgeUserRepository) GetAll(documentID uint) (*documentDTOs.DocumentAcknowledgementListResponseDTO, error) {
	var acknowledgementList []documentDTOs.DocumentAcknowledgementListDTO

	// Construct raw SQL query with UNION to combine shared users and department users
	rawSQL := `
        SELECT users.id, users.first_name_en, users.first_name_ar, users.email,
            CASE WHEN document_acknowledge_users.id IS NOT NULL THEN 1 ELSE 0 END AS acknowledged
        FROM document_shared_users
        JOIN users ON users.id = document_shared_users.user_id
        LEFT JOIN document_acknowledge_users ON document_acknowledge_users.user_id = users.id AND document_acknowledge_users.document_id = ?
        WHERE document_shared_users.document_id = ?

        UNION

        SELECT users.id, users.first_name_en, users.first_name_ar, users.email,
            CASE WHEN document_acknowledge_users.id IS NOT NULL THEN 1 ELSE 0 END AS acknowledged
        FROM document_shared_units
        JOIN organization_units ON organization_units.id = document_shared_units.organization_unit_id
        JOIN users ON users.organization_unit_1_id = organization_units.id
        LEFT JOIN document_acknowledge_users ON document_acknowledge_users.user_id = users.id AND document_acknowledge_users.document_id = ?
        WHERE document_shared_units.document_id = ?
    `

	// Execute the raw SQL query with GORM
	query := repo.db.Raw(rawSQL, documentID, documentID, documentID, documentID).Scan(&acknowledgementList)

	if query.Error != nil {
		return nil, query.Error
	}

	response := &documentDTOs.DocumentAcknowledgementListResponseDTO{
		Acknowledgements: acknowledgementList,
	}

	return response, nil
}

func (repo *DocumentAcknowledgeUserRepository) Create(documentAcknowledgeUser *models.DocumentAcknowledgeUser, tx ...utility.DatabaseTransaction) (*models.DocumentAcknowledgeUser, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	result := db.Create(&documentAcknowledgeUser)

	return documentAcknowledgeUser, result.Error
}

func (repo *DocumentAcknowledgeUserRepository) Exists(documentID uint, userID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.DocumentAcknowledgeUser{}).
		Where("document_id = ? AND user_id = ?", documentID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
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
