package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDocumentRequestRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditID uint, params *auditDTOs.AuditDocumentRequestListQueryParamsDTO) (*auditDTOs.AuditDocumentRequestListResponseDTO, error)
	GetByID(id uint) (*models.AuditDocumentRequest, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditDocumentRequest, error)
	Create(auditDocumentRequest *models.AuditDocumentRequest, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditDocumentRequest *models.AuditDocumentRequest, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateAuditableItems(auditDocumentRequestAuditableItems []*models.AuditDocumentRequestAuditableItem, tx ...utility.DatabaseTransaction) error
	DeleteAuditableItems(itemIds []*uint, tx ...utility.DatabaseTransaction) error
}

type AuditDocumentRequestRepository struct {
	db *gorm.DB
}

func NewAuditDocumentRequestRepository(db *gorm.DB) *AuditDocumentRequestRepository {
	return &AuditDocumentRequestRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditDocumentRequestRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditDocumentRequestRepository) GetAll(auditID uint, params *auditDTOs.AuditDocumentRequestListQueryParamsDTO) (*auditDTOs.AuditDocumentRequestListResponseDTO, error) {
	var auditDocumentRequests []*auditDTOs.AuditDocumentRequestListDTO
	var totalCount int64

	query := repo.db.Table("audit_document_requests").
		Select(`
			audit_document_requests.*,
			CONVERT(VARCHAR, audit_document_requests.requested_at, 23) AS requested_at,
			audit_document_request_statuses.title_en AS audit_document_request_status_title_en,
			audit_document_request_statuses.title_ar AS audit_document_request_status_title_ar,
			audit_document_request_statuses.type AS audit_document_request_status_type,
			audit_document_request_statuses.color_code AS audit_document_request_status_color_code,
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
			LEFT JOIN audit_document_request_statuses ON audit_document_request_statuses.id = audit_document_requests.audit_document_request_status_id
			LEFT JOIN organization_units ON organization_units.id = audit_document_requests.responsible_unit_id
        	LEFT JOIN users AS responsible_user ON responsible_user.id = audit_document_requests.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = audit_document_requests.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_document_requests.updated_by_id
		`).Where(`audit_document_requests.audit_id = ?`, auditID)

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
			audit_document_requests.title LIKE ? OR
			audit_document_requests.description LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			audit_document_request_statuses.title_en LIKE ? OR
			audit_document_request_statuses.title_ar LIKE ?`,
			keyword, keyword, keyword, keyword, keyword, keyword)
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

func (repo *AuditDocumentRequestRepository) GetByID(id uint) (*models.AuditDocumentRequest, error) {
	var auditDocumentRequest *models.AuditDocumentRequest

	result := repo.db.First(&auditDocumentRequest, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditDocumentRequest, nil
}

func (repo *AuditDocumentRequestRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditDocumentRequest, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditDocumentRequest *models.AuditDocumentRequest

	result := db.
		Preload("Audit").
		Preload("AuditDocumentRequestStatus").
		Preload("AuditableItems").
		Preload("AuditableItems.AuditableItem.AuditStatus").
		Preload("AuditableItems.AuditableItem.AuditorRiskRating").
		Preload("AuditableItems.AuditableItem.Process").
		Preload("AuditableItems.AuditableItem.Control").
		Preload("AuditableItems.AuditableItem.Document").
		Preload("AuditableItems.AuditableItem.IncidentRegister").
		Preload("AuditableItems.AuditableItem.ComplianceRequirement").
		Preload("AuditableItems.AuditableItem.RiskRegister").
		Preload("AuditableItems.AuditableItem.User").
		Preload("AuditableItems.AuditableItem.ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditDocumentRequest, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditDocumentRequest, nil
}

func (repo *AuditDocumentRequestRepository) Create(auditDocumentRequest *models.AuditDocumentRequest, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditDocumentRequest).Error; err != nil {
		return 0, err
	}

	return auditDocumentRequest.ID, nil
}

func (repo *AuditDocumentRequestRepository) Update(auditDocumentRequest *models.AuditDocumentRequest, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditDocumentRequest).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDocumentRequestRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit Document Request
	if err := db.Delete(&models.AuditDocumentRequest{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDocumentRequestRepository) CreateAuditableItems(auditDocumentRequestAuditableItems []*models.AuditDocumentRequestAuditableItem, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&auditDocumentRequestAuditableItems).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDocumentRequestRepository) DeleteAuditableItems(itemIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", itemIds).Delete(&models.AuditDocumentRequestAuditableItem{}).Error; err != nil {
		return err
	}

	return nil
}
