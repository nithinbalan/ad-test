package repositories

import (
	"errors"
	"jethur/grc/organization/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ControlCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlCategory, error)
	GetByID(id uint) (*models.ControlCategory, error)
	Create(controlCategory *models.ControlCategory) (*models.ControlCategory, error)
	Update(id uint, updatedControlCategory *models.ControlCategory) (*models.ControlCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ControlCategoryRepository struct {
	db *gorm.DB
}

func NewControlCategoryRepository(db *gorm.DB) *ControlCategoryRepository {
	return &ControlCategoryRepository{
		db: db,
	}
}

func (repo *ControlCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlCategory, error) {

	var controlCategories []*models.ControlCategory

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
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
	result := query.Find(&controlCategories)

	return controlCategories, result.Error
}

func (repo *ControlCategoryRepository) GetByID(id uint) (*models.ControlCategory, error) {
	var controlCategory *models.ControlCategory
	result := repo.db.First(&controlCategory, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return controlCategory, result.Error
}

func (repo *ControlCategoryRepository) Create(controlCategory *models.ControlCategory) (*models.ControlCategory, error) {
	result := repo.db.Create(&controlCategory)

	return controlCategory, result.Error
}

func (repo *ControlCategoryRepository) Update(id uint, updatedControlCategory *models.ControlCategory) (*models.ControlCategory, error) {
	var controlCategory *models.ControlCategory

	if err := repo.db.First(&controlCategory, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&controlCategory).Updates(updatedControlCategory).Error; err != nil {
		return nil, err
	}

	return controlCategory, nil
}

func (repo *ControlCategoryRepository) Delete(id uint) error {

	var controlCategory *models.ControlCategory
	if err := repo.db.First(&controlCategory, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ControlCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ControlCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ControlCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ControlCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ControlCategoryRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ControlCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ControlCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *ControlCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.ControlCategory{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *ControlCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ControlCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
