package repositories

import (
	bcmDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type BCMProcessTypeRepositoryInterface interface {
	GetAll(params *bcmDTOs.BCMProcessTypeListQueryParamsDTO) ([]*bcmDTOs.BCMProcessTypeDTO, error)
	GetByID(id uint) (*models.BCMProcessType, error)
	Create(tool *models.BCMProcessType) (uint, error)
	Update(tool *models.BCMProcessType) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

// Implementation
type BCMProcessTypeRepository struct {
	db *gorm.DB
}

func NewBCMProcessTypeRepository(db *gorm.DB) *BCMProcessTypeRepository {
	return &BCMProcessTypeRepository{
		db: db,
	}
}

func (repo *BCMProcessTypeRepository) GetAll(params *bcmDTOs.BCMProcessTypeListQueryParamsDTO) ([]*bcmDTOs.BCMProcessTypeDTO, error) {
	var tools []*bcmDTOs.BCMProcessTypeDTO

	query := repo.db.Table("bcm_process_types").
		Select(`
            bcm_process_types.*`)

	if params.Keyword != "" {
		query = query.Where(`
				(bcm_process_types.title_en LIKE ? OR
				bcm_process_types.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	result := query.Find(&tools)

	return tools, result.Error
}

func (repo *BCMProcessTypeRepository) GetByID(id uint) (*models.BCMProcessType, error) {
	var tool *models.BCMProcessType
	result := repo.db.Preload("CreatedByUser").Preload("UpdatedByUser").First(&tool, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return tool, nil
}

func (repo *BCMProcessTypeRepository) Create(tool *models.BCMProcessType) (uint, error) {
	if err := repo.db.Create(tool).Error; err != nil {
		return 0, err
	}

	return tool.ID, nil
}

func (repo *BCMProcessTypeRepository) Update(updatedTool *models.BCMProcessType) error {
	if err := repo.db.Save(updatedTool).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BCMProcessTypeRepository) Delete(id uint) error {
	var tool *models.BCMProcessType
	if err := repo.db.First(&tool, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.BCMProcessType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *BCMProcessTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.BCMProcessType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.BCMProcessType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *BCMProcessTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.BCMProcessType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.BCMProcessType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}
