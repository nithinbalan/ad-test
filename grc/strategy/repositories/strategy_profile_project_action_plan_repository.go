package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileProjectActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectActionPlanListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectActionPlanListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectActionPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectActionPlan, error)
	Create(strategyProfileProjectActionPlan *models.StrategyProfileProjectActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type StrategyProfileProjectActionPlanRepository struct {
	db *gorm.DB
}

func NewStrategyProfileProjectActionPlanRepository(db *gorm.DB) *StrategyProfileProjectActionPlanRepository {
	return &StrategyProfileProjectActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction
func (repo *StrategyProfileProjectActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

// GetAll returns list of objectives with focus area titles and user info
func (repo *StrategyProfileProjectActionPlanRepository) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectActionPlanListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectActionPlanListDTO, error) {
	var strategyProfileProjectActionPlans []*strategyDTOs.StrategyProfileProjectActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("strategy_profile_project_action_plans").
		Select(`
			strategy_profile_project_action_plans.*,
			strategy_profiles.title AS strategy_profile_title,
			strategy_profile_projects.id AS strategy_profile_project_id,
			projects.title_en AS strategy_profile_project_title_en,
			projects.title_ar AS strategy_profile_project_title_ar,
			strategy_profile_project_milestones.id AS strategy_profile_project_milestone_id,
			project_milestones.title_en AS strategy_profile_project_milestone_title_en,
			project_milestones.title_ar AS strategy_profile_project_milestone_title_ar,
			project_action_plans.id AS project_action_plan_id,
			project_action_plans.title AS project_action_plan_title,
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
			LEFT JOIN strategy_profiles ON strategy_profiles.id = strategy_profile_project_action_plans.strategy_profile_id
			LEFT JOIN strategy_profile_projects ON strategy_profile_projects.id = strategy_profile_project_action_plans.strategy_profile_project_id
			LEFT JOIN projects ON projects.id = strategy_profile_projects.project_id
			LEFT JOIN strategy_profile_project_milestones ON strategy_profile_project_milestones.id = strategy_profile_project_action_plans.strategy_profile_project_milestone_id
			LEFT JOIN project_milestones ON project_milestones.id = strategy_profile_project_milestones.project_milestone_id
			LEFT JOIN project_action_plans ON project_action_plans.id = strategy_profile_project_action_plans.project_action_plan_id
			LEFT JOIN users AS created_user ON created_user.id = strategy_profile_project_action_plans.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategy_profile_project_action_plans.updated_by_id
		`).
		Where("strategy_profile_project_action_plans.strategy_profile_id = ?", strategyProfileID)

	// Apply filters
	for column, values := range params.Filters {
		// Map API param to DB column
		if column == "project_milestone_id" {
			column = "strategy_profile_project_action_plans.strategy_profile_project_milestone_id"
		}
		if column == "project_id" {
			column = "strategy_profile_project_action_plans.strategy_profile_project_id"
		}
		if len(values) > 0 {
			query = query.Where(column+" IN (?)", values)
		}
	}

	// Apply keyword search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			project_action_plans.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?
		`, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Pagination
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&strategyProfileProjectActionPlans).Error
	if err != nil {
		return nil, err
	}

	return strategyProfileProjectActionPlans, nil
}

// GetByID fetches a single project action plan by ID
func (repo *StrategyProfileProjectActionPlanRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectActionPlan, error) {
	var projectActionPlan *models.StrategyProfileProjectActionPlan

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.First(&projectActionPlan, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return projectActionPlan, nil
}
func (repo *StrategyProfileProjectActionPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectActionPlan, error) {
	var projectActionPlan *models.StrategyProfileProjectActionPlan

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Preload("StrategyProfileProject").
		Preload("ProjectActionPlan").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&projectActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return projectActionPlan, nil
}

// Create inserts a new project action plan record
func (repo *StrategyProfileProjectActionPlanRepository) Create(projectActionPlan *models.StrategyProfileProjectActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(projectActionPlan).Error; err != nil {
		return 0, err
	}

	return projectActionPlan.ID, nil
}

// Delete removes a project action plan record by ID
func (repo *StrategyProfileProjectActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.StrategyProfileProjectActionPlan{}, id).Error; err != nil {
		return err
	}
	return nil
}
