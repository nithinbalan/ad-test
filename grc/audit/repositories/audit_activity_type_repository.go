package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditActivityTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditActivityType, error)
	GetByID(id uint) (*models.AuditActivityType, error)
	Create(article *models.AuditActivityType) (*models.AuditActivityType, error)
	Update(id uint, auditActivityType *models.AuditActivityType) (*models.AuditActivityType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type AuditActivityTypeRepository struct {
	db *gorm.DB
}

func NewAuditActivityTypeRepository(db *gorm.DB) *AuditActivityTypeRepository {
	return &AuditActivityTypeRepository{
		db: db,
	}
}

func (repo *AuditActivityTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditActivityType, error) {

	var auditActivityTypes []*models.AuditActivityType

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
	result := query.Find(&auditActivityTypes)

	return auditActivityTypes, result.Error
}

func (repo *AuditActivityTypeRepository) GetByID(id uint) (*models.AuditActivityType, error) {
	var auditActivityType *models.AuditActivityType
	result := repo.db.First(&auditActivityType, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditActivityType, result.Error
}

func (repo *AuditActivityTypeRepository) Create(auditActivityType *models.AuditActivityType) (*models.AuditActivityType, error) {
	result := repo.db.Create(&auditActivityType)

	return auditActivityType, result.Error
}

func (repo *AuditActivityTypeRepository) Update(id uint, updatedAuditActivityType *models.AuditActivityType) (*models.AuditActivityType, error) {
	var auditActivityType *models.AuditActivityType

	if err := repo.db.First(&auditActivityType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&auditActivityType).Updates(updatedAuditActivityType).Error; err != nil {
		return nil, err
	}

	return auditActivityType, nil
}

func (repo *AuditActivityTypeRepository) Delete(id uint) error {

	var auditActivityType *models.AuditActivityType
	if err := repo.db.First(&auditActivityType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.AuditActivityType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *AuditActivityTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.AuditActivityType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AuditActivityType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *AuditActivityTypeRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.AuditActivityType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AuditActivityType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *AuditActivityTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.AuditActivityType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *AuditActivityTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.AuditActivityType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
