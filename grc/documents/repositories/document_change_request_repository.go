package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentChangeRequestRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.DocumentChangeRequestQueryParamsDTO) (*dtos.DocumentChangeRequestListResponseDTO, error)
	GetByID(id uint) (*models.DocumentChangeRequest, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequest, error)
	Create(changeRequest *models.DocumentChangeRequest, tx ...utility.DatabaseTransaction) (uint, error)
	Update(changeRequest *models.DocumentChangeRequest, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	GetDocumentByID(id uint) (*models.DocumentChangeRequestDocument, error)
	GetDocumentByChangeRequestID(id uint) (*models.DocumentChangeRequestDocument, error)
	CreateDocument(doc *models.DocumentChangeRequestDocument, tx ...utility.DatabaseTransaction) error
	UpdateDocument(doc *models.DocumentChangeRequestDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocument(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

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

func (repo *DocumentChangeRequestRepository) GetAll(params *dtos.DocumentChangeRequestQueryParamsDTO) (*dtos.DocumentChangeRequestListResponseDTO, error) {
	var documentChangeRequests []*dtos.DocumentChangeRequestListDTO
	var totalCount int64

	query := repo.db.Table("documents_change_requests").
		Select(`
			documents_change_requests.id,
			documents_change_requests.reference_number,
			documents_change_requests.document_id,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			documents_change_requests.document_change_request_type_id,
			change_request_types.title_en AS document_change_request_type_en,
			change_request_types.title_ar AS document_change_request_type_ar,
			documents_change_requests.document_change_request_status_id,
			change_request_statuses.title_en AS document_change_request_status_en,
			change_request_statuses.title_ar AS document_change_request_status_ar,
			documents_change_requests.reason_for_change,
			documents_change_requests.version,
			created_by_user.id AS created_by_user_id,
			created_by_user.first_name_en AS created_by_user_first_name_en,
			created_by_user.first_name_ar AS created_by_user_first_name_ar,
			created_by_user.last_name_en AS created_by_user_last_name_en,
			created_by_user.last_name_ar AS created_by_user_last_name_ar,
			created_by_user.email AS created_by_user_email,
			created_by_user.profile_image_path AS created_by_user_profile_image_path,
			updated_by_user.id AS updated_by_user_id,
			updated_by_user.first_name_en AS updated_by_user_first_name_en,
			updated_by_user.first_name_ar AS updated_by_user_first_name_ar,
			updated_by_user.last_name_en AS updated_by_user_last_name_en,
			updated_by_user.last_name_ar AS updated_by_user_last_name_ar,
			updated_by_user.email AS updated_by_user_email,
			updated_by_user.profile_image_path AS updated_by_user_profile_image_path,
			documents_change_requests.created_at,
			documents_change_requests.updated_at
		`).
		Joins("JOIN documents ON documents.id = documents_change_requests.document_id").
		Joins("JOIN document_change_request_types AS change_request_types ON change_request_types.id = documents_change_requests.document_change_request_type_id").
		Joins("JOIN document_change_request_statuses AS change_request_statuses ON change_request_statuses.id = documents_change_requests.document_change_request_status_id").
		Joins("LEFT JOIN users AS created_by_user ON created_by_user.id = documents_change_requests.created_by_id").
		Joins("LEFT JOIN users AS updated_by_user ON updated_by_user.id = documents_change_requests.updated_by_id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where(`
        (
            documents_change_requests.reason_for_change LIKE ? OR
            change_request_types.title_en LIKE ? OR
            change_request_types.title_ar LIKE ? OR
            change_request_statuses.title_en LIKE ? OR
            change_request_statuses.title_ar LIKE ?
        )`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.SortBy == "" {
		params.SortBy = "reference_number"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	err := query.Offset(offset).Limit(params.PageSize).Find(&documentChangeRequests).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.DocumentChangeRequestListResponseDTO{
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
	var changeRequest models.DocumentChangeRequest
	result := repo.db.First(&changeRequest, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &changeRequest, nil
}

func (repo *DocumentChangeRequestRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequest, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var changeRequest models.DocumentChangeRequest
	result := db.Preload("Documents").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&changeRequest, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &changeRequest, nil
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

func (repo *DocumentChangeRequestRepository) GetDocumentByID(id uint) (*models.DocumentChangeRequestDocument, error) {
	var document models.DocumentChangeRequestDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *DocumentChangeRequestRepository) GetDocumentByChangeRequestID(id uint) (*models.DocumentChangeRequestDocument, error) {
	var document models.DocumentChangeRequestDocument
	if err := repo.db.Where("document_change_request_id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *DocumentChangeRequestRepository) CreateDocument(doc *models.DocumentChangeRequestDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&doc).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentChangeRequestRepository) UpdateDocument(doc *models.DocumentChangeRequestDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(&doc).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentChangeRequestRepository) DeleteDocument(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", docIds).Delete(&models.DocumentChangeRequestDocument{}).Error; err != nil {
		return err
	}
	return nil
}
