package repositories

import (
	"jethur/grc/models"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ViolationReportActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(violationReportActionPlanID uint, sort string, sortBy string, keyword string) ([]*violationReportDTOs.ViolationReportActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.ViolationReportActionPlanStatusUpdate, error)
	Create(ViolationReportActionPlanStatusUpdate *models.ViolationReportActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(ViolationReportActionPlanStatusUpdate *models.ViolationReportActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error
	GetLatestViolationReportActionPlanStatusUpdate(violationReportActionPlanID uint) (*violationReportDTOs.ViolationReportActionPlanStatusUpdateListDTO, error)

	// Document
	CreateDocuments(docs []*models.ViolationReportActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ViolationReportActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	ViolationReportActionPlanUpdateStatus(ViolationReportActionPlan *models.ViolationReportActionPlan, tx ...utility.DatabaseTransaction) error
}

type ViolationReportActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewViolationReportActionPlanStatusUpdateRepository(db *gorm.DB) *ViolationReportActionPlanStatusUpdateRepository {
	return &ViolationReportActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ViolationReportActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ViolationReportActionPlanStatusUpdateRepository) GetAll(violationReportActionPlanID uint, sort string, sortBy string, keyword string) ([]*violationReportDTOs.ViolationReportActionPlanStatusUpdateListDTO, error) {
	var updates []*violationReportDTOs.ViolationReportActionPlanStatusUpdateListDTO

	query := repo.db.Table("violation_report_action_plan_status_updates").
		Select(`
            violation_report_action_plan_status_updates.*,
            violation_report_action_plans.title AS violation_report_action_plan_title,
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
            LEFT JOIN violation_report_action_plans ON violation_report_action_plans.id = violation_report_action_plan_status_updates.violation_report_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = violation_report_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = violation_report_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = violation_report_action_plan_status_updates.updated_by_id
        `).
		Where("violation_report_action_plan_status_updates.violation_report_action_plan_id = ?", violationReportActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (violation_report_action_plans.title LIKE ? OR
            violation_report_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "violation_report_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&updates)

	// Fetch documents for each status update
	for _, update := range updates {
		var documents []*models.ViolationReportActionPlanStatusUpdateDocument
		repo.db.Model(&models.ViolationReportActionPlanStatusUpdateDocument{}).
			Where("violation_report_action_plan_status_update_id = ?", update.ID).
			Find(&documents)
		update.Documents = documents
	}

	return updates, result.Error
}

func (repo *ViolationReportActionPlanStatusUpdateRepository) GetByID(id uint) (*models.ViolationReportActionPlanStatusUpdate, error) {
	var statusUpdate *models.ViolationReportActionPlanStatusUpdate
	result := repo.db.Preload("Documents").First(&statusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return statusUpdate, nil
}

func (repo *ViolationReportActionPlanStatusUpdateRepository) Create(statusUpdate *models.ViolationReportActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(statusUpdate).Error; err != nil {
		return 0, err
	}

	return statusUpdate.ID, nil
}

func (repo *ViolationReportActionPlanStatusUpdateRepository) Update(statusUpdate *models.ViolationReportActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(statusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ViolationReportActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.ViolationReportActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ViolationReportActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.ViolationReportActionPlanStatusUpdateDocument, error) {
	var document models.ViolationReportActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ViolationReportActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).
		Delete(&models.ViolationReportActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ViolationReportActionPlanStatusUpdateRepository) GetLatestViolationReportActionPlanStatusUpdate(violationReportActionPlanID uint) (*violationReportDTOs.ViolationReportActionPlanStatusUpdateListDTO, error) {
	var update *violationReportDTOs.ViolationReportActionPlanStatusUpdateListDTO

	query := repo.db.Table("violation_report_action_plan_status_updates").
		Select(`
            violation_report_action_plan_status_updates.*,
            violation_report_action_plans.title AS violation_report_action_plan_title,
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
            LEFT JOIN violation_report_action_plans ON violation_report_action_plans.id = violation_report_action_plan_status_updates.violation_report_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = violation_report_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = violation_report_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = violation_report_action_plan_status_updates.updated_by_id
        `).
		Where("violation_report_action_plan_status_updates.violation_report_action_plan_id = ?", violationReportActionPlanID).
		Order("violation_report_action_plan_status_updates.created_at DESC").
		First(&update)

	if query.Error != nil {
		return nil, query.Error
	}

	return update, nil
}

func (repo *ViolationReportActionPlanStatusUpdateRepository) ViolationReportActionPlanUpdateStatus(plan *models.ViolationReportActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(plan).Error; err != nil {
		return err
	}
	return nil
}
