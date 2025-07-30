package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileProjectRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProject, error)
	Create(strategyProfileProject *models.StrategyProfileProject, tx ...utility.DatabaseTransaction) (uint, error)
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type StrategyProfileProjectRepository struct {
	db *gorm.DB
}

func NewStrategyProfileProjectRepository(db *gorm.DB) *StrategyProfileProjectRepository {
	return &StrategyProfileProjectRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StrategyProfileProjectRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategyProfileProjectRepository) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectListDTO, error) {
	var strategyProfileProjects []*strategyDTOs.StrategyProfileProjectListDTO
	var totalCount int64

	query := repo.db.Table("strategy_profile_projects").
		Select(`
            strategy_profile_projects.*,
			strategy_profiles.title AS strategy_profile_title,
			projects.title_en AS project_title_en,
			projects.title_ar AS project_title_ar,
			projects.description_en AS project_description_en,
			projects.description_ar AS project_description_ar,
			strategic_objectives.title_en AS strategy_profile_strategic_objective_title_en,
			strategic_objectives.title_ar AS strategy_profile_strategic_objective_title_ar,
			strategic_intiatives.title_en AS strategy_profile_strategic_intiative_title_en,
			strategic_intiatives.title_ar AS strategy_profile_strategic_intiative_title_ar,
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
			LEFT JOIN strategy_profiles ON strategy_profiles.id = strategy_profile_projects.strategy_profile_id
			LEFT JOIN projects ON projects.id = strategy_profile_projects.project_id
			LEFT JOIN strategy_profile_strategic_objectives ON strategy_profile_strategic_objectives.id = strategy_profile_projects.strategy_profile_strategic_objective_id
			LEFT JOIN strategic_objectives ON strategic_objectives.id = strategy_profile_strategic_objectives.strategic_objective_id
			LEFT JOIN strategy_profile_strategic_intiatives ON strategy_profile_strategic_intiatives.id = strategy_profile_projects.strategy_profile_strategic_intiative_id
			LEFT JOIN strategic_intiatives ON strategic_intiatives.id = strategy_profile_strategic_intiatives.strategic_intiative_id
			LEFT JOIN users AS created_user ON created_user.id = strategy_profile_projects.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategy_profile_projects.updated_by_id
        `).Where("strategy_profile_projects.strategy_profile_id = ?", strategyProfileID)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			// Map API param to DB column
			if columnName == "strategic_objective_id" {
				columnName = "strategy_profile_projects.strategy_profile_strategic_objective_id"
			}
			if columnName == "strategic_intiative_id" {
				columnName = "strategy_profile_projects.strategy_profile_strategic_intiative_id"
			}
			query = query.Where(columnName+" IN (?)", values)
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
			(projects.title_en LIKE ? OR
			projects.title_ar LIKE ? OR
			projects.description_en LIKE ? OR
			projects.description_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Find(&strategyProfileProjects).Error
	if err != nil {
		return nil, err
	}
	return strategyProfileProjects, nil
}

func (repo *StrategyProfileProjectRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProject, error) {
	var strategyProfileProject *models.StrategyProfileProject

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.First(&strategyProfileProject, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileProject, nil
}

func (repo *StrategyProfileProjectRepository) Create(strategyProfileProject *models.StrategyProfileProject, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(strategyProfileProject).Error; err != nil {
		return 0, err
	}

	return strategyProfileProject.ID, nil
}

func (repo *StrategyProfileProjectRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Strategy profile
	if err := db.Delete(&models.StrategyProfileProject{}, id).Error; err != nil {
		return err
	}

	return nil
}
