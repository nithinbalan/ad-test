package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type ProjectActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(projectID uint, params *projectDTOs.ProjectActionPlanListQueryParamsDTO) (*projectDTOs.ProjectActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.ProjectActionPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectActionPlan, error)
	Create(projectActionPlan *models.ProjectActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(projectActionPlan *models.ProjectActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
}

type ProjectActionPlanRepository struct {
	db *gorm.DB
}

func NewProjectActionPlanRepository(db *gorm.DB) *ProjectActionPlanRepository {
	return &ProjectActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ProjectActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProjectActionPlanRepository) GetAll(projectID uint, params *projectDTOs.ProjectActionPlanListQueryParamsDTO) (*projectDTOs.ProjectActionPlanListResponseDTO, error) {

	var projectActionPlans []*projectDTOs.ProjectActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("project_action_plans").
		Select(`
        project_action_plans.*,
		CONVERT(VARCHAR, project_action_plans.target_date, 23) AS target_date,
		projects.title_en AS project_title_en,
		projects.title_ar AS project_title_ar,
        project_milestones.title_en AS project_milestone_title_en,
        project_milestones.title_ar AS project_milestone_title_ar,
		project_milestones.description_en AS project_milestone_description_en,
        project_milestones.description_ar AS project_milestone_description_ar,
		severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
		action_plan_statuses.type AS action_plan_status_type,
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
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = project_action_plans.action_plan_status_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = project_action_plans.responsible_user_id
        LEFT JOIN users AS review_user ON review_user.id = project_action_plans.review_user_id
		LEFT JOIN severity_levels ON severity_levels.id = project_action_plans.severity_level_id
        LEFT JOIN users AS created_user ON created_user.id = project_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = project_action_plans.updated_by_id
        LEFT JOIN project_milestones ON project_milestones.id = project_action_plans.project_milestone_id
		LEFT JOIN projects ON projects.id = project_action_plans.project_id
    `)

		// Conditionally apply project ID filter
	if projectID != 0 {
		query = query.Where("project_action_plans.project_id = ?", projectID)
	}

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("project_action_plans.%s IN (?)", columnName), values)
		}
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
            (project_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            project_action_plans.target_date LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
            responsible_user.email LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
            responsible_user.last_name_ar LIKE ?)`,
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

	// Execute the query
	err := query.Find(&projectActionPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &projectDTOs.ProjectActionPlanListResponseDTO{
		Items: projectActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ProjectActionPlanRepository) GetByID(id uint) (*models.ProjectActionPlan, error) {
	var projectActionPlan *models.ProjectActionPlan
	result := repo.db.First(&projectActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return projectActionPlan, nil
}

func (repo *ProjectActionPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectActionPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var projectActionPlan models.ProjectActionPlan

	// Load the project action plan with all related data
	result := db.
		Preload("Project").
		Preload("ProjectMilestone").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("SeverityLevel").
		Preload("ReviewUser").
		Preload("ReviewUser.OrganizationUnit1").
		Preload("ActionPlanStatus").
		Preload("ActionPlanStatusUpdates").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&projectActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &projectActionPlan, nil
}

func (repo *ProjectActionPlanRepository) Create(projectActionPlan *models.ProjectActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(projectActionPlan).Error; err != nil {
		return 0, err
	}
	return projectActionPlan.ID, nil
}

func (repo *ProjectActionPlanRepository) Update(projectActionPlan *models.ProjectActionPlan, tx ...utility.DatabaseTransaction) error {
	if err := repo.db.Save(projectActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProjectActionPlanRepository) Delete(id uint) error {

	var projectActionPlan *models.ProjectActionPlan
	if err := repo.db.First(&projectActionPlan, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ProjectActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
