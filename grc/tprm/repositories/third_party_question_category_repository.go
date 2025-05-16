package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ThirdPartyQuestionCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyQuestionCategory, error)
	GetByID(id uint) (*models.ThirdPartyQuestionCategory, error)
	Create(article *models.ThirdPartyQuestionCategory) (*models.ThirdPartyQuestionCategory, error)
	Update(id uint, thirdPartyQuestionCategory *models.ThirdPartyQuestionCategory) (*models.ThirdPartyQuestionCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ThirdPartyQuestionCategoryRepository struct {
	db *gorm.DB
}

func NewThirdPartyQuestionCategoryRepository(db *gorm.DB) *ThirdPartyQuestionCategoryRepository {
	return &ThirdPartyQuestionCategoryRepository{
		db: db,
	}
}

func (repo *ThirdPartyQuestionCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyQuestionCategory, error) {
	var thirdPartyQuestionCategories []*models.ThirdPartyQuestionCategory
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
	result := query.Find(&thirdPartyQuestionCategories)

	return thirdPartyQuestionCategories, result.Error
}

func (repo *ThirdPartyQuestionCategoryRepository) GetByID(id uint) (*models.ThirdPartyQuestionCategory, error) {
	var thirdPartyQuestionCategory *models.ThirdPartyQuestionCategory
	result := repo.db.First(&thirdPartyQuestionCategory, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return thirdPartyQuestionCategory, nil
}

func (repo *ThirdPartyQuestionCategoryRepository) Create(thirdPartyQuestionCategory *models.ThirdPartyQuestionCategory) (*models.ThirdPartyQuestionCategory, error) {
	result := repo.db.Create(&thirdPartyQuestionCategory)

	return thirdPartyQuestionCategory, result.Error
}

func (repo *ThirdPartyQuestionCategoryRepository) Update(id uint, updatedType *models.ThirdPartyQuestionCategory) (*models.ThirdPartyQuestionCategory, error) {
	var thirdPartyQuestionCategory *models.ThirdPartyQuestionCategory
	if err := repo.db.First(&thirdPartyQuestionCategory, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&thirdPartyQuestionCategory).Updates(updatedType).Error; err != nil {
		return nil, err
	}

	return thirdPartyQuestionCategory, nil
}

func (repo *ThirdPartyQuestionCategoryRepository) Delete(id uint) error {
	var thirdPartyQuestionCategory *models.ThirdPartyQuestionCategory
	if err := repo.db.First(&thirdPartyQuestionCategory, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ThirdPartyQuestionCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ThirdPartyQuestionCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ThirdPartyQuestionCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThirdPartyQuestionCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ThirdPartyQuestionCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ThirdPartyQuestionCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThirdPartyQuestionCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *ThirdPartyQuestionCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.ThirdPartyQuestionCategory{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *ThirdPartyQuestionCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ThirdPartyQuestionCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
