package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type MeetingMinuteRepositoryInterface interface {
	GetAll(currentUserID uint, params *meetingDTOs.MeetingMinuteListQueryParamsDTO) (*meetingDTOs.MeetingMinuteListResponseDTO, error)
	GetByID(id uint) (*models.MeetingMinute, error)
	GetFullDetailsByID(id uint) (*models.MeetingMinute, error)
	Create(meetingMinute *models.MeetingMinute) (uint, error)
	Update(meetingMinute *models.MeetingMinute) error
	Delete(id uint) error
}

type MeetingMinuteRepository struct {
	db *gorm.DB
}

func NewMeetingMinuteRepository(db *gorm.DB) *MeetingMinuteRepository {
	return &MeetingMinuteRepository{
		db: db,
	}
}

func (repo *MeetingMinuteRepository) GetAll(currentUserID uint, params *meetingDTOs.MeetingMinuteListQueryParamsDTO) (*meetingDTOs.MeetingMinuteListResponseDTO, error) {
	var meetingMinutes []*meetingDTOs.MeetingMinuteListDTO
	var totalCount int64

	query := repo.db.Table("meeting_minutes").
		Select(`
			meeting_minutes.*,
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
            LEFT JOIN users AS created_user ON created_user.id = meeting_minutes.created_by
            LEFT JOIN users AS updated_user ON updated_user.id = meeting_minutes.updated_by
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
			(meeting_minutes.minutes LIKE ?)`,
			keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "minutes" // Default sort by minutes if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&meetingMinutes).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &meetingDTOs.MeetingMinuteListResponseDTO{
		Items: meetingMinutes,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *MeetingMinuteRepository) GetByID(id uint) (*models.MeetingMinute, error) {
	var meetingMinute *models.MeetingMinute

	result := repo.db.First(&meetingMinute, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return meetingMinute, nil

}

func (repo *MeetingMinuteRepository) GetFullDetailsByID(id uint) (*models.MeetingMinute, error) {
	var meetingMinute *models.MeetingMinute

	result := repo.db.
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&meetingMinute, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return meetingMinute, nil
}

func (repo *MeetingMinuteRepository) Create(meetingMinute *models.MeetingMinute) (uint, error) {
	if err := repo.db.Create(meetingMinute).Error; err != nil {
		return 0, err
	}

	return meetingMinute.ID, nil
}

func (repo *MeetingMinuteRepository) Update(meetingMinute *models.MeetingMinute) error {
	if err := repo.db.Save(meetingMinute).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingMinuteRepository) Delete(id uint) error {
	// Delete the MeetingMinute
	if err := repo.db.Delete(&models.MeetingMinute{}, id).Error; err != nil {
		return err
	}

	return nil
}
