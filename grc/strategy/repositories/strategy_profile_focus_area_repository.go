package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileFocusAreaRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileFocusAreaListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileFocusAreaListDTO, error)
	GetByID(id uint) (*models.StrategyProfileFocusArea, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileFocusArea, error)
	Create(strategyProfileFocusArea *models.StrategyProfileFocusArea, tx ...utility.DatabaseTransaction) (uint, error)
	Update(strategyProfileFocusArea *models.StrategyProfileFocusArea, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type StrategyProfileFocusAreaRepository struct {
	db *gorm.DB
}

func NewStrategyProfileFocusAreaRepository(db *gorm.DB) *StrategyProfileFocusAreaRepository {
	return &StrategyProfileFocusAreaRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StrategyProfileFocusAreaRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategyProfileFocusAreaRepository) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileFocusAreaListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileFocusAreaListDTO, error) {
	var strategyProfileFocusAreas []*strategyDTOs.StrategyProfileFocusAreaListDTO
	var totalCount int64

	query := repo.db.Table("strategy_profile_focus_areas").
		Select(`
            strategy_profile_focus_areas.*,
			strategy_profiles.title AS strategy_profile_title,
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
			LEFT JOIN strategy_profiles ON strategy_profiles.id = strategy_profile_focus_areas.strategy_profile_id
			LEFT JOIN users AS created_user ON created_user.id = strategy_profile_focus_areas.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategy_profile_focus_areas.updated_by_id
        `).Where("strategy_profile_focus_areas.strategy_profile_id = ?", strategyProfileID)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
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
			(strategy_profile_focus_areas.title LIKE ? OR
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

	err := query.Find(&strategyProfileFocusAreas).Error
	if err != nil {
		return nil, err
	}
	return strategyProfileFocusAreas, nil
}

func (repo *StrategyProfileFocusAreaRepository) GetByID(id uint) (*models.StrategyProfileFocusArea, error) {
	var strategyProfileFocusArea *models.StrategyProfileFocusArea

	result := repo.db.First(&strategyProfileFocusArea, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileFocusArea, nil

}

func (repo *StrategyProfileFocusAreaRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileFocusArea, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var strategyProfileFocusArea *models.StrategyProfileFocusArea

	result := db.
		Preload("StrategyProfile").
		Preload("StrategyProfile.ResponsibleUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&strategyProfileFocusArea, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileFocusArea, nil
}

func (repo *StrategyProfileFocusAreaRepository) Create(strategyProfileFocusArea *models.StrategyProfileFocusArea, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(strategyProfileFocusArea).Error; err != nil {
		return 0, err
	}

	return strategyProfileFocusArea.ID, nil
}

func (repo *StrategyProfileFocusAreaRepository) Update(strategyProfileFocusArea *models.StrategyProfileFocusArea, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(strategyProfileFocusArea).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategyProfileFocusAreaRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Strategy profile
	if err := db.Delete(&models.StrategyProfileFocusArea{}, id).Error; err != nil {
		return err
	}

	return nil
}
