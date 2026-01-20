package repositories

import (
	"errors"
	"jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type ViolationReportSubCategoryRepositoryInterface interface {
	GetAll(params *dtos.ViolationReportSubCategoryListQueryParamsDTO) ([]*dtos.ViolationReportSubCategoryListDTO, error)
	GetByID(id uint) (*models.ViolationReportSubCategory, error)
	Create(vrsc *models.ViolationReportSubCategory) (uint, error)
	Update(vrsc *models.ViolationReportSubCategory) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn, titleAr string, categoryID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string, categoryID uint) (bool, error)
}

// Implementation
type ViolationReportSubCategoryRepository struct {
	db *gorm.DB
}

func NewViolationReportSubCategoryRepository(db *gorm.DB) *ViolationReportSubCategoryRepository {
	return &ViolationReportSubCategoryRepository{db: db}
}

func (repo *ViolationReportSubCategoryRepository) GetAll(params *dtos.ViolationReportSubCategoryListQueryParamsDTO) ([]*dtos.ViolationReportSubCategoryListDTO, error) {
	var vrscs []*dtos.ViolationReportSubCategoryListDTO

	query := repo.db.Table("violation_report_sub_categories").
		Select(`
            violation_report_sub_categories.*, 
            violation_report_categories.title_en as violation_report_category_title_en, 
            violation_report_categories.title_ar as violation_report_category_title_ar`).
		Joins(`
		    LEFT JOIN violation_report_categories ON violation_report_categories.id = violation_report_sub_categories.violation_report_category_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where(`
				(violation_report_sub_categories.title_en LIKE ? OR
				violation_report_sub_categories.title_ar LIKE ?)`,
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

	result := query.Find(&vrscs)

	return vrscs, result.Error
}

func (repo *ViolationReportSubCategoryRepository) GetByID(id uint) (*models.ViolationReportSubCategory, error) {
	var vrsc *models.ViolationReportSubCategory
	result := repo.db.Preload("ViolationReportCategory").First(&vrsc, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return vrsc, result.Error
}

func (repo *ViolationReportSubCategoryRepository) Create(vrsc *models.ViolationReportSubCategory) (uint, error) {
	if err := repo.db.Create(vrsc).Error; err != nil {
		return 0, err
	}
	return vrsc.ID, nil
}

func (repo *ViolationReportSubCategoryRepository) Update(updatedVRSC *models.ViolationReportSubCategory) error {
	if err := repo.db.Save(updatedVRSC).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ViolationReportSubCategoryRepository) Delete(id uint) error {
	var vrsc *models.ViolationReportSubCategory
	if err := repo.db.First(&vrsc, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ViolationReportSubCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ViolationReportSubCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ViolationReportSubCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ViolationReportSubCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ViolationReportSubCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ViolationReportSubCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ViolationReportSubCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *ViolationReportSubCategoryRepository) ExistsByTitle(titleEn, titleAr string, categoryID uint) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.ViolationReportSubCategory{}).
		Where("(title_en = ? OR title_ar = ?) AND violation_report_category_id = ?", titleEn, titleAr, categoryID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *ViolationReportSubCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string, categoryID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ViolationReportSubCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?) AND violation_report_category_id = ?", id, titleEn, titleAr, categoryID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
