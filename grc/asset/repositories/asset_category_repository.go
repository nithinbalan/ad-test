package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AssetCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetCategory, error)
	GetByID(id uint) (*models.AssetCategory, error)
	Create(assetCategory *models.AssetCategory) (*models.AssetCategory, error)
	Update(id uint, assetCategory *models.AssetCategory) (*models.AssetCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type AssetCategoryRepository struct {
	db *gorm.DB
}

func NewAssetCategoryRepository(db *gorm.DB) *AssetCategoryRepository {
	return &AssetCategoryRepository{
		db: db,
	}
}

func (repo *AssetCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetCategory, error) {
	var assetCategories []*models.AssetCategory
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
	result := query.Find(&assetCategories)

	return assetCategories, result.Error
}

func (repo *AssetCategoryRepository) GetByID(id uint) (*models.AssetCategory, error) {
	var assetCategory *models.AssetCategory
	result := repo.db.First(&assetCategory, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return assetCategory, nil
}

func (repo *AssetCategoryRepository) Create(assetCategory *models.AssetCategory) (*models.AssetCategory, error) {
	result := repo.db.Create(&assetCategory)

	return assetCategory, result.Error
}

func (repo *AssetCategoryRepository) Update(id uint, updatedCategory *models.AssetCategory) (*models.AssetCategory, error) {
	var assetCategory *models.AssetCategory
	if err := repo.db.First(&assetCategory, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&assetCategory).Updates(updatedCategory).Error; err != nil {
		return nil, err
	}

	return assetCategory, nil
}

func (repo *AssetCategoryRepository) Delete(id uint) error {
	var assetCategory *models.AssetCategory
	if err := repo.db.First(&assetCategory, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.AssetCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *AssetCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.AssetCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AssetCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *AssetCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.AssetCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AssetCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *AssetCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.AssetCategory{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *AssetCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.AssetCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
