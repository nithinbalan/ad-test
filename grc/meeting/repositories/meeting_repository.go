package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type MeetingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *meetingDTOs.MeetingListQueryParamsDTO) (*meetingDTOs.MeetingListResponseDTO, error)
	GetByID(id uint) (*models.Meeting, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Meeting, error)
	Create(meeting *models.Meeting, tx ...utility.DatabaseTransaction) (uint, error)
	Update(meeting *models.Meeting, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateParticipants(participants []*models.MeetingParticipant, tx ...utility.DatabaseTransaction) error
	DeleteParticipants(participantIDs []uint, tx ...utility.DatabaseTransaction) error
}

type MeetingRepository struct {
	db *gorm.DB
}

func NewMeetingRepository(db *gorm.DB) *MeetingRepository {
	return &MeetingRepository{
		db: db,
	}
}

func (repo *MeetingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *MeetingRepository) GetAll(params *meetingDTOs.MeetingListQueryParamsDTO) (*meetingDTOs.MeetingListResponseDTO, error) {
	var meetings []*meetingDTOs.MeetingListDTO
	var totalCount int64

	query := repo.db.Table("meeting").
		Select(`
            meeting.id,
            meeting.reference_number,
            meeting.title,
            meeting.description,
            meeting.start_time,
            meeting.end_time,
            meeting.created_at,
            meeting.updated_at
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
			(meeting.title LIKE ? OR
			meeting.description LIKE ? OR
			meeting.reference_number LIKE ?)`,
			keyword, keyword, keyword)
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&meetings).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &meetingDTOs.MeetingListResponseDTO{
		Items: meetings,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *MeetingRepository) GetByID(id uint) (*models.Meeting, error) {
	var meeting *models.Meeting

	result := repo.db.First(&meeting, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return meeting, nil

}

func (repo *MeetingRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Meeting, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var meeting *models.Meeting

	result := db.
		Preload("Meeting").
		Preload("MeetingType").
		Preload("MeetingFrequency").
		Preload("Organizer").
		Preload("MeetingStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Participants").
		Preload("Participants.User").
		First(&meeting, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return meeting, nil
}

func (repo *MeetingRepository) Create(meeting *models.Meeting, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(meeting).Error; err != nil {
		return 0, err
	}

	return meeting.ID, nil
}

func (repo *MeetingRepository) Update(meeting *models.Meeting, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(meeting).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Meeting
	if err := db.Delete(&models.Meeting{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingRepository) CreateParticipants(participants []*models.MeetingParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(participants).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingRepository) DeleteParticipants(participantIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete the participants
	if err := db.Where("id IN (?)", participantIDs).Delete(&models.MeetingParticipant{}).Error; err != nil {
		return err
	}

	return nil
}
