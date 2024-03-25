package repositories

import (
	"errors"
	"jethur/grc/compliance/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ComplianceRequirementCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceRequirementCategory, error)
	GetByID(id uint) (*models.ComplianceRequirementCategory, error)
	Create(article *models.ComplianceRequirementCategory) (*models.ComplianceRequirementCategory, error)
	Update(id uint, complianceRequirementCategory *models.ComplianceRequirementCategory) (*models.ComplianceRequirementCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ComplianceRequirementCategoryRepository struct {
	db *gorm.DB
}

func NewComplianceRequirementCategoryRepository(db *gorm.DB) *ComplianceRequirementCategoryRepository {
	return &ComplianceRequirementCategoryRepository{
		db: db,
	}
}

func (repo *ComplianceRequirementCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceRequirementCategory, error) {

	var complianceRequirementCategories []*models.ComplianceRequirementCategory

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
	result := query.Find(&complianceRequirementCategories)

	return complianceRequirementCategories, result.Error
}

func (repo *ComplianceRequirementCategoryRepository) GetByID(id uint) (*models.ComplianceRequirementCategory, error) {
	var complianceRequirementCategory *models.ComplianceRequirementCategory
	result := repo.db.First(&complianceRequirementCategory, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return complianceRequirementCategory, result.Error
}

func (repo *ComplianceRequirementCategoryRepository) Create(complianceRequirementCategory *models.ComplianceRequirementCategory) (*models.ComplianceRequirementCategory, error) {
	result := repo.db.Create(&complianceRequirementCategory)

	return complianceRequirementCategory, result.Error
}

func (repo *ComplianceRequirementCategoryRepository) Update(id uint, updatedCategory *models.ComplianceRequirementCategory) (*models.ComplianceRequirementCategory, error) {
	var complianceRequirementCategory *models.ComplianceRequirementCategory

	if err := repo.db.First(&complianceRequirementCategory, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&complianceRequirementCategory).Updates(updatedCategory).Error; err != nil {
		return nil, err
	}

	return complianceRequirementCategory, nil
}

func (repo *ComplianceRequirementCategoryRepository) Delete(id uint) error {

	var complianceRequirementCategory *models.ComplianceRequirementCategory
	if err := repo.db.First(&complianceRequirementCategory, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ComplianceRequirementCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceRequirementCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ComplianceRequirementCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ComplianceRequirementCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ComplianceRequirementCategoryRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ComplianceRequirementCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ComplianceRequirementCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *ComplianceRequirementCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.ComplianceRequirementCategory{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *ComplianceRequirementCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ComplianceRequirementCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
