package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	logDTOs "jethur/grc/logs/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ActivityLogRepositoryInteface interface {
	GetAll(params *logDTOs.ActivityLogListQueryParamsDTO) (*logDTOs.ActivityLogListResponseDTO, error)
	GetFullDetailsByID(id uint) (*models.ActivityLog, error)
	SaveActivityLogs(logs []*models.ActivityLog, tx ...*gorm.DB) error
}

type ActivityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{
		db: db,
	}
}

func (repo *ActivityLogRepository) GetAll(params *logDTOs.ActivityLogListQueryParamsDTO) (*logDTOs.ActivityLogListResponseDTO, error) {
	var totalCount int64
	var activityLogs []*logDTOs.ActivityLogListDTO

	// Build query
	query := repo.db.Table("activity_logs").
		Select(`
			activity_logs.*,
			users.id AS user_id,
			users.first_name_en AS user_first_name_en,
			users.first_name_ar AS user_first_name_ar,
			users.last_name_en AS user_last_name_en,
			users.last_name_ar AS user_last_name_ar,
			users.email AS user_email,
			users.profile_image_path AS user_profile_image_path,
			users.initial AS user_initial,
			users.color AS user_color
		`).
		Joins(`
			LEFT JOIN users ON users.id = activity_logs.user_id
		`)

	// Apply filters for keyword search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("table_name LIKE ? OR changed_data LIKE ?", keyword, keyword)
	}

	// Count total records for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	if params.SortBy == "" {
		params.SortBy = "created_at" // Default sort column
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default sort order
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	err := query.Offset(offset).Limit(params.PageSize).Find(&activityLogs).Error
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	// Build response
	response := &logDTOs.ActivityLogListResponseDTO{
		Items: activityLogs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ActivityLogRepository) GetFullDetailsByID(id uint) (*models.ActivityLog, error) {
	var activityLog *models.ActivityLog

	result := repo.db.
		Preload("User").
		First(&activityLog, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return activityLog, nil
}

func (repo *ActivityLogRepository) SaveActivityLogs(logs []*models.ActivityLog, tx ...*gorm.DB) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Create(&logs).Error
}
