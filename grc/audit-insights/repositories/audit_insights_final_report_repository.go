package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"

	"gorm.io/gorm"
)

type AuditInsightsFinalReportRepositoryInterface interface {
	GetAll(params *auditDTOs.AuditFinalReportListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditFinalReportListResponseDTO, error)
}

type AuditInsightsFinalReportRepository struct {
	db *gorm.DB
}

func NewAuditInsightsFinalReportRepository(db *gorm.DB) *AuditInsightsFinalReportRepository {
	return &AuditInsightsFinalReportRepository{
		db: db,
	}
}

func (repo *AuditInsightsFinalReportRepository) GetAll(params *auditDTOs.AuditFinalReportListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditFinalReportListResponseDTO, error) {
	var auditFinalReports []*auditDTOs.AuditFinalReportListDTO
	var totalCount int64

	query := repo.db.Table("audit_plan_auditable_items").
		Select(`
            audit_final_reports.id,
            audit_final_reports.audit_id,
            audit_final_reports.submitted_at,
            audit_final_reports.reference_number,
            audit_final_reports.next_review_level,
            audit_final_reports.audit_final_report_status_id,
            audit_final_report_statuses.title_en AS audit_final_report_status_title_en,
            audit_final_report_statuses.title_ar AS audit_final_report_status_title_ar,
            audit_final_report_statuses.color_code AS audit_final_report_status_color_code,
            audit_final_report_statuses.type AS audit_final_report_status_type,
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
            submitted_user.id AS submitted_user_id,
            submitted_user.first_name_en AS submitted_user_first_name_en,
            submitted_user.first_name_ar AS submitted_user_first_name_ar,
            submitted_user.last_name_en AS submitted_user_last_name_en,
            submitted_user.last_name_ar AS submitted_user_last_name_ar,
            submitted_user.initial AS submitted_user_initial,
            submitted_user.color AS submitted_user_color,
            submitted_user.email AS submitted_user_email,
            submitted_user.profile_image_path AS submitted_user_profile_image_path,
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
            LEFT JOIN audit_final_reports ON audit_final_reports.audit_id = audits.id
            LEFT JOIN audit_final_report_statuses ON audit_final_report_statuses.id = audit_final_reports.audit_final_report_status_id
            LEFT JOIN users AS submitted_user ON submitted_user.id = audit_final_reports.submitted_by_id
            LEFT JOIN users AS created_user ON created_user.id = audit_final_reports.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = audit_final_reports.updated_by_id
        `).
		Where("auditable_items.responsible_unit_id = ? AND audit_final_report_statuses.type ? AND audit_final_reports.id IS NOT NULL", currentUserResponsibleUnitID, "published").
		Group(`
            audit_final_reports.id,
            audit_final_reports.audit_id,
            audit_final_reports.submitted_at,
            audit_final_reports.reference_number,
            audit_final_reports.next_review_level,
            audit_final_reports.audit_final_report_status_id,
            audit_final_report_statuses.id,
            audit_final_report_statuses.title_en,
            audit_final_report_statuses.title_ar,
            audit_final_report_statuses.color_code,
            audit_final_report_statuses.type,
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
            submitted_user.id,
            submitted_user.first_name_en,
            submitted_user.first_name_ar,
            submitted_user.last_name_en,
            submitted_user.last_name_ar,
            submitted_user.initial,
            submitted_user.color,
            submitted_user.email,
            submitted_user.profile_image_path,
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
	// Apply dynamic filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply date filters
	if params.ReportDate != "" {
		query = query.Where("audit_final_reports.submitted_at = ?", params.ReportDate)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audit_final_reports.reference_number LIKE ? OR
			audit_plans.title LIKE ? OR
			team_lead.first_name_en LIKE ? OR
			team_lead.first_name_ar LIKE ? OR
			team_lead.last_name_en LIKE ? OR
			team_lead.last_name_ar LIKE ? OR
			submitted_user.first_name_en LIKE ? OR
			submitted_user.first_name_ar LIKE ? OR
			submitted_user.last_name_en LIKE ? OR
			submitted_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditFinalReports).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditFinalReportListResponseDTO{
		Items: auditFinalReports,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}
