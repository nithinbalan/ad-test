package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type IncidentEvidenceTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentEvidenceType, error)
	GetByID(id uint) (*models.IncidentEvidenceType, error)
	Create(evidenceType *models.IncidentEvidenceType) (*models.IncidentEvidenceType, error)
	Update(id uint, evidenceType *models.IncidentEvidenceType) (*models.IncidentEvidenceType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type IncidentEvidenceTypeRepository struct {
	db *gorm.DB
}

func NewIncidentEvidenceTypeRepository(db *gorm.DB) *IncidentEvidenceTypeRepository {
	return &IncidentEvidenceTypeRepository{
		db: db,
	}
}

func (repo *IncidentEvidenceTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentEvidenceType, error) {

	var incidentEvidenceTypes []*models.IncidentEvidenceType

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
	result := query.Find(&incidentEvidenceTypes)

	return incidentEvidenceTypes, result.Error
}

func (repo *IncidentEvidenceTypeRepository) GetByID(id uint) (*models.IncidentEvidenceType, error) {
	var evidenceType *models.IncidentEvidenceType
	result := repo.db.First(&evidenceType, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return evidenceType, result.Error
}

func (repo *IncidentEvidenceTypeRepository) Create(evidenceType *models.IncidentEvidenceType) (*models.IncidentEvidenceType, error) {
	result := repo.db.Create(&evidenceType)

	return evidenceType, result.Error
}

func (repo *IncidentEvidenceTypeRepository) Update(id uint, updatedEvidenceType *models.IncidentEvidenceType) (*models.IncidentEvidenceType, error) {
	var evidenceType *models.IncidentEvidenceType

	if err := repo.db.First(&evidenceType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&evidenceType).Updates(updatedEvidenceType).Error; err != nil {
		return nil, err
	}

	return evidenceType, nil
}

func (repo *IncidentEvidenceTypeRepository) Delete(id uint) error {

	var evidenceType *models.IncidentEvidenceType
	if err := repo.db.First(&evidenceType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.IncidentEvidenceType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentEvidenceTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.IncidentEvidenceType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentEvidenceType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *IncidentEvidenceTypeRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.IncidentEvidenceType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentEvidenceType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *IncidentEvidenceTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.IncidentEvidenceType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *IncidentEvidenceTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.IncidentEvidenceType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
