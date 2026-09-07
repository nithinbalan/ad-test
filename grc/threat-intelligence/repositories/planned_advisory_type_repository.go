package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type PlannedAdvisoryTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryType, error)
	GetByID(id uint) (*models.PlannedAdvisoryType, error)
	Create(advisoryType *models.PlannedAdvisoryType) (*models.PlannedAdvisoryType, error)
	Update(id uint, advisoryType *models.PlannedAdvisoryType) (*models.PlannedAdvisoryType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

// implementation
type PlannedAdvisoryTypeRepository struct {
	db *gorm.DB
}

func NewPlannedAdvisoryTypeRepository(db *gorm.DB) *PlannedAdvisoryTypeRepository {
	return &PlannedAdvisoryTypeRepository{
		db: db,
	}
}

func (repo *PlannedAdvisoryTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryType, error) {

	var plannedAdvisoryTypes []*models.PlannedAdvisoryType

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
	result := query.Find(&plannedAdvisoryTypes)

	return plannedAdvisoryTypes, result.Error
}

func (repo *PlannedAdvisoryTypeRepository) GetByID(id uint) (*models.PlannedAdvisoryType, error) {
	var advisoryType *models.PlannedAdvisoryType
	result := repo.db.First(&advisoryType, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return advisoryType, result.Error
}

func (repo *PlannedAdvisoryTypeRepository) Create(advisoryType *models.PlannedAdvisoryType) (*models.PlannedAdvisoryType, error) {
	result := repo.db.Create(&advisoryType)

	return advisoryType, result.Error
}

func (repo *PlannedAdvisoryTypeRepository) Update(id uint, updatedAdvisoryType *models.PlannedAdvisoryType) (*models.PlannedAdvisoryType, error) {
	var advisoryType *models.PlannedAdvisoryType

	if err := repo.db.First(&advisoryType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&advisoryType).Updates(updatedAdvisoryType).Error; err != nil {
		return nil, err
	}

	return advisoryType, nil
}

func (repo *PlannedAdvisoryTypeRepository) Delete(id uint) error {

	var advisoryType *models.PlannedAdvisoryType
	if err := repo.db.First(&advisoryType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.PlannedAdvisoryType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *PlannedAdvisoryTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.PlannedAdvisoryType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.PlannedAdvisoryType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *PlannedAdvisoryTypeRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.PlannedAdvisoryType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.PlannedAdvisoryType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}
