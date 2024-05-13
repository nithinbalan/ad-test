package repositories

import (
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(incidentActionPlanID uint, sort string, sortBy string, keyword string) ([]*incidentDTOs.IncidentActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.IncidentActionPlanStatusUpdate, error)
	Create(IncidentActionPlanStatusUpdate *models.IncidentActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(IncidentActionPlanStatusUpdate *models.IncidentActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.IncidentActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.IncidentActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type IncidentActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewIncidentActionPlanStatusUpdateRepository(db *gorm.DB) *IncidentActionPlanStatusUpdateRepository {
	return &IncidentActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *IncidentActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentActionPlanStatusUpdateRepository) GetAll(incidentActionPlanID uint, sort string, sortBy string, keyword string) ([]*incidentDTOs.IncidentActionPlanStatusUpdateListDTO, error) {

	var incidentActionPlanUpdates []*incidentDTOs.IncidentActionPlanStatusUpdateListDTO

	query := repo.db.Table("incident_action_plan_status_updates").
		Select(`
            incident_action_plan_status_updates.*,
            incident_action_plans.title AS incident_action_plan_title,
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
            LEFT JOIN incident_action_plans ON incident_action_plans.id = incident_action_plan_status_updates.incident_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = incident_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = incident_action_plan_status_updates.created_by
            LEFT JOIN users AS updated_user ON updated_user.id = incident_action_plan_status_updates.updated_by
        `).
		Where("incident_action_plan_status_updates.incident_action_plan_id = ?", incidentActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (incident_action_plans.title LIKE ? OR
            incident_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "incident_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&incidentActionPlanUpdates)

	return incidentActionPlanUpdates, result.Error
}

func (repo *IncidentActionPlanStatusUpdateRepository) GetByID(id uint) (*models.IncidentActionPlanStatusUpdate, error) {
	var incidentActionPlanStatusUpdate *models.IncidentActionPlanStatusUpdate
	result := repo.db.First(&incidentActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentActionPlanStatusUpdate, nil
}

func (repo *IncidentActionPlanStatusUpdateRepository) Create(IncidentActionPlanStatusUpdate *models.IncidentActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(IncidentActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return IncidentActionPlanStatusUpdate.ID, nil
}

func (repo *IncidentActionPlanStatusUpdateRepository) Update(IncidentActionPlanStatusUpdate *models.IncidentActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(IncidentActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.IncidentActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *IncidentActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.IncidentActionPlanStatusUpdateDocument, error) {
	var document models.IncidentActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *IncidentActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.IncidentActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}
