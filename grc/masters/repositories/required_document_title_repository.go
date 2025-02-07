package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RequiredDocumentTitleRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RequiredDocumentTitle, error)
	GetByID(id uint) (*models.RequiredDocumentTitle, error)
	Create(document *models.RequiredDocumentTitle) (*models.RequiredDocumentTitle, error)
	Update(id uint, document *models.RequiredDocumentTitle) (*models.RequiredDocumentTitle, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type RequiredDocumentTitleRepository struct {
	db *gorm.DB
}

func NewRequiredDocumentTitleRepository(db *gorm.DB) *RequiredDocumentTitleRepository {
	return &RequiredDocumentTitleRepository{
		db: db,
	}
}

func (repo *RequiredDocumentTitleRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RequiredDocumentTitle, error) {
	var documents []*models.RequiredDocumentTitle

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
	result := query.Find(&documents)

	return documents, result.Error
}

func (repo *RequiredDocumentTitleRepository) GetByID(id uint) (*models.RequiredDocumentTitle, error) {
	var document *models.RequiredDocumentTitle
	result := repo.db.First(&document, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return document, result.Error
}

func (repo *RequiredDocumentTitleRepository) Create(document *models.RequiredDocumentTitle) (*models.RequiredDocumentTitle, error) {
	result := repo.db.Create(&document)

	return document, result.Error
}

func (repo *RequiredDocumentTitleRepository) Update(id uint, updatedDocument *models.RequiredDocumentTitle) (*models.RequiredDocumentTitle, error) {
	var document *models.RequiredDocumentTitle

	if err := repo.db.First(&document, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&document).Updates(updatedDocument).Error; err != nil {
		return nil, err
	}

	return document, nil
}

func (repo *RequiredDocumentTitleRepository) Delete(id uint) error {
	var document *models.RequiredDocumentTitle
	if err := repo.db.First(&document, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.RequiredDocumentTitle{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RequiredDocumentTitleRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RequiredDocumentTitle{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RequiredDocumentTitle{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *RequiredDocumentTitleRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RequiredDocumentTitle{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RequiredDocumentTitle{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	// Return the IDs of the records that were deactivated
	return idsToUpdate, nil
}

func (repo *RequiredDocumentTitleRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.RequiredDocumentTitle{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *RequiredDocumentTitleRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.RequiredDocumentTitle{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
