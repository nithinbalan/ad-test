package repositories

import (
	"errors"
	"jethur/grc/organization/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type StrategicObjectiveRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StrategicObjective, error)
	GetByID(id uint) (*models.StrategicObjective, error)
	Create(strategicObjective *models.StrategicObjective) (*models.StrategicObjective, error)
	Update(id uint, strategicObjective *models.StrategicObjective) (*models.StrategicObjective, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type StrategicObjectiveRepository struct {
	db *gorm.DB
}

func NewStrategicObjectiveRepository(db *gorm.DB) *StrategicObjectiveRepository {
	return &StrategicObjectiveRepository{
		db: db,
	}
}

func (repo *StrategicObjectiveRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StrategicObjective, error) {

	var strategicObjectives []*models.StrategicObjective

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
	result := query.Find(&strategicObjectives)

	return strategicObjectives, result.Error
}

func (repo *StrategicObjectiveRepository) GetByID(id uint) (*models.StrategicObjective, error) {
	var strategicObjective *models.StrategicObjective
	result := repo.db.First(&strategicObjective, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return strategicObjective, result.Error
}

func (repo *StrategicObjectiveRepository) Create(strategicObjective *models.StrategicObjective) (*models.StrategicObjective, error) {
	result := repo.db.Create(&strategicObjective)

	return strategicObjective, result.Error
}

func (repo *StrategicObjectiveRepository) Update(id uint, updatedStrategicObjective *models.StrategicObjective) (*models.StrategicObjective, error) {
	var strategicObjective *models.StrategicObjective

	if err := repo.db.First(&strategicObjective, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&strategicObjective).Updates(updatedStrategicObjective).Error; err != nil {
		return nil, err
	}

	return strategicObjective, nil
}

func (repo *StrategicObjectiveRepository) Delete(id uint) error {

	var strategicObjective *models.StrategicObjective
	if err := repo.db.First(&strategicObjective, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.StrategicObjective{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *StrategicObjectiveRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.StrategicObjective{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.StrategicObjective{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *StrategicObjectiveRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.StrategicObjective{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.StrategicObjective{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *StrategicObjectiveRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.StrategicObjective{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *StrategicObjectiveRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.StrategicObjective{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
