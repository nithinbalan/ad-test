package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"

	"gorm.io/gorm"
)

type AuditInsightsAuditFindingRepositoryInterface interface {
	GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditFindingListResponseDTO, error)
}

type AuditInsightsAuditFindingRepository struct {
	db *gorm.DB
}

func NewAuditInsightsAuditFindingRepository(db *gorm.DB) *AuditInsightsAuditFindingRepository {
	return &AuditInsightsAuditFindingRepository{
		db: db,
	}
}

func (repo *AuditInsightsAuditFindingRepository) GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditFindingListResponseDTO, error) {
	var auditFindings []*auditDTOs.AuditFindingListDTO
	var totalCount int64

	query := repo.db.Table("audit_findings").
		Select(`
		audit_findings.id,
		audit_findings.audit_id,
		audit_findings.title,
		audit_findings.reference_number,
		audit_findings.description,
		audit_findings.auditable_item_id,
		audit_findings.audit_finding_category_id,
		audit_findings.severity_level_id,
		audit_findings.recommendations,
		audit_findings.audit_finding_status_id,
		audit_findings.created_at,
		audit_findings.updated_at,
		audit_findings.created_by_id,
		audit_findings.updated_by_id,
		audits.id AS audit_id,
		audit_plans.title AS audit_title,
		auditable_items.title_en AS auditable_item_title_en,
		auditable_items.title_ar AS auditable_item_title_ar,
		audit_finding_categories.title_en AS audit_finding_category_title_en,
		audit_finding_categories.title_ar AS audit_finding_category_title_ar,
		severity_levels.title_en AS severity_level_title_en,
		severity_levels.title_ar AS severity_level_title_ar,
		severity_levels.color_code AS severity_level_color_code,
		audit_finding_statuses.title_en AS audit_finding_status_title_en,
		audit_finding_statuses.title_ar AS audit_finding_status_title_ar,
		audit_finding_statuses.color_code AS audit_finding_status_color_code,
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
		LEFT JOIN audits ON audits.id = audit_findings.audit_id
        LEFT JOIN audit_statuses ON audit_statuses.id = audits.audit_status_id
		LEFT JOIN audit_plans ON audit_plans.id = audits.audit_plan_id
		LEFT JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id
		LEFT JOIN audit_finding_categories ON audit_finding_categories.id = audit_findings.audit_finding_category_id
		LEFT JOIN severity_levels ON severity_levels.id = audit_findings.severity_level_id
		LEFT JOIN audit_finding_statuses ON audit_finding_statuses.id = audit_findings.audit_finding_status_id
		LEFT JOIN users AS created_user ON created_user.id = audit_findings.created_by_id
		LEFT JOIN users AS updated_user ON updated_user.id = audit_findings.updated_by_id
	`).
		Where("auditable_items.responsible_unit_id = ? AND audit_statuses.type = ?", currentUserResponsibleUnitID, "completed").
		Group(`
		audit_findings.id,
		audit_findings.audit_id,
        audit_findings.title,
		audit_findings.reference_number,
		audit_findings.description,
		audit_findings.auditable_item_id,
		audit_findings.audit_finding_category_id,
		audit_findings.severity_level_id,
		audit_findings.recommendations,
		audit_findings.audit_finding_status_id,
		audit_findings.created_at,
		audit_findings.updated_at,
		audit_findings.created_by_id,
		audit_findings.updated_by_id,
		audits.id,
		audit_plans.title,
		auditable_items.title_en,
		auditable_items.title_ar,
		audit_finding_categories.title_en,
		audit_finding_categories.title_ar,
		severity_levels.title_en,
		severity_levels.title_ar,
		severity_levels.color_code,
		audit_finding_statuses.title_en,
		audit_finding_statuses.title_ar,
		audit_finding_statuses.color_code,
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

	for _, filter := range params.QuickFilter {
		switch filter {
		case "new":
			query = query.Where("audit_finding_statuses.type = ?", "new")
		case "in_progress":
			query = query.Where("audit_finding_statuses.type = ?", "in-progress")
		case "resolved":
			query = query.Where("audit_finding_statuses.type = ?", "resolved")
		case "reopened":
			query = query.Where("audit_finding_statuses.type = ?", "reopened")
		case "closed":
			query = query.Where("audit_finding_statuses.type = ?", "closed")
		case "deferred":
			query = query.Where("audit_finding_statuses.type = ?", "deferred")
		case "is_closed":
			query = query.Where("audit_finding_statuses.flag = ?", 1) // Example flag-based filter
		default:
			// Optional: Add a default filter or ignore unknown filters
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
			(auditable_items.title_en LIKE ? OR
			audit_findings.title LIKE ? OR
			audit_findings.reference_number LIKE ? OR
			audit_finding_categories.title_en LIKE ? OR
			severity_levels.title_en LIKE ? OR
			audit_finding_statuses.title_en LIKE ? OR
			audit_finding_statuses.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			updated_user.first_name_en LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audit_findings.created_at" // Default sort by created_at if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditFindings).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditFindingListResponseDTO{
		Items: auditFindings,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}
