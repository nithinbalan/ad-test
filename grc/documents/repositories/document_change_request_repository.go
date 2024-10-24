package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface for the DocumentChangeRequestRepository
type DocumentChangeRequestRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *documentDTOs.DocumentChangeRequestQueryParamsDTO) (*documentDTOs.DocumentChangeRequestListResponseDTO, error)
	GetByID(id uint) (*models.DocumentChangeRequest, error)
	GetByDocumentID(documentID uint) (*models.DocumentChangeRequest, error)
	GetDocumentByID(documentID uint) (*models.DocumentChangeRequestSupportingDocument, error) // New method
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequest, error)
	Create(request *models.DocumentChangeRequest, tx ...utility.DatabaseTransaction) (uint, error)
	Update(request *models.DocumentChangeRequest, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateSupportingDocuments(docs []*models.DocumentChangeRequestSupportingDocument, tx ...utility.DatabaseTransaction) error
	DeleteSupportingDocuments(docIDs []*uint, tx ...utility.DatabaseTransaction) error
}

// Implementation of the DocumentChangeRequestRepository
type DocumentChangeRequestRepository struct {
	db *gorm.DB
}

func NewDocumentChangeRequestRepository(db *gorm.DB) *DocumentChangeRequestRepository {
	return &DocumentChangeRequestRepository{
		db: db,
	}
}

func (repo *DocumentChangeRequestRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentChangeRequestRepository) GetAll(params *documentDTOs.DocumentChangeRequestQueryParamsDTO) (*documentDTOs.DocumentChangeRequestListResponseDTO, error) {
	var documentChangeRequests []*documentDTOs.DocumentChangeRequestListDTO
	var totalCount int64

	query := repo.db.Table("document_change_requests").
		Select(`
			document_change_requests.id,
			document_change_requests.documents_id,
			documents.title_en AS document_title,
			document_change_requests.document_change_request_type_id,
			change_request_types.title_en AS document_change_request_type_en,
			change_request_types.title_ar AS document_change_request_type_ar,
			document_change_requests.reason_for_change,
			document_change_requests.document_change_request_statuses_id,
			change_request_statuses.title_en AS document_change_request_status_en,
			change_request_statuses.title_ar AS document_change_request_status_ar,
			change_request_statuses.color_code AS document_change_request_status_color,
			created_by_user.id AS created_by_user_id,
			created_by_user.first_name_en AS created_by_user_first_name_en,
			created_by_user.last_name_en AS created_by_user_last_name_en,
			created_by_user.email AS created_by_user_email,
			created_by_user.profile_image_path AS created_by_user_profile_image_path,
			document_change_requests.created_at,
			document_change_requests.updated_at
		`).
		Joins("JOIN documents ON documents.id = document_change_requests.documents_id").
		Joins("JOIN document_change_request_types AS change_request_types ON change_request_types.id = document_change_requests.document_change_request_type_id").
		Joins("JOIN document_change_request_statuses AS change_request_statuses ON change_request_statuses.id = document_change_requests.document_change_request_statuses_id").
		Joins("LEFT JOIN users AS created_by_user ON created_by_user.id = document_change_requests.created_by_id")

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting, default sort by document title
	if params.SortBy == "" {
		params.SortBy = "document_title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Pagination
	offset := (params.Page - 1) * params.PageSize
	err := query.Offset(offset).Limit(params.PageSize).Find(&documentChangeRequests).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	// Build the response
	response := &documentDTOs.DocumentChangeRequestListResponseDTO{
		Items: documentChangeRequests,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *DocumentChangeRequestRepository) GetByID(id uint) (*models.DocumentChangeRequest, error) {
	var changeRequest *models.DocumentChangeRequest
	result := repo.db.First(&changeRequest, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return changeRequest, nil
}

// GetByDocumentID retrieves a document change request by document ID
func (repo *DocumentChangeRequestRepository) GetByDocumentID(documentID uint) (*models.DocumentChangeRequest, error) {
	var changeRequest *models.DocumentChangeRequest
	result := repo.db.Where("documents_id = ?", documentID).First(&changeRequest)
	if result.Error != nil {
		return nil, result.Error
	}
	return changeRequest, nil
}

// New Method: GetDocumentByID retrieves a document by its ID
func (repo *DocumentChangeRequestRepository) GetDocumentByID(documentID uint) (*models.DocumentChangeRequestSupportingDocument, error) {
	var document *models.DocumentChangeRequestSupportingDocument
	result := repo.db.First(&document, documentID)
	if result.Error != nil {
		return nil, result.Error
	}
	return document, nil
}

func (repo *DocumentChangeRequestRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequest, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var changeRequest *models.DocumentChangeRequest

	result := db.Preload("Documents").
		Preload("DocumentChangeRequestType").
		Preload("DocumentChangeRequestStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&changeRequest, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return changeRequest, nil
}

func (repo *DocumentChangeRequestRepository) Create(changeRequest *models.DocumentChangeRequest, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(changeRequest).Error; err != nil {
		return 0, err
	}
	return changeRequest.ID, nil
}

func (repo *DocumentChangeRequestRepository) Update(changeRequest *models.DocumentChangeRequest, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(changeRequest).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentChangeRequestRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.DocumentChangeRequest{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentChangeRequestRepository) CreateSupportingDocuments(docs []*models.DocumentChangeRequestSupportingDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentChangeRequestRepository) DeleteSupportingDocuments(docIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", docIDs).Delete(&models.DocumentChangeRequestSupportingDocument{}).Error; err != nil {
		return err
	}
	return nil
}
