package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type IncidentTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentType, error)
	GetByID(id uint) (*models.IncidentType, error)
	Create(article *models.IncidentType) (*models.IncidentType, error)
	Update(id uint, incidentType *models.IncidentType) (*models.IncidentType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type IncidentTypeRepository struct {
	db *gorm.DB
}

func NewIncidentTypeRepository(db *gorm.DB) *IncidentTypeRepository {
	return &IncidentTypeRepository{
		db: db,
	}
}

func (repo *IncidentTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentType, error) {
	var incidentTypes []*models.IncidentType
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
	result := query.Find(&incidentTypes)

	return incidentTypes, result.Error
}

func (repo *IncidentTypeRepository) GetByID(id uint) (*models.IncidentType, error) {
	var incidentType *models.IncidentType
	result := repo.db.First(&incidentType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentType, nil
}

func (repo *IncidentTypeRepository) Create(incidentType *models.IncidentType) (*models.IncidentType, error) {
	result := repo.db.Create(&incidentType)

	return incidentType, result.Error
}

func (repo *IncidentTypeRepository) Update(id uint, updatedType *models.IncidentType) (*models.IncidentType, error) {
	var incidentType *models.IncidentType
	if err := repo.db.First(&incidentType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&incidentType).Updates(updatedType).Error; err != nil {
		return nil, err
	}

	return incidentType, nil
}

func (repo *IncidentTypeRepository) Delete(id uint) error {
	var incidentType *models.IncidentType
	if err := repo.db.First(&incidentType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.IncidentType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.IncidentType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *IncidentTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.IncidentType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *IncidentTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.IncidentType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *IncidentTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.IncidentType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
