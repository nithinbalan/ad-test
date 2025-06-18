package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategicObjectivePeriodicReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.StrategicObjectivePeriodicReviewQueryParamsDTO, strategicObjectiveID uint) (*dtos.StrategicObjectivePeriodicReviewListResponseDTO, error)
	GetByID(id uint) (*models.StrategicObjectiveReview, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategicObjectiveReview, error)
	Create(review *models.StrategicObjectiveReview, tx ...utility.DatabaseTransaction) (uint, error)
	Update(review *models.StrategicObjectiveReview, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	GetDocumentByID(id uint) (*models.StrategicObjectiveReviewDocument, error)
	CreateDocuments(docs []*models.StrategicObjectiveReviewDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type StrategicObjectivePeriodicReviewRepository struct {
	db *gorm.DB
}

func NewStrategicObjectivePeriodicReviewRepository(db *gorm.DB) *StrategicObjectivePeriodicReviewRepository {
	return &StrategicObjectivePeriodicReviewRepository{
		db: db,
	}
}

func (repo *StrategicObjectivePeriodicReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategicObjectivePeriodicReviewRepository) GetAll(params *dtos.StrategicObjectivePeriodicReviewQueryParamsDTO, strategicObjectiveID uint) (*dtos.StrategicObjectivePeriodicReviewListResponseDTO, error) {
	var reviews []*dtos.StrategicObjectivePeriodicReviewListDTO
	var totalCount int64

	query := repo.db.Table("strategic_objective_reviews").
		Select(`
			strategic_objective_reviews.id,
			strategic_objective_reviews.strategic_objective_id,
			strategic_objectives.title_en AS strategic_objective_review_title_en,
			strategic_objectives.title_ar AS strategic_objective_review_title_ar,
			strategic_objective_reviews.review_date,
			strategic_objective_reviews.description,
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
			strategic_objective_reviews.created_at,
			strategic_objective_reviews.updated_at
		`).
		Joins("JOIN strategic_objectives ON strategic_objectives.id = strategic_objective_reviews.strategic_objective_id").
		Joins("LEFT JOIN users AS created_by_user ON created_by_user.id = strategic_objective_reviews.created_by_id").
		Joins("LEFT JOIN users AS updated_by_user ON updated_by_user.id = strategic_objective_reviews.updated_by_id")

	// Add strategicObjectiveID filter
	query = query.Where("strategic_objective_reviews.strategic_objective_id = ?", strategicObjectiveID)

	// Apply dynamic filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
        (
            strategic_objective_reviews.description LIKE ? OR
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

	// Fetch supportive strategic_objective_reviews for each review
	for _, review := range reviews {
		var supportDocs []*dtos.StrategicObjectivePeriodicReviewSupportDocumentDTO
		err := repo.db.Table("strategic_objective_review_documents").
			Select("id, url, name, size, token, original_name, type").
			Where("strategic_objective_review_id = ?", review.ID).
			Find(&supportDocs).Error
		if err != nil {
			return nil, err
		}
		review.SupportDocuments = supportDocs
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	// Build the response
	response := &dtos.StrategicObjectivePeriodicReviewListResponseDTO{
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

func (repo *StrategicObjectivePeriodicReviewRepository) GetByID(id uint) (*models.StrategicObjectiveReview, error) {
	var review models.StrategicObjectiveReview
	result := repo.db.First(&review, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &review, nil
}

func (repo *StrategicObjectivePeriodicReviewRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategicObjectiveReview, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var review models.StrategicObjectiveReview
	result := db.Preload("StrategicObjective").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&review, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &review, nil
}

func (repo *StrategicObjectivePeriodicReviewRepository) Create(review *models.StrategicObjectiveReview, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(review).Error; err != nil {
		return 0, err
	}
	return review.ID, nil
}

func (repo *StrategicObjectivePeriodicReviewRepository) Update(review *models.StrategicObjectiveReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(review).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StrategicObjectivePeriodicReviewRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.StrategicObjectiveReview{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StrategicObjectivePeriodicReviewRepository) GetDocumentByID(id uint) (*models.StrategicObjectiveReviewDocument, error) {
	var document models.StrategicObjectiveReviewDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *StrategicObjectivePeriodicReviewRepository) CreateDocuments(docs []*models.StrategicObjectiveReviewDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StrategicObjectivePeriodicReviewRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", docIds).Delete(&models.StrategicObjectiveReviewDocument{}).Error; err != nil {
		return err
	}
	return nil
}
