package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(complianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.ComplianceActionPlanStatusUpdate, error)
	Create(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.ComplianceActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ComplianceActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	GetLatestComplianceActionPlanStatusUpdate(complianceActionPlanID uint) (*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error)
}

type ComplianceActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewComplianceActionPlanStatusUpdateRepository(db *gorm.DB) *ComplianceActionPlanStatusUpdateRepository {
	return &ComplianceActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceActionPlanStatusUpdateRepository) GetAll(complianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error) {

	var complianceActionPlanUpdates []*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO

	query := repo.db.Table("compliance_action_plan_status_updates").
		Select(`
            compliance_action_plan_status_updates.*,
            compliance_action_plans.title AS compliance_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
        	created_user.id AS created_user_id,
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
            LEFT JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = compliance_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = compliance_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = compliance_action_plan_status_updates.updated_by_id
        `).
		Where("compliance_action_plan_status_updates.compliance_action_plan_id = ?", complianceActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (compliance_action_plans.title LIKE ? OR
            compliance_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "compliance_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&complianceActionPlanUpdates)

	for _, statusUpdate := range complianceActionPlanUpdates {
		var documents []*complianceDTOs.ComplianceActionPlanStatusUpdateDocumentDTO

		repo.db.Model(&models.ComplianceActionPlanStatusUpdateDocument{}).Where("compliance_action_plan_status_update_id = ?", statusUpdate.ID).Find(&documents)
		statusUpdate.Documents = documents
	}

	return complianceActionPlanUpdates, result.Error
}

func (repo *ComplianceActionPlanStatusUpdateRepository) GetByID(id uint) (*models.ComplianceActionPlanStatusUpdate, error) {
	var complianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate
	result := repo.db.First(&complianceActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceActionPlanStatusUpdate, nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) Create(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(ComplianceActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return ComplianceActionPlanStatusUpdate.ID, nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) Update(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(ComplianceActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.ComplianceActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *ComplianceActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.ComplianceActionPlanStatusUpdateDocument, error) {
	var document models.ComplianceActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ComplianceActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) GetLatestComplianceActionPlanStatusUpdate(complianceActionPlanID uint) (*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error) {

	var complianceActionPlanUpdate *complianceDTOs.ComplianceActionPlanStatusUpdateListDTO

	query := repo.db.Table("compliance_action_plan_status_updates").
		Select(`
            compliance_action_plan_status_updates.*,
            compliance_action_plans.title AS compliance_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
        	created_user.id AS created_user_id,
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
            LEFT JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = compliance_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = compliance_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = compliance_action_plan_status_updates.updated_by_id
        `).
		Where("compliance_action_plan_status_updates.compliance_action_plan_id = ?", complianceActionPlanID).
		Order("compliance_action_plan_status_updates.created_at DESC").
		First(&complianceActionPlanUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return complianceActionPlanUpdate, nil

}
