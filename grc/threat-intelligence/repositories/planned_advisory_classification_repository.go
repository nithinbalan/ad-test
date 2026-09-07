package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type PlannedAdvisoryClassificationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryClassification, error)
	GetByID(id uint) (*models.PlannedAdvisoryClassification, error)
	Create(classification *models.PlannedAdvisoryClassification) (*models.PlannedAdvisoryClassification, error)
	Update(id uint, classification *models.PlannedAdvisoryClassification) (*models.PlannedAdvisoryClassification, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

// implementation
type PlannedAdvisoryClassificationRepository struct {
	db *gorm.DB
}

func NewPlannedAdvisoryClassificationRepository(db *gorm.DB) *PlannedAdvisoryClassificationRepository {
	return &PlannedAdvisoryClassificationRepository{
		db: db,
	}
}

func (repo *PlannedAdvisoryClassificationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryClassification, error) {

	var plannedAdvisoryClassifications []*models.PlannedAdvisoryClassification

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
	result := query.Find(&plannedAdvisoryClassifications)

	return plannedAdvisoryClassifications, result.Error
}

func (repo *PlannedAdvisoryClassificationRepository) GetByID(id uint) (*models.PlannedAdvisoryClassification, error) {
	var classification *models.PlannedAdvisoryClassification
	result := repo.db.First(&classification, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return classification, result.Error
}

func (repo *PlannedAdvisoryClassificationRepository) Create(classification *models.PlannedAdvisoryClassification) (*models.PlannedAdvisoryClassification, error) {
	result := repo.db.Create(&classification)

	return classification, result.Error
}

func (repo *PlannedAdvisoryClassificationRepository) Update(id uint, updatedClassification *models.PlannedAdvisoryClassification) (*models.PlannedAdvisoryClassification, error) {
	var classification *models.PlannedAdvisoryClassification

	if err := repo.db.First(&classification, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&classification).Updates(updatedClassification).Error; err != nil {
		return nil, err
	}

	return classification, nil
}

func (repo *PlannedAdvisoryClassificationRepository) Delete(id uint) error {

	var classification *models.PlannedAdvisoryClassification
	if err := repo.db.First(&classification, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.PlannedAdvisoryClassification{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *PlannedAdvisoryClassificationRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.PlannedAdvisoryClassification{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.PlannedAdvisoryClassification{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *PlannedAdvisoryClassificationRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.PlannedAdvisoryClassification{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.PlannedAdvisoryClassification{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}
