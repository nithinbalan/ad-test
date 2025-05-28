package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type DesignationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Designation, error)
	GetByID(id uint) (*models.Designation, error)
	Create(article *models.Designation) (*models.Designation, error)
	Update(id uint, category *models.Designation) (*models.Designation, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)

	GetByTitleEn(titleEn string) (*models.Designation, error)
}

// implementation
type DesignationRepository struct {
	db *gorm.DB
}

func NewDesignationRepository(db *gorm.DB) *DesignationRepository {
	return &DesignationRepository{
		db: db,
	}
}

func (repo *DesignationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Designation, error) {
	var designations []*models.Designation

	query := repo.db.Preload("JobDescription")

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
	result := query.Find(&designations)

	return designations, result.Error
}

func (repo *DesignationRepository) GetByID(id uint) (*models.Designation, error) {
	var designation *models.Designation
	result := repo.db.First(&designation, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return designation, nil
}

func (repo *DesignationRepository) Create(designation *models.Designation) (*models.Designation, error) {
	result := repo.db.Create(&designation)

	return designation, result.Error
}

func (repo *DesignationRepository) Update(id uint, updatedDesignation *models.Designation) (*models.Designation, error) {
	var designation *models.Designation

	if err := repo.db.First(&designation, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&designation).Updates(updatedDesignation).Error; err != nil {
		return nil, err
	}

	return designation, nil
}

func (repo *DesignationRepository) Delete(id uint) error {

	var designation *models.Designation
	if err := repo.db.First(&designation, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.Designation{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *DesignationRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.Designation{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Designation{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *DesignationRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.Designation{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Designation{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *DesignationRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.Designation{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *DesignationRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.Designation{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *DesignationRepository) GetByTitleEn(titleEn string) (*models.Designation, error) {
	var designation models.Designation

	result := repo.db.Where("title_en = ?", titleEn).First(&designation)
	if result.Error != nil {
		return nil, result.Error
	}

	return &designation, nil
}
