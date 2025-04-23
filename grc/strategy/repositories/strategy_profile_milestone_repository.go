package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileMilestoneRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileMilestoneListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileMilestoneListDTO, error)
	GetByID(id uint) (*models.StrategyProfileMilestone, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileMilestone, error)
	Create(strategyProfileMilestone *models.StrategyProfileMilestone, tx ...utility.DatabaseTransaction) (uint, error)
	Update(strategyProfileMilestone *models.StrategyProfileMilestone, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type StrategyProfileMilestoneRepository struct {
	db *gorm.DB
}

func NewStrategyProfileMilestoneRepository(db *gorm.DB) *StrategyProfileMilestoneRepository {
	return &StrategyProfileMilestoneRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StrategyProfileMilestoneRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategyProfileMilestoneRepository) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileMilestoneListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileMilestoneListDTO, error) {
	var strategyProfileMilestones []*strategyDTOs.StrategyProfileMilestoneListDTO
	var totalCount int64

	query := repo.db.Table("strategy_profile_milestones").
		Select(`
            strategy_profile_milestones.*,
			strategy_profiles.title AS strategy_profile_title,
			strategy_profile_focus_areas.title AS strategy_profile_focus_area_title,
			strategy_profile_goals.title AS strategy_profile_goal_title,
			strategy_profile_objectives.title AS strategy_profile_objective_title,
			strategy_profile_milestones.title AS strategy_profile_milestone_title,
			strategy_profile_projects.title AS strategy_profile_project_title,
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
			LEFT JOIN strategy_profile_goals ON strategy_profile_goals.id = strategy_profile_milestones.strategy_profile_goal_id
			LEFT JOIN strategy_profiles ON strategy_profiles.id = strategy_profile_milestones.strategy_profile_id
			LEFT JOIN strategy_profile_focus_areas ON strategy_profile_focus_areas.id = strategy_profile_milestones.strategy_profile_focus_area_id
			LEFT JOIN strategy_profile_objectives ON strategy_profile_objectives.id = strategy_profile_milestones.strategy_profile_objective_id
			LEFT JOIN strategy_profile_initiatives ON strategy_profile_initiatives.id = strategy_profile_milestones.strategy_profile_initiative_id
			LEFT JOIN strategy_profile_projects ON strategy_profile_projects.id = strategy_profile_milestones.strategy_profile_project_id
			LEFT JOIN users AS created_user ON created_user.id = strategy_profile_milestones.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategy_profile_milestones.updated_by_id
        `).Where("strategy_profile_milestones.strategy_profile_id = ?", strategyProfileID)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("strategy_profile_milestones."+columnName+" IN (?)", values)
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
			(strategy_profile_milestones.title LIKE ? OR
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

	err := query.Find(&strategyProfileMilestones).Error
	if err != nil {
		return nil, err
	}
	return strategyProfileMilestones, nil
}

func (repo *StrategyProfileMilestoneRepository) GetByID(id uint) (*models.StrategyProfileMilestone, error) {
	var strategyProfileMilestone *models.StrategyProfileMilestone

	result := repo.db.First(&strategyProfileMilestone, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileMilestone, nil

}

func (repo *StrategyProfileMilestoneRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileMilestone, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var strategyProfileMilestone *models.StrategyProfileMilestone

	result := db.
		Preload("StrategyProfile").
		Preload("StrategyProfile.ResponsibleUser").
		Preload("StrategyProfileFocusArea").
		Preload("StrategyProfileGoal").
		Preload("StrategyProfileObjective").
		Preload("StrategyProfileInitiative").
		Preload("StrategyProfileProject").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&strategyProfileMilestone, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileMilestone, nil
}

func (repo *StrategyProfileMilestoneRepository) Create(strategyProfileMilestone *models.StrategyProfileMilestone, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(strategyProfileMilestone).Error; err != nil {
		return 0, err
	}

	return strategyProfileMilestone.ID, nil
}

func (repo *StrategyProfileMilestoneRepository) Update(strategyProfileMilestone *models.StrategyProfileMilestone, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(strategyProfileMilestone).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategyProfileMilestoneRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Strategy profile
	if err := db.Delete(&models.StrategyProfileMilestone{}, id).Error; err != nil {
		return err
	}

	return nil
}
