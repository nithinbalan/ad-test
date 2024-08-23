package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditEngagementTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditEngagementType, error)
	GetByID(id uint) (*models.AuditEngagementType, error)
	Create(article *models.AuditEngagementType) (*models.AuditEngagementType, error)
	Update(id uint, auditEngagementType *models.AuditEngagementType) (*models.AuditEngagementType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type AuditEngagementTypeRepository struct {
	db *gorm.DB
}

func NewAuditEngagementTypeRepository(db *gorm.DB) *AuditEngagementTypeRepository {
	return &AuditEngagementTypeRepository{
		db: db,
	}
}

func (repo *AuditEngagementTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditEngagementType, error) {

	var auditEngagementTypes []*models.AuditEngagementType

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
	result := query.Find(&auditEngagementTypes)

	return auditEngagementTypes, result.Error
}

func (repo *AuditEngagementTypeRepository) GetByID(id uint) (*models.AuditEngagementType, error) {
	var auditEngagementType *models.AuditEngagementType
	result := repo.db.First(&auditEngagementType, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditEngagementType, result.Error
}

func (repo *AuditEngagementTypeRepository) Create(auditEngagementType *models.AuditEngagementType) (*models.AuditEngagementType, error) {
	result := repo.db.Create(&auditEngagementType)

	return auditEngagementType, result.Error
}

func (repo *AuditEngagementTypeRepository) Update(id uint, updatedEngagementType *models.AuditEngagementType) (*models.AuditEngagementType, error) {
	var auditEngagementType *models.AuditEngagementType

	if err := repo.db.First(&auditEngagementType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&auditEngagementType).Updates(updatedEngagementType).Error; err != nil {
		return nil, err
	}

	return auditEngagementType, nil
}

func (repo *AuditEngagementTypeRepository) Delete(id uint) error {

	var auditEngagementType *models.AuditEngagementType
	if err := repo.db.First(&auditEngagementType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.AuditEngagementType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *AuditEngagementTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.AuditEngagementType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AuditEngagementType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *AuditEngagementTypeRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.AuditEngagementType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AuditEngagementType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *AuditEngagementTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.AuditEngagementType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *AuditEngagementTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.AuditEngagementType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
