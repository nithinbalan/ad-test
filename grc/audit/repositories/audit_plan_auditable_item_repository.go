package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditPlanAuditableItemRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditPlanID uint, params *auditDTOs.AuditPlanAuditableItemListQueryParamsDTO) (*auditDTOs.AuditPlanAuditableItemListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanAuditableItem, error)
	Create(auditPlanAuditableItems []*models.AuditPlanAuditableItem, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanAuditableItem, error)
	GetByAuditPlanAndAuditableItemID(auditPlanID, auditableItemID uint, auditPlanAuditableItem *models.AuditPlanAuditableItem, tx ...utility.DatabaseTransaction) error
	GetByAuditableItemID(auditableItemID uint) (*models.AuditPlanAuditableItem, error)

}

type AuditPlanAuditableItemRepository struct {
	db *gorm.DB
}

func NewAuditPlanAuditableItemRepository(db *gorm.DB) *AuditPlanAuditableItemRepository {
	return &AuditPlanAuditableItemRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditPlanAuditableItemRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditPlanAuditableItemRepository) GetAll(auditPlanID uint, params *auditDTOs.AuditPlanAuditableItemListQueryParamsDTO) (*auditDTOs.AuditPlanAuditableItemListResponseDTO, error) {
	var auditPlanAuditableItems []*auditDTOs.AuditPlanAuditableItemListDTO
	var totalCount int64

	query := repo.db.Table("audit_plan_auditable_items").
		Select(`
			audit_plan_auditable_items.id,
			audit_plan_auditable_items.audit_plan_id,
			audit_plan_auditable_items.auditable_item_id,
			audit_plan_auditable_items.created_by_id,
			audit_plan_auditable_items.created_at,
			processes.title_en AS process_title_en,
			processes.title_ar AS process_title_ar,
			controls.title_en AS control_title_en,
			controls.title_ar AS control_title_ar,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			incident_register.title AS incident_register_title,
			compliance_requirements.title AS compliance_requirement_title,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			risk_register.title AS risk_register_title,
			users.id AS user_id,
			users.first_name_en AS user_first_name_en,
			users.first_name_ar AS user_first_name_ar,
			users.last_name_en AS user_last_name_en,
			users.last_name_ar AS user_last_name_ar,
			users.email AS user_email,
			users.initial AS user_initial,
			users.color AS user_color,
			users.profile_image_path AS user_profile_image_path,
			audit_plans.reference_number,
			audit_plans.title AS audit_plan_title,
			auditable_items.title_en AS auditable_item_title_en,
			auditable_items.title_ar AS auditable_item_title_ar,
			auditable_items.type AS auditable_item_type,
			auditable_items.comment AS auditable_item_comment,
			audit_categories.title_en AS audit_category_title_en,
			audit_categories.title_ar AS audit_category_title_ar,
			audit_teams.title_en AS audit_team_title_en,
			audit_teams.title_ar AS audit_team_title_ar,
			team_lead.id AS team_lead_id,
			team_lead.first_name_en AS team_lead_first_name_en,
			team_lead.first_name_ar AS team_lead_first_name_ar,
			team_lead.last_name_en AS team_lead_last_name_en,
			team_lead.last_name_ar AS team_lead_last_name_ar,
			team_lead.email AS team_lead_email,
			team_lead.profile_image_path AS team_lead_profile_image_path,
			audit_statuses.title_en AS audit_status_title_en,
			audit_statuses.title_ar AS audit_status_title_ar,
			audit_statuses.type AS audit_status_type,
			audit_statuses.color_code AS audit_status_color_code,
			auditor_risk_ratings.title_en AS auditor_risk_rating_title_en,
			auditor_risk_ratings.title_ar AS auditor_risk_rating_title_ar,
			auditor_risk_ratings.color_code AS auditor_risk_rating_color_code,
			auditor_risk_ratings.type AS auditor_risk_rating_type,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN audit_plans ON audit_plans.id = audit_plan_auditable_items.audit_plan_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
			LEFT JOIN processes ON processes.id = auditable_items.process_id
			LEFT JOIN controls ON controls.id = auditable_items.control_id
			LEFT JOIN documents ON documents.id = auditable_items.document_id
			LEFT JOIN incident_register ON incident_register.id = auditable_items.incident_register_id
			LEFT JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
			LEFT JOIN organization_units ON organization_units.id = auditable_items.responsible_unit_id
			LEFT JOIN risk_register ON risk_register.id = auditable_items.risk_register_id
			LEFT JOIN users ON users.id = auditable_items.user_id
			LEFT JOIN audit_categories ON audit_categories.id = audit_plans.audit_category_id
			LEFT JOIN audit_teams ON audit_teams.id = audit_plans.audit_team_id
			LEFT JOIN users AS team_lead ON team_lead.id = audit_plans.team_lead_id
			LEFT JOIN audit_statuses ON audit_statuses.id = audit_plans.audit_status_id
        	LEFT JOIN auditor_risk_ratings ON auditor_risk_ratings.id = auditable_items.auditor_risk_rating_id
			LEFT JOIN users AS created_user ON created_user.id = audit_plan_auditable_items.created_by_id
		`).Where("audit_plan_auditable_items.audit_plan_id = ?", auditPlanID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audit_plans.title LIKE ? OR
			audit_plans.description LIKE ? OR
			audit_plans.reference_number LIKE ? OR
			audit_categories.title_en LIKE ? OR
			audit_categories.title_ar LIKE ? OR
			auditable_items.title_en LIKE ? OR
			auditable_items.title_ar LIKE ? OR
			audit_teams.title_en LIKE ? OR
			audit_teams.title_ar LIKE ? OR
			audit_statuses.title_en LIKE ? OR
			audit_statuses.title_ar LIKE ? OR
			team_lead.first_name_en LIKE ? OR
			team_lead.first_name_ar LIKE ? OR
			team_lead.last_name_en LIKE ? OR
			team_lead.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audit_plans.title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditPlanAuditableItems).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditPlanAuditableItemListResponseDTO{
		Items: auditPlanAuditableItems,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditPlanAuditableItemRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanAuditableItem, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditPlan *models.AuditPlanAuditableItem

	result := db.First(&auditPlan, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlan, nil
}

func (repo *AuditPlanAuditableItemRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanAuditableItem, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var auditableItem models.AuditPlanAuditableItem

	// Load the auditable item with all related data
	result := db.
		Preload("AuditPlan").
		Preload("AuditableItem").
		Preload("AuditableItem.AuditStatus").
		Preload("AuditableItem.AuditorRiskRating").
		Preload("AuditableItem.Process").
		Preload("AuditableItem.Control").
		Preload("AuditableItem.Document").
		Preload("AuditableItem.IncidentRegister").
		Preload("AuditableItem.ComplianceRequirement").
		Preload("AuditableItem.RiskRegister").
		Preload("AuditableItem.User").
		Preload("AuditableItem.ResponsibleUnit").
		Preload("CreatedByUser").
		First(&auditableItem, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &auditableItem, nil
}

func (repo *AuditPlanAuditableItemRepository) Create(auditPlanAuditableItems []*models.AuditPlanAuditableItem, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&auditPlanAuditableItems).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanAuditableItemRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var auditPlanAuditableItem *models.AuditPlanAuditableItem
	if err := db.First(&auditPlanAuditableItem, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.AuditPlanAuditableItem{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *AuditPlanAuditableItemRepository) GetByAuditPlanAndAuditableItemID(auditPlanID, auditableItemID uint, auditPlanAuditableItem *models.AuditPlanAuditableItem, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("audit_plan_id = ? AND auditable_item_id = ?", auditPlanID, auditableItemID).First(auditPlanAuditableItem).Error
}


func (repo *AuditPlanAuditableItemRepository) GetByAuditableItemID(auditableItemID uint) (*models.AuditPlanAuditableItem, error) {
	var auditPlanAuditableItem *models.AuditPlanAuditableItem

	result := repo.db.Where("auditable_item_id = ?", auditableItemID).First(&auditPlanAuditableItem)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlanAuditableItem, nil
}
