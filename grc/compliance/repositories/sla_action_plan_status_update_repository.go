package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type SLAActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(slaActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.SLAActionPlanStatusUpdate, error)
	Create(SLAActionPlanStatusUpdate *models.SLAActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(SLAActionPlanStatusUpdate *models.SLAActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error
	GetLatestSLAActionPlanStatusUpdate(slaActionPlanID uint) (*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error)

	//Document
	CreateDocuments(docs []*models.SLAActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.SLAActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type SLAActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewSLAActionPlanStatusUpdateRepository(db *gorm.DB) *SLAActionPlanStatusUpdateRepository {
	return &SLAActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *SLAActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *SLAActionPlanStatusUpdateRepository) GetAll(slaActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error) {

	var slaActionPlanUpdates []*complianceDTOs.SLAActionPlanStatusUpdateListDTO

	query := repo.db.Table("sla_action_plan_status_updates").
		Select(`
            sla_action_plan_status_updates.*,
            sla_action_plans.title AS sla_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
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
            LEFT JOIN sla_action_plans ON sla_action_plans.id = sla_action_plan_status_updates.sla_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = sla_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = sla_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = sla_action_plan_status_updates.updated_by_id
        `).
		Where("sla_action_plan_status_updates.sla_action_plan_id = ?", slaActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (sla_action_plans.title LIKE ? OR
            sla_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "sla_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&slaActionPlanUpdates)

	for _, statusUpdate := range slaActionPlanUpdates {
		var documents []*complianceDTOs.SLAActionPlanStatusUpdateDocumentDTO

		repo.db.Model(&models.SLAActionPlanStatusUpdateDocument{}).Where("sla_action_plan_status_update_id = ?", statusUpdate.ID).Find(&documents)
		statusUpdate.Documents = documents
	}

	return slaActionPlanUpdates, result.Error
}

func (repo *SLAActionPlanStatusUpdateRepository) GetByID(id uint) (*models.SLAActionPlanStatusUpdate, error) {
	var slaActionPlanStatusUpdate *models.SLAActionPlanStatusUpdate
	result := repo.db.First(&slaActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return slaActionPlanStatusUpdate, nil
}

func (repo *SLAActionPlanStatusUpdateRepository) Create(SLAActionPlanStatusUpdate *models.SLAActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(SLAActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return SLAActionPlanStatusUpdate.ID, nil
}

func (repo *SLAActionPlanStatusUpdateRepository) Update(SLAActionPlanStatusUpdate *models.SLAActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(SLAActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SLAActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.SLAActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *SLAActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.SLAActionPlanStatusUpdateDocument, error) {
	var document models.SLAActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *SLAActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.SLAActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLAActionPlanStatusUpdateRepository) GetLatestSLAActionPlanStatusUpdate(slaActionPlanID uint) (*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error) {

	var slaActionPlanUpdate *complianceDTOs.SLAActionPlanStatusUpdateListDTO

	query := repo.db.Table("sla_action_plan_status_updates").
		Select(`
            sla_action_plan_status_updates.*,
            sla_action_plans.title AS sla_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
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
            LEFT JOIN sla_action_plans ON sla_action_plans.id = sla_action_plan_status_updates.sla_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = sla_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = sla_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = sla_action_plan_status_updates.updated_by_id
        `).
		Where("sla_action_plan_status_updates.sla_action_plan_id = ?", slaActionPlanID).
		Order("sla_action_plan_status_updates.created_at DESC").
		First(&slaActionPlanUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return slaActionPlanUpdate, nil
}
