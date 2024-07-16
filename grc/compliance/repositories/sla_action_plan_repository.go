package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type SLAActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(slaID uint, params *complianceDTOs.SLAActionPlanListQueryParamsDTO) (*complianceDTOs.SLAActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.SLAActionPlan, error)
	GetFullDetailsByID(id uint) (*models.SLAActionPlan, error)
	Create(slaActionPlan *models.SLAActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(slaActionPlan *models.SLAActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateInitialSLAActionPlanStatusUpdate(slaActionPlanStatusUpdate *models.SLAActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteSLAActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error
}

type SLAActionPlanRepository struct {
	db *gorm.DB
}

func NewSLAActionPlanRepository(db *gorm.DB) *SLAActionPlanRepository {
	return &SLAActionPlanRepository{
		db: db,
	}
}

func (repo *SLAActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *SLAActionPlanRepository) GetAll(slaID uint, params *complianceDTOs.SLAActionPlanListQueryParamsDTO) (*complianceDTOs.SLAActionPlanListResponseDTO, error) {

	var slaActionPlans []*complianceDTOs.SLAActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("sla_action_plans").
		Select(`
        sla_action_plans.*,
        sla_check_questions.title_en AS sla_check_question_title_en,
		sla_check_questions.title_ar AS sla_check_question_title_ar,
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
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = sla_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = sla_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = sla_action_plans.responsible_user_id
        LEFT JOIN users AS created_user ON created_user.id = sla_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = sla_action_plans.updated_by_id
        LEFT JOIN sla_check_questions ON sla_check_questions.id = sla_action_plans.sla_check_question_id
    `).Where(`sla_action_plans.sla_register_id`, slaID)

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
            (sla_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            sla_action_plans.target_date LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
            responsible_user.email LIKE ?)`,
			keyword, keyword, keyword, keyword,
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&slaActionPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.SLAActionPlanListResponseDTO{
		Items: slaActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *SLAActionPlanRepository) GetByID(id uint) (*models.SLAActionPlan, error) {
	var slaActionPlan *models.SLAActionPlan
	result := repo.db.First(&slaActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return slaActionPlan, nil
}

func (repo *SLAActionPlanRepository) GetFullDetailsByID(id uint) (*models.SLAActionPlan, error) {
	var slaActionPlan *models.SLAActionPlan

	result := repo.db.
		Preload("SLACheckQuestion").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("SeverityLevel").
		Preload("ActionPlanStatus").
		Preload("ActionPlanStatusUpdates").
		First(&slaActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return slaActionPlan, nil
}

func (repo *SLAActionPlanRepository) Create(slaActionPlan *models.SLAActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	if err := repo.db.Create(slaActionPlan).Error; err != nil {
		return 0, err
	}
	return slaActionPlan.ID, nil
}

func (repo *SLAActionPlanRepository) Update(slaActionPlan *models.SLAActionPlan, tx ...utility.DatabaseTransaction) error {
	if err := repo.db.Save(slaActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SLAActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var slaActionPlan *models.SLAActionPlan
	if err := db.First(&slaActionPlan, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.SLAActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SLAActionPlanRepository) CreateInitialSLAActionPlanStatusUpdate(slaActionPlanStatusUpdate *models.SLAActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(slaActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}
	return slaActionPlanStatusUpdate.ID, nil
}

func (repo *SLAActionPlanRepository) DeleteSLAActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.SLAActionPlanStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil
}
