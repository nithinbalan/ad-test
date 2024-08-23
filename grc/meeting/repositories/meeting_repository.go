package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type MeetingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(currentUserID uint, params *meetingDTOs.MeetingListQueryParamsDTO) (*meetingDTOs.MeetingListResponseDTO, error)
	GetByID(id uint) (*models.Meeting, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Meeting, error)
	Create(meeting *models.Meeting, tx ...utility.DatabaseTransaction) (uint, error)
	Update(meeting *models.Meeting, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateParticipants(participants []*models.MeetingParticipant, tx ...utility.DatabaseTransaction) error
	DeleteParticipants(participantIDs []uint, tx ...utility.DatabaseTransaction) error
	CreateMeetingActionPlans(actionPlans []*models.MeetingActionPlan, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	CreateMeetingComplianceRequirements(complianceRequirements []*models.MeetingComplianceRequirement, tx ...utility.DatabaseTransaction) error
	CreateMeetingControls(controls []*models.MeetingControl, tx ...utility.DatabaseTransaction) error
	DeleteMeetingControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteMeetingComplianceRequirements(requirementIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateMeetingRisks(risks []*models.MeetingRisk, tx ...utility.DatabaseTransaction) error
	DeleteMeetingRisks(riskIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateMeetingProcesses(processes []*models.MeetingProcess, tx ...utility.DatabaseTransaction) error
	DeleteMeetingProcesses(processIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateMeetingIncidents(incidents []*models.MeetingIncident, tx ...utility.DatabaseTransaction) error
	DeleteMeetingIncidents(incidentIds []*uint, tx ...utility.DatabaseTransaction) error

	GetAllCounts(meetingID uint) (*meetingDTOs.MeetingsCountResponseDTO, error)
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

func (repo *MeetingRepository) GetAll(currentUserID uint, params *meetingDTOs.MeetingListQueryParamsDTO) (*meetingDTOs.MeetingListResponseDTO, error) {
	var meetings []*meetingDTOs.MeetingListDTO
	var totalCount int64

	query := repo.db.Table("meetings").
		Select(`
			meetings.id,
			meetings.meeting_id,
			meeting_types.title_en AS meeting_type_title_en,
			meeting_types.title_ar AS meeting_type_title_er,
			meeting_frequencies.title_en AS meeting_frequency_title_en,
			meeting_frequencies.title_ar AS meeting_frequency_title_ar,
			meetings.reference_number,
			meetings.title,
			meetings.description,
			meetings.location,
			CONVERT(VARCHAR, meetings.date, 23) AS date,
			meetings.start_time,
			meetings.end_time,
			users.id AS organizer_user_id,
			users.first_name_en AS organizer_user_first_name_en,
			users.first_name_ar AS organizer_user_first_name_ar,
			users.last_name_en AS organizer_user_last_name_en,
			users.last_name_ar AS organizer_user_last_name_ar,
			users.initial AS organizer_user_initial,
        	users.color AS organizer_user_color,
        	users.email AS organizer_user_email,
        	users.profile_image_path AS organizer_user_profile_image_path,
			meetings.created_at,
			meetings.updated_at,
			meetings.is_public,
			meeting_statuses.title_en AS meeting_status_title_en,
			meeting_statuses.title_ar AS meeting_status_title_ar,
			meeting_statuses.color_code AS meeting_status_color_code,
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
			LEFT JOIN meeting_types ON meetings.meeting_type_id = meeting_types.id
			LEFT JOIN meeting_frequencies ON meetings.meeting_frequency_id = meeting_frequencies.id
			LEFT JOIN users ON meetings.organizer_id = users.id
			LEFT JOIN meeting_statuses ON meetings.meeting_status_id = meeting_statuses.id
            LEFT JOIN users AS created_user ON created_user.id = meetings.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = meetings.updated_by_id
        `).Where("(meetings.is_public = ? OR meetings.organizer_id = ? OR meetings.id IN (SELECT meeting_id FROM meeting_participants WHERE user_id = ?))", true, currentUserID, currentUserID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "today":
			query = query.Where("CONVERT(date, meetings.date) = CONVERT(date, GETDATE())")
		case "this_week":
			query = query.Where("DATEPART(week, meetings.date) = DATEPART(week, GETDATE())")
		case "this_month":
			query = query.Where("DATEPART(month, meetings.date) = DATEPART(month, GETDATE())")
		case "scheduled":
			query = query.Where("meeting_statuses.type = ?", "scheduled")
		case "completed":
			query = query.Where("meeting_statuses.type = ?", "completed")
		case "cancelled":
			query = query.Where("meeting_statuses.type = ?", "cancelled")
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
			meetings.reference_number LIKE ? OR
			meetings.location LIKE ? OR
			meeting_types.title_en LIKE ? OR
			meeting_types.title_ar LIKE ? OR
			meeting_frequencies.title_en LIKE ? OR
			meeting_frequencies.title_ar LIKE ? OR
			users.first_name_en LIKE ? OR
			users.last_name_en LIKE ? OR
			users.first_name_ar LIKE ? OR
			users.last_name_ar LIKE ? OR
			meeting_statuses.title_en LIKE ? OR
			meeting_statuses.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
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
		Preload("Meeting.MeetingType").
		Preload("Meeting.MeetingFrequency").
		Preload("Meeting.Organizer").
		Preload("Meeting.MeetingStatus").
		Preload("Meeting.CreatedByUser").
		Preload("MeetingType").
		Preload("MeetingFrequency").
		Preload("Organizer").
		Preload("MeetingStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Agendas").
		Preload("Agendas.MeetingMinutes").
		Preload("Participants").
		Preload("Participants.User").
		Preload("Controls").
		Preload("ComplianceRequirements").
		Preload("Risks").
		Preload("Processes").
		Preload("Incidents").
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

func (repo *MeetingRepository) CreateMeetingActionPlans(actionPlans []*models.MeetingActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(actionPlans).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingRepository) CreateMeetingComplianceRequirements(complianceRequirements []*models.MeetingComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&complianceRequirements).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MeetingRepository) CreateMeetingControls(controls []*models.MeetingControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MeetingRepository) DeleteMeetingControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete controls by IDs
	if err := db.Where("id IN (?)", controlIds).Delete(&models.MeetingControl{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingRepository) DeleteMeetingComplianceRequirements(requirementIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete compliance requirements by IDs
	if err := db.Where("id IN (?)", requirementIds).Delete(&models.MeetingComplianceRequirement{}).Error; err != nil {
		return err
	}

	return nil
}
func (repo *MeetingRepository) CreateMeetingRisks(risks []*models.MeetingRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&risks).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MeetingRepository) DeleteMeetingRisks(riskIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete risks by IDs
	if err := db.Where("id IN (?)", riskIds).Delete(&models.MeetingRisk{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingRepository) CreateMeetingProcesses(processes []*models.MeetingProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&processes).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MeetingRepository) DeleteMeetingProcesses(processIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete process by IDs
	if err := db.Where("id IN (?)", processIds).Delete(&models.MeetingProcess{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingRepository) CreateMeetingIncidents(incidents []*models.MeetingIncident, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&incidents).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MeetingRepository) DeleteMeetingIncidents(incidentIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete incidents by IDs
	if err := db.Where("id IN (?)", incidentIds).Delete(&models.MeetingIncident{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber *string
	err := repo.db.Model(&models.Meeting{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Scan(&lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != nil && len(*lastReferenceNumber) >= 4 {
		lastDigits = (*lastReferenceNumber)[len(*lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}

	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "MT", currentYear, lastReferenceID+1)
	return newReferenceNumber, nil
}

func (repo *MeetingRepository) GetAllCounts(meetingID uint) (*meetingDTOs.MeetingsCountResponseDTO, error) {
	var agendaCount, momCount, actionPlanCount int64
	var err error

	// Query for agenda count
	err = repo.db.Model(&models.MeetingAgenda{}).
		Where("meeting_id = ?", meetingID).
		Count(&agendaCount).Error
	if err != nil {
		return nil, err
	}

	// Query for minutes of meeting (MoM) count through meeting_agendas
	err = repo.db.Model(&models.MeetingMinute{}).
		Joins("JOIN meeting_agendas ON meeting_agendas.id = meeting_minutes.meeting_agenda_id").
		Where("meeting_agendas.meeting_id = ?", meetingID).
		Count(&momCount).Error
	if err != nil {
		return nil, err
	}

	// Query for action plan count
	err = repo.db.Model(&models.MeetingActionPlan{}).
		Where("meeting_id = ?", meetingID).
		Count(&actionPlanCount).Error
	if err != nil {
		return nil, err
	}

	return &meetingDTOs.MeetingsCountResponseDTO{
		Agenda:     int(agendaCount),
		Mom:        int(momCount),
		ActionPlan: int(actionPlanCount),
	}, nil
}
