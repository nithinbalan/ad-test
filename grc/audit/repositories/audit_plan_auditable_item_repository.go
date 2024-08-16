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
	Delete(auditPlanID uint, auditableItemID uint, tx ...utility.DatabaseTransaction) error
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
			audit_plans.reference_number,
			audit_plans.title AS audit_plan_title,
			auditable_items.title_en AS auditable_item_title_en,
			auditable_items.title_ar AS auditable_item_title_ar,
			audit_categories.title_en AS audit_category_title_en,
			audit_categories.title_ar AS audit_category_title_ar,
			audit_teams.title_en AS audit_team_title_en,
			audit_teams.title_ar AS audit_team_title_ar,
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
			LEFT JOIN audit_categories ON audit_categories.id = audit_plans.audit_category_id
			LEFT JOIN audit_teams ON audit_teams.id = audit_plans.audit_team_id
			LEFT JOIN users AS team_lead ON team_lead.id = audit_plans.team_lead_id
			LEFT JOIN audit_statuses ON audit_statuses.id = audit_plans.audit_status_id
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

func (repo *AuditPlanAuditableItemRepository) Delete(auditPlanID uint, auditableItemID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete the records based on audit_plan_id and auditable_item_id
	if err := db.Where("audit_plan_id = ? AND auditable_item_id = ?", auditPlanID, auditableItemID).Delete(&models.AuditPlanAuditableItem{}).Error; err != nil {
		return err
	}

	return nil
}
