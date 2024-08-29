package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditFindingActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditFindingActionPlanID uint, sort string, sortBy string, keyword string) ([]*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.AuditFindingActionPlanStatusUpdate, error)
	Create(auditFindingActionPlanStatusUpdate *models.AuditFindingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditFindingActionPlanStatusUpdate *models.AuditFindingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error
	GetLatestAuditFindingActionPlanStatusUpdate(auditFindingActionPlanID uint) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)

	// Document methods
	CreateDocuments(docs []*models.AuditFindingActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.AuditFindingActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type AuditFindingActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewAuditFindingActionPlanStatusUpdateRepository(db *gorm.DB) *AuditFindingActionPlanStatusUpdateRepository {
	return &AuditFindingActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditFindingActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditFindingActionPlanStatusUpdateRepository) GetAll(auditFindingActionPlanID uint, sort string, sortBy string, keyword string) ([]*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {

	var auditFindingActionPlanUpdates []*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO

	query := repo.db.Table("audit_finding_action_plan_status_updates").
		Select(`
            audit_finding_action_plan_status_updates.*,
            audit_finding_action_plans.title AS audit_finding_action_plan_title,
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
            LEFT JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_status_updates.audit_finding_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = audit_finding_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = audit_finding_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = audit_finding_action_plan_status_updates.updated_by_id
        `).
		Where("audit_finding_action_plan_status_updates.audit_finding_action_plan_id = ?", auditFindingActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (audit_finding_action_plans.title LIKE ? OR
            audit_finding_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "audit_finding_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&auditFindingActionPlanUpdates)

	return auditFindingActionPlanUpdates, result.Error
}

func (repo *AuditFindingActionPlanStatusUpdateRepository) GetByID(id uint) (*models.AuditFindingActionPlanStatusUpdate, error) {
	var auditFindingActionPlanStatusUpdate *models.AuditFindingActionPlanStatusUpdate
	result := repo.db.First(&auditFindingActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditFindingActionPlanStatusUpdate, nil
}

func (repo *AuditFindingActionPlanStatusUpdateRepository) Create(auditFindingActionPlanStatusUpdate *models.AuditFindingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditFindingActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return auditFindingActionPlanStatusUpdate.ID, nil
}

func (repo *AuditFindingActionPlanStatusUpdateRepository) Update(auditFindingActionPlanStatusUpdate *models.AuditFindingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditFindingActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AuditFindingActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.AuditFindingActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *AuditFindingActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.AuditFindingActionPlanStatusUpdateDocument, error) {
	var document models.AuditFindingActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *AuditFindingActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.AuditFindingActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFindingActionPlanStatusUpdateRepository) GetLatestAuditFindingActionPlanStatusUpdate(audit_findingActionPlanID uint) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {

	var audit_findingActionPlanUpdate *auditDTOs.AuditFindingActionPlanStatusUpdateListDTO

	query := repo.db.Table("audit_finding_action_plan_status_updates").
		Select(`
            audit_finding_action_plan_status_updates.*,
            audit_finding_action_plans.title AS audit_finding_action_plan_title,
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
            LEFT JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_status_updates.audit_finding_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = audit_finding_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = audit_finding_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = audit_finding_action_plan_status_updates.updated_by_id
        `).
		Where("audit_finding_action_plan_status_updates.audit_finding_action_plan_id = ?", audit_findingActionPlanID).
		Order("audit_finding_action_plan_status_updates.created_at DESC").
		First(&audit_findingActionPlanUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return audit_findingActionPlanUpdate, nil

}

// func (repo *AuditFindingActionPlanStatusUpdateRepository) GetLatestAuditFindingActionPlanStatusUpdate(auditFindingActionPlanID uint) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {

// 	var auditFindingActionPlanUpdate *auditDTOs.AuditFindingActionPlanStatusUpdateListDTO

// 	query := repo.db.Table("audit_finding_action_plan_status_updates").
// 		Select(`
//             audit_finding_action_plan_status_updates.*,
//             audit_finding_action_plans.title AS audit_finding_action_plan_title,
//             action_plan_statuses.title_en AS action_plan_status_title_en,
//             action_plan_statuses.title_ar AS action_plan_status_title_ar,
//             action_plan_statuses.color_code AS action_plan_status_color_code,
//         	created_user.id AS created_user_id,
//         	created_user.first_name_en AS created_user_first_name_en,
//         	created_user.first_name_ar AS created_user_first_name_ar,
//         	created_user.last_name_en AS created_user_last_name_en,
//         	created_user.last_name_ar AS created_user_last_name_ar,
//         	created_user.initial AS created_user_initial,
//         	created_user.color AS created_user_color,
//         	created_user.email AS created_user_email,
//         	created_user.profile_image_path AS created_user_profile_image_path,
// 			updated_user.id AS updated_user_id,
//     		updated_user.first_name_en AS updated_user_first_name_en,
//         	updated_user.first_name_ar AS updated_user_first_name_ar,
//         	updated_user.last_name_en AS updated_user.last_name_en,
//         	updated_user.last_name_ar AS updated_user.last_name_ar,
// 			updated_user.initial AS updated_user_initial,
//         	updated_user.color AS updated_user.color,
//         	updated_user.email AS updated_user.email,
//         	updated_user.profile_image_path AS updated_user.profile_image_path
//         `).
// 		Joins(`
//             LEFT JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_status_updates.audit_finding_action_plan_id
//             LEFT JOIN action_plan_statuses ON action_plan_statuses.id = audit_finding_action_plan_status_updates.action_plan_status_id
//             LEFT JOIN users AS created_user ON created_user.id = audit_finding_action_plan_status_updates.created_by_id
//             LEFT JOIN users AS updated_user ON updated_user.id = audit_finding_action_plan_status_updates.updated_by_id
//         `).
// 		Where("audit_finding_action_plan_status_updates.audit_finding_action_plan_id = ?", auditFindingActionPlanID).
// 		Order("audit_finding_action_plan_status_updates.created_at DESC").
// 		First(&auditFindingActionPlanUpdate)

// 	if query.Error != nil {
// 		return nil, query.Error
// 	}

// 	return auditFindingActionPlanUpdate, nil

// }
