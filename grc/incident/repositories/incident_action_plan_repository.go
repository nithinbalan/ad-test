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
	Create(incidentActionPlan *models.IncidentActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(incidentActionPlan *models.IncidentActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetCountByRiskID(incidentID uint, tx ...utility.DatabaseTransaction) (int64, error)

	CreateInitialIncidentActionPlanStatusUpdate(incidentActionPlanStatusUpdate *models.IncidentActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteIncidentActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error
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
		incident_register.reference_number AS incident_reference_number,
		incident_register.title AS incident_title,
		CONVERT(VARCHAR, incident_action_plans.target_date, 23) AS target_date,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
        severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
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
		LEFT JOIN incident_register ON incident_register.id = incident_action_plans.incident_register_id
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = incident_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = incident_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = incident_action_plans.responsible_user_id
		LEFT JOIN users AS review_user ON review_user.id = incident_action_plans.review_user_id
        LEFT JOIN users AS created_user ON created_user.id = incident_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = incident_action_plans.updated_by_id
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
	// Conditionally apply incident ID filter
	if incidentID != 0 {
		query = query.Where("incident_action_plans.incident_register_id = ?", incidentID)
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&incidentActionPlans).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &incidentDTOs.IncidentActionPlanListResponseDTO{
			Items: incidentActionPlans,
		}, nil
	}
	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}
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
		Preload("ActionPlanStatusUpdates").
		Preload("ReviewUser").
		First(&incidentActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentActionPlan, nil
}

func (repo *IncidentActionPlanRepository) Create(incidentActionPlan *models.IncidentActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentActionPlan).Error; err != nil {
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

func (repo *IncidentActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentActionPlan *models.IncidentActionPlan
	if err := db.First(&incidentActionPlan, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.IncidentActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentActionPlanRepository) CreateInitialIncidentActionPlanStatusUpdate(incidentActionPlanStatusUpdate *models.IncidentActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}
	return incidentActionPlanStatusUpdate.ID, nil
}

func (repo *IncidentActionPlanRepository) DeleteIncidentActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.IncidentActionPlanStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil

}

func (repo *IncidentActionPlanRepository) GetCountByRiskID(incidentID uint, tx ...utility.DatabaseTransaction) (int64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var count int64
	if err := db.Model(&models.IncidentActionPlan{}).
		Where("incident_register_id = ?", incidentID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
