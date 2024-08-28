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
	GetAll(auditID uint, params *auditDTOs.AuditMeetingListQueryParamsDTO) ([]*models.AuditMeeting, *commonDTOs.PaginationInfoDTO, error)
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

func (repo *AuditMeetingRepository) GetAll(auditID uint, params *auditDTOs.AuditMeetingListQueryParamsDTO) ([]*models.AuditMeeting, *commonDTOs.PaginationInfoDTO, error) {
	var auditMeetings []*models.AuditMeeting
	var totalCount int64

	// Set the model for the query
	query := repo.db.Model(&models.AuditMeeting{}).
		Preload("Participants").
		Preload("Participants.User").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("audit_id = ?", auditID)

	// Apply date filters if provided
	if params.StartDate != "" {
		query = query.Where("date >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("date <= ?", params.EndDate)
	}

	// Apply keyword search filter if provided
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audit_meetings.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "date" // Default sort by date if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	// Apply sorting, pagination offset, and limit
	query = query.Order(params.SortBy + " " + params.Sort).
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize)

	// Execute the query and fetch the audit meetings
	err := query.Find(&auditMeetings).Error
	if err != nil {
		return nil, nil, err
	}

	// Prepare pagination info
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	pagination := &commonDTOs.PaginationInfoDTO{
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
	}

	return auditMeetings, pagination, nil
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
