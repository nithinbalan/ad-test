package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type IncidentEvidenceSourceRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentEvidenceSource, error)
	GetByID(id uint) (*models.IncidentEvidenceSource, error)
	Create(evidenceSource *models.IncidentEvidenceSource) (*models.IncidentEvidenceSource, error)
	Update(id uint, evidenceSource *models.IncidentEvidenceSource) (*models.IncidentEvidenceSource, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type IncidentEvidenceSourceRepository struct {
	db *gorm.DB
}

func NewIncidentEvidenceSourceRepository(db *gorm.DB) *IncidentEvidenceSourceRepository {
	return &IncidentEvidenceSourceRepository{
		db: db,
	}
}

func (repo *IncidentEvidenceSourceRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentEvidenceSource, error) {

	var incidentEvidenceSources []*models.IncidentEvidenceSource

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
	result := query.Find(&incidentEvidenceSources)

	return incidentEvidenceSources, result.Error
}

func (repo *IncidentEvidenceSourceRepository) GetByID(id uint) (*models.IncidentEvidenceSource, error) {
	var evidenceSource *models.IncidentEvidenceSource
	result := repo.db.First(&evidenceSource, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return evidenceSource, result.Error
}

func (repo *IncidentEvidenceSourceRepository) Create(evidenceSource *models.IncidentEvidenceSource) (*models.IncidentEvidenceSource, error) {
	result := repo.db.Create(&evidenceSource)

	return evidenceSource, result.Error
}

func (repo *IncidentEvidenceSourceRepository) Update(id uint, updatedEvidenceSource *models.IncidentEvidenceSource) (*models.IncidentEvidenceSource, error) {
	var evidenceSource *models.IncidentEvidenceSource

	if err := repo.db.First(&evidenceSource, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&evidenceSource).Updates(updatedEvidenceSource).Error; err != nil {
		return nil, err
	}

	return evidenceSource, nil
}

func (repo *IncidentEvidenceSourceRepository) Delete(id uint) error {

	var evidenceSource *models.IncidentEvidenceSource
	if err := repo.db.First(&evidenceSource, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.IncidentEvidenceSource{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentEvidenceSourceRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.IncidentEvidenceSource{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentEvidenceSource{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *IncidentEvidenceSourceRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.IncidentEvidenceSource{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentEvidenceSource{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *IncidentEvidenceSourceRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.IncidentEvidenceSource{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *IncidentEvidenceSourceRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.IncidentEvidenceSource{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
