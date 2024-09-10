package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type UnitRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Unit, error)
	GetByID(id uint) (*models.Unit, error)
	Create(article *models.Unit) (*models.Unit, error)
	Update(id uint, unit *models.Unit) (*models.Unit, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
	ExistsByCode(codeEn string, codeAr string) (bool, error)
	ExistsByCodeExcludingID(id uint, codeEn, codeAr string) (bool, error)
}

// implementation
type UnitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) *UnitRepository {
	return &UnitRepository{
		db: db,
	}
}

func (repo *UnitRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Unit, error) {

	var units []*models.Unit

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ? OR code_en LIKE ? OR code_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
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
	result := query.Find(&units)

	return units, result.Error
}

func (repo *UnitRepository) GetByID(id uint) (*models.Unit, error) {
	var unit *models.Unit
	result := repo.db.First(&unit, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return unit, result.Error
}

func (repo *UnitRepository) Create(unit *models.Unit) (*models.Unit, error) {
	result := repo.db.Create(&unit)

	return unit, result.Error
}

func (repo *UnitRepository) Update(id uint, updatedUnit *models.Unit) (*models.Unit, error) {
	var unit *models.Unit

	if err := repo.db.First(&unit, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&unit).Updates(updatedUnit).Error; err != nil {
		return nil, err
	}

	return unit, nil
}

func (repo *UnitRepository) Delete(id uint) error {

	var unit *models.Unit
	if err := repo.db.First(&unit, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.Unit{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *UnitRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.Unit{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Unit{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *UnitRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.Unit{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Unit{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *UnitRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.Unit{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *UnitRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.Unit{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *UnitRepository) ExistsByCode(codeEn, codeAr string) (bool, error) {
	var count int64

	// Check for the existence of codeEn or codeAr in a single query
	err := repo.db.Model(&models.Unit{}).
		Where("code_en = ? OR code_ar = ?", codeEn, codeAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either code exists, otherwise false
	return count > 0, nil
}

func (repo *UnitRepository) ExistsByCodeExcludingID(id uint, codeEn, codeAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.Unit{}).
		Where("id <> ? AND (code_en = ? OR code_ar = ?)", id, codeEn, codeAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
