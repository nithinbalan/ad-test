package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileInitiativeRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileInitiativeListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileInitiativeListDTO, error)
	GetByID(id uint) (*models.StrategyProfileInitiative, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileInitiative, error)
	Create(strategyProfileInitiative *models.StrategyProfileInitiative, tx ...utility.DatabaseTransaction) (uint, error)
	Update(strategyProfileInitiative *models.StrategyProfileInitiative, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type StrategyProfileInitiativeRepository struct {
	db *gorm.DB
}

func NewStrategyProfileInitiativeRepository(db *gorm.DB) *StrategyProfileInitiativeRepository {
	return &StrategyProfileInitiativeRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StrategyProfileInitiativeRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategyProfileInitiativeRepository) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileInitiativeListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileInitiativeListDTO, error) {
	var strategyProfileInitiatives []*strategyDTOs.StrategyProfileInitiativeListDTO
	var totalCount int64

	query := repo.db.Table("strategy_profile_initiatives").
		Select(`
            strategy_profile_initiatives.*,
			strategy_profiles.title AS strategy_profile_title,
			strategy_profile_focus_areas.title AS strategy_profile_focus_area_title,
			strategy_profile_goals.title AS strategy_profile_goal_title,
			strategy_profile_objectives.title AS strategy_profile_objective_title,
			CONVERT(VARCHAR, strategy_profile_initiatives.start_date, 23) AS start_date,
			CONVERT(VARCHAR, strategy_profile_initiatives.end_date, 23) AS end_date,
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
			LEFT JOIN users AS responsible_user ON responsible_user.id = strategy_profile_initiatives.responsible_user_id
			LEFT JOIN strategy_profile_goals ON strategy_profile_goals.id = strategy_profile_initiatives.strategy_profile_goal_id
			LEFT JOIN strategy_profiles ON strategy_profiles.id = strategy_profile_initiatives.strategy_profile_id
			LEFT JOIN strategy_profile_focus_areas ON strategy_profile_focus_areas.id = strategy_profile_initiatives.strategy_profile_focus_area_id
			LEFT JOIN strategy_profile_objectives ON strategy_profile_objectives.id = strategy_profile_initiatives.strategy_profile_objective_id
			LEFT JOIN users AS created_user ON created_user.id = strategy_profile_initiatives.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategy_profile_initiatives.updated_by_id
        `).Where("strategy_profile_initiatives.strategy_profile_id = ?", strategyProfileID)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("strategy_profile_initiatives."+columnName+" IN (?)", values)
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
			(strategy_profile_initiatives.title LIKE ? OR
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

	err := query.Find(&strategyProfileInitiatives).Error
	if err != nil {
		return nil, err
	}
	return strategyProfileInitiatives, nil
}

func (repo *StrategyProfileInitiativeRepository) GetByID(id uint) (*models.StrategyProfileInitiative, error) {
	var strategyProfileInitiative *models.StrategyProfileInitiative

	result := repo.db.First(&strategyProfileInitiative, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileInitiative, nil

}

func (repo *StrategyProfileInitiativeRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileInitiative, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var strategyProfileInitiative *models.StrategyProfileInitiative

	result := db.
		Preload("StrategyProfile").
		Preload("StrategyProfile.ResponsibleUser").
		Preload("StrategyProfileFocusArea").
		Preload("StrategyProfileGoal").
		Preload("StrategyProfileObjective").
		Preload("ResponsibleUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&strategyProfileInitiative, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileInitiative, nil
}

func (repo *StrategyProfileInitiativeRepository) Create(strategyProfileInitiative *models.StrategyProfileInitiative, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(strategyProfileInitiative).Error; err != nil {
		return 0, err
	}

	return strategyProfileInitiative.ID, nil
}

func (repo *StrategyProfileInitiativeRepository) Update(strategyProfileInitiative *models.StrategyProfileInitiative, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(strategyProfileInitiative).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategyProfileInitiativeRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Strategy profile
	if err := db.Delete(&models.StrategyProfileInitiative{}, id).Error; err != nil {
		return err
	}

	return nil
}
