package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentChangeRequestInitiateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.DocumentChangeRequestInitiateQueryParamsDTO) (*dtos.DocumentChangeRequestInitiateListResponseDTO, error)
	GetByID(id uint) (*models.DocumentChangeRequestInitiate, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestInitiate, error)
	Create(changeRequest *models.DocumentChangeRequestInitiate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(changeRequest *models.DocumentChangeRequestInitiate, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	GetDocumentByID(id uint) (*models.DocumentChangeRequestInitiateDocument, error)
	CreateDocuments(docs []*models.DocumentChangeRequestInitiateDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type DocumentChangeRequestInitiateRepository struct {
	db *gorm.DB
}

func NewDocumentChangeRequestInitiateRepository(db *gorm.DB) *DocumentChangeRequestInitiateRepository {
	return &DocumentChangeRequestInitiateRepository{
		db: db,
	}
}

func (repo *DocumentChangeRequestInitiateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentChangeRequestInitiateRepository) GetAll(params *dtos.DocumentChangeRequestInitiateQueryParamsDTO) (*dtos.DocumentChangeRequestInitiateListResponseDTO, error) {
	var documentChangeRequests []*dtos.DocumentChangeRequestInitiateListDTO
	var totalCount int64

	query := repo.db.Table("document_change_request_initiates").
		Select(`
			document_change_request_initiates.id,
			document_change_request_initiates.documents_id,
			documents.title_en AS document_title,
			documents.title_ar AS document_title_ar,
			document_change_request_initiates.document_change_request_initiate_type_id,
			initiate_types.title_en AS document_change_request_initiate_type_title_en,
			initiate_types.title_ar AS document_change_request_initiate_type_title_ar,
			document_change_request_initiates.reason_for_change,
			document_change_request_initiates.document_change_request_initiate_status_id,
			initiate_statuses.title_en AS document_change_request_initiate_status_title_en,
			initiate_statuses.title_ar AS document_change_request_initiate_status_title_ar,
			initiate_statuses.color_code AS document_change_request_initiate_status_color_code,
			created_by_user.id AS created_by_user_id,
			created_by_user.first_name_en AS created_by_user_first_name_en,
			created_by_user.first_name_ar AS created_by_user_first_name_ar,
			created_by_user.last_name_en AS created_by_user_last_name_en,
			created_by_user.last_name_ar AS created_by_user_last_name_ar,
			created_by_user.email AS created_by_user_email,
			created_by_user.initial AS created_user_initial,
	        created_by_user.color AS created_user_color,
			created_by_user.profile_image_path AS created_by_user_profile_image_path,
			updated_by_user.id AS updated_by_user_id,
			updated_by_user.first_name_en AS updated_by_user_first_name_en,
			updated_by_user.first_name_ar AS updated_by_user_first_name_ar,
			updated_by_user.last_name_en AS updated_by_user_last_name_en,
			updated_by_user.last_name_ar AS updated_by_user_last_name_ar,
			updated_by_user.email AS updated_by_user_email,
			updated_by_user.initial AS updated_user_initial,
	        updated_by_user.color AS updated_user_color,
			updated_by_user.profile_image_path AS updated_by_user_profile_image_path,
			document_change_request_initiates.created_at,
			document_change_request_initiates.updated_at
		`).
		Joins("JOIN documents ON documents.id = document_change_request_initiates.documents_id").
		Joins("JOIN document_change_request_initiate_types AS initiate_types ON initiate_types.id = document_change_request_initiates.document_change_request_initiate_type_id").
		Joins("JOIN document_change_request_initiate_statuses AS initiate_statuses ON initiate_statuses.id = document_change_request_initiates.document_change_request_initiate_status_id").
		Joins("LEFT JOIN users AS created_by_user ON created_by_user.id = document_change_request_initiates.created_by_id").
		Joins("LEFT JOIN users AS updated_by_user ON updated_by_user.id = document_change_request_initiates.updated_by_id")

	// Apply dynamic filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters for initiate status
	for _, filter := range params.QuickFilter {
		switch filter {
		case "initiated":
			query = query.Where("initiate_statuses.type = ?", "initiated")
		case "approved":
			query = query.Where("initiate_statuses.type = ?", "approved")
		case "rejected":
			query = query.Where("initiate_statuses.type = ?", "rejected")
		case "modify":
			query = query.Where("initiate_types.type = ?", "modify")
		case "archive":
			query = query.Where("initiate_types.type = ?", "archive")
		case "restore":
			query = query.Where("initiate_types.type = ?", "restore")
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
        (
            document_change_request_initiates.reason_for_change LIKE ? OR
            initiate_types.title_en LIKE ? OR
            initiate_types.title_ar LIKE ? OR
            initiate_statuses.title_en LIKE ? OR
            initiate_statuses.title_ar LIKE ? OR
            created_by_user.first_name_en LIKE ? OR
            created_by_user.first_name_ar LIKE ?
        )`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
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
	response := &dtos.DocumentChangeRequestInitiateListResponseDTO{
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

func (repo *DocumentChangeRequestInitiateRepository) GetByID(id uint) (*models.DocumentChangeRequestInitiate, error) {
	var changeRequest models.DocumentChangeRequestInitiate
	result := repo.db.First(&changeRequest, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &changeRequest, nil
}

func (repo *DocumentChangeRequestInitiateRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestInitiate, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var changeRequest models.DocumentChangeRequestInitiate
	result := db.Preload("Documents").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&changeRequest, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &changeRequest, nil
}

func (repo *DocumentChangeRequestInitiateRepository) Create(changeRequest *models.DocumentChangeRequestInitiate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(changeRequest).Error; err != nil {
		return 0, err
	}
	return changeRequest.ID, nil
}

func (repo *DocumentChangeRequestInitiateRepository) Update(changeRequest *models.DocumentChangeRequestInitiate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(changeRequest).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentChangeRequestInitiateRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.DocumentChangeRequestInitiate{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentChangeRequestInitiateRepository) GetDocumentByID(id uint) (*models.DocumentChangeRequestInitiateDocument, error) {
	var document models.DocumentChangeRequestInitiateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *DocumentChangeRequestInitiateRepository) CreateDocuments(docs []*models.DocumentChangeRequestInitiateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentChangeRequestInitiateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", docIds).Delete(&models.DocumentChangeRequestInitiateDocument{}).Error; err != nil {
		return err
	}
	return nil
}
