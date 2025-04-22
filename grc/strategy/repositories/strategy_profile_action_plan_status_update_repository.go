package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(recoveryActionPlanID uint, sort string, sortBy string, keyword string) ([]*strategyDTOs.StrategyProfileActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.StrategyProfileActionPlanStatusUpdate, error)
	Create(StrategyProfileActionPlanStatusUpdate *models.StrategyProfileActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(StrategyProfileActionPlanStatusUpdate *models.StrategyProfileActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error
	GetLatestStrategyProfileActionPlanStatusUpdate(recoveryActionPlanID uint) (*strategyDTOs.StrategyProfileActionPlanStatusUpdateListDTO, error)

	//Document
	CreateDocuments(docs []*models.StrategyProfileActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.StrategyProfileActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	ActionPlanUpdateStatus(actionPlan *models.StrategyProfileActionPlan, tx ...utility.DatabaseTransaction) error
}

type StrategyProfileActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewStrategyProfileActionPlanStatusUpdateRepository(db *gorm.DB) *StrategyProfileActionPlanStatusUpdateRepository {
	return &StrategyProfileActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StrategyProfileActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategyProfileActionPlanStatusUpdateRepository) GetAll(recoveryActionPlanID uint, sort string, sortBy string, keyword string) ([]*strategyDTOs.StrategyProfileActionPlanStatusUpdateListDTO, error) {

	var strategyProfileActionPlanUpdates []*strategyDTOs.StrategyProfileActionPlanStatusUpdateListDTO

	query := repo.db.Table("strategy_profile_action_plan_status_updates").
		Select(`
            strategy_profile_action_plan_status_updates.*,
            strategy_profile_action_plans.title AS strategy_profile_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
			action_plan_statuses.type AS action_plan_status_type,
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
            LEFT JOIN strategy_profile_action_plans ON strategy_profile_action_plans.id = strategy_profile_action_plan_status_updates.strategy_profile_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = strategy_profile_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = strategy_profile_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = strategy_profile_action_plan_status_updates.updated_by_id
        `).
		Where("strategy_profile_action_plan_status_updates.strategy_profile_action_plan_id = ?", recoveryActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (strategy_profile_action_plans.title LIKE ? OR
            strategy_profile_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "strategy_profile_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&strategyProfileActionPlanUpdates)

	for _, statusUpdate := range strategyProfileActionPlanUpdates {
		var documents []*strategyDTOs.StrategyProfileActionPlanStatusUpdateDocumentDTO

		repo.db.Model(&models.StrategyProfileActionPlanStatusUpdateDocument{}).Where("strategy_profile_action_plan_status_update_id = ?", statusUpdate.ID).Find(&documents)
		statusUpdate.Documents = documents
	}

	return strategyProfileActionPlanUpdates, result.Error
}

func (repo *StrategyProfileActionPlanStatusUpdateRepository) GetByID(id uint) (*models.StrategyProfileActionPlanStatusUpdate, error) {
	var strategyProfileActionPlanStatusUpdate *models.StrategyProfileActionPlanStatusUpdate
	result := repo.db.Preload("Documents").First(&strategyProfileActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileActionPlanStatusUpdate, nil
}

func (repo *StrategyProfileActionPlanStatusUpdateRepository) Create(StrategyProfileActionPlanStatusUpdate *models.StrategyProfileActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(StrategyProfileActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return StrategyProfileActionPlanStatusUpdate.ID, nil
}

func (repo *StrategyProfileActionPlanStatusUpdateRepository) Update(StrategyProfileActionPlanStatusUpdate *models.StrategyProfileActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(StrategyProfileActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StrategyProfileActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.StrategyProfileActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *StrategyProfileActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.StrategyProfileActionPlanStatusUpdateDocument, error) {
	var document models.StrategyProfileActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *StrategyProfileActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.StrategyProfileActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategyProfileActionPlanStatusUpdateRepository) GetLatestStrategyProfileActionPlanStatusUpdate(recoveryActionPlanID uint) (*strategyDTOs.StrategyProfileActionPlanStatusUpdateListDTO, error) {

	var strategyProfileActionPlanUpdate *strategyDTOs.StrategyProfileActionPlanStatusUpdateListDTO

	query := repo.db.Table("strategy_profile_action_plan_status_updates").
		Select(`
            strategy_profile_action_plan_status_updates.*,
            strategy_profile_action_plans.title AS strategy_profile_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
			action_plan_statuses.type AS action_plan_status_type,
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
            LEFT JOIN strategy_profile_action_plans ON strategy_profile_action_plans.id = strategy_profile_action_plan_status_updates.strategy_profile_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = strategy_profile_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = strategy_profile_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = strategy_profile_action_plan_status_updates.updated_by_id
        `).
		Where("strategy_profile_action_plan_status_updates.strategy_profile_action_plan_id = ?", recoveryActionPlanID).
		Order("strategy_profile_action_plan_status_updates.created_at DESC").
		First(&strategyProfileActionPlanUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return strategyProfileActionPlanUpdate, nil
}

func (repo *StrategyProfileActionPlanStatusUpdateRepository) ActionPlanUpdateStatus(actionPlan *models.StrategyProfileActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(actionPlan).Error; err != nil {
		return err
	}
	return nil
}
