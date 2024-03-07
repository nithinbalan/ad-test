package repositories

import (
	"errors"
	"jethur/grc/compliance/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ComplianceCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceCategory, error)
	GetByID(id uint) (*models.ComplianceCategory, error)
	Create(article *models.ComplianceCategory) (*models.ComplianceCategory, error)
	Update(id uint, category *models.ComplianceCategory) (*models.ComplianceCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ComplianceCategoryRepository struct {
	db *gorm.DB
}

func NewComplianceCategoryRepository(db *gorm.DB) *ComplianceCategoryRepository {
	return &ComplianceCategoryRepository{
		db: db,
	}
}

func (repo *ComplianceCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceCategory, error) {

	var categories []*models.ComplianceCategory

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

func (repo *ComplianceCategoryRepository) GetByID(id uint) (*models.ComplianceCategory, error) {
	var category *models.ComplianceCategory
	result := repo.db.First(&category, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return category, result.Error
}

func (repo *ComplianceCategoryRepository) Create(category *models.ComplianceCategory) (*models.ComplianceCategory, error) {
	result := repo.db.Create(&category)

	return category, result.Error
}

func (repo *ComplianceCategoryRepository) Update(id uint, updatedCategory *models.ComplianceCategory) (*models.ComplianceCategory, error) {
	var category *models.ComplianceCategory

	if err := repo.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&category).Updates(updatedCategory).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (repo *ComplianceCategoryRepository) Delete(id uint) error {

	var category *models.ComplianceCategory
	if err := repo.db.First(&category, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ComplianceCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ComplianceCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ComplianceCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ComplianceCategoryRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ComplianceCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ComplianceCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *ComplianceCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.ComplianceCategory{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *ComplianceCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ComplianceCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
