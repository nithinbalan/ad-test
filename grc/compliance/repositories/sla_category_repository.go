package repositories

import (
	"errors"
	"jethur/grc/compliance/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SLACategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SLACategory, error)
	GetByID(id uint) (*models.SLACategory, error)
	Create(article *models.SLACategory) (*models.SLACategory, error)
	Update(id uint, category *models.SLACategory) (*models.SLACategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type SLACategoryRepository struct {
	db *gorm.DB
}

func NewSLACategoryRepository(db *gorm.DB) *SLACategoryRepository {
	return &SLACategoryRepository{
		db: db,
	}
}

func (repo *SLACategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SLACategory, error) {

	var categories []*models.SLACategory

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
	result := query.Find(&categories)

	return categories, result.Error
}

func (repo *SLACategoryRepository) GetByID(id uint) (*models.SLACategory, error) {
	var category *models.SLACategory
	result := repo.db.First(&category, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return category, result.Error
}

func (repo *SLACategoryRepository) Create(category *models.SLACategory) (*models.SLACategory, error) {
	result := repo.db.Create(&category)

	return category, result.Error
}

func (repo *SLACategoryRepository) Update(id uint, updatedCategory *models.SLACategory) (*models.SLACategory, error) {
	var category *models.SLACategory

	if err := repo.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&category).Updates(updatedCategory).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (repo *SLACategoryRepository) Delete(id uint) error {

	var category *models.SLACategory
	if err := repo.db.First(&category, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SLACategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SLACategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.SLACategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SLACategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *SLACategoryRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.SLACategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SLACategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *SLACategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.SLACategory{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *SLACategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SLACategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
