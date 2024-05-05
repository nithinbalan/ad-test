package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ChecklistRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Checklist, error)
	GetByID(id uint) (*models.Checklist, error)
	Create(article *models.Checklist) (*models.Checklist, error)
	Update(id uint, checklist *models.Checklist) (*models.Checklist, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ChecklistRepository struct {
	db *gorm.DB
}

func NewChecklistRepository(db *gorm.DB) *ChecklistRepository {
	return &ChecklistRepository{
		db: db,
	}
}

func (repo *ChecklistRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Checklist, error) {

	var checklists []*models.Checklist

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
	result := query.Find(&checklists)

	return checklists, result.Error
}

func (repo *ChecklistRepository) GetByID(id uint) (*models.Checklist, error) {
	var checklist *models.Checklist
	result := repo.db.First(&checklist, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return checklist, result.Error
}

func (repo *ChecklistRepository) Create(checklist *models.Checklist) (*models.Checklist, error) {
	result := repo.db.Create(&checklist)

	return checklist, result.Error
}

func (repo *ChecklistRepository) Update(id uint, updatedChecklist *models.Checklist) (*models.Checklist, error) {
	var checklist *models.Checklist

	if err := repo.db.First(&checklist, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&checklist).Updates(updatedChecklist).Error; err != nil {
		return nil, err
	}

	return checklist, nil
}

func (repo *ChecklistRepository) Delete(id uint) error {

	var checklist *models.Checklist
	if err := repo.db.First(&checklist, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.Checklist{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ChecklistRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.Checklist{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Checklist{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ChecklistRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.Checklist{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Checklist{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *ChecklistRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.Checklist{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *ChecklistRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.Checklist{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
