package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AssetTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetType, error)
	GetByID(id uint) (*models.AssetType, error)
	Create(assetType *models.AssetType) (*models.AssetType, error)
	Update(id uint, assetType *models.AssetType) (*models.AssetType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type AssetTypeRepository struct {
	db *gorm.DB
}

func NewAssetTypeRepository(db *gorm.DB) *AssetTypeRepository {
	return &AssetTypeRepository{
		db: db,
	}
}

func (repo *AssetTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetType, error) {
	var assetTypes []*models.AssetType
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
	result := query.Find(&assetTypes)

	return assetTypes, result.Error
}

func (repo *AssetTypeRepository) GetByID(id uint) (*models.AssetType, error) {
	var assetType *models.AssetType
	result := repo.db.First(&assetType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return assetType, nil
}

func (repo *AssetTypeRepository) Create(assetType *models.AssetType) (*models.AssetType, error) {
	result := repo.db.Create(&assetType)

	return assetType, result.Error
}

func (repo *AssetTypeRepository) Update(id uint, updatedType *models.AssetType) (*models.AssetType, error) {
	var assetType *models.AssetType
	if err := repo.db.First(&assetType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&assetType).Updates(updatedType).Error; err != nil {
		return nil, err
	}

	return assetType, nil
}

func (repo *AssetTypeRepository) Delete(id uint) error {
	var assetType *models.AssetType
	if err := repo.db.First(&assetType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.AssetType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *AssetTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are 'inactive' and within the provided 'ids'
	if err := repo.db.Model(&models.AssetType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AssetType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *AssetTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are 'active' and within the provided 'ids'
	if err := repo.db.Model(&models.AssetType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AssetType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *AssetTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.AssetType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *AssetTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.AssetType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
