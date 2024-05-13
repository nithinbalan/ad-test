package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type IncidentActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(incidentID uint, params *incidentDTOs.IncidentActionPlanListQueryParamsDTO) (*incidentDTOs.IncidentActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.IncidentActionPlan, error)
	GetFullDetailsByID(id uint) (*models.IncidentActionPlan, error)
	Create(incidentActionPlan *models.IncidentActionPlan) (uint, error)
	Update(incidentActionPlan *models.IncidentActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
}

type IncidentActionPlanRepository struct {
	db *gorm.DB
}

func NewIncidentActionPlanRepository(db *gorm.DB) *IncidentActionPlanRepository {
	return &IncidentActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *IncidentActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentActionPlanRepository) GetAll(incidentID uint, params *incidentDTOs.IncidentActionPlanListQueryParamsDTO) (*incidentDTOs.IncidentActionPlanListResponseDTO, error) {

	var incidentActionPlans []*incidentDTOs.IncidentActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("incident_action_plans").
		Select(`
        incident_action_plans.*,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
        severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
        responsible_user.first_name_en AS responsible_user_first_name_en,
        responsible_user.first_name_ar AS responsible_user_first_name_ar,
        responsible_user.last_name_en AS responsible_user_last_name_en,
        responsible_user.last_name_ar AS responsible_user_last_name_ar,
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
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = incident_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = incident_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = incident_action_plans.responsible_user_id
        LEFT JOIN users AS created_user ON created_user.id = incident_action_plans.created_by
        LEFT JOIN users AS updated_user ON updated_user.id = incident_action_plans.updated_by
    `).Where(`incident_action_plans.incident_register_id`, incidentID)

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
            (incident_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            incident_action_plans.target_date LIKE ? OR
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentActionPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &incidentDTOs.IncidentActionPlanListResponseDTO{
		Items: incidentActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *IncidentActionPlanRepository) GetByID(id uint) (*models.IncidentActionPlan, error) {
	var incidentActionPlan *models.IncidentActionPlan
	result := repo.db.First(&incidentActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentActionPlan, nil
}

func (repo *IncidentActionPlanRepository) GetFullDetailsByID(id uint) (*models.IncidentActionPlan, error) {
	var incidentActionPlan *models.IncidentActionPlan

	result := repo.db.
		Preload("IncidentRegister").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("SeverityLevel").
		Preload("ActionPlanStatus").
		First(&incidentActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentActionPlan, nil
}

func (repo *IncidentActionPlanRepository) Create(incidentActionPlan *models.IncidentActionPlan) (uint, error) {
	if err := repo.db.Create(incidentActionPlan).Error; err != nil {
		return 0, err
	}
	return incidentActionPlan.ID, nil
}

func (repo *IncidentActionPlanRepository) Update(incidentActionPlan *models.IncidentActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentActionPlanRepository) Delete(id uint) error {

	var incidentActionPlan *models.IncidentActionPlan
	if err := repo.db.First(&incidentActionPlan, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.IncidentActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
