package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"

	"gorm.io/gorm"
)

type AuditInsightsAuditRepositoryInterface interface {
	GetAll(params *auditDTOs.AuditListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditListResponseDTO, error)
}

type AuditInsightsAuditRepository struct {
	db *gorm.DB
}

func NewAuditInsightsAuditRepository(db *gorm.DB) *AuditInsightsAuditRepository {
	return &AuditInsightsAuditRepository{
		db: db,
	}
}

func (repo *AuditInsightsAuditRepository) GetAll(params *auditDTOs.AuditListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditListResponseDTO, error) {
	var audits []*auditDTOs.AuditListDTO
	var totalCount int64

	query := repo.db.Table("audit_plan_auditable_items").
		Select(`
        audits.id,
        audits.audit_plan_id,
        CONVERT(VARCHAR, audits.start_date, 23) AS start_date,
        CONVERT(VARCHAR, audits.end_date, 23) AS end_date,
        audit_plans.title AS audit_plan_title,
        audit_plans.description AS audit_plan_description,
        audit_plans.purpose AS audit_plan_purpose,
        audit_plans.objective AS audit_plan_objective,
        audit_plans.audit_status_id AS audit_status_id,
        audit_statuses.title_en AS audit_status_title_en,
        audit_statuses.title_ar AS audit_status_title_ar,
        audit_statuses.type AS audit_status_type,
        audit_statuses.color_code AS audit_status_color_code,
        audit_categories.title_en AS audit_category_title_en,
        audit_categories.title_ar AS audit_category_title_ar,
        team_lead.id AS team_lead_id,
        team_lead.first_name_en AS team_lead_first_name_en,
        team_lead.first_name_ar AS team_lead_first_name_ar,
        team_lead.last_name_en AS team_lead_last_name_en,
        team_lead.last_name_ar AS team_lead_last_name_ar,
        team_lead.email AS team_lead_email,
        team_lead.initial AS team_lead_initial,
        team_lead.color AS team_lead_color,
        team_lead.profile_image_path AS team_lead_profile_image_path,
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
        LEFT JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
        LEFT JOIN audit_plans ON audit_plans.id = audit_plan_auditable_items.audit_plan_id
        LEFT JOIN audits ON audits.audit_plan_id = audit_plans.id
        LEFT JOIN audit_statuses ON audit_statuses.id = audit_plans.audit_status_id
        LEFT JOIN audit_categories ON audit_categories.id = audit_plans.audit_category_id
        LEFT JOIN users AS team_lead ON team_lead.id = audit_plans.team_lead_id
        LEFT JOIN users AS created_user ON created_user.id = audits.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = audits.updated_by_id
    `).
		Where(`auditable_items.responsible_unit_id = ? AND audits.id IS NOT NULL`, currentUserResponsibleUnitID).
		Group(`
        audits.id,
        audits.audit_plan_id,
        audits.start_date,
        audits.end_date,
        audit_plans.title,
        audit_plans.description,
        audit_plans.purpose,
        audit_plans.objective,
        audit_plans.audit_status_id,
        audit_statuses.title_en,
        audit_statuses.title_ar,
        audit_statuses.type,
        audit_statuses.color_code,
        audit_categories.title_en,
        audit_categories.title_ar,
        team_lead.id,
        team_lead.first_name_en,
        team_lead.first_name_ar,
        team_lead.last_name_en,
        team_lead.last_name_ar,
        team_lead.email,
        team_lead.initial,
        team_lead.color,
        team_lead.profile_image_path,
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
        updated_user.profile_image_path
    `)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			//audits and audits plans using same fields as audit_status_id, so attaching table name in to the field
			if columnName == "audit_status_id" {
				// Prepend the table name for "audit_status_id"
				query = query.Where("audits.audit_status_id IN (?)", values)
			} else {
				// Use the column name as-is for other filters
				query = query.Where(columnName+" IN (?)", values)
			}
		}
	}

	// Apply other filters as needed (e.g., date filters)
	if params.StartDate != "" {
		query = query.Where("audits.start_date >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("audits.end_date <= ?", params.EndDate)
	}
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audits.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audits.start_date LIKE ? OR
			audits.end_date LIKE ? OR
			audit_plans.title LIKE ? OR
			audit_categories.title_en LIKE ? OR
			audit_categories.title_ar LIKE ? OR
			team_lead.first_name_en LIKE ? OR
			team_lead.first_name_ar LIKE ? OR
			team_lead.last_name_en LIKE ? OR
			team_lead.last_name_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}
	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audits.start_date" // Default sort by start_date if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&audits).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditListResponseDTO{
		Items: audits,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}
