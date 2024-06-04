package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceStatusUpdateListDTO, error)
	GetByID(id uint) (*models.ComplianceStatusUpdate, error)
	Create(ComplianceStatusUpdate *models.ComplianceStatusUpdate) (uint, error)
	Update(ComplianceStatusUpdate *models.ComplianceStatusUpdate) error

	//Document
	CreateDocuments(docs []*models.ComplianceStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ComplianceStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceStatusUpdateRepository struct {
	db *gorm.DB
}

func NewComplianceStatusUpdateRepository(db *gorm.DB) *ComplianceStatusUpdateRepository {
	return &ComplianceStatusUpdateRepository{
		db: db,
	}
}

func (repo *ComplianceStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceStatusUpdateRepository) GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceStatusUpdateListDTO, error) {

	var complianceUpdates []*complianceDTOs.ComplianceStatusUpdateListDTO

	query := repo.db.Table("compliance_status_updates").
		Select(`
            compliance_status_updates.*,
			compliance_register.title_en AS compliance_title_en,
			compliance_register.title_ar AS compliance_title_ar,
			compliance_register.description_en AS compliance_description_en,
			compliance_register.description_ar AS compliance_description_ar,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
			created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar 
        `).
		Joins(`
		    LEFT JOIN compliance_register ON compliance_register.id = compliance_status_updates.compliance_register_id
		    LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_status_updates.compliance_status_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_status_updates.created_by
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_status_updates.updated_by
		`).
		Where("compliance_status_updates.compliance_register_id = ?", complianceID)

	if keyword != "" {
		query = query.Where(`
			(compliance_register.title LIKE ? OR
			compliance_status_updates.comment LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "compliance_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&complianceUpdates)

    if result.Error != nil {
        return nil, result.Error
    }

	// Fetching documents for each compliance status updates
	for _, update := range complianceUpdates {
		var documents []*models.ComplianceStatusUpdateDocument
		repo.db.Model(&models.ComplianceStatusUpdateDocument{}).Where("compliance_status_update_id = ?", update.ID).Find(&documents)
		update.Documents = documents
	}

	return complianceUpdates, nil

}

func (repo *ComplianceStatusUpdateRepository) GetByID(id uint) (*models.ComplianceStatusUpdate, error) {
	var complianceStatusUpdate *models.ComplianceStatusUpdate
	result := repo.db.First(&complianceStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceStatusUpdate, nil
}

func (repo *ComplianceStatusUpdateRepository) Create(ComplianceStatusUpdate *models.ComplianceStatusUpdate) (uint, error) {

	if err := repo.db.Create(ComplianceStatusUpdate).Error; err != nil {
		return 0, err
	}

	return ComplianceStatusUpdate.ID, nil
}

func (repo *ComplianceStatusUpdateRepository) Update(ComplianceStatusUpdate *models.ComplianceStatusUpdate) error {
	if err := repo.db.Save(ComplianceStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceStatusUpdateRepository) CreateDocuments(docs []*models.ComplianceStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceStatusUpdateRepository) GetDocumentByID(id uint) (*models.ComplianceStatusUpdateDocument, error) {
	var document models.ComplianceStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ComplianceStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ComplianceStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}
