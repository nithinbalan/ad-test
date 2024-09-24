package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"

	"gorm.io/gorm"
)

type AuditInsightsDocumentRequestRepositoryInterface interface {
	GetAll(params *auditDTOs.AuditDocumentRequestListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditDocumentRequestListResponseDTO, error)
}

type AuditInsightsDocumentRequestRepository struct {
	db *gorm.DB
}

func NewAuditInsightsDocumentRequestRepository(db *gorm.DB) *AuditInsightsDocumentRequestRepository {
	return &AuditInsightsDocumentRequestRepository{
		db: db,
	}
}

func (repo *AuditInsightsDocumentRequestRepository) GetAll(params *auditDTOs.AuditDocumentRequestListQueryParamsDTO, currentUserResponsibleUnitID uint) (*auditDTOs.AuditDocumentRequestListResponseDTO, error) {
	var auditDocumentRequests []*auditDTOs.AuditDocumentRequestListDTO
	var totalCount int64

	query := repo.db.Table("audit_document_requests").
		Select(`
			audit_document_requests.*,
			audit_document_request_statuses.title_en AS audit_document_request_status_title_en,
			audit_document_request_statuses.title_ar AS audit_document_request_status_title_ar,
			audit_document_request_statuses.type AS audit_document_request_status_type,
			audit_document_request_statuses.color_code AS audit_document_request_status_color_code,
			audit_plans.id AS audit_plan_id,
			audit_plans.title AS audit_plan_title,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
			requested_user.id AS requested_user_id,
			requested_user.first_name_en AS requested_user_first_name_en,
        	requested_user.first_name_ar AS requested_user_first_name_ar,
        	requested_user.last_name_en AS requested_user_last_name_en,
        	requested_user.last_name_ar AS requested_user_last_name_ar,   
        	requested_user.initial AS requested_user_initial,
        	requested_user.color AS requested_user_color,
        	requested_user.email AS requested_user_email,
        	requested_user.profile_image_path AS requested_user_profile_image_path,
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
			LEFT JOIN audit_document_request_statuses ON audit_document_request_statuses.id = audit_document_requests.audit_document_request_status_id
			LEFT JOIN organization_units ON organization_units.id = audit_document_requests.responsible_unit_id
			LEFT JOIN audits ON audits.id = audit_document_requests.audit_id
			LEFT JOIN audit_plans ON audit_plans.id = audits.audit_plan_id
        	LEFT JOIN users AS responsible_user ON responsible_user.id = audit_document_requests.responsible_user_id
			LEFT JOIN users AS requested_user ON requested_user.id = audit_document_requests.requested_by_id
			LEFT JOIN users AS created_user ON created_user.id = audit_document_requests.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_document_requests.updated_by_id
		`).Where(`audit_document_requests.responsible_unit_id = ?`, currentUserResponsibleUnitID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Apply other filters as needed (e.g., date filters)
	if params.RequestedAt != "" {
		query = query.Where("requested_at = ?", params.RequestedAt)
	}
	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("audit_document_request_statuses.type = ?", "open")
		case "draft":
			query = query.Where("audit_document_request_statuses.type = ?", "draft")
		case "closed":
			query = query.Where("audit_document_request_statuses.type = ?", "closed")
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
			audit_plans.title LIKE ? OR
			audit_document_requests.title LIKE ? OR
			audit_document_requests.description LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			audit_document_request_statuses.title_en LIKE ? OR
			audit_document_request_statuses.title_ar LIKE ?`,
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditDocumentRequests).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditDocumentRequestListResponseDTO{
		Items: auditDocumentRequests,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}
