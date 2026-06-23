package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AwarenessReportSubCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessReportSubCategory, error)
	GetByID(id uint) (*models.AwarenessReportSubCategory, error)
	Create(subCategory *models.AwarenessReportSubCategory) (*models.AwarenessReportSubCategory, error)
	Update(id uint, subCategory *models.AwarenessReportSubCategory) (*models.AwarenessReportSubCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type AwarenessReportSubCategoryRepository struct {
	db *gorm.DB
}

func NewAwarenessReportSubCategoryRepository(db *gorm.DB) *AwarenessReportSubCategoryRepository {
	return &AwarenessReportSubCategoryRepository{
		db: db,
	}
}

func (repo *AwarenessReportSubCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessReportSubCategory, error) {
	var subCategories []*models.AwarenessReportSubCategory

	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
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

	result := query.Order(sortBy + " " + sort).Find(&subCategories)

	return subCategories, result.Error
}

func (repo *AwarenessReportSubCategoryRepository) GetByID(id uint) (*models.AwarenessReportSubCategory, error) {
	var subCategory models.AwarenessReportSubCategory

	if err := repo.db.First(&subCategory, id).Error; err != nil {
		return nil, err
	}

	return &subCategory, nil
}

func (repo *AwarenessReportSubCategoryRepository) Create(subCategory *models.AwarenessReportSubCategory) (*models.AwarenessReportSubCategory, error) {
	if err := repo.db.Create(subCategory).Error; err != nil {
		return nil, err
	}

	return subCategory, nil
}

func (repo *AwarenessReportSubCategoryRepository) Update(id uint, updatedSubCategory *models.AwarenessReportSubCategory) (*models.AwarenessReportSubCategory, error) {
	var subCategory models.AwarenessReportSubCategory

	if err := repo.db.First(&subCategory, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&subCategory).Updates(updatedSubCategory).Error; err != nil {
		return nil, err
	}

	return &subCategory, nil
}

func (repo *AwarenessReportSubCategoryRepository) Delete(id uint) error {
	var subCategory models.AwarenessReportSubCategory

	if err := repo.db.First(&subCategory, id).Error; err != nil {
		return err
	}

	return repo.db.Delete(&subCategory).Error
}

func (repo *AwarenessReportSubCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.AwarenessReportSubCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AwarenessReportSubCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *AwarenessReportSubCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.AwarenessReportSubCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AwarenessReportSubCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *AwarenessReportSubCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.AwarenessReportSubCategory{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *AwarenessReportSubCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.AwarenessReportSubCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}