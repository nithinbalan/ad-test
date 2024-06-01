package repositories

import (
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type IncidentSubCategoryRepositoryInterface interface {
	GetAll(params *incidentDTOs.IncidentSubCategoryListQueryParamsDTO) ([]*incidentDTOs.IncidentSubCategoryDTO, error)
	GetByID(id uint) (*models.IncidentSubCategory, error)
	Create(subCategory *models.IncidentSubCategory) (uint, error)
	Update(subCategory *models.IncidentSubCategory) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string, categoryID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string, categoryID uint) (bool, error)
}

// implementation
type IncidentSubCategoryRepository struct {
	db *gorm.DB
}

func NewIncidentSubCategoryRepository(db *gorm.DB) *IncidentSubCategoryRepository {
	return &IncidentSubCategoryRepository{
		db: db,
	}
}

func (repo *IncidentSubCategoryRepository) GetAll(params *incidentDTOs.IncidentSubCategoryListQueryParamsDTO) ([]*incidentDTOs.IncidentSubCategoryDTO, error) {
	var incidentSubCategories []*incidentDTOs.IncidentSubCategoryDTO

	query := repo.db.Table("incident_sub_categories").
		Select(`
            incident_sub_categories.*, incident_categories.title_en as incident_category_title_en, incident_categories.title_ar as incident_category_title_ar`).
		Joins(`
		    LEFT JOIN incident_categories ON incident_categories.id = incident_sub_categories.incident_category_id
		`)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		query = query.Where(`
				(incident_sub_categories.title_en LIKE ? OR
				incident_sub_categories.title_ar LIKE ?)`,
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

	result := query.Find(&incidentSubCategories)

	return incidentSubCategories, result.Error

}

func (repo *IncidentSubCategoryRepository) GetByID(id uint) (*models.IncidentSubCategory, error) {
	var subCategory *models.IncidentSubCategory
	result := repo.db.Preload("IncidentCategory").First(&subCategory, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return subCategory, nil
}

func (repo *IncidentSubCategoryRepository) Create(subCategory *models.IncidentSubCategory) (uint, error) {
	if err := repo.db.Create(subCategory).Error; err != nil {
		return 0, err
	}

	return subCategory.ID, nil
}

func (repo *IncidentSubCategoryRepository) Update(updatedSubCategory *models.IncidentSubCategory) error {
	if err := repo.db.Save(updatedSubCategory).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentSubCategoryRepository) Delete(id uint) error {
	var subCategory *models.IncidentSubCategory
	if err := repo.db.First(&subCategory, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.IncidentSubCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentSubCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.IncidentSubCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentSubCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *IncidentSubCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.IncidentSubCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentSubCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *IncidentSubCategoryRepository) ExistsByTitle(titleEn, titleAr string, CategoryID uint) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.IncidentSubCategory{}).
		Where("(title_en = ? OR title_ar = ?) AND incident_category_id = ?", titleEn, titleAr, CategoryID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *IncidentSubCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string, CategoryID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.IncidentSubCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?) AND incident_category_id = ?", id, titleEn, titleAr, CategoryID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
