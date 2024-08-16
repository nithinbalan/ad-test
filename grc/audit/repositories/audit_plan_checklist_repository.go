package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditPlanChecklistRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditPlanId uint, params *auditDTOs.AuditPlanChecklistListQueryParamsDTO) ([]*auditDTOs.AuditPlanChecklistListDTO, error)
	GetByID(id uint) (*models.AuditPlanChecklist, error)
	Create(auditPlanChecklist *models.AuditPlanChecklist, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditPlanChecklist *models.AuditPlanChecklist, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AuditPlanChecklistRepository struct {
	db *gorm.DB
}

func NewAuditPlanChecklistRepository(db *gorm.DB) *AuditPlanChecklistRepository {
	return &AuditPlanChecklistRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditPlanChecklistRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditPlanChecklistRepository) GetAll(auditPlanId uint, params *auditDTOs.AuditPlanChecklistListQueryParamsDTO) ([]*auditDTOs.AuditPlanChecklistListDTO, error) {
	var auditPlanChecklists []*auditDTOs.AuditPlanChecklistListDTO

	query := repo.db.Table("audit_plan_checklists").
		Select(`
			audit_plan_checklists.id,
			audit_plan_checklists.audit_plan_id,
			audit_plan_checklists.auditable_item_id,
			auditable_items.title_en AS auditable_item_title_en,
			auditable_items.title_ar AS auditable_item_title_ar,
			auditable_items.process_id,
			auditable_items.control_id,
			auditable_items.document_id,
			auditable_items.incident_register_id,
			auditable_items.compliance_requirement_id,
			auditable_items.risk_register_id,
			auditable_items.user_id,
			audit_plan_checklists.question,
			audit_plan_checklists.created_at,
			audit_plan_checklists.updated_at,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
			LEFT JOIN users AS created_user ON created_user.id = audit_plan_checklists.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_plan_checklists.updated_by_id
		`).Where(`audit_plan_checklists.audit_plan_id`, auditPlanId)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audit_plan_checklists.question LIKE ? OR
			auditable_items.title_en LIKE ? OR
			auditable_items.title_ar LIKE ?)`,
			keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "created_at" // Default sort by created_at if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	err := query.Find(&auditPlanChecklists).Error
	if err != nil {
		return nil, err
	}

	return auditPlanChecklists, nil
}

func (repo *AuditPlanChecklistRepository) GetByID(id uint) (*models.AuditPlanChecklist, error) {
	var auditPlanChecklist *models.AuditPlanChecklist

	result := repo.db.First(&auditPlanChecklist, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlanChecklist, nil
}

func (repo *AuditPlanChecklistRepository) Create(auditPlanChecklist *models.AuditPlanChecklist, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditPlanChecklist).Error; err != nil {
		return 0, err
	}

	return auditPlanChecklist.ID, nil
}

func (repo *AuditPlanChecklistRepository) Update(auditPlanChecklist *models.AuditPlanChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditPlanChecklist).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanChecklistRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.AuditPlanChecklist{}, id).Error; err != nil {
		return err
	}

	return nil
}
