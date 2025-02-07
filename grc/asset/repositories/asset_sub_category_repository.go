package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AssetSubCategoryRepositoryInterface interface {
	GetAll(params *assetDTOs.AssetSubCategoryListQueryParamsDTO) ([]*assetDTOs.AssetSubCategoryDTO, error)
	GetByID(id uint) (*models.AssetSubCategory, error)
	Create(category *models.AssetSubCategory) (uint, error)
	Update(category *models.AssetSubCategory) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type AssetSubCategoryRepository struct {
	db *gorm.DB
}

func NewAssetSubCategoryRepository(db *gorm.DB) *AssetSubCategoryRepository {
	return &AssetSubCategoryRepository{
		db: db,
	}
}

func (repo *AssetSubCategoryRepository) GetAll(params *assetDTOs.AssetSubCategoryListQueryParamsDTO) ([]*assetDTOs.AssetSubCategoryDTO, error) {
	var assetSubCategories []*assetDTOs.AssetSubCategoryDTO

	query := repo.db.Table("asset_sub_categories").
		Select(`
            asset_sub_categories.*, asset_categories.title_en as asset_category_title_en, asset_categories.title_ar as asset_category_title_ar`).
		Joins(`
		    LEFT JOIN asset_categories ON asset_categories.id = asset_sub_categories.asset_category_id
		`)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		query = query.Where(`
				(asset_sub_categories.title_en LIKE ? OR
				asset_sub_categories.title_ar LIKE ?)`,
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

	result := query.Find(&assetSubCategories)

	return assetSubCategories, result.Error
}

func (repo *AssetSubCategoryRepository) GetByID(id uint) (*models.AssetSubCategory, error) {
	var category *models.AssetSubCategory
	result := repo.db.Preload("AssetCategory").First(&category, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (repo *AssetSubCategoryRepository) Create(category *models.AssetSubCategory) (uint, error) {
	if err := repo.db.Create(category).Error; err != nil {
		return 0, err
	}

	return category.ID, nil
}

func (repo *AssetSubCategoryRepository) Update(updatedCategory *models.AssetSubCategory) error {
	if err := repo.db.Save(updatedCategory).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetSubCategoryRepository) Delete(id uint) error {
	var category *models.AssetSubCategory
	if err := repo.db.First(&category, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.AssetSubCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *AssetSubCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.AssetSubCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AssetSubCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *AssetSubCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.AssetSubCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AssetSubCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *AssetSubCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.AssetSubCategory{}).
		Where("(title_en = ? OR title_ar = ?)", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *AssetSubCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.AssetSubCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
