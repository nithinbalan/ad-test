package repositories

import (
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ControlActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(controlActionPlanID uint, sort string, sortBy string, keyword string) ([]*organizationDTOs.ControlActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.ControlActionPlanStatusUpdate, error)
	Create(controlActionPlanStatusUpdate *models.ControlActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(controlActionPlanStatusUpdate *models.ControlActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.ControlActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ControlActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type ControlActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewControlActionPlanStatusUpdateRepository(db *gorm.DB) *ControlActionPlanStatusUpdateRepository {
	return &ControlActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ControlActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ControlActionPlanStatusUpdateRepository) GetAll(controlActionPlanID uint, sort string, sortBy string, keyword string) ([]*organizationDTOs.ControlActionPlanStatusUpdateListDTO, error) {

	var controlActionPlanUpdates []*organizationDTOs.ControlActionPlanStatusUpdateListDTO

	query := repo.db.Table("control_action_plan_status_updates").
		Select(`
            control_action_plan_status_updates.*,
            control_action_plans.title AS control_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
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
            LEFT JOIN control_action_plans ON control_action_plans.id = control_action_plan_status_updates.control_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = control_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = control_action_plan_status_updates.created_by
            LEFT JOIN users AS updated_user ON updated_user.id = control_action_plan_status_updates.updated_by
        `).
		Where("control_action_plan_status_updates.control_action_plan_id = ?", controlActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (control_action_plans.title LIKE ? OR
            control_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "control_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)

	result := query.Find(&controlActionPlanUpdates)
	return controlActionPlanUpdates, result.Error
}

func (repo *ControlActionPlanStatusUpdateRepository) GetByID(id uint) (*models.ControlActionPlanStatusUpdate, error) {
	var controlActionPlanStatusUpdate *models.ControlActionPlanStatusUpdate
	result := repo.db.Preload("Documents").First(&controlActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlActionPlanStatusUpdate, nil
}

func (repo *ControlActionPlanStatusUpdateRepository) Create(controlActionPlanStatusUpdate *models.ControlActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(controlActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return controlActionPlanStatusUpdate.ID, nil
}

func (repo *ControlActionPlanStatusUpdateRepository) Update(controlActionPlanStatusUpdate *models.ControlActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(controlActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ControlActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.ControlActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *ControlActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.ControlActionPlanStatusUpdateDocument, error) {
	var document models.ControlActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ControlActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ControlActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}
