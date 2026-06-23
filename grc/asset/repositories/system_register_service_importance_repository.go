package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type SystemRegisterServiceImportanceRepositoryInterface interface {
	GetAll(params *assetDTOs.SystemRegisterServiceImportanceListQueryParamsDTO) ([]*assetDTOs.SystemRegisterServiceImportanceDTO, error)
	GetByID(id uint) (*models.SystemRegisterServiceImportance, error)
	Create(tool *models.SystemRegisterServiceImportance) (uint, error)
	Update(tool *models.SystemRegisterServiceImportance) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

// Implementation
type SystemRegisterServiceImportanceRepository struct {
	db *gorm.DB
}

func NewSystemRegisterServiceImportanceRepository(db *gorm.DB) *SystemRegisterServiceImportanceRepository {
	return &SystemRegisterServiceImportanceRepository{
		db: db,
	}
}

func (repo *SystemRegisterServiceImportanceRepository) GetAll(params *assetDTOs.SystemRegisterServiceImportanceListQueryParamsDTO) ([]*assetDTOs.SystemRegisterServiceImportanceDTO, error) {
	var tools []*assetDTOs.SystemRegisterServiceImportanceDTO

	query := repo.db.Table("system_register_service_importances").
		Select(`
            system_register_service_importances.*`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		query = query.Where(`
				(system_register_service_importances.title_en LIKE ? OR
				system_register_service_importances.title_ar LIKE ?)`,
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

func (repo *SystemRegisterServiceImportanceRepository) GetByID(id uint) (*models.SystemRegisterServiceImportance, error) {
	var tool *models.SystemRegisterServiceImportance
	result := repo.db.Preload("CreatedByUser").Preload("UpdatedByUser").First(&tool, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return tool, nil
}

func (repo *SystemRegisterServiceImportanceRepository) Create(tool *models.SystemRegisterServiceImportance) (uint, error) {
	if err := repo.db.Create(tool).Error; err != nil {
		return 0, err
	}

	return tool.ID, nil
}

func (repo *SystemRegisterServiceImportanceRepository) Update(updatedTool *models.SystemRegisterServiceImportance) error {
	if err := repo.db.Save(updatedTool).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SystemRegisterServiceImportanceRepository) Delete(id uint) error {
	var tool *models.SystemRegisterServiceImportance
	if err := repo.db.First(&tool, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SystemRegisterServiceImportance{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SystemRegisterServiceImportanceRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.SystemRegisterServiceImportance{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SystemRegisterServiceImportance{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *SystemRegisterServiceImportanceRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.SystemRegisterServiceImportance{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SystemRegisterServiceImportance{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}
