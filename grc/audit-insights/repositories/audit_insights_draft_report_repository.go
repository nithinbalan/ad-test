package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"

	"gorm.io/gorm"
)

type AuditInsightsDraftReportRepositoryInterface interface {
	GetAll(params *auditDTOs.AuditDraftReportListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditDraftReportListResponseDTO, error)
}

type AuditInsightsDraftReportRepository struct {
	db *gorm.DB
}

func NewAuditInsightsDraftReportRepository(db *gorm.DB) *AuditInsightsDraftReportRepository {
	return &AuditInsightsDraftReportRepository{
		db: db,
	}
}

func (repo *AuditInsightsDraftReportRepository) GetAll(params *auditDTOs.AuditDraftReportListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditDraftReportListResponseDTO, error) {
	var auditDraftReports []*auditDTOs.AuditDraftReportListDTO
	var totalCount int64

	query := repo.db.Table("audit_plan_auditable_items").
		Select(`
            audit_draft_reports.id,
            audit_draft_reports.audit_id,
            audit_draft_reports.report_date,
            audit_draft_reports.executive_summary,
            audit_draft_reports.conclusion,
            audit_draft_report_statuses.id AS audit_draft_report_status_id,
            audit_draft_report_statuses.title_en AS audit_draft_report_status_title_en,
            audit_draft_report_statuses.title_ar AS audit_draft_report_status_title_ar,
            audit_draft_report_statuses.color_code AS audit_draft_report_status_color_code,
            audit_draft_report_statuses.type AS audit_draft_report_status_type,
            audits.audit_plan_id,
            CONVERT(VARCHAR, audits.start_date, 23) AS start_date,
            CONVERT(VARCHAR, audits.end_date, 23) AS end_date,
            audit_plans.title AS audit_plan_title,
            audit_plans.description AS audit_plan_description,
            audit_plans.purpose AS audit_plan_purpose,
            audit_plans.objective AS audit_plan_objective,
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
            audit_draft_reports.internal_review_submitted_at,
            internal_review_user.id AS internal_review_user_id,
            internal_review_user.first_name_en AS internal_review_user_first_name_en,
            internal_review_user.first_name_ar AS internal_review_user_first_name_ar,
            internal_review_user.last_name_en AS internal_review_user_last_name_en,
            internal_review_user.last_name_ar AS internal_review_user_last_name_ar,
            internal_review_user.initial AS internal_review_user_initial,
            internal_review_user.color AS internal_review_user_color,
            internal_review_user.email AS internal_review_user_email,
            internal_review_user.profile_image_path AS internal_review_user_profile_image_path,
            audit_draft_reports.auditee_review_submitted_at,
            auditee_review_user.id AS auditee_review_user_id,
            auditee_review_user.first_name_en AS auditee_review_user_first_name_en,
            auditee_review_user.first_name_ar AS auditee_review_user_first_name_ar,
            auditee_review_user.last_name_en AS auditee_review_user_last_name_en,
            auditee_review_user.last_name_ar AS auditee_review_user_last_name_ar,
            auditee_review_user.initial AS auditee_review_user_initial,
            auditee_review_user.color AS auditee_review_user_color,
            auditee_review_user.email AS auditee_review_user_email,
            auditee_review_user.profile_image_path AS auditee_review_user_profile_image_path,
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
            LEFT JOIN audit_draft_reports ON audit_draft_reports.audit_id = audits.id
            LEFT JOIN audit_draft_report_statuses ON audit_draft_report_statuses.id = audit_draft_reports.audit_draft_report_status_id
            LEFT JOIN users AS internal_review_user ON internal_review_user.id = audit_draft_reports.internal_review_submitted_by_id
		    LEFT JOIN users AS auditee_review_user ON auditee_review_user.id = audit_draft_reports.auditee_review_submitted_by_id
            LEFT JOIN users AS created_user ON created_user.id = audit_draft_reports.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = audit_draft_reports.updated_by_id
        `).
		Where(`auditable_items.responsible_unit_id = ? AND audit_draft_reports.id IS NOT NULL`, currentUserResponsibleUnitID).
		Group(`
            audit_draft_reports.id,
            audit_draft_reports.audit_id,
            audit_draft_reports.report_date,
            audit_draft_reports.executive_summary,
            audit_draft_reports.conclusion,
            audit_draft_report_statuses.id,
            audit_draft_report_statuses.title_en,
            audit_draft_report_statuses.title_ar,
            audit_draft_report_statuses.color_code,
            audit_draft_report_statuses.type,
            audits.audit_plan_id,
            audits.start_date,
            audits.end_date,
            audit_plans.title,
            audit_plans.description,
            audit_plans.purpose,
            audit_plans.objective,
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
            audit_draft_reports.internal_review_submitted_at,
	        audit_draft_reports.auditee_review_submitted_at,
            internal_review_user.id,
            internal_review_user.first_name_en,
            internal_review_user.first_name_ar,
            internal_review_user.last_name_en,
            internal_review_user.last_name_ar,
            internal_review_user.initial,
            internal_review_user.color,
            internal_review_user.email,
            internal_review_user.profile_image_path,
            auditee_review_user.id,
            auditee_review_user.first_name_en,
            auditee_review_user.first_name_ar,
            auditee_review_user.last_name_en,
            auditee_review_user.last_name_ar,
            auditee_review_user.initial,
            auditee_review_user.color,
            auditee_review_user.email,
            auditee_review_user.profile_image_path,
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
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply other filters as needed (e.g., date filters)
	if params.ReportDate != "" {
		query = query.Where("audits.report_date = ?", params.ReportDate)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audit_draft_reports.report_date LIKE ? OR
			audits.start_date LIKE ? OR
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
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditDraftReports).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditDraftReportListResponseDTO{
		Items: auditDraftReports,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}
