package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RCMRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *auditDTOs.RCMListQueryParamsDTO, auditPlanID uint) ([]map[string]interface{}, int64, error)
}

type RCMRepository struct {
	db *gorm.DB
}

func NewRCMRepository(db *gorm.DB) *RCMRepository {
	return &RCMRepository{
		db: db,
	}
}

func (repo *RCMRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RCMRepository) GetAll(params *auditDTOs.RCMListQueryParamsDTO, auditPlanID uint) ([]map[string]interface{}, int64, error) {
	var rcmList []map[string]interface{}
	var totalCount int64

	query := repo.db.Table("audit_plan_auditable_item_risks").
		Select(`
			audit_plan_auditable_item_risks.id AS id,
			audit_plan_auditable_item_risks.audit_plan_auditable_item_id AS audit_plan_auditable_item_id,
			audit_plan_auditable_item_risks.risk_register_id AS risk_register_id,
			risk_register.title AS risk_title,
			audit_plan_auditable_item_risks.created_at AS created_at,
			audit_plan_auditable_item_risks.updated_at AS updated_at,
			audit_plan_auditable_item_risk_controls.controls_id AS control_id,
			controls.title_en AS control_title_en,
			controls.title_ar AS control_title_ar,
			control_implementation_statuses.title_en AS control_implementation_status_title_en,
			control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
			control_implementation_statuses.color_code AS control_implementation_status_color_code,
			control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
			control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
			control_efficiency_statuses.color_code AS control_efficiency_status_color_code
		`).
		Joins(`
		    LEFT JOIN audit_plan_auditable_items ON audit_plan_auditable_items.id = audit_plan_auditable_item_risks.audit_plan_auditable_item_id
			LEFT JOIN risk_register ON risk_register.id = audit_plan_auditable_item_risks.risk_register_id
			LEFT JOIN users AS created_user ON created_user.id = audit_plan_auditable_item_risks.created_by_id
			LEFT JOIN audit_plan_auditable_item_risk_controls ON audit_plan_auditable_item_risk_controls.audit_plan_auditable_item_risks_id = audit_plan_auditable_item_risks.id
			LEFT JOIN controls ON controls.id = audit_plan_auditable_item_risk_controls.controls_id
			LEFT JOIN control_implementation_statuses ON control_implementation_statuses.id = controls.control_implementation_status_id
		    LEFT JOIN control_efficiency_statuses ON control_efficiency_statuses.id = controls.control_efficiency_status_id
			LEFT JOIN users AS control_user ON control_user.id = audit_plan_auditable_item_risk_controls.created_by_id
		`).Where(`audit_plan_auditable_items.audit_plan_id = ?`, auditPlanID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(risk_registers.title_en LIKE ? OR
			risk_registers.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			controls.title_en LIKE ? OR
			controls.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audit_plan_auditable_item_risks.id" // Default sort by ID if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&rcmList).Error
	if err != nil {
		return nil, 0, err
	}

	return rcmList, totalCount, nil
}
