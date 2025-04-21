package repositories

import (
	"jethur/grc/models"
	"jethur/grc/strategy/dtos"

	"gorm.io/gorm"
)

// Interfaces
type StrategyProfileSettingRepositoryInterface interface {
	GetAll(params *dtos.StrategyProfileSettingListQueryParamsDTO) (*[]dtos.StrategyProfileSettingListDTO, error)
	GetByID(id uint) (*models.StrategyProfileSetting, error)
	GetFullDetailsByID(id uint) (*models.StrategyProfileSetting, error)
	Update(moduleId uint, module *models.StrategyProfileSetting) (uint, error)
}

// Implementation
type StrategyProfileSettingRepository struct {
	db *gorm.DB
}

func NewStrategyProfileSettingRepository(db *gorm.DB) *StrategyProfileSettingRepository {
	return &StrategyProfileSettingRepository{
		db: db,
	}
}

func (repo *StrategyProfileSettingRepository) GetAll(params *dtos.StrategyProfileSettingListQueryParamsDTO) (*[]dtos.StrategyProfileSettingListDTO, error) {
	var settings *[]dtos.StrategyProfileSettingListDTO

	query := repo.db.Table("strategy_profile_settings").
		Select(`
            strategy_profile_settings.*`)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		query = query.Where(`
				(strategy_profile_settings.title_en LIKE ? OR
				strategy_profile_settings.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title_en" // Default sort by title_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "is_enable":
			query = query.Where("strategy_profile_settings.is_enable = ?", true)
		}
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	result := query.Find(&settings)

	return settings, result.Error

}

func (repo *StrategyProfileSettingRepository) GetByID(id uint) (*models.StrategyProfileSetting, error) {
	var module *models.StrategyProfileSetting
	result := repo.db.First(&module, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return module, nil
}

func (repo *StrategyProfileSettingRepository) GetFullDetailsByID(id uint) (*models.StrategyProfileSetting, error) {
	var module *models.StrategyProfileSetting // Declare the variable inside the function

	result := repo.db.
		First(&module, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return module, nil
}

func (repo *StrategyProfileSettingRepository) Update(moduleId uint, module *models.StrategyProfileSetting) (uint, error) {
	db := repo.db
	if err := db.Save(module).Error; err != nil {
		return 0, err
	}

	return moduleId, nil
}
