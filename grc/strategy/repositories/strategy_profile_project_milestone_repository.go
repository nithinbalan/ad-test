package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileProjectMilestoneRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectMilestoneListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectMilestoneListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectMilestone, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectMilestone, error)
	Create(strategyProfileProjectMilestone *models.StrategyProfileProjectMilestone, tx ...utility.DatabaseTransaction) (uint, error)
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type StrategyProfileProjectMilestoneRepository struct {
	db *gorm.DB
}

func NewStrategyProfileProjectMilestoneRepository(db *gorm.DB) *StrategyProfileProjectMilestoneRepository {
	return &StrategyProfileProjectMilestoneRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction
func (repo *StrategyProfileProjectMilestoneRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

// GetAll returns list of objectives with focus area titles and user info
func (repo *StrategyProfileProjectMilestoneRepository) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectMilestoneListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectMilestoneListDTO, error) {
	var strategyProfileProjectMilestones []*strategyDTOs.StrategyProfileProjectMilestoneListDTO
	var totalCount int64

	query := repo.db.Table("strategy_profile_project_milestones").
		Select(`
			strategy_profile_project_milestones.*,
			strategy_profiles.title AS strategy_profile_title,
			strategy_profile_projects.id AS strategy_profile_project_id,
			projects.title_en AS strategy_profile_project_title_en,
			projects.title_ar AS strategy_profile_project_title_ar,
			project_milestones.id AS project_milestone_id,
			project_milestones.title_en AS project_milestone_title_en,
			project_milestones.title_ar AS project_milestone_title_ar,
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
			LEFT JOIN strategy_profiles ON strategy_profiles.id = strategy_profile_project_milestones.strategy_profile_id
			LEFT JOIN strategy_profile_projects ON strategy_profile_projects.id = strategy_profile_project_milestones.strategy_profile_project_id
			LEFT JOIN projects ON projects.id = strategy_profile_projects.project_id
			LEFT JOIN project_milestones ON project_milestones.id = strategy_profile_project_milestones.project_milestone_id
			LEFT JOIN users AS created_user ON created_user.id = strategy_profile_project_milestones.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategy_profile_project_milestones.updated_by_id
		`).
		Where("strategy_profile_project_milestones.strategy_profile_id = ?", strategyProfileID)

	// Apply filters
	for column, values := range params.Filters {
		// Map API param to DB column
		if column == "project_id" {
			column = "strategy_profile_project_milestones.strategy_profile_project_id"
		}
		if len(values) > 0 {
			query = query.Where(column+" IN (?)", values)
		}
	}

	// Apply keyword search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			project_milestones.title_en LIKE ? OR
			project_milestones.title_ar LIKE ? OR
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

	err := query.Find(&strategyProfileProjectMilestones).Error
	if err != nil {
		return nil, err
	}

	return strategyProfileProjectMilestones, nil
}

// GetByID fetches a single project milestone by ID
func (repo *StrategyProfileProjectMilestoneRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectMilestone, error) {
	var projectMilestone *models.StrategyProfileProjectMilestone

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.First(&projectMilestone, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return projectMilestone, nil
}

func (repo *StrategyProfileProjectMilestoneRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileProjectMilestone, error) {
	var projectMilestone *models.StrategyProfileProjectMilestone

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Preload("StrategyProfileProject").
		Preload("ProjectMilestone").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&projectMilestone, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return projectMilestone, nil
}

// Create inserts a new project milestone record
func (repo *StrategyProfileProjectMilestoneRepository) Create(projectMilestone *models.StrategyProfileProjectMilestone, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(projectMilestone).Error; err != nil {
		return 0, err
	}

	return projectMilestone.ID, nil
}

// Delete removes a project milestone record by ID
func (repo *StrategyProfileProjectMilestoneRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.StrategyProfileProjectMilestone{}, id).Error; err != nil {
		return err
	}
	return nil
}
