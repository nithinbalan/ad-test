package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceRequirementStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(complianceRequirementID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceRequirementStatusUpdateListDTO, error)
	GetByID(id uint) (*models.ComplianceRequirementStatusUpdate, error)
	Create(ComplianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(ComplianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate, tx ...utility.DatabaseTransaction) error
	//Document
	CreateDocuments(docs []*models.ComplianceRequirementStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ComplianceRequirementStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	GetLatestComplianceRequirementStatusUpdate(complianceRequirementID uint) (*complianceDTOs.ComplianceRequirementStatusUpdateListDTO, error)
}

type ComplianceRequirementStatusUpdateRepository struct {
	db *gorm.DB
}

func NewComplianceRequirementStatusUpdateRepository(db *gorm.DB) *ComplianceRequirementStatusUpdateRepository {
	return &ComplianceRequirementStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceRequirementStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceRequirementStatusUpdateRepository) GetAll(complianceRequirementID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceRequirementStatusUpdateListDTO, error) {

	var complianceRequirementStatusUpdates []*complianceDTOs.ComplianceRequirementStatusUpdateListDTO

	query := repo.db.Table("compliance_requirement_status_updates").
		Select(`
            compliance_requirement_status_updates.*,
            compliance_requirements.title AS compliance_requirement_title,
            compliance_statuses.title_en AS compliance_status_title_en,
            compliance_statuses.title_ar AS compliance_status_title_ar,
            compliance_statuses.color_code AS compliance_status_color_code,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,   
			updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
            updated_user.color AS updated_user_color,
            updated_user.email AS updated_user_email,
            updated_user.profile_image_path AS updated_user_profile_image_path 
        `).
		Joins(`
            LEFT JOIN compliance_requirements ON compliance_requirements.id = compliance_requirement_status_updates.compliance_requirement_id
            LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_requirement_status_updates.compliance_status_id
            LEFT JOIN users AS created_user ON created_user.id = compliance_requirement_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = compliance_requirement_status_updates.updated_by_id
        `).
		Where("compliance_requirement_status_updates.compliance_requirement_id = ?", complianceRequirementID)

	if keyword != "" {
		query = query.Where(`
            (compliance_requirements.title LIKE ? OR
            compliance_requirement_status_updates.comment LIKE ? OR
            compliance_statuses.title_en LIKE ? OR
            compliance_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "compliance_requirement_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&complianceRequirementStatusUpdates)

	return complianceRequirementStatusUpdates, result.Error
}

func (repo *ComplianceRequirementStatusUpdateRepository) GetByID(id uint) (*models.ComplianceRequirementStatusUpdate, error) {
	var complianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate
	result := repo.db.First(&complianceRequirementStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceRequirementStatusUpdate, nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) Create(ComplianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(ComplianceRequirementStatusUpdate).Error; err != nil {
		return 0, err
	}

	return ComplianceRequirementStatusUpdate.ID, nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) Update(ComplianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(ComplianceRequirementStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) CreateDocuments(docs []*models.ComplianceRequirementStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *ComplianceRequirementStatusUpdateRepository) GetDocumentByID(id uint) (*models.ComplianceRequirementStatusUpdateDocument, error) {
	var document models.ComplianceRequirementStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ComplianceRequirementStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) GetLatestComplianceRequirementStatusUpdate(complianceRequirementID uint) (*complianceDTOs.ComplianceRequirementStatusUpdateListDTO, error) {

	var complianceRequirementStatusUpdate *complianceDTOs.ComplianceRequirementStatusUpdateListDTO

	query := repo.db.Table("compliance_requirement_status_updates").
		Select(`
            compliance_requirement_status_updates.*,
            compliance_requirements.title AS compliance_requirement_title,
            compliance_statuses.title_en AS compliance_status_title_en,
            compliance_statuses.title_ar AS compliance_status_title_ar,
            compliance_statuses.color_code AS compliance_status_color_code,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,   
			updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
            updated_user.color AS updated_user_color,
            updated_user.email AS updated_user_email,
            updated_user.profile_image_path AS updated_user_profile_image_path 
        `).
		Joins(`
            LEFT JOIN compliance_requirements ON compliance_requirements.id = compliance_requirement_status_updates.compliance_requirement_id
            LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_requirement_status_updates.compliance_status_id
            LEFT JOIN users AS created_user ON created_user.id = compliance_requirement_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = compliance_requirement_status_updates.updated_by_id
        `).
		Where("compliance_requirement_status_updates.compliance_requirement_id = ?", complianceRequirementID).
		Order("compliance_requirement_status_updates.created_at DESC").
		First(&complianceRequirementStatusUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return complianceRequirementStatusUpdate, nil
}
