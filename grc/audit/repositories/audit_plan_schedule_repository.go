package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditPlanScheduleRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditPlanID uint, params *auditDTOs.AuditPlanScheduleListQueryParamsDTO) (*auditDTOs.AuditPlanScheduleListResponseDTO, error)
	GetByID(id uint) (*models.AuditPlanSchedule, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanSchedule, error)
	Create(auditPlanSchedule *models.AuditPlanSchedule, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditPlanSchedule *models.AuditPlanSchedule, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateAuditableItems(auditPlanScheduleAuditableItems []*models.AuditPlanScheduleAuditableItem, tx ...utility.DatabaseTransaction) error
	DeleteAuditableItems(auditPlanScheduleAuditableItemIds []*uint, tx ...utility.DatabaseTransaction) error

	GetAllSchedulesByAuditPlanID(auditPlanID uint, tx ...utility.DatabaseTransaction) ([]*models.AuditPlanSchedule, error)
}

type AuditPlanScheduleRepository struct {
	db *gorm.DB
}

func NewAuditPlanScheduleRepository(db *gorm.DB) *AuditPlanScheduleRepository {
	return &AuditPlanScheduleRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditPlanScheduleRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditPlanScheduleRepository) GetAll(auditPlanID uint, params *auditDTOs.AuditPlanScheduleListQueryParamsDTO) (*auditDTOs.AuditPlanScheduleListResponseDTO, error) {
	var auditPlanSchedules []*auditDTOs.AuditPlanScheduleListDTO
	var totalCount int64

	query := repo.db.Table("audit_plan_schedules").
		Select(`
			audit_plan_schedules.id,
			audit_plan_schedules.audit_plan_id,
			audit_plan_schedules.audit_activity_type_id,
			audit_activity_types.title_en AS audit_activity_type_title_en,
			audit_activity_types.title_ar AS audit_activity_type_title_ar,
			audit_plan_schedules.responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			audit_plan_schedules.title,
			audit_plan_schedules.description,
			audit_plans.title AS audit_plan_title,
			CONVERT(VARCHAR, audit_plan_schedules.start_date, 23) AS start_date,
			CONVERT(VARCHAR, audit_plan_schedules.target_date, 23) AS target_date,
			audit_plan_schedules.audit_schedule_status_id,
			audit_schedule_statuses.title_en AS audit_schedule_status_title_en,
			audit_schedule_statuses.title_ar AS audit_schedule_status_title_ar,
			audit_schedule_statuses.color_code AS audit_schedule_status_color_code,
			audit_plan_schedules.created_by_id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			audit_plan_schedules.updated_by_id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			CONVERT(VARCHAR, audit_plan_schedules.created_at, 23) AS created_at,
			CONVERT(VARCHAR, audit_plan_schedules.updated_at, 23) AS updated_at
		`).
		Joins(`
			LEFT JOIN audit_activity_types ON audit_activity_types.id = audit_plan_schedules.audit_activity_type_id
			LEFT JOIN audit_plans ON audit_plans.id = audit_plan_schedules.audit_plan_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = audit_plan_schedules.responsible_user_id
			LEFT JOIN audit_schedule_statuses ON audit_schedule_statuses.id = audit_plan_schedules.audit_schedule_status_id
			LEFT JOIN users AS created_user ON created_user.id = audit_plan_schedules.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_plan_schedules.updated_by_id
		`).Where(`audit_plan_schedules.audit_plan_id = ?`, auditPlanID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Apply other filters as needed (e.g., date filters)
	if params.StartDate != "" {
		query = query.Where("audit_plan_schedules.start_date >= ?", params.StartDate)
	}
	if params.TargetDate != "" {
		query = query.Where("audit_plan_schedules.target_date <= ?", params.TargetDate)
	}
	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "scheduled":
			query = query.Where("audit_schedule_statuses.type = ?", "scheduled")
		case "completed":
			query = query.Where("audit_schedule_statuses.type = ?", "completed")
		case "in_progress":
			query = query.Where("audit_schedule_statuses.type = ?", "in_progress")
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
			audit_plan_schedules.title LIKE ? OR
			audit_plan_schedules.description LIKE ? OR
			audit_plans.title LIKE ? OR
			audit_activity_types.title_en LIKE ? OR
			audit_activity_types.title_ar LIKE ? OR
			audit_schedule_statuses.title_en LIKE ? OR
			audit_schedule_statuses.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?`,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword)
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditPlanSchedules).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditPlanScheduleListResponseDTO{
		Items: auditPlanSchedules,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditPlanScheduleRepository) GetByID(id uint) (*models.AuditPlanSchedule, error) {
	var auditPlanSchedule *models.AuditPlanSchedule

	result := repo.db.First(&auditPlanSchedule, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlanSchedule, nil
}

func (repo *AuditPlanScheduleRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanSchedule, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditPlanSchedule *models.AuditPlanSchedule

	result := db.
		Preload("AuditPlan").
		Preload("AuditActivityType").
		Preload("ResponsibleUser").
		Preload("AuditScheduleStatus").
		Preload("AuditableItems").
		Preload("AuditableItems.AuditableItem").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditPlanSchedule, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlanSchedule, nil
}

func (repo *AuditPlanScheduleRepository) Create(auditPlanSchedule *models.AuditPlanSchedule, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditPlanSchedule).Error; err != nil {
		return 0, err
	}

	return auditPlanSchedule.ID, nil
}

func (repo *AuditPlanScheduleRepository) Update(auditPlanSchedule *models.AuditPlanSchedule, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditPlanSchedule).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanScheduleRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit Schedule
	if err := db.Delete(&models.AuditPlanSchedule{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanScheduleRepository) CreateAuditableItems(auditPlanScheduleAuditableItems []*models.AuditPlanScheduleAuditableItem, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&auditPlanScheduleAuditableItems).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanScheduleRepository) DeleteAuditableItems(auditPlanScheduleAuditableItemIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete controls by IDs
	if err := db.Where("id IN (?)", auditPlanScheduleAuditableItemIds).Delete(&models.AuditPlanScheduleAuditableItem{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanScheduleRepository) GetAllSchedulesByAuditPlanID(auditPlanID uint, tx ...utility.DatabaseTransaction) ([]*models.AuditPlanSchedule, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var schedules []*models.AuditPlanSchedule

	err := db.Where("audit_plan_id = ?", auditPlanID).Find(&schedules).Error
	if err != nil {
		return nil, err
	}

	return schedules, nil
}
