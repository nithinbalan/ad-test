package repositories

import (
	"errors"
	"jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type ViolationReportCategoryRepositoryInterface interface {
	GetAll(params *dtos.ViolationReportCategoryListQueryParamsDTO) ([]*dtos.ViolationReportCategoryListDTO, error)
	GetByID(id uint) (*models.ViolationReportCategory, error)
	Create(vrc *models.ViolationReportCategory) (uint, error)
	Update(vrc *models.ViolationReportCategory) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn, titleAr string, departmentID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string, departmentID uint) (bool, error)
}

// Implementation
type ViolationReportCategoryRepository struct {
	db *gorm.DB
}

func NewViolationReportCategoryRepository(db *gorm.DB) *ViolationReportCategoryRepository {
	return &ViolationReportCategoryRepository{db: db}
}

func (repo *ViolationReportCategoryRepository) GetAll(params *dtos.ViolationReportCategoryListQueryParamsDTO) ([]*dtos.ViolationReportCategoryListDTO, error) {
	var vrcs []*dtos.ViolationReportCategoryListDTO

	query := repo.db.Table("violation_report_categories").
		Select(`
            violation_report_categories.*, 
            department.title_en as department_title_en, 
            department.title_ar as department_title_ar`).
		Joins(`
			LEFT JOIN organization_units AS department ON department.id = violation_report_categories.department_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where(`
				(violation_report_categories.title_en LIKE ? OR
				violation_report_categories.title_ar LIKE ?)`,
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

	result := query.Find(&vrcs)

	return vrcs, result.Error
}

func (repo *ViolationReportCategoryRepository) GetByID(id uint) (*models.ViolationReportCategory, error) {
	var vrc *models.ViolationReportCategory
	result := repo.db.Preload("Department").First(&vrc, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return vrc, result.Error
}

func (repo *ViolationReportCategoryRepository) Create(vrc *models.ViolationReportCategory) (uint, error) {
	if err := repo.db.Create(vrc).Error; err != nil {
		return 0, err
	}
	return vrc.ID, nil
}

func (repo *ViolationReportCategoryRepository) Update(updatedVRC *models.ViolationReportCategory) error {
	if err := repo.db.Save(updatedVRC).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ViolationReportCategoryRepository) Delete(id uint) error {
	var vrc *models.ViolationReportCategory
	if err := repo.db.First(&vrc, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ViolationReportCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ViolationReportCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ViolationReportCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ViolationReportCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ViolationReportCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ViolationReportCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ViolationReportCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *ViolationReportCategoryRepository) ExistsByTitle(titleEn, titleAr string, departmentID uint) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.ViolationReportCategory{}).
		Where("(title_en = ? OR title_ar = ?) AND department_id = ?", titleEn, titleAr, departmentID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *ViolationReportCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string, departmentID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ViolationReportCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?) AND department_id = ?", id, titleEn, titleAr, departmentID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
