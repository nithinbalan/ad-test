package repositories

import (
	bcmDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type BCMExecutionRecordTypeRepositoryInterface interface {
	GetAll(params *bcmDTOs.BCMExecutionRecordTypeListQueryParamsDTO) ([]*bcmDTOs.BCMExecutionRecordTypeDTO, error)
	GetByID(id uint) (*models.BCMExecutionRecordType, error)
	Create(tool *models.BCMExecutionRecordType) (uint, error)
	Update(tool *models.BCMExecutionRecordType) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

// Implementation
type BCMExecutionRecordTypeRepository struct {
	db *gorm.DB
}

func NewBCMExecutionRecordTypeRepository(db *gorm.DB) *BCMExecutionRecordTypeRepository {
	return &BCMExecutionRecordTypeRepository{
		db: db,
	}
}

func (repo *BCMExecutionRecordTypeRepository) GetAll(params *bcmDTOs.BCMExecutionRecordTypeListQueryParamsDTO) ([]*bcmDTOs.BCMExecutionRecordTypeDTO, error) {
	var tools []*bcmDTOs.BCMExecutionRecordTypeDTO

	query := repo.db.Table("bcm_execution_record_types").
		Select(`
            bcm_execution_record_types.*`)

	if params.Keyword != "" {
		query = query.Where(`
				(bcm_execution_record_types.title_en LIKE ? OR
				bcm_execution_record_types.title_ar LIKE ?)`,
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

func (repo *BCMExecutionRecordTypeRepository) GetByID(id uint) (*models.BCMExecutionRecordType, error) {
	var tool *models.BCMExecutionRecordType
	result := repo.db.Preload("CreatedByUser").Preload("UpdatedByUser").First(&tool, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return tool, nil
}

func (repo *BCMExecutionRecordTypeRepository) Create(tool *models.BCMExecutionRecordType) (uint, error) {
	if err := repo.db.Create(tool).Error; err != nil {
		return 0, err
	}

	return tool.ID, nil
}

func (repo *BCMExecutionRecordTypeRepository) Update(updatedTool *models.BCMExecutionRecordType) error {
	if err := repo.db.Save(updatedTool).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BCMExecutionRecordTypeRepository) Delete(id uint) error {
	var tool *models.BCMExecutionRecordType
	if err := repo.db.First(&tool, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.BCMExecutionRecordType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *BCMExecutionRecordTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.BCMExecutionRecordType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.BCMExecutionRecordType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *BCMExecutionRecordTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.BCMExecutionRecordType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.BCMExecutionRecordType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}
