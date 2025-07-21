package repositories

import (
	"jethur/grc/models"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(projectActionPlanID uint, sort string, sortBy string, keyword string) ([]*projectDTOs.ProjectActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.ProjectActionPlanStatusUpdate, error)
	Create(ProjectActionPlanStatusUpdate *models.ProjectActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(ProjectActionPlanStatusUpdate *models.ProjectActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.ProjectActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ProjectActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	GetLatestProjectActionPlanStatusUpdate(projectActionPlanID uint) (*projectDTOs.ProjectActionPlanStatusUpdateListDTO, error)
	ProjectActionPlanUpdateStatus(projectActionPlan *models.ProjectActionPlan, tx ...utility.DatabaseTransaction) error
}

type ProjectActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewProjectActionPlanStatusUpdateRepository(db *gorm.DB) *ProjectActionPlanStatusUpdateRepository {
	return &ProjectActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ProjectActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProjectActionPlanStatusUpdateRepository) GetAll(projectActionPlanID uint, sort string, sortBy string, keyword string) ([]*projectDTOs.ProjectActionPlanStatusUpdateListDTO, error) {

	var projectActionPlanUpdates []*projectDTOs.ProjectActionPlanStatusUpdateListDTO

	query := repo.db.Table("project_action_plan_status_updates").
		Select(`
            project_action_plan_status_updates.*,
            project_action_plans.title AS project_action_plan_title,
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
            LEFT JOIN project_action_plans ON project_action_plans.id = project_action_plan_status_updates.project_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = project_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = project_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = project_action_plan_status_updates.updated_by_id
        `).
		Where("project_action_plan_status_updates.project_action_plan_id = ?", projectActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (project_action_plans.title LIKE ? OR
            project_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "project_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&projectActionPlanUpdates)

	for _, statusUpdate := range projectActionPlanUpdates {
		var documents []*projectDTOs.ProjectActionPlanStatusUpdateDocumentDTO

		repo.db.Model(&models.ProjectActionPlanStatusUpdateDocument{}).Where("project_action_plan_status_update_id = ?", statusUpdate.ID).Find(&documents)
		statusUpdate.Documents = documents
	}

	return projectActionPlanUpdates, result.Error
}

func (repo *ProjectActionPlanStatusUpdateRepository) GetByID(id uint) (*models.ProjectActionPlanStatusUpdate, error) {
	var projectActionPlanStatusUpdate *models.ProjectActionPlanStatusUpdate
	result := repo.db.Preload("Documents").First(&projectActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return projectActionPlanStatusUpdate, nil
}

func (repo *ProjectActionPlanStatusUpdateRepository) Create(ProjectActionPlanStatusUpdate *models.ProjectActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(ProjectActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return ProjectActionPlanStatusUpdate.ID, nil
}

func (repo *ProjectActionPlanStatusUpdateRepository) Update(ProjectActionPlanStatusUpdate *models.ProjectActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(ProjectActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProjectActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.ProjectActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *ProjectActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.ProjectActionPlanStatusUpdateDocument, error) {
	var document models.ProjectActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ProjectActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ProjectActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ProjectActionPlanStatusUpdateRepository) GetLatestProjectActionPlanStatusUpdate(projectActionPlanID uint) (*projectDTOs.ProjectActionPlanStatusUpdateListDTO, error) {

	var projectActionPlanUpdate *projectDTOs.ProjectActionPlanStatusUpdateListDTO

	query := repo.db.Table("project_action_plan_status_updates").
		Select(`
            project_action_plan_status_updates.*,
            project_action_plans.title AS project_action_plan_title,
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
            LEFT JOIN project_action_plans ON project_action_plans.id = project_action_plan_status_updates.project_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = project_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = project_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = project_action_plan_status_updates.updated_by_id
        `).
		Where("project_action_plan_status_updates.project_action_plan_id = ?", projectActionPlanID).
		Order("project_action_plan_status_updates.created_at DESC").
		First(&projectActionPlanUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	var documents []*projectDTOs.ProjectActionPlanStatusUpdateDocumentDTO

	repo.db.Model(&models.ProjectActionPlanStatusUpdateDocument{}).Where("project_action_plan_status_update_id = ?", projectActionPlanUpdate.ID).Find(&documents)
	projectActionPlanUpdate.Documents = documents

	return projectActionPlanUpdate, nil

}

func (repo *ProjectActionPlanStatusUpdateRepository) ProjectActionPlanUpdateStatus(ProjectActionPlan *models.ProjectActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(ProjectActionPlan).Error; err != nil {
		return err
	}
	return nil
}
