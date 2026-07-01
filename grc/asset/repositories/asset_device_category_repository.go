package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetDeviceCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetDeviceCategory, error)
	GetByID(id uint) (*models.AssetDeviceCategory, error)
	Create(assetDeviceCategory *models.AssetDeviceCategory) (*models.AssetDeviceCategory, error)
	Update(id uint, assetDeviceCategory *models.AssetDeviceCategory) (*models.AssetDeviceCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type AssetDeviceCategoryRepository struct {
	db *gorm.DB
}

func NewAssetDeviceCategoryRepository(db *gorm.DB) *AssetDeviceCategoryRepository {
	return &AssetDeviceCategoryRepository{
		db: db,
	}
}

func (repo *AssetDeviceCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetDeviceCategory, error) {
	var assetDeviceCategories []*models.AssetDeviceCategory
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

	result := query.Order(sortBy + " " + sort).Find(&assetDeviceCategories)

	return assetDeviceCategories, result.Error
}

func (repo *AssetDeviceCategoryRepository) GetByID(id uint) (*models.AssetDeviceCategory, error) {
	var assetDeviceCategory models.AssetDeviceCategory

	if err := repo.db.First(&assetDeviceCategory, id).Error; err != nil {
		return nil, err
	}

	return &assetDeviceCategory, nil
}

func (repo *AssetDeviceCategoryRepository) Create(assetDeviceCategory *models.AssetDeviceCategory) (*models.AssetDeviceCategory, error) {
	if err := repo.db.Create(assetDeviceCategory).Error; err != nil {
		return nil, err
	}

	return assetDeviceCategory, nil
}

func (repo *AssetDeviceCategoryRepository) Update(id uint, updatedCategory *models.AssetDeviceCategory) (*models.AssetDeviceCategory, error) {
	var assetDeviceCategory models.AssetDeviceCategory

	if err := repo.db.First(&assetDeviceCategory, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&assetDeviceCategory).Updates(updatedCategory).Error; err != nil {
		return nil, err
	}

	return &assetDeviceCategory, nil
}

func (repo *AssetDeviceCategoryRepository) Delete(id uint) error {
	var assetDeviceCategory models.AssetDeviceCategory

	if err := repo.db.First(&assetDeviceCategory, id).Error; err != nil {
		return err
	}

	return repo.db.Delete(&assetDeviceCategory).Error
}

func (repo *AssetDeviceCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.AssetDeviceCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AssetDeviceCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *AssetDeviceCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.AssetDeviceCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AssetDeviceCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *AssetDeviceCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.AssetDeviceCategory{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *AssetDeviceCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.AssetDeviceCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}