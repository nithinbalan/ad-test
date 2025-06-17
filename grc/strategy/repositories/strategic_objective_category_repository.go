package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type StrategicObjectiveCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StrategicObjectiveCategory, error)
	GetByID(id uint) (*models.StrategicObjectiveCategory, error)
	Create(article *models.StrategicObjectiveCategory) (*models.StrategicObjectiveCategory, error)
	Update(id uint, category *models.StrategicObjectiveCategory) (*models.StrategicObjectiveCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type StrategicObjectiveCategoryRepository struct {
	db *gorm.DB
}

func NewStrategicObjectiveCategoryRepository(db *gorm.DB) *StrategicObjectiveCategoryRepository {
	return &StrategicObjectiveCategoryRepository{
		db: db,
	}
}

func (repo *StrategicObjectiveCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StrategicObjectiveCategory, error) {
	var strategicObjectiveCategories []*models.StrategicObjectiveCategory
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
	result := query.Find(&strategicObjectiveCategories)

	return strategicObjectiveCategories, result.Error
}

func (repo *StrategicObjectiveCategoryRepository) GetByID(id uint) (*models.StrategicObjectiveCategory, error) {
	var category *models.StrategicObjectiveCategory
	result := repo.db.First(&category, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (repo *StrategicObjectiveCategoryRepository) Create(category *models.StrategicObjectiveCategory) (*models.StrategicObjectiveCategory, error) {
	result := repo.db.Create(&category)

	return category, result.Error
}

func (repo *StrategicObjectiveCategoryRepository) Update(id uint, updatedCategory *models.StrategicObjectiveCategory) (*models.StrategicObjectiveCategory, error) {
	var category *models.StrategicObjectiveCategory
	if err := repo.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&category).Updates(updatedCategory).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (repo *StrategicObjectiveCategoryRepository) Delete(id uint) error {
	var category *models.StrategicObjectiveCategory
	if err := repo.db.First(&category, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.StrategicObjectiveCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *StrategicObjectiveCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.StrategicObjectiveCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.StrategicObjectiveCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *StrategicObjectiveCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.StrategicObjectiveCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.StrategicObjectiveCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *StrategicObjectiveCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.StrategicObjectiveCategory{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *StrategicObjectiveCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.StrategicObjectiveCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
