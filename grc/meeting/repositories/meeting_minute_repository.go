package repositories

import (
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type MeetingMinuteRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(meetingID uint, currentUserID uint, params *meetingDTOs.MeetingMinuteListQueryParamsDTO) ([]*meetingDTOs.MeetingMinuteListDTO, error)
	GetByID(id uint) (*models.MeetingMinute, error)
	GetFullDetailsByID(id uint) (*models.MeetingMinute, error)
	Create(meetingMinute *models.MeetingMinute) (uint, error)
	Update(meetingMinute *models.MeetingMinute, tx ...utility.DatabaseTransaction) error
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

// BeginTransaction starts a new transaction.
func (repo *MeetingMinuteRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *MeetingMinuteRepository) GetAll(meetingID uint, currentUserID uint, params *meetingDTOs.MeetingMinuteListQueryParamsDTO) ([]*meetingDTOs.MeetingMinuteListDTO, error) {
	var meetingMinutes []*meetingDTOs.MeetingMinuteListDTO
	var totalCount int64

	query := repo.db.Table("meeting_minutes").
		Select(`
		meeting_minutes.*,
		meeting_agendas.title AS meeting_agenda_title,
		meetings.title AS meeting_title,
		meeting_minute_statuses.title_en AS meeting_minute_status_title_en,
		meeting_minute_statuses.title_ar AS meeting_minute_status_title_ar,
		meeting_minute_statuses.type AS meeting_minute_status_type,
		meeting_minute_statuses.color_code AS meeting_minute_status_color_code,
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
		LEFT JOIN meeting_agendas ON meeting_minutes.meeting_agenda_id = meeting_agendas.id
		LEFT JOIN meeting_minute_statuses ON meeting_minutes.meeting_minute_status_id = meeting_minute_statuses.id
		LEFT JOIN meetings ON meeting_minutes.meeting_id = meetings.id
		LEFT JOIN users AS created_user ON created_user.id = meeting_minutes.created_by_id
		LEFT JOIN users AS updated_user ON updated_user.id = meeting_minutes.updated_by_id
	`).
		Where(`meetings.id = ?`, meetingID)

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
			(meeting_minutes.minutes LIKE ? OR
			(meetings.title_en LIKE ? OR
			(meeting_minute_statuses.title_en LIKE ? OR
			(meeting_minute_statuses.title_ar LIKE ? OR
			meeting_agendas.title LIKE ?)`,
			keyword,keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "minutes" // Default sort by minutes if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&meetingMinutes).Error; err != nil {
		return nil, err
	}

	return meetingMinutes, nil
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
		Preload("MeetingAgenda").
		Preload("Meeting").
		Preload("MeetingMinuteStatus").
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

func (repo *MeetingMinuteRepository) Update(meetingMinute *models.MeetingMinute, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(meetingMinute).Error; err != nil {
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
