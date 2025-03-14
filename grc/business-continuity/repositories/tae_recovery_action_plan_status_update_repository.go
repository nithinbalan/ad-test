package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type TaeRecoveryActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(recoveryActionPlanID uint, sort string, sortBy string, keyword string) ([]*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.TaeRecoveryActionPlanStatusUpdate, error)
	Create(TaeRecoveryActionPlanStatusUpdate *models.TaeRecoveryActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(TaeRecoveryActionPlanStatusUpdate *models.TaeRecoveryActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error
	GetLatestTaeRecoveryActionPlanStatusUpdate(recoveryActionPlanID uint) (*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error)

	//Document
	CreateDocuments(docs []*models.TaeRecoveryActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.TaeRecoveryActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	RecoveryActionPlanUpdateStatus(actionPlan *models.TestAndExerciseRecoveryActionPlan, tx ...utility.DatabaseTransaction) error
}

type TaeRecoveryActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewTaeRecoveryActionPlanStatusUpdateRepository(db *gorm.DB) *TaeRecoveryActionPlanStatusUpdateRepository {
	return &TaeRecoveryActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *TaeRecoveryActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *TaeRecoveryActionPlanStatusUpdateRepository) GetAll(recoveryActionPlanID uint, sort string, sortBy string, keyword string) ([]*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error) {

	var taeRecoveryActionPlanUpdates []*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO

	query := repo.db.Table("tae_recovery_action_plan_status_updates").
		Select(`
            tae_recovery_action_plan_status_updates.*,
            test_and_exercise_recovery_action_plans.title AS test_and_exercise_recovery_action_plan_title,
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
            LEFT JOIN test_and_exercise_recovery_action_plans ON test_and_exercise_recovery_action_plans.id = tae_recovery_action_plan_status_updates.test_and_exercise_recovery_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = tae_recovery_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = tae_recovery_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = tae_recovery_action_plan_status_updates.updated_by_id
        `).
		Where("tae_recovery_action_plan_status_updates.test_and_exercise_recovery_action_plan_id = ?", recoveryActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (test_and_exercise_recovery_action_plans.title LIKE ? OR
            tae_recovery_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "tae_recovery_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&taeRecoveryActionPlanUpdates)

	for _, statusUpdate := range taeRecoveryActionPlanUpdates {
		var documents []*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateDocumentDTO

		repo.db.Model(&models.TaeRecoveryActionPlanStatusUpdateDocument{}).Where("tae_recovery_action_plan_status_update_id = ?", statusUpdate.ID).Find(&documents)
		statusUpdate.Documents = documents
	}

	return taeRecoveryActionPlanUpdates, result.Error
}

func (repo *TaeRecoveryActionPlanStatusUpdateRepository) GetByID(id uint) (*models.TaeRecoveryActionPlanStatusUpdate, error) {
	var taeRecoveryActionPlanStatusUpdate *models.TaeRecoveryActionPlanStatusUpdate
	result := repo.db.Preload("Documents").First(&taeRecoveryActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return taeRecoveryActionPlanStatusUpdate, nil
}

func (repo *TaeRecoveryActionPlanStatusUpdateRepository) Create(TaeRecoveryActionPlanStatusUpdate *models.TaeRecoveryActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(TaeRecoveryActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return TaeRecoveryActionPlanStatusUpdate.ID, nil
}

func (repo *TaeRecoveryActionPlanStatusUpdateRepository) Update(TaeRecoveryActionPlanStatusUpdate *models.TaeRecoveryActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(TaeRecoveryActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TaeRecoveryActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.TaeRecoveryActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *TaeRecoveryActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.TaeRecoveryActionPlanStatusUpdateDocument, error) {
	var document models.TaeRecoveryActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *TaeRecoveryActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.TaeRecoveryActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TaeRecoveryActionPlanStatusUpdateRepository) GetLatestTaeRecoveryActionPlanStatusUpdate(recoveryActionPlanID uint) (*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error) {

	var taeRecoveryActionPlanUpdate *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO

	query := repo.db.Table("tae_recovery_action_plan_status_updates").
		Select(`
            tae_recovery_action_plan_status_updates.*,
            test_and_exercise_recovery_action_plans.title AS test_and_exercise_recovery_action_plan_title,
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
            LEFT JOIN test_and_exercise_recovery_action_plans ON test_and_exercise_recovery_action_plans.id = tae_recovery_action_plan_status_updates.test_and_exercise_recovery_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = tae_recovery_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = tae_recovery_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = tae_recovery_action_plan_status_updates.updated_by_id
        `).
		Where("tae_recovery_action_plan_status_updates.test_and_exercise_recovery_action_plan_id = ?", recoveryActionPlanID).
		Order("tae_recovery_action_plan_status_updates.created_at DESC").
		First(&taeRecoveryActionPlanUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return taeRecoveryActionPlanUpdate, nil
}

func (repo *TaeRecoveryActionPlanStatusUpdateRepository) RecoveryActionPlanUpdateStatus(actionPlan *models.TestAndExerciseRecoveryActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(actionPlan).Error; err != nil {
		return err
	}
	return nil
}
