package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type SystemRegisterRecoveryToolRepositoryInterface interface {
	GetAll(params *assetDTOs.SystemRegisterRecoveryToolListQueryParamsDTO) ([]*assetDTOs.SystemRegisterRecoveryToolDTO, error)
	GetByID(id uint) (*models.SystemRegisterRecoveryTool, error)
	Create(tool *models.SystemRegisterRecoveryTool) (uint, error)
	Update(tool *models.SystemRegisterRecoveryTool) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

// Implementation
type SystemRegisterRecoveryToolRepository struct {
	db *gorm.DB
}

func NewSystemRegisterRecoveryToolRepository(db *gorm.DB) *SystemRegisterRecoveryToolRepository {
	return &SystemRegisterRecoveryToolRepository{
		db: db,
	}
}

func (repo *SystemRegisterRecoveryToolRepository) GetAll(params *assetDTOs.SystemRegisterRecoveryToolListQueryParamsDTO) ([]*assetDTOs.SystemRegisterRecoveryToolDTO, error) {
	var tools []*assetDTOs.SystemRegisterRecoveryToolDTO

	query := repo.db.Table("system_register_recovery_tools").
		Select(`
            system_register_recovery_tools.*`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		query = query.Where(`
				(system_register_recovery_tools.title_en LIKE ? OR
				system_register_recovery_tools.title_ar LIKE ?)`,
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

	result := query.Find(&tools)

	return tools, result.Error
}

func (repo *SystemRegisterRecoveryToolRepository) GetByID(id uint) (*models.SystemRegisterRecoveryTool, error) {
	var tool *models.SystemRegisterRecoveryTool
	result := repo.db.Preload("CreatedByUser").Preload("UpdatedByUser").First(&tool, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return tool, nil
}

func (repo *SystemRegisterRecoveryToolRepository) Create(tool *models.SystemRegisterRecoveryTool) (uint, error) {
	if err := repo.db.Create(tool).Error; err != nil {
		return 0, err
	}

	return tool.ID, nil
}

func (repo *SystemRegisterRecoveryToolRepository) Update(updatedTool *models.SystemRegisterRecoveryTool) error {
	if err := repo.db.Save(updatedTool).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SystemRegisterRecoveryToolRepository) Delete(id uint) error {
	var tool *models.SystemRegisterRecoveryTool
	if err := repo.db.First(&tool, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SystemRegisterRecoveryTool{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SystemRegisterRecoveryToolRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.SystemRegisterRecoveryTool{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SystemRegisterRecoveryTool{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *SystemRegisterRecoveryToolRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.SystemRegisterRecoveryTool{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SystemRegisterRecoveryTool{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}
