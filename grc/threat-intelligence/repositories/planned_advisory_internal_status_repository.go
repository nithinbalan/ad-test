package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type PlannedAdvisoryInternalStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryInternalStatus, error)
	GetByID(id uint) (*models.PlannedAdvisoryInternalStatus, error)
	Create(internalStatus *models.PlannedAdvisoryInternalStatus) (*models.PlannedAdvisoryInternalStatus, error)
	Update(id uint, internalStatus *models.PlannedAdvisoryInternalStatus) (*models.PlannedAdvisoryInternalStatus, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

// implementation
type PlannedAdvisoryInternalStatusRepository struct {
	db *gorm.DB
}

func NewPlannedAdvisoryInternalStatusRepository(db *gorm.DB) *PlannedAdvisoryInternalStatusRepository {
	return &PlannedAdvisoryInternalStatusRepository{
		db: db,
	}
}

func (repo *PlannedAdvisoryInternalStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryInternalStatus, error) {

	var plannedAdvisoryInternalStatuses []*models.PlannedAdvisoryInternalStatus

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
	result := query.Find(&plannedAdvisoryInternalStatuses)

	return plannedAdvisoryInternalStatuses, result.Error
}

func (repo *PlannedAdvisoryInternalStatusRepository) GetByID(id uint) (*models.PlannedAdvisoryInternalStatus, error) {
	var internalStatus *models.PlannedAdvisoryInternalStatus
	result := repo.db.First(&internalStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return internalStatus, result.Error
}

func (repo *PlannedAdvisoryInternalStatusRepository) Create(internalStatus *models.PlannedAdvisoryInternalStatus) (*models.PlannedAdvisoryInternalStatus, error) {
	result := repo.db.Create(&internalStatus)

	return internalStatus, result.Error
}

func (repo *PlannedAdvisoryInternalStatusRepository) Update(id uint, updatedInternalStatus *models.PlannedAdvisoryInternalStatus) (*models.PlannedAdvisoryInternalStatus, error) {
	var internalStatus *models.PlannedAdvisoryInternalStatus

	if err := repo.db.First(&internalStatus, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&internalStatus).Updates(updatedInternalStatus).Error; err != nil {
		return nil, err
	}

	return internalStatus, nil
}

func (repo *PlannedAdvisoryInternalStatusRepository) Delete(id uint) error {

	var internalStatus *models.PlannedAdvisoryInternalStatus
	if err := repo.db.First(&internalStatus, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.PlannedAdvisoryInternalStatus{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *PlannedAdvisoryInternalStatusRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.PlannedAdvisoryInternalStatus{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.PlannedAdvisoryInternalStatus{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *PlannedAdvisoryInternalStatusRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.PlannedAdvisoryInternalStatus{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.PlannedAdvisoryInternalStatus{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *PlannedAdvisoryInternalStatusRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.PlannedAdvisoryInternalStatus{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *PlannedAdvisoryInternalStatusRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.PlannedAdvisoryInternalStatus{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
