package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentPostReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(incidentID uint, params *incidentDTOs.IncidentPostReviewListQueryParamsDTO) (*incidentDTOs.IncidentPostReviewListResponseDTO, error)
	GetByID(id uint) (*models.IncidentPostReview, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentPostReview, error)
	Create(incidentPostReview *models.IncidentPostReview, tx ...utility.DatabaseTransaction) (uint, error)
	Update(incidentPostReview *models.IncidentPostReview, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateLessons(lessons []*models.IncidentPostReviewLesson, tx ...utility.DatabaseTransaction) error
	DeleteLessons(lessonsIDs []*uint, tx ...utility.DatabaseTransaction) error
	CreateParticipants(users []*models.IncidentPostReviewUser, tx ...utility.DatabaseTransaction) error
	DeleteParticipants(userIDs []*uint, tx ...utility.DatabaseTransaction) error
	CreateRecommendations(recommendations []*models.IncidentPostReviewRecommendation, tx ...utility.DatabaseTransaction) error
	DeleteRecommendations(recommendationIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type IncidentPostReviewRepository struct {
	db *gorm.DB
}

func NewIncidentPostReviewRepository(db *gorm.DB) *IncidentPostReviewRepository {
	return &IncidentPostReviewRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *IncidentPostReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentPostReviewRepository) GetAll(incidentID uint, params *incidentDTOs.IncidentPostReviewListQueryParamsDTO) (*incidentDTOs.IncidentPostReviewListResponseDTO, error) {

	var incidentPostReviews []*incidentDTOs.IncidentPostReviewListDTO
	var totalCount int64

	query := repo.db.Table("incident_post_reviews").
		Select(`
        incident_post_reviews.*,
        created_user.first_name_en AS created_user_first_name_en,
        created_user.first_name_ar AS created_user_first_name_ar,
        created_user.last_name_en AS created_user_last_name_en,
        created_user.last_name_ar AS created_user_last_name_ar,   
        updated_user.first_name_en AS updated_user_first_name_en,
        updated_user.first_name_ar AS updated_user_first_name_ar,
        updated_user.last_name_en AS updated_user_last_name_en,
        updated_user.last_name_ar AS updated_user_last_name_ar
    `).
		Joins(`
		LEFT JOIN users AS created_user ON created_user.id = incident_post_reviews.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = incident_post_reviews.updated_by_id
    `)

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
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (incident_post_reviews.summary LIKE ?)`,
			keyword)
	}
	// Conditionally apply incident ID filter
	if incidentID != 0 {
		query = query.Where("incident_post_reviews.incident_register_id = ?", incidentID)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "summary" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentPostReviews).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &incidentDTOs.IncidentPostReviewListResponseDTO{
		Items: incidentPostReviews,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *IncidentPostReviewRepository) GetByID(id uint) (*models.IncidentPostReview, error) {
	var incidentPostReview *models.IncidentPostReview
	result := repo.db.First(&incidentPostReview, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentPostReview, nil
}

func (repo *IncidentPostReviewRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentPostReview, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentPostReview *models.IncidentPostReview

	result := db.
		Preload("IncidentRegister").
		Preload("IncidentPostReviewEffectiveness").
		Preload("Recommendations").
		Preload("Users.User").
		Preload("Lessons").
		First(&incidentPostReview, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentPostReview, nil
}

func (repo *IncidentPostReviewRepository) Create(incidentPostReview *models.IncidentPostReview, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentPostReview).Error; err != nil {
		return 0, err
	}
	return incidentPostReview.ID, nil
}

func (repo *IncidentPostReviewRepository) Update(incidentPostReview *models.IncidentPostReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentPostReview).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentPostReviewRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentPostReview *models.IncidentPostReview
	if err := db.First(&incidentPostReview, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.IncidentPostReview{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentPostReviewRepository) CreateLessons(lessons []*models.IncidentPostReviewLesson, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&lessons).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentPostReviewRepository) DeleteLessons(lessonsIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete compliance requirements by IDs
	if err := db.Where("id IN (?)", lessonsIDs).Delete(&models.IncidentPostReviewLesson{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentPostReviewRepository) CreateParticipants(users []*models.IncidentPostReviewUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentPostReviewRepository) DeleteParticipants(userIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete users by IDs
	if err := db.Where("id IN (?)", userIDs).Delete(&models.IncidentPostReviewUser{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentPostReviewRepository) CreateRecommendations(recommendations []*models.IncidentPostReviewRecommendation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&recommendations).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentPostReviewRepository) DeleteRecommendations(recommendationIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete Incident PostReview Recommendation by IDs
	if err := db.Where("id IN (?)", recommendationIDs).Delete(&models.IncidentPostReviewRecommendation{}).Error; err != nil {
		return err
	}

	return nil
}
