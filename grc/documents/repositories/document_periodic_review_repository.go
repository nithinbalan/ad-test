package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentPeriodicReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.DocumentPeriodicReviewQueryParamsDTO, documentId uint) (*dtos.DocumentPeriodicReviewListResponseDTO, error)
	GetByID(id uint) (*models.DocumentReview, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentReview, error)
	Create(review *models.DocumentReview, tx ...utility.DatabaseTransaction) (uint, error)
	Update(review *models.DocumentReview, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	GetDocumentByID(id uint) (*models.DocumentReviewSupportDocument, error)
	CreateDocuments(docs []*models.DocumentReviewSupportDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type DocumentPeriodicReviewRepository struct {
	db *gorm.DB
}

func NewDocumentPeriodicReviewRepository(db *gorm.DB) *DocumentPeriodicReviewRepository {
	return &DocumentPeriodicReviewRepository{
		db: db,
	}
}

func (repo *DocumentPeriodicReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentPeriodicReviewRepository) GetAll(params *dtos.DocumentPeriodicReviewQueryParamsDTO, documentId uint) (*dtos.DocumentPeriodicReviewListResponseDTO, error) {
	var reviews []*dtos.DocumentPeriodicReviewListDTO
	var totalCount int64

	query := repo.db.Table("document_reviews").
		Select(`
			document_reviews.id,
			document_reviews.document_id,
			documents.title_en AS document_title,
			document_reviews.review_date,
			document_reviews.is_change_required,
			document_reviews.description,
			created_by_user.id AS created_by_user_id,
			created_by_user.first_name_en AS created_by_user_first_name_en,
			created_by_user.first_name_ar AS created_by_user_first_name_ar,
			created_by_user.last_name_en AS created_by_user_last_name_en,
			created_by_user.last_name_ar AS created_by_user_last_name_ar,
			created_by_user.email AS created_by_user_email,
			created_by_user.initial AS created_by_user_initial,
	        created_by_user.color AS created_by_user_color,
			created_by_user.profile_image_path AS created_by_user_profile_image_path,
			updated_by_user.id AS updated_by_user_id,
			updated_by_user.first_name_en AS updated_by_user_first_name_en,
			updated_by_user.first_name_ar AS updated_by_user_first_name_ar,
			updated_by_user.last_name_en AS updated_by_user_last_name_en,
			updated_by_user.last_name_ar AS updated_by_user_last_name_ar,
			updated_by_user.email AS updated_by_user_email,
			updated_by_user.initial AS updated_by_user_initial,
	        updated_by_user.color AS updated_by_user_color,
			updated_by_user.profile_image_path AS updated_by_user_profile_image_path,
			document_reviews.created_at,
			document_reviews.updated_at
		`).
		Joins("JOIN documents ON documents.id = document_reviews.document_id").
		Joins("LEFT JOIN users AS created_by_user ON created_by_user.id = document_reviews.created_by_id").
		Joins("LEFT JOIN users AS updated_by_user ON updated_by_user.id = document_reviews.updated_by_id")

	// Add documentId filter
	query = query.Where("document_reviews.document_id = ?", documentId)

	// Apply dynamic filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters if applicable
	for _, filter := range params.QuickFilter {
		if filter == "is_change_required" {
			query = query.Where("document_reviews.is_change_required = ?", true)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
        (
            document_reviews.description LIKE ? OR
            created_by_user.first_name_en LIKE ? OR
            created_by_user.first_name_ar LIKE ? OR
            updated_by_user.first_name_en LIKE ? OR
            updated_by_user.first_name_ar LIKE ?
        )`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Sort by the specified column and order
	query = query.Order(params.SortBy + " " + params.Sort)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Pagination
	offset := (params.Page - 1) * params.PageSize
	err := query.Offset(offset).Limit(params.PageSize).Find(&reviews).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	// Build the response
	response := &dtos.DocumentPeriodicReviewListResponseDTO{
		Items: reviews,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *DocumentPeriodicReviewRepository) GetByID(id uint) (*models.DocumentReview, error) {
	var review models.DocumentReview
	result := repo.db.First(&review, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &review, nil
}

func (repo *DocumentPeriodicReviewRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentReview, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var review models.DocumentReview
	result := db.Preload("SupportDocuments").
		Preload("DocumentChangeRequest").
		Preload("DocumentChangeRequest.DocumentChangeRequestStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&review, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &review, nil
}

func (repo *DocumentPeriodicReviewRepository) Create(review *models.DocumentReview, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(review).Error; err != nil {
		return 0, err
	}
	return review.ID, nil
}

func (repo *DocumentPeriodicReviewRepository) Update(review *models.DocumentReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(review).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentPeriodicReviewRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.DocumentReview{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentPeriodicReviewRepository) GetDocumentByID(id uint) (*models.DocumentReviewSupportDocument, error) {
	var document models.DocumentReviewSupportDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *DocumentPeriodicReviewRepository) CreateDocuments(docs []*models.DocumentReviewSupportDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentPeriodicReviewRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", docIds).Delete(&models.DocumentReviewSupportDocument{}).Error; err != nil {
		return err
	}
	return nil
}
