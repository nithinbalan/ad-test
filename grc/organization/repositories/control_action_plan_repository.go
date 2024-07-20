package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type ControlActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(controlID uint, params *organizationDTOs.ControlActionPlanListQueryParamsDTO) (*organizationDTOs.ControlActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.ControlActionPlan, error)
	GetFullDetailsByID(id uint) (*models.ControlActionPlan, error)
	Create(controlActionPlan *models.ControlActionPlan) (uint, error)
	Update(controlActionPlan *models.ControlActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error

	CreateInitialControlActionPlanStatusUpdate(controlActionPlanStatusUpdate *models.ControlActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteControlActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error
}

type ControlActionPlanRepository struct {
	db *gorm.DB
}

func NewControlActionPlanRepository(db *gorm.DB) *ControlActionPlanRepository {
	return &ControlActionPlanRepository{
		db: db,
	}
}

func (repo *ControlActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ControlActionPlanRepository) GetAll(controlID uint, params *organizationDTOs.ControlActionPlanListQueryParamsDTO) (*organizationDTOs.ControlActionPlanListResponseDTO, error) {

	var controlActionPlans []*organizationDTOs.ControlActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("control_action_plans").
		Select(`
        control_action_plans.*,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
        severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
        responsible_user.first_name_en AS responsible_user_first_name_en,
        responsible_user.first_name_ar AS responsible_user_first_name_ar,
        responsible_user.last_name_en AS responsible_user_last_name_en,
        responsible_user.last_name_ar AS responsible_user_last_name_ar,
		responsible_user.initial AS responsible_user_initial,
        responsible_user.color AS responsible_user_color,
        responsible_user.email AS responsible_user_email,
        responsible_user.profile_image_path AS responsible_user_profile_image_path,
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
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = control_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = control_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = control_action_plans.responsible_user_id
        LEFT JOIN users AS created_user ON created_user.id = control_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = control_action_plans.updated_by_id
    `).Where(`control_action_plans.control_id`, controlID)

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

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (control_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            control_action_plans.target_date LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
            responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword)
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&controlActionPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &organizationDTOs.ControlActionPlanListResponseDTO{
		Items: controlActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ControlActionPlanRepository) GetByID(id uint) (*models.ControlActionPlan, error) {
	var controlActionPlan *models.ControlActionPlan
	result := repo.db.First(&controlActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlActionPlan, nil
}

func (repo *ControlActionPlanRepository) GetFullDetailsByID(id uint) (*models.ControlActionPlan, error) {
	var controlActionPlan *models.ControlActionPlan

	result := repo.db.
		Preload("Control").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("SeverityLevel").
		Preload("ActionPlanStatus").
		Preload("CreatedByUser").
		Preload("CreatedByUser.OrganizationUnit1").
		Preload("UpdatedByUser").
		Preload("UpdatedByUser.OrganizationUnit1").
		Preload("ActionPlanStatusUpdates").
		First(&controlActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlActionPlan, nil
}

func (repo *ControlActionPlanRepository) Create(controlActionPlan *models.ControlActionPlan) (uint, error) {
	if err := repo.db.Create(controlActionPlan).Error; err != nil {
		return 0, err
	}
	return controlActionPlan.ID, nil
}

func (repo *ControlActionPlanRepository) Update(controlActionPlan *models.ControlActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(controlActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ControlActionPlanRepository) Delete(id uint) error {

	var controlActionPlan *models.ControlActionPlan
	if err := repo.db.First(&controlActionPlan, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ControlActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ControlActionPlanRepository) CreateInitialControlActionPlanStatusUpdate(controlActionPlanStatusUpdate *models.ControlActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(controlActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}
	return controlActionPlanStatusUpdate.ID, nil
}

func (repo *ControlActionPlanRepository) DeleteControlActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ControlActionPlanStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil
}
