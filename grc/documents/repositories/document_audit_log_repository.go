package repositories

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentAuditLogRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(documentID uint, documentVersionID uint, params *documentDTOs.DocumentAuditLogListQueryParamsDTO) (*documentDTOs.DocumentAuditLogListResponseDTO, error)
	GetByType(docuemntAuditLog string) (*models.DocumentAuditLogType, error)
	Create(documentAuditLog *models.DocumentAuditLog, tx ...utility.DatabaseTransaction) error
}

type DocumentAuditLogRepository struct {
	db *gorm.DB
}

func NewDocumentAuditLogRepository(db *gorm.DB) *DocumentAuditLogRepository {
	return &DocumentAuditLogRepository{
		db: db,
	}
}

func (repo *DocumentAuditLogRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentAuditLogRepository) GetAll(documentID uint, documentVersionID uint, params *documentDTOs.DocumentAuditLogListQueryParamsDTO) (*documentDTOs.DocumentAuditLogListResponseDTO, error) {
	var auditLogList []*documentDTOs.DocumentAuditLogListDTO
	var totalCount int64

	query := repo.db.Table("document_audit_logs").
		Select(`
			document_audit_logs.id AS id,
			document_audit_logs.document_id AS document_id,
			document_audit_logs.document_version_id AS document_version_id,
			document_audit_logs.user_id AS user_id,
			document_audit_logs.created_by_id AS created_by_id,
			document_audit_logs.created_at AS created_at,
			document_audit_log_types.title_en AS document_audit_log_type_title_en,
			document_audit_log_types.title_ar AS document_audit_log_type_title_ar,
			document_audit_log_types.color_code AS document_audit_log_type_color_code,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			documents.reference_number AS document_reference_number,
			document_versions.url AS document_url,
			document_versions.name AS document_name,
			document_versions.size AS document_size,
			document_versions.token AS document_token,
			document_versions.original_name AS document_original_name,
			document_versions.type AS document_type,
			document_versions.version AS document_version,
			users.first_name_en AS user_first_name_en,
			users.first_name_ar AS user_first_name_ar,
			users.last_name_en AS user_last_name_en,
			users.last_name_ar AS user_last_name_ar,
			users.email AS user_email,
			users.initial AS user_initial,
			users.color AS user_color,
			users.profile_image_path AS user_profile_image_path,
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
			LEFT JOIN document_audit_log_types ON document_audit_log_types.id = document_audit_logs.document_audit_log_type_id
			LEFT JOIN documents ON documents.id = document_audit_logs.document_id 
			LEFT JOIN document_versions ON document_versions.document_id = documents.id 
			LEFT JOIN users ON users.id = document_audit_logs.user_id
			LEFT JOIN users AS created_user ON created_user.id = document_audit_logs.created_by_id
		`).
		Where("document_audit_logs.document_id = ? AND document_audit_logs.document_version_id = ?", documentID, documentVersionID)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(document_audit_log_types.title_en LIKE ? OR
			document_audit_log_types.title_ar LIKE ? OR
			documents.title_en LIKE ? OR
			documents.title_ar LIKE ? OR
			users.first_name_en LIKE ? OR
			users.first_name_ar LIKE ? OR
			users.last_name_en LIKE ? OR
			users.last_name_ar LIKE ? OR 
			users.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Apply sorting
	if params.SortBy == "" {
		params.SortBy = "document_audit_logs.created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditLogList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &documentDTOs.DocumentAuditLogListResponseDTO{
		Items: auditLogList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *DocumentAuditLogRepository) GetByType(docuemntAuditLog string) (*models.DocumentAuditLogType, error) {
	var documentAuditLogType *models.DocumentAuditLogType

	result := repo.db.Where("type = ?", docuemntAuditLog).First(&documentAuditLogType)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Document Audit Log type with type "+docuemntAuditLog)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return documentAuditLogType, nil
}

func (repo *DocumentAuditLogRepository) Create(documentAuditLog *models.DocumentAuditLog, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(documentAuditLog).Error
}
