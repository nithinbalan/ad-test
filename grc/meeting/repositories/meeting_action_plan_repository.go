package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type MeetingActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(meetingID uint, params *meetingDTOs.MeetingActionPlanListQueryParamsDTO) (*meetingDTOs.MeetingActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.MeetingActionPlan, error)
	GetFullDetailsByID(id uint) (*models.MeetingActionPlan, error)
	Create(meetingActionPlan *models.MeetingActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(meetingActionPlan *models.MeetingActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateInitialMeetingActionPlanStatusUpdate(meetingActionPlanStatusUpdate *models.MeetingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteMeetingActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error

	GetMeetingActionPlanSummary(meetingID uint) (openCount int, closedCount int, err error)
}

type MeetingActionPlanRepository struct {
	db *gorm.DB
}

func NewMeetingActionPlanRepository(db *gorm.DB) *MeetingActionPlanRepository {
	return &MeetingActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *MeetingActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *MeetingActionPlanRepository) GetAll(meetingID uint, params *meetingDTOs.MeetingActionPlanListQueryParamsDTO) (*meetingDTOs.MeetingActionPlanListResponseDTO, error) {

	var meetingActionPlans []*meetingDTOs.MeetingActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("meeting_action_plans").
		Select(`
        meeting_action_plans.*,
		meetings.id AS meeting_id,
		meetings.reference_number AS meeting_reference_number,
		meetings.title AS meeting_title,
		CONVERT(VARCHAR, meeting_action_plans.target_date, 23) AS target_date,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
        action_plan_statuses.color_code AS action_plan_status_color_code,
        severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
		meeting_minutes.minutes AS meeting_minutes,
		responsible_user.id AS responsible_user_id,
        responsible_user.first_name_en AS responsible_user_first_name_en,
        responsible_user.first_name_ar AS responsible_user_first_name_ar,
        responsible_user.last_name_en AS responsible_user_last_name_en,
        responsible_user.last_name_ar AS responsible_user_last_name_ar,
		responsible_user.initial AS responsible_user_initial,
        responsible_user.color AS responsible_user_color,
        responsible_user.email AS responsible_user_email,
        responsible_user.profile_image_path AS responsible_user_profile_image_path,
		review_user.id AS review_user_id,
    	review_user.first_name_en AS review_user_first_name_en,
		review_user.first_name_ar AS review_user_first_name_ar,
		review_user.last_name_en AS review_user_last_name_en,
		review_user.last_name_ar AS review_user_last_name_ar,
		review_user.initial AS review_user_initial,
        review_user.color AS review_user_color,
        review_user.email AS review_user_email,
        review_user.profile_image_path AS review_user_profile_image_path,
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
		LEFT JOIN meeting_minutes ON meeting_minutes.id = meeting_action_plans.meeting_minute_id
		LEFT JOIN meetings ON meetings.id = meeting_action_plans.meeting_id
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = meeting_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = meeting_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = meeting_action_plans.responsible_user_id
		LEFT JOIN users AS review_user ON review_user.id = meeting_action_plans.review_user_id
        LEFT JOIN users AS created_user ON created_user.id = meeting_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = meeting_action_plans.updated_by_id
    `)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply target_date filters
	if params.TargetDateFrom != "" {
		query = query.Where("target_date >= ?", params.TargetDateFrom)
	}
	if params.TargetDateTo != "" {
		query = query.Where("target_date <= ?", params.TargetDateTo)
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "overdue":
			query = query.Where("target_date < ?", time.Now())
		case "completed":
			query = query.Where("action_plan_statuses.type = ?", "completed")
		case "in_progress":
			query = query.Where("action_plan_statuses.type = ?", "in_progress")
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (meeting_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            meeting_action_plans.target_date LIKE ? OR
            meeting_minutes.minutes LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
            responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword)
	}

	if meetingID != 0 {
		query = query.Where("meeting_action_plans.meeting_id = ?", meetingID)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&meetingActionPlans).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &meetingDTOs.MeetingActionPlanListResponseDTO{
			Items: meetingActionPlans,
		}, nil
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &meetingDTOs.MeetingActionPlanListResponseDTO{
		Items: meetingActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *MeetingActionPlanRepository) GetByID(id uint) (*models.MeetingActionPlan, error) {
	var meetingActionPlan *models.MeetingActionPlan
	result := repo.db.First(&meetingActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return meetingActionPlan, nil
}

func (repo *MeetingActionPlanRepository) GetFullDetailsByID(id uint) (*models.MeetingActionPlan, error) {
	var meetingActionPlan *models.MeetingActionPlan

	result := repo.db.
		Preload("Meeting").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("SeverityLevel").
		Preload("ActionPlanStatus").
		Preload("MeetingMinute").
		Preload("CreatedByUser").
		Preload("CreatedByUser.OrganizationUnit1").
		Preload("UpdatedByUser").
		Preload("UpdatedByUser.OrganizationUnit1").
		Preload("ActionPlanStatusUpdates").
		Preload("ReviewUser").
		First(&meetingActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return meetingActionPlan, nil
}

func (repo *MeetingActionPlanRepository) Create(meetingActionPlan *models.MeetingActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(meetingActionPlan).Error; err != nil {
		fmt.Println(err)
		return 0, err
	}

	return meetingActionPlan.ID, nil
}

func (repo *MeetingActionPlanRepository) Update(meetingActionPlan *models.MeetingActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(meetingActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MeetingActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var meetingActionPlan *models.MeetingActionPlan
	if err := db.First(&meetingActionPlan, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.MeetingActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *MeetingActionPlanRepository) CreateInitialMeetingActionPlanStatusUpdate(meetingActionPlanStatusUpdate *models.MeetingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(meetingActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}
	return meetingActionPlanStatusUpdate.ID, nil
}

func (repo *MeetingActionPlanRepository) DeleteMeetingActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.MeetingActionPlanStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingActionPlanRepository) GetMeetingActionPlanSummary(meetingID uint) (openCount int, closedCount int, err error) {
	var open int64
	var closed int64

	err = repo.db.Model(&models.MeetingActionPlan{}).
		Where("meeting_id = ? AND action_plan_status_id IN (SELECT id FROM action_plan_statuses WHERE type IN ('new', 'in_progress', 'on_hold'))", meetingID).
		Count(&open).Error
	if err != nil {
		return 0, 0, err
	}

	err = repo.db.Model(&models.MeetingActionPlan{}).
		Where("meeting_id = ? AND action_plan_status_id IN (SELECT id FROM action_plan_statuses WHERE type IN ('completed', 'terminated'))", meetingID).
		Count(&closed).Error
	if err != nil {
		return 0, 0, err
	}

	openCount = int(open)
	closedCount = int(closed)

	return openCount, closedCount, nil
}
