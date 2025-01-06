package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"time"

	"gorm.io/gorm"
)

type AuditInsightsFindingActionPlanRepositoryInterface interface {
	GetAll(params *auditDTOs.AuditFindingActionPlanListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditFindingActionPlanListResponseDTO, error)
	GetFullDetailsByID(id uint) (*models.AuditFindingActionPlan, error)
}

type AuditInsightsFindingActionPlanRepository struct {
	db *gorm.DB
}

func NewAuditInsightsFindingActionPlanRepository(db *gorm.DB) *AuditInsightsFindingActionPlanRepository {
	return &AuditInsightsFindingActionPlanRepository{
		db: db,
	}
}

func (repo *AuditInsightsFindingActionPlanRepository) GetAll(params *auditDTOs.AuditFindingActionPlanListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditFindingActionPlanListResponseDTO, error) {
	var auditFindingActionPlans []*auditDTOs.AuditFindingActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("audit_finding_action_plans").
		Select(`
        audit_finding_action_plans.id,
		audit_finding_action_plans.audit_finding_id,
		audit_finding_action_plans.action_plan_status_id,
		audit_finding_action_plans.severity_level_id,
		audit_finding_action_plans.responsible_user_id,
		audit_finding_action_plans.review_user_id,
		audit_finding_action_plans.created_by_id,
		audit_finding_action_plans.updated_by_id,
		audit_finding_action_plans.title,
		audit_finding_action_plans.description,
		audit_finding_action_plans.percentage,
		audit_finding_action_plans.comments,
		audit_finding_action_plans.target_date,
		audit_finding_action_plans.created_at,
		audit_finding_action_plans.updated_at,
		audit_findings.title AS audit_finding_title,
		audit_findings.reference_number AS audit_finding_reference_number,
		audit_findings.description AS audit_finding_description,
		CONVERT(VARCHAR, audit_finding_action_plans.target_date, 23) AS target_date,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
        severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
		responsible_user.id AS responsible_user_id,
        responsible_user.first_name_en AS responsible_user_first_name_en,
        responsible_user.first_name_ar AS responsible_user_first_name_ar,
        responsible_user.last_name_en AS responsible_user_last_name_en,
        responsible_user.last_name_ar AS responsible_user_last_name_ar,
		responsible_user.initial AS responsible_user_initial,
        responsible_user.color AS responsible_user_color,
        responsible_user.email AS responsible_user_email,
        responsible_user.profile_image_path AS responsible_user_profile_image_path, 
		review_user.id AS review_user_id,
    	review_user.first_name_en AS review_user_first_name_en,
		review_user.first_name_ar AS review_user_first_name_ar,
		review_user.last_name_en AS review_user_last_name_en,
		review_user.last_name_ar AS review_user_last_name_ar,
		review_user.initial AS review_user_initial,
        review_user.color AS review_user_color,
        review_user.email AS review_user_email,
        review_user.profile_image_path AS review_user_profile_image_path,
        created_user.first_name_en AS created_user_first_name_en,
        created_user.first_name_ar AS created_user_first_name_ar,
        created_user.last_name_en AS created_user_last_name_en,
        created_user.last_name_ar AS created_user_last_name_ar,  
		created_user.initial AS created_user_initial,
        created_user.color AS created_user_color,
        created_user.email AS created_user_email,
        created_user.profile_image_path AS created_user_profile_image_path,
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
        LEFT JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
        LEFT JOIN audits ON audits.id = audit_findings.audit_id
        LEFT JOIN audit_statuses ON audit_statuses.id = audits.audit_status_id
        LEFT JOIN audit_plans ON audit_plans.id = audits.audit_plan_id
        LEFT JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id
        LEFT JOIN audit_finding_categories ON audit_finding_categories.id = audit_findings.audit_finding_category_id
        LEFT JOIN audit_finding_statuses ON audit_finding_statuses.id = audit_findings.audit_finding_status_id
        LEFT JOIN users AS created_user ON created_user.id = audit_findings.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = audit_findings.updated_by_id
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = audit_finding_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = audit_finding_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = audit_finding_action_plans.responsible_user_id
		LEFT JOIN users AS review_user ON review_user.id = audit_finding_action_plans.review_user_id
    `).
		Where("auditable_items.responsible_unit_id = ? AND audit_statuses.type = ?", currentUserResponsibleUnitID, "completed").
		Group(`
			audit_finding_action_plans.id,
			audit_finding_action_plans.audit_finding_id,
			audit_finding_action_plans.action_plan_status_id,
			audit_finding_action_plans.severity_level_id,
			audit_finding_action_plans.responsible_user_id,
			audit_finding_action_plans.review_user_id,
			audit_finding_action_plans.created_by_id,
			audit_finding_action_plans.updated_by_id,
			audit_finding_action_plans.title,
			audit_finding_action_plans.description,
			audit_finding_action_plans.percentage,
			audit_finding_action_plans.comments,
			audit_finding_action_plans.target_date,
			audit_finding_action_plans.created_at,
			audit_finding_action_plans.updated_at,
			audit_findings.title,
			audit_findings.reference_number,
			audit_findings.description,
			action_plan_statuses.title_en,
			action_plan_statuses.title_ar,
			action_plan_statuses.color_code,
			severity_levels.title_en,
			severity_levels.title_ar,
			severity_levels.color_code,
			responsible_user.id,
			responsible_user.first_name_en,
			responsible_user.first_name_ar,
			responsible_user.last_name_en,
			responsible_user.last_name_ar,
			responsible_user.initial,
			responsible_user.color,
			responsible_user.email,
			responsible_user.profile_image_path,
			review_user.id,
			review_user.first_name_en,
			review_user.first_name_ar,
			review_user.last_name_en,
			review_user.last_name_ar,
			review_user.initial,
			review_user.color,
			review_user.email,
			review_user.profile_image_path,
			created_user.id,
			created_user.first_name_en,
			created_user.first_name_ar,
			created_user.last_name_en,
			created_user.last_name_ar,
			created_user.initial,
			created_user.color,
			created_user.email,
			created_user.profile_image_path,
			updated_user.id,
			updated_user.first_name_en,
			updated_user.first_name_ar,
			updated_user.last_name_en,
			updated_user.last_name_ar,
			updated_user.initial,
			updated_user.color,
			updated_user.email,
			updated_user.profile_image_path`,
		)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply target_date filters
	if params.TargetDateFrom != "" {
		query = query.Where("target_date >= ?", params.TargetDateFrom)
	}
	if params.TargetDateTo != "" {
		query = query.Where("target_date <= ?", params.TargetDateTo)
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "overdue":
			query = query.Where("target_date < ?", time.Now())
		case "completed":
			query = query.Where("action_plan_statuses.type = ?", "completed")
		case "in_progress":
			query = query.Where("action_plan_statuses.type = ?", "in_progress")
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
            (audit_finding_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            audit_finding_action_plans.target_date LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
            responsible_user.last_name_ar LIKE ?)`,

			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword)
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditFindingActionPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditFindingActionPlanListResponseDTO{
		Items: auditFindingActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditInsightsFindingActionPlanRepository) GetFullDetailsByID(id uint) (*models.AuditFindingActionPlan, error) {
	var auditFindingActionPlan *models.AuditFindingActionPlan

	result := repo.db.
		Preload("AuditFinding").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("SeverityLevel").
		Preload("ActionPlanStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ActionPlanStatusUpdates").
		Preload("ReviewUser").
		First(&auditFindingActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditFindingActionPlan, nil
}
