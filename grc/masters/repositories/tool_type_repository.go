package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ToolTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ToolType, error)
	GetByID(id uint) (*models.ToolType, error)
	Create(article *models.ToolType) (*models.ToolType, error)
	Update(id uint, toolType *models.ToolType) (*models.ToolType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ToolTypeRepository struct {
	db *gorm.DB
}

func NewToolTypeRepository(db *gorm.DB) *ToolTypeRepository {
	return &ToolTypeRepository{
		db: db,
	}
}

func (repo *ToolTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ToolType, error) {

	var toolTypes []*models.ToolType

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE  ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&toolTypes)

	return toolTypes, result.Error
}

func (repo *ToolTypeRepository) GetByID(id uint) (*models.ToolType, error) {
	var toolType *models.ToolType
	result := repo.db.First(&toolType, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return toolType, result.Error
}

func (repo *ToolTypeRepository) Create(toolType *models.ToolType) (*models.ToolType, error) {
	result := repo.db.Create(&toolType)

	return toolType, result.Error
}

func (repo *ToolTypeRepository) Update(id uint, updatedToolType *models.ToolType) (*models.ToolType, error) {
	var toolType *models.ToolType

	if err := repo.db.First(&toolType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&toolType).Updates(updatedToolType).Error; err != nil {
		return nil, err
	}

	return toolType, nil
}

func (repo *ToolTypeRepository) Delete(id uint) error {

	var toolType *models.ToolType
	if err := repo.db.First(&toolType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ToolType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ToolTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ToolType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ToolType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ToolTypeRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ToolType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ToolType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *ToolTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.ToolType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *ToolTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ToolType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
