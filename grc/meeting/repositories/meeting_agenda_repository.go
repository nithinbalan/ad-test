package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type MeetingAgendaRepositoryInterface interface {
	GetAll(meetingID uint, currentUserID uint, params *meetingDTOs.MeetingAgendaListQueryParamsDTO) (*meetingDTOs.MeetingAgendaListResponseDTO, error)
	GetByID(id uint) (*models.MeetingAgenda, error)
	GetFullDetailsByID(id uint) (*models.MeetingAgenda, error)
	Create(meeting *models.MeetingAgenda) (uint, error)
	Update(meeting *models.MeetingAgenda) error
	Delete(id uint) error
}

type MeetingAgendaRepository struct {
	db *gorm.DB
}

func NewMeetingAgendaRepository(db *gorm.DB) *MeetingAgendaRepository {
	return &MeetingAgendaRepository{
		db: db,
	}
}

func (repo *MeetingAgendaRepository) GetAll(meetingID uint, currentUserID uint, params *meetingDTOs.MeetingAgendaListQueryParamsDTO) (*meetingDTOs.MeetingAgendaListResponseDTO, error) {
	var meetingAgendas []*meetingDTOs.MeetingAgendaListDTO
	var totalCount int64

	query := repo.db.Table("meeting_agendas").
		Select(`
			meeting_agendas.*,
			meetings.id as meeting_id,
			meetings.title as meeting_title,
			meetings.location as meeting_location,
			meetings.date as meeting_date,
			meetings.start_time as meeting_start_time,
			meetings.end_time as meeting_end_time,
        	created_user.first_name_en AS created_user_first_name_en,
        	created_user.first_name_ar AS created_user_first_name_ar,
        	created_user.last_name_en AS created_user_last_name_en,
        	created_user.last_name_ar AS created_user_last_name_ar,   
        	created_user.initial AS created_user_initial,
        	created_user.color AS created_user_color,
        	created_user.email AS created_user_email,
        	created_user.profile_image_path AS created_user_profile_image_path,   
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
			LEFT JOIN meetings ON meeting_agendas.meeting_id = meetings.id
            LEFT JOIN users AS created_user ON created_user.id = meeting_agendas.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = meeting_agendas.updated_by_id
        `).Where(`meeting_agendas.meeting_id`, meetingID)

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
			(meetings.title LIKE ? OR
			meetings.location LIKE ? OR
			meeting_agendas.title LIKE ? OR
			meeting_agendas.time_from LIKE ? OR
			meeting_agendas.time_to LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&meetingAgendas).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &meetingDTOs.MeetingAgendaListResponseDTO{
		Items: meetingAgendas,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *MeetingAgendaRepository) GetByID(id uint) (*models.MeetingAgenda, error) {
	var meetingAgenda *models.MeetingAgenda

	result := repo.db.First(&meetingAgenda, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return meetingAgenda, nil

}

func (repo *MeetingAgendaRepository) GetFullDetailsByID(id uint) (*models.MeetingAgenda, error) {
	var meetingAgenda *models.MeetingAgenda

	result := repo.db.
		Preload("MeetingMinutes").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&meetingAgenda, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return meetingAgenda, nil
}

func (repo *MeetingAgendaRepository) Create(meetingAgenda *models.MeetingAgenda) (uint, error) {
	if err := repo.db.Create(meetingAgenda).Error; err != nil {
		return 0, err
	}

	return meetingAgenda.ID, nil
}

func (repo *MeetingAgendaRepository) Update(meetingAgenda *models.MeetingAgenda) error {
	if err := repo.db.Save(meetingAgenda).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingAgendaRepository) Delete(id uint) error {
	// Delete the MeetingAgenda
	if err := repo.db.Delete(&models.MeetingAgenda{}, id).Error; err != nil {
		return err
	}

	return nil
}
