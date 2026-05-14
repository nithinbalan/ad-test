package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type LocationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Location, error)
	GetByID(id uint) (*models.Location, error)
	Create(location *models.Location) (*models.Location, error)
	Update(id uint, location *models.Location) (*models.Location, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

// implementation
type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{
		db: db,
	}
}

func (repo *LocationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Location, error) {

	var locations []*models.Location

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&locations)

	return locations, result.Error
}

func (repo *LocationRepository) GetByID(id uint) (*models.Location, error) {
	var location *models.Location
	result := repo.db.Preload("CreatedByUser").Preload("UpdatedByUser").First(&location, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return location, result.Error
}

func (repo *LocationRepository) Create(location *models.Location) (*models.Location, error) {
	result := repo.db.Create(&location)

	return location, result.Error
}

func (repo *LocationRepository) Update(id uint, updatedLocation *models.Location) (*models.Location, error) {
	var location *models.Location

	if err := repo.db.First(&location, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&location).Updates(updatedLocation).Error; err != nil {
		return nil, err
	}

	return location, nil
}

func (repo *LocationRepository) Delete(id uint) error {

	var location *models.Location
	if err := repo.db.First(&location, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.Location{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *LocationRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.Location{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Location{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *LocationRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.Location{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Location{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}
