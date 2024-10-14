package repositories

import (
	auditFindingDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"

	"gorm.io/gorm"
)

type AuditInsightsDraftFindingActionPlanRepositoryInterface interface {
	GetAll(findingID uint, currentUserResponsibleUnitID uint, params *auditFindingDTOs.AuditDraftFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditDraftFindingActionPlanListResponseDTO, error)

}

type AuditInsightsDraftFindingActionPlanRepository struct {
	db *gorm.DB
}

func NewAuditInsightsDraftFindingActionPlanRepository(db *gorm.DB) *AuditInsightsDraftFindingActionPlanRepository {
	return &AuditInsightsDraftFindingActionPlanRepository{
		db: db,
	}
}

func (repo *AuditInsightsDraftFindingActionPlanRepository) GetAll(findingID uint, currentUserResponsibleUnitID uint,params *auditFindingDTOs.AuditDraftFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditDraftFindingActionPlanListResponseDTO, error) {

	var auditInsightsDraftFindingActionPlans []*auditFindingDTOs.AuditDraftFindingActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("audit_draft_finding_action_plans").
		Select(`
        audit_draft_finding_action_plans.*,
		audit_draft_findings.title AS audit_finding_title,
		audit_draft_findings.description AS audit_finding_description,
		CONVERT(VARCHAR, audit_draft_finding_action_plans.target_date, 23) AS target_date,
        responsible_user.id AS responsible_user_id,
        responsible_user.first_name_en AS responsible_user_first_name_en,
        responsible_user.first_name_ar AS responsible_user_first_name_ar,
        responsible_user.last_name_en AS responsible_user_last_name_en,
        responsible_user.last_name_ar AS responsible_user_last_name_ar,
		responsible_user.initial AS responsible_user_initial,
        responsible_user.color AS responsible_user_color,
        responsible_user.email AS responsible_user_email,
        responsible_user.profile_image_path AS responsible_user_profile_image_path, 
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
		LEFT JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_action_plans.audit_draft_finding_id
		LEFT JOIN auditable_items ON auditable_items.id = audit_draft_findings.auditable_item_id
		LEFT JOIN users AS responsible_user ON responsible_user.id = audit_draft_finding_action_plans.responsible_user_id
        LEFT JOIN users AS created_user ON created_user.id = audit_draft_finding_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = audit_draft_finding_action_plans.updated_by_id
    `).Where("audit_draft_finding_action_plans.audit_draft_finding_id = ? AND auditable_items.responsible_unit_id = ?", findingID, currentUserResponsibleUnitID)

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

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (audit_draft_finding_action_plans.title LIKE ? OR
            audit_draft_finding_action_plans.target_date LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
            responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword,
			keyword)
	}
	// Conditionally apply finding ID filter
	if findingID != 0 {
		query = query.Where("audit_draft_finding_action_plans.audit_draft_finding_id = ?", findingID)
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditInsightsDraftFindingActionPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditFindingDTOs.AuditDraftFindingActionPlanListResponseDTO{
		Items: auditInsightsDraftFindingActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

