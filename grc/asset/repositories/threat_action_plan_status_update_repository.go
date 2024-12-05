package repositories

import (
	threatDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThreatActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(threatActionPlanID uint, sort string, sortBy string, keyword string) ([]*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.ThreatActionPlanStatusUpdate, error)
	Create(threatActionPlanStatusUpdate *models.ThreatActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(threatActionPlanStatusUpdate *models.ThreatActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error
	GetLatestThreatActionPlanStatusUpdate(threatActionPlanID uint) (*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error)

	// Document operations
	CreateDocuments(docs []*models.ThreatActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ThreatActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type ThreatActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewThreatActionPlanStatusUpdateRepository(db *gorm.DB) *ThreatActionPlanStatusUpdateRepository {
	return &ThreatActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ThreatActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThreatActionPlanStatusUpdateRepository) GetAll(threatActionPlanID uint, sort string, sortBy string, keyword string) ([]*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error) {
	var threatActionPlanUpdates []*threatDTOs.ThreatActionPlanStatusUpdateListDTO

	query := repo.db.Table("threat_action_plan_status_updates").
		Select(`
            threat_action_plan_status_updates.*,
            threat_action_plans.title AS threat_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
            created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
            updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar
        `).
		Joins(`
            LEFT JOIN threat_action_plans ON threat_action_plans.id = threat_action_plan_status_updates.threat_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = threat_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = threat_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = threat_action_plan_status_updates.updated_by_id
        `).
		Where("threat_action_plan_status_updates.threat_action_plan_id = ?", threatActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (threat_action_plans.title LIKE ? OR
            threat_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,

			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "threat_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&threatActionPlanUpdates)

	for _, statusUpdate := range threatActionPlanUpdates {
		var documents []*threatDTOs.ThreatActionPlanStatusUpdateDocumentDTO
		repo.db.Model(&models.ThreatActionPlanStatusUpdateDocument{}).Where("threat_action_plan_status_update_id = ?", statusUpdate.ID).Find(&documents)
		statusUpdate.Documents = documents
	}

	return threatActionPlanUpdates, result.Error
}

func (repo *ThreatActionPlanStatusUpdateRepository) GetByID(id uint) (*models.ThreatActionPlanStatusUpdate, error) {
	var threatActionPlanStatusUpdate *models.ThreatActionPlanStatusUpdate
	result := repo.db.Preload("Documents").First(&threatActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return threatActionPlanStatusUpdate, nil
}

func (repo *ThreatActionPlanStatusUpdateRepository) Create(threatActionPlanStatusUpdate *models.ThreatActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(threatActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return threatActionPlanStatusUpdate.ID, nil
}

func (repo *ThreatActionPlanStatusUpdateRepository) Update(threatActionPlanStatusUpdate *models.ThreatActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(threatActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThreatActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.ThreatActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThreatActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.ThreatActionPlanStatusUpdateDocument, error) {
	var document models.ThreatActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ThreatActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", docIds).Delete(&models.ThreatActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThreatActionPlanStatusUpdateRepository) GetLatestThreatActionPlanStatusUpdate(threatActionPlanID uint) (*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error) {
	var threatActionPlanUpdate *threatDTOs.ThreatActionPlanStatusUpdateListDTO

	query := repo.db.Table("threat_action_plan_status_updates").
		Select(`
            threat_action_plan_status_updates.*,
            threat_action_plans.title AS threat_action_plan_title,
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
			updated_user.id AS updated_user_id,
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
            LEFT JOIN threat_action_plans ON threat_action_plans.id = threat_action_plan_status_updates.threat_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = threat_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = threat_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = threat_action_plan_status_updates.updated_by_id
        `).
		Where("threat_action_plan_status_updates.threat_action_plan_id = ?", threatActionPlanID).
		Order("threat_action_plan_status_updates.created_at DESC").
		First(&threatActionPlanUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return threatActionPlanUpdate, nil
}
