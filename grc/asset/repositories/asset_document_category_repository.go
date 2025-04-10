package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AssetDocumentCategoryRepositoryInterface interface {
	GetAll(params *assetDTOs.AssetDocumentCategoryListQueryParamsDTO) ([]*assetDTOs.AssetDocumentCategoryDTO, error)
	GetByID(id uint) (*models.AssetDocumentCategory, error)
	Create(category *models.AssetDocumentCategory) (uint, error)
	Update(category *models.AssetDocumentCategory) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type AssetDocumentCategoryRepository struct {
	db *gorm.DB
}

func NewAssetDocumentCategoryRepository(db *gorm.DB) *AssetDocumentCategoryRepository {
	return &AssetDocumentCategoryRepository{
		db: db,
	}
}

func (repo *AssetDocumentCategoryRepository) GetAll(params *assetDTOs.AssetDocumentCategoryListQueryParamsDTO) ([]*assetDTOs.AssetDocumentCategoryDTO, error) {
	var assetDocumentCategories []*assetDTOs.AssetDocumentCategoryDTO

	query := repo.db.Table("asset_document_categories").
		Select(`
            asset_document_categories.*`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		query = query.Where(`
				(asset_document_categories.title_en LIKE ? OR
				asset_document_categories.title_ar LIKE ?)`,
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

	result := query.Find(&assetDocumentCategories)

	return assetDocumentCategories, result.Error
}

func (repo *AssetDocumentCategoryRepository) GetByID(id uint) (*models.AssetDocumentCategory, error) {
	var category *models.AssetDocumentCategory
	result := repo.db.First(&category, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (repo *AssetDocumentCategoryRepository) Create(category *models.AssetDocumentCategory) (uint, error) {
	if err := repo.db.Create(category).Error; err != nil {
		return 0, err
	}

	return category.ID, nil
}

func (repo *AssetDocumentCategoryRepository) Update(updatedCategory *models.AssetDocumentCategory) error {
	if err := repo.db.Save(updatedCategory).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetDocumentCategoryRepository) Delete(id uint) error {
	var category *models.AssetDocumentCategory
	if err := repo.db.First(&category, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.AssetDocumentCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *AssetDocumentCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.AssetDocumentCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AssetDocumentCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *AssetDocumentCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Fetch IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.AssetDocumentCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AssetDocumentCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *AssetDocumentCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.AssetDocumentCategory{}).
		Where("(title_en = ? OR title_ar = ?)", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *AssetDocumentCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.AssetDocumentCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
