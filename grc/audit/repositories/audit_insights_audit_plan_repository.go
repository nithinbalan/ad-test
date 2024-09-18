package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditInsightsAuditPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *auditDTOs.AuditPlanListQueryParamsDTO, currentUserID uint, currentUserResponsibleUnitID uint) (*auditDTOs.AuditPlanListResponseDTO, error)
	GetByID(id uint) (*models.AuditPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlan, error)
}

type AuditInsightsAuditPlanRepository struct {
	db *gorm.DB
}

func NewAuditInsightsAuditPlanRepository(db *gorm.DB) *AuditInsightsAuditPlanRepository {
	return &AuditInsightsAuditPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditInsightsAuditPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditInsightsAuditPlanRepository) GetAll(params *auditDTOs.AuditPlanListQueryParamsDTO, currentUserID uint, currentUserResponsibleUnitID uint) (*auditDTOs.AuditPlanListResponseDTO, error) {
	var auditPlans []*auditDTOs.AuditPlanListDTO
	var totalCount int64

	query := repo.db.Table("audit_plans").
		Select(`
			audit_plans.*,
			audit_categories.title_en AS audit_category_title_en,
			audit_categories.title_ar AS audit_category_title_ar,
			CONVERT(VARCHAR, audit_plans.start_date, 23) AS start_date,
			CONVERT(VARCHAR, audit_plans.target_date, 23) AS target_date,
			audit_statuses.title_en AS audit_status_title_en,
			audit_statuses.title_ar AS audit_status_title_ar,
			audit_statuses.type AS audit_status_type,
			audit_statuses.color_code AS audit_status_color_code,
			team_lead.id AS team_lead_id,
			team_lead.first_name_en AS team_lead_first_name_en,
			team_lead.first_name_ar AS team_lead_first_name_ar,
			team_lead.last_name_en AS team_lead_last_name_en,
			team_lead.last_name_ar AS team_lead_last_name_ar,
			team_lead.email AS team_lead_email,
			team_lead.initial AS team_lead_initial,
			team_lead.profile_image_path AS team_lead_profile_image_path,
			audit_teams.title_en AS audit_team_title_en,
			audit_teams.title_en AS audit_team_title_en,
			audit_teams_team_lead.id AS audit_team_team_lead_id,
			audit_teams_team_lead.first_name_en AS audit_team_team_lead_first_name_en,
			audit_teams_team_lead.first_name_ar AS audit_team_team_lead_first_name_ar,
			audit_teams_team_lead.last_name_en AS audit_team_team_lead_last_name_en,
			audit_teams_team_lead.last_name_ar AS audit_team_team_lead_last_name_ar,
			audit_teams_team_lead.email AS audit_team_team_lead_email,
			audit_teams_team_lead.initial AS audit_team_team_lead_initial,
			audit_teams_team_lead.color AS audit_team_team_lead_color,
			audit_teams_team_lead.profile_image_path AS audit_team_team_lead_profile_image_path,
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
			LEFT JOIN audit_categories ON audit_categories.id = audit_plans.audit_category_id
			LEFT JOIN audit_teams ON audit_teams.id = audit_plans.audit_team_id
        	LEFT JOIN users AS audit_teams_team_lead ON audit_teams_team_lead.id = audit_teams.team_lead_id
			LEFT JOIN audit_statuses ON audit_statuses.id = audit_plans.audit_status_id
			LEFT JOIN users AS team_lead ON team_lead.id = audit_plans.team_lead_id
			LEFT JOIN users AS created_user ON created_user.id = audit_plans.created_by_id
        	LEFT JOIN users AS updated_user ON updated_user.id = audit_plans.updated_by_id
			LEFT JOIN audit_plan_auditable_items ON audit_plan_auditable_items.audit_plan_id = audit_plans.id
			LEFT JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
			LEFT JOIN audit_plan_audit_members ON audit_plan_audit_members.audit_plan_id = audit_plans.id
	
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply other filters as needed (e.g., date filters)
	if params.StartDate != "" {
		query = query.Where("start_date >= ?", params.StartDate)
	}
	if params.TargetDate != "" {
		query = query.Where("target_date <= ?", params.TargetDate)
	}
	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "pending":
			query = query.Where("audit_statuses.type = ?", "pending")
		case "completed":
			query = query.Where("audit_statuses.type = ?", "completed")
		case "in_progress":
			query = query.Where("audit_statuses.type = ?", "in_progress")
		}
	}
	// **Apply user-specific filtering**:
	query = query.Where(`
		team_lead.id = ? OR
		audit_plan_audit_members.user_id = ? OR
		auditable_items.responsible_unit_id = ?`,
		currentUserID, currentUserID, currentUserResponsibleUnitID)

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
			audit_teams.title_en LIKE ? OR
			audit_teams.title_ar LIKE ? OR
			audit_statuses.title_en LIKE ? OR
			audit_statuses.title_ar LIKE ? OR
			team_lead.first_name_en LIKE ? OR
			team_lead.first_name_ar LIKE ? OR
			team_lead.last_name_en LIKE ? OR
			team_lead.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditPlanListResponseDTO{
		Items: auditPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditInsightsAuditPlanRepository) GetByID(id uint) (*models.AuditPlan, error) {
	var auditPlan *models.AuditPlan

	result := repo.db.First(&auditPlan, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlan, nil
}

func (repo *AuditInsightsAuditPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditPlan *models.AuditPlan

	result := db.
		Preload("AuditCategory").
		Preload("AuditTeam").
		Preload("AuditTeam.Users").
		Preload("AuditTeam.Users.User").
		Preload("AuditTeam.TeamLead").
		Preload("AuditStatus").
		Preload("TeamLead").
		Preload("Members").
		Preload("AuditInsightsAuditPlanChecklists").
		Preload("Members.User").
		Preload("AuditableItems").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlan, nil
}