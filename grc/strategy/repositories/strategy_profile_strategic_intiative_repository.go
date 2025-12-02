package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileStrategicIntiativeRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileStrategicIntiativeListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileStrategicIntiativeListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileStrategicIntiative, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileStrategicIntiative, error)
	Create(strategicIntiative *models.StrategyProfileStrategicIntiative, tx ...utility.DatabaseTransaction) (uint, error)
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type StrategyProfileStrategicIntiativeRepository struct {
	db *gorm.DB
}

func NewStrategyProfileStrategicIntiativeRepository(db *gorm.DB) *StrategyProfileStrategicIntiativeRepository {
	return &StrategyProfileStrategicIntiativeRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction
func (repo *StrategyProfileStrategicIntiativeRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

// GetAll returns list of objectives with focus area titles and user info
func (repo *StrategyProfileStrategicIntiativeRepository) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileStrategicIntiativeListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileStrategicIntiativeListDTO, error) {
	var strategicIntiatives []*strategyDTOs.StrategyProfileStrategicIntiativeListDTO
	var totalCount int64

	query := repo.db.Table("strategy_profile_strategic_intiatives").
		Select(`
			strategy_profile_strategic_intiatives.*,
			strategy_profiles.title AS strategy_profile_title,
			strategy_profile_strategic_objectives.id AS strategy_profile_strategic_objective_id,
			strategic_objectives.title_en AS strategy_profile_strategic_objective_title_en,
			strategic_objectives.title_ar AS strategy_profile_strategic_objective_title_ar,
			strategic_intiatives.id AS strategic_intiative_id,
			strategic_intiatives.title_en AS strategic_intiative_title_en,
			strategic_intiatives.title_ar AS strategic_intiative_title_ar,
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
			LEFT JOIN strategy_profiles ON strategy_profiles.id = strategy_profile_strategic_intiatives.strategy_profile_id
			LEFT JOIN strategy_profile_strategic_objectives ON strategy_profile_strategic_objectives.id = strategy_profile_strategic_intiatives.strategy_profile_strategic_objective_id
			LEFT JOIN strategic_objectives ON strategic_objectives.id = strategy_profile_strategic_objectives.strategic_objective_id
			LEFT JOIN strategic_intiatives ON strategic_intiatives.id = strategy_profile_strategic_intiatives.strategic_intiative_id
			LEFT JOIN users AS created_user ON created_user.id = strategy_profile_strategic_intiatives.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategy_profile_strategic_intiatives.updated_by_id
		`).
		Where("strategy_profile_strategic_intiatives.strategy_profile_id = ?", strategyProfileID)

	// Apply filters
	for column, values := range params.Filters {
		// Map API param to DB column
		if column == "strategic_objective_id" {
			column = "strategy_profile_strategic_intiatives.strategy_profile_strategic_objective_id"
		}
		if len(values) > 0 {
			query = query.Where(column+" IN (?)", values)
		}
	}

	// Apply keyword search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			strategic_intiatives.title_en LIKE ? OR
			strategic_intiatives.title_ar LIKE ? OR
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

	err := query.Find(&strategicIntiatives).Error
	if err != nil {
		return nil, err
	}

	return strategicIntiatives, nil
}

// GetByID fetches a single objective by ID
func (repo *StrategyProfileStrategicIntiativeRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileStrategicIntiative, error) {
	var strategicIntiative *models.StrategyProfileStrategicIntiative

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.First(&strategicIntiative, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return strategicIntiative, nil
}

func (repo *StrategyProfileStrategicIntiativeRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileStrategicIntiative, error) {
	var strategyProfileStrategicIntiative *models.StrategyProfileStrategicIntiative

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Preload("StrategyProfileStrategicObjective").
		Preload("StrategicIntiative").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&strategyProfileStrategicIntiative, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileStrategicIntiative, nil
}

// Create inserts a new objective record
func (repo *StrategyProfileStrategicIntiativeRepository) Create(strategicIntiative *models.StrategyProfileStrategicIntiative, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(strategicIntiative).Error; err != nil {
		return 0, err
	}

	return strategicIntiative.ID, nil
}

// Delete removes an objective record by ID
func (repo *StrategyProfileStrategicIntiativeRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.StrategyProfileStrategicIntiative{}, id).Error; err != nil {
		return err
	}
	return nil
}
