package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type SystemRegisterDisruptionLevelRepositoryInterface interface {
	GetAll(params *assetDTOs.SystemRegisterDisruptionLevelListQueryParamsDTO) ([]*assetDTOs.SystemRegisterDisruptionLevelDTO, error)
	GetByID(id uint) (*models.SystemRegisterDisruptionLevel, error)
	Create(level *models.SystemRegisterDisruptionLevel) (uint, error)
	Update(level *models.SystemRegisterDisruptionLevel) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

// Implementation
type SystemRegisterDisruptionLevelRepository struct {
	db *gorm.DB
}

func NewSystemRegisterDisruptionLevelRepository(db *gorm.DB) *SystemRegisterDisruptionLevelRepository {
	return &SystemRegisterDisruptionLevelRepository{
		db: db,
	}
}

func (repo *SystemRegisterDisruptionLevelRepository) GetAll(params *assetDTOs.SystemRegisterDisruptionLevelListQueryParamsDTO) ([]*assetDTOs.SystemRegisterDisruptionLevelDTO, error) {
	var levels []*assetDTOs.SystemRegisterDisruptionLevelDTO

	query := repo.db.Table("system_register_disruption_levels").
		Select(`
            system_register_disruption_levels.*`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		query = query.Where(`
				(system_register_disruption_levels.title_en LIKE ? OR
				system_register_disruption_levels.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "created_at" // Default sort by created_at if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc" // Default order is descending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	result := query.Find(&levels)

	return levels, result.Error
}

func (repo *SystemRegisterDisruptionLevelRepository) GetByID(id uint) (*models.SystemRegisterDisruptionLevel, error) {
	var level *models.SystemRegisterDisruptionLevel
	result := repo.db.Preload("CreatedByUser").Preload("UpdatedByUser").First(&level, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return level, nil
}

func (repo *SystemRegisterDisruptionLevelRepository) Create(level *models.SystemRegisterDisruptionLevel) (uint, error) {
	if err := repo.db.Create(level).Error; err != nil {
		return 0, err
	}

	return level.ID, nil
}

func (repo *SystemRegisterDisruptionLevelRepository) Update(updatedLevel *models.SystemRegisterDisruptionLevel) error {
	if err := repo.db.Save(updatedLevel).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SystemRegisterDisruptionLevelRepository) Delete(id uint) error {
	var level *models.SystemRegisterDisruptionLevel
	if err := repo.db.First(&level, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SystemRegisterDisruptionLevel{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SystemRegisterDisruptionLevelRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.SystemRegisterDisruptionLevel{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SystemRegisterDisruptionLevel{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *SystemRegisterDisruptionLevelRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.SystemRegisterDisruptionLevel{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SystemRegisterDisruptionLevel{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}
