package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditMeetingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *auditDTOs.AuditMeetingListQueryParamsDTO) (*auditDTOs.AuditMeetingListResponseDTO, error)
	GetByID(id uint) (*models.AuditMeeting, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditMeeting, error)
	Create(auditMeeting *models.AuditMeeting, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditMeeting *models.AuditMeeting, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	// Participant methods
	CreateParticipants(participants []*models.AuditMeetingParticipant, tx ...utility.DatabaseTransaction) error
	DeleteParticipant(auditMeetingId uint, participantIds []*uint, tx ...utility.DatabaseTransaction) error
}

type AuditMeetingRepository struct {
	db *gorm.DB
}

func NewAuditMeetingRepository(db *gorm.DB) *AuditMeetingRepository {
	return &AuditMeetingRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditMeetingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditMeetingRepository) GetAll(params *auditDTOs.AuditMeetingListQueryParamsDTO) (*auditDTOs.AuditMeetingListResponseDTO, error) {
	var auditMeetings []*auditDTOs.AuditMeetingListDTO
	var totalCount int64

	query := repo.db.Table("audit_meetings").
		Select(`
			audit_meetings.*,
			CONVERT(VARCHAR, audit_meetings.date, 23) AS date,
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
			LEFT JOIN audits ON audits.id = audit_meetings.audit_id
			LEFT JOIN users AS created_user ON created_user.id = audit_meetings.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_meetings.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply other filters as needed (e.g., date filters)
	if params.StartDate != "" {
		query = query.Where("audit_meetings.date >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("audit_meetings.date <= ?", params.EndDate)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audit_meetings.title LIKE ? OR
			audits.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audit_meetings.date" // Default sort by date if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditMeetings).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditMeetingListResponseDTO{
		Items: auditMeetings,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditMeetingRepository) GetByID(id uint) (*models.AuditMeeting, error) {
	var auditMeeting *models.AuditMeeting

	result := repo.db.First(&auditMeeting, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditMeeting, nil
}

func (repo *AuditMeetingRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditMeeting, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditMeeting *models.AuditMeeting

	result := db.
		Preload("Audit").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Participants").
		Preload("Participants.User").
		First(&auditMeeting, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditMeeting, nil
}

func (repo *AuditMeetingRepository) Create(auditMeeting *models.AuditMeeting, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditMeeting).Error; err != nil {
		return 0, err
	}

	return auditMeeting.ID, nil
}

func (repo *AuditMeetingRepository) Update(auditMeeting *models.AuditMeeting, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditMeeting).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditMeetingRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the AuditMeeting
	if err := db.Delete(&models.AuditMeeting{}, id).Error; err != nil {
		return err
	}

	return nil
}

// Participant methods

func (repo *AuditMeetingRepository) CreateParticipants(participants []*models.AuditMeetingParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&participants).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditMeetingRepository) DeleteParticipant(auditMeetingId uint, participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("audit_meeting_id = ? AND user_id IN (?)", auditMeetingId, participantIds).Delete(&models.AuditMeetingParticipant{}).Error; err != nil {
		return err
	}

	return nil
}
