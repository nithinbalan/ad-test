package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

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

	GenerateReferenceNumber() (string, error)

	GetDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestDocument, error)
	GetDocumentByChangeRequestID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestDocument, error)
	CreateDocument(doc *models.DocumentChangeRequestDocument, tx ...utility.DatabaseTransaction) error
	UpdateDocument(doc *models.DocumentChangeRequestDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocument(docId uint, tx ...utility.DatabaseTransaction) error
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

	query := repo.db.Table("document_change_requests").
		Select(`
			document_change_requests.id,
			document_change_requests.reference_number,
			document_change_requests.document_id,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			document_change_requests.document_change_request_type_id,
			change_request_types.title_en AS document_change_request_type_en,
			change_request_types.title_ar AS document_change_request_type_ar,
			document_change_requests.document_change_request_status_id,
			change_request_statuses.title_en AS document_change_request_status_en,
			change_request_statuses.title_ar AS document_change_request_status_ar,
			change_request_statuses.color_code AS document_change_request_status_color_code,
			document_change_requests.reason_for_change,
			document_change_requests.version,
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
			document_change_requests.created_at,
			document_change_requests.updated_at
		`).
		Joins("JOIN documents ON documents.id = document_change_requests.document_id").
		Joins("JOIN document_change_request_types AS change_request_types ON change_request_types.id = document_change_requests.document_change_request_type_id").
		Joins("JOIN document_change_request_statuses AS change_request_statuses ON change_request_statuses.id = document_change_requests.document_change_request_status_id").
		Joins("LEFT JOIN users AS created_by_user ON created_by_user.id = document_change_requests.created_by_id").
		Joins("LEFT JOIN users AS updated_by_user ON updated_by_user.id = document_change_requests.updated_by_id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where(`
        (
            document_change_requests.reason_for_change LIKE ? OR
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
	result := db.Preload("DocumentChangeRequestType").
		Preload("DocumentChangeRequestStatus").
		Preload("ChangeRequestDocument").
		Preload("DocumentChangeRequestInitiate").
		Preload("Document").
		Preload("Document.ResponsibleUser").
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

func (repo *DocumentChangeRequestRepository) GetDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var document models.DocumentChangeRequestDocument
	if err := db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *DocumentChangeRequestRepository) GetDocumentByChangeRequestID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var document models.DocumentChangeRequestDocument
	if err := db.Where("document_change_request_id = ?", id).First(&document).Error; err != nil {
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

func (repo *DocumentChangeRequestRepository) DeleteDocument(docId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id = ? ", docId).Delete(&models.DocumentChangeRequestDocument{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentChangeRequestRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.DocumentChangeRequest{}).
		Select("reference_number").
		Where("reference_number IS NOT NULL"). // Ensure to filter out NULL values
		Order("reference_number DESC").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "CR", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}
