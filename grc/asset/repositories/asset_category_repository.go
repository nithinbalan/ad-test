package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AssetCategoryRepositoryInterface interface {
	GetAll(params *assetDTOs.AssetCategoryListQueryParamsDTO) ([]*assetDTOs.AssetCategoryDTO, error)
	GetByID(id uint) (*models.AssetCategory, error)
	Create(category *models.AssetCategory) (uint, error)
	Update(category *models.AssetCategory) error
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

func (repo *AssetCategoryRepository) GetAll(params *assetDTOs.AssetCategoryListQueryParamsDTO) ([]*assetDTOs.AssetCategoryDTO, error) {
	var assetCategories []*assetDTOs.AssetCategoryDTO

	query := repo.db.Table("asset_categories").
		Select(`
            asset_categories.*, asset_types.title_en as asset_type_title_en, asset_types.title_ar as asset_type_title_ar`)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		query = query.Where(`
				(asset_categories.title_en LIKE ? OR
				asset_categories.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title_en" // Default sort by title_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	result := query.Find(&assetCategories)

	return assetCategories, result.Error
}

func (repo *AssetCategoryRepository) GetByID(id uint) (*models.AssetCategory, error) {
	var category *models.AssetCategory
	result := repo.db.First(&category, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (repo *AssetCategoryRepository) Create(category *models.AssetCategory) (uint, error) {
	if err := repo.db.Create(category).Error; err != nil {
		return 0, err
	}

	return category.ID, nil
}

func (repo *AssetCategoryRepository) Update(updatedCategory *models.AssetCategory) error {
	if err := repo.db.Save(updatedCategory).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetCategoryRepository) Delete(id uint) error {
	var category *models.AssetCategory
	if err := repo.db.First(&category, id).Error; err != nil {
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
	// Fetch IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.AssetCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AssetCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *AssetCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.AssetCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
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

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.AssetCategory{}).
		Where("(title_en = ? OR title_ar = ?)", titleEn, titleAr).
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
