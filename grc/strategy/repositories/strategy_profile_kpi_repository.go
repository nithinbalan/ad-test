package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileKpiRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileKpiListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileKpiListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileKpi, error)
	Create(kpi *models.StrategyProfileKpi, tx ...utility.DatabaseTransaction) (uint, error)
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetObjectivesForKPI(kpiID uint, tx ...utility.DatabaseTransaction) ([]*models.StrategyProfileStrategicObjective, error)
	GetKPIsForObjective(objectiveID uint, tx ...utility.DatabaseTransaction) ([]*strategyDTOs.KpiRegisterMinimalDetails, error)
}

type StrategyProfileKpiRepository struct {
	db *gorm.DB
}

func NewStrategyProfileKpiRepository(db *gorm.DB) *StrategyProfileKpiRepository {
	return &StrategyProfileKpiRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction
func (repo *StrategyProfileKpiRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

// GetAll returns list of objectives with focus area titles and user info
func (repo *StrategyProfileKpiRepository) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileKpiListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileKpiListDTO, error) {
	var kpis []*strategyDTOs.StrategyProfileKpiListDTO
	var totalCount int64

	query := repo.db.Table("strategy_profile_kpis").
		Select(`
			strategy_profile_kpis.*,
			strategy_profiles.title AS strategy_profile_title,
			strategy_profile_strategic_objectives.id AS strategy_profile_strategic_objective_id,
			strategic_objectives.title_en AS strategy_profile_strategic_objective_title_en,
			strategic_objectives.title_ar AS strategy_profile_strategic_objective_title_ar,
			kpi_register.id AS kpi_id,
			kpi_register.title AS kpi_title,
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
			LEFT JOIN strategy_profiles ON strategy_profiles.id = strategy_profile_kpis.strategy_profile_id
			LEFT JOIN strategy_profile_strategic_objectives ON strategy_profile_strategic_objectives.id = strategy_profile_kpis.strategy_profile_strategic_objective_id
			LEFT JOIN strategic_objectives ON strategic_objectives.id = strategy_profile_strategic_objectives.strategic_objective_id
			LEFT JOIN kpi_register ON kpi_register.id = strategy_profile_kpis.kpi_id
			LEFT JOIN users AS created_user ON created_user.id = strategy_profile_kpis.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategy_profile_kpis.updated_by_id
		`).
		Where("strategy_profile_kpis.strategy_profile_id = ?", strategyProfileID)

	// Apply filters
	for column, values := range params.Filters {
		// Map API param to DB column
		if column == "strategic_objective_id" {
			column = "strategy_profile_kpis.strategy_profile_strategic_objective_id"
		}
		if len(values) > 0 {
			query = query.Where(column+" IN (?)", values)
		}
	}

	// Apply keyword search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			kpi_register.title LIKE ? OR
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

	err := query.Find(&kpis).Error
	if err != nil {
		return nil, err
	}

	return kpis, nil
}

// GetByID fetches a single objective by ID
func (repo *StrategyProfileKpiRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileKpi, error) {
	var kpi models.StrategyProfileKpi
	db := repo.db

	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.First(&kpi, id).Error; err != nil {
		return nil, err
	}

	return &kpi, nil
}

// Create inserts a new objective record
func (repo *StrategyProfileKpiRepository) Create(kpi *models.StrategyProfileKpi, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(kpi).Error; err != nil {
		return 0, err
	}

	return kpi.ID, nil
}

// Delete removes an objective record by ID
func (repo *StrategyProfileKpiRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.StrategyProfileKpi{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StrategyProfileKpiRepository) GetObjectivesForKPI(kpiID uint, tx ...utility.DatabaseTransaction) ([]*models.StrategyProfileStrategicObjective, error) {
	var objectives []*models.StrategyProfileStrategicObjective

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.
		Joins("JOIN strategy_profile_kpis ON strategy_profile_kpis.strategy_profile_strategic_objective_id = strategy_profile_strategic_objectives.id").
		Where("strategy_profile_kpis.kpi_id = ?", kpiID).
		Find(&objectives).Error

	if err != nil {
		return nil, err
	}
	return objectives, nil
}
func (repo *StrategyProfileKpiRepository) GetKPIsForObjective(objectiveID uint, tx ...utility.DatabaseTransaction) ([]*strategyDTOs.KpiRegisterMinimalDetails, error) {
	var kpis []*strategyDTOs.KpiRegisterMinimalDetails
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Select only the required fields from both tables
	err := db.
		Table("strategy_profile_kpis").
		Select(`
            kpi_register.id as "kpi_id",
            kpi_register.achievement as "achievement",
            kpi_register.weightage as "weightage"
        `).
		Joins("LEFT JOIN kpi_register ON kpi_register.id = strategy_profile_kpis.kpi_id").
		Where("strategy_profile_kpis.strategy_profile_strategic_objective_id = ?", objectiveID).
		Find(&kpis).Error

	if err != nil {
		return nil, err
	}
	return kpis, nil
}
