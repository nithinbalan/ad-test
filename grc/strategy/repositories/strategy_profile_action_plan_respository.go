package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileActionPlanListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileActionPlanListDTO, error)
	GetByID(id uint) (*models.StrategyProfileActionPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileActionPlan, error)
	Create(strategyProfileActionPlan *models.StrategyProfileActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(strategyProfileActionPlan *models.StrategyProfileActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type StrategyProfileActionPlanRepository struct {
	db *gorm.DB
}

func NewStrategyProfileActionPlanRepository(db *gorm.DB) *StrategyProfileActionPlanRepository {
	return &StrategyProfileActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StrategyProfileActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategyProfileActionPlanRepository) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileActionPlanListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileActionPlanListDTO, error) {
	var strategyProfileActionPlans []*strategyDTOs.StrategyProfileActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("strategy_profile_action_plans").
		Select(`
            strategy_profile_action_plans.*,
			strategy_profiles.title AS strategy_profile_title,
			strategy_profile_focus_areas.title AS strategy_profile_focus_area_title,
			strategy_profile_goals.title AS strategy_profile_goal_title,
			strategy_profile_objectives.title AS strategy_profile_objective_title,
			strategy_profile_action_plans.title AS strategy_profile_action_plan_title,
			strategy_profile_projects.title AS strategy_profile_project_title,
			strategy_profile_milestones.title AS strategy_profile_milestone_title,
			CONVERT(VARCHAR, strategy_profile_goals.target_date, 23) AS target_date,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
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
			LEFT JOIN users AS responsible_user ON responsible_user.id = strategy_profile_action_plans.responsible_user_id
			LEFT JOIN strategy_profile_goals ON strategy_profile_goals.id = strategy_profile_action_plans.strategy_profile_goal_id
			LEFT JOIN strategy_profiles ON strategy_profiles.id = strategy_profile_action_plans.strategy_profile_id
			LEFT JOIN strategy_profile_focus_areas ON strategy_profile_focus_areas.id = strategy_profile_action_plans.strategy_profile_focus_area_id
			LEFT JOIN strategy_profile_objectives ON strategy_profile_objectives.id = strategy_profile_action_plans.strategy_profile_objective_id
			LEFT JOIN strategy_profile_initiatives ON strategy_profile_initiatives.id = strategy_profile_action_plans.strategy_profile_initiative_id
			LEFT JOIN strategy_profile_projects ON strategy_profile_projects.id = strategy_profile_action_plans.strategy_profile_project_id
			LEFT JOIN strategy_profile_milestones ON strategy_profile_milestones.id = strategy_profile_action_plans.strategy_profile_milestone_id
			LEFT JOIN users AS created_user ON created_user.id = strategy_profile_action_plans.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategy_profile_action_plans.updated_by_id
        `).Where("strategy_profile_action_plans.strategy_profile_id = ?", strategyProfileID)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("strategy_profile_action_plans."+columnName+" IN (?)", values)
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
			(strategy_profile_action_plans.title LIKE ? OR
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

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&strategyProfileActionPlans).Error
	if err != nil {
		return nil, err
	}
	return strategyProfileActionPlans, nil
}

func (repo *StrategyProfileActionPlanRepository) GetByID(id uint) (*models.StrategyProfileActionPlan, error) {
	var strategyProfileActionPlan *models.StrategyProfileActionPlan

	result := repo.db.First(&strategyProfileActionPlan, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileActionPlan, nil

}

func (repo *StrategyProfileActionPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileActionPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var strategyProfileActionPlan *models.StrategyProfileActionPlan

	result := db.
		Preload("StrategyProfile").
		Preload("StrategyProfile.ResponsibleUser").
		Preload("StrategyProfileFocusArea").
		Preload("StrategyProfileGoal").
		Preload("StrategyProfileObjective").
		Preload("StrategyProfileInitiative").
		Preload("StrategyProfileProject").
		Preload("StrategyProfileMilestone").
		Preload("ResponsibleUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&strategyProfileActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileActionPlan, nil
}

func (repo *StrategyProfileActionPlanRepository) Create(strategyProfileActionPlan *models.StrategyProfileActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(strategyProfileActionPlan).Error; err != nil {
		return 0, err
	}

	return strategyProfileActionPlan.ID, nil
}

func (repo *StrategyProfileActionPlanRepository) Update(strategyProfileActionPlan *models.StrategyProfileActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(strategyProfileActionPlan).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategyProfileActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Strategy profile
	if err := db.Delete(&models.StrategyProfileActionPlan{}, id).Error; err != nil {
		return err
	}

	return nil
}
