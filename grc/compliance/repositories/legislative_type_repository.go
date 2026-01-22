package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type LegislationTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LegislationType, error)
	GetByID(id uint) (*models.LegislationType, error)
	Create(body *models.LegislationType) (*models.LegislationType, error)
	Update(id uint, body *models.LegislationType) (*models.LegislationType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type LegislationTypeRepository struct {
	db *gorm.DB
}

func NewLegislationTypeRepository(db *gorm.DB) *LegislationTypeRepository {
	return &LegislationTypeRepository{db: db}
}

// GetAll
func (repo *LegislationTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LegislationType, error) {
	var items []*models.LegislationType
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
	result := query.Find(&items)

	return items, result.Error
}

// GetByID
func (repo *LegislationTypeRepository) GetByID(id uint) (*models.LegislationType, error) {
	var item *models.LegislationType
	result := repo.db.First(&item, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

// Create
func (repo *LegislationTypeRepository) Create(body *models.LegislationType) (*models.LegislationType, error) {
	result := repo.db.Create(&body)
	return body, result.Error
}

// Update
func (repo *LegislationTypeRepository) Update(id uint, updatedBody *models.LegislationType) (*models.LegislationType, error) {
	var body *models.LegislationType

	if err := repo.db.First(&body, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&body).Updates(updatedBody).Error; err != nil {
		return nil, err
	}

	return body, nil
}

// Delete
func (repo *LegislationTypeRepository) Delete(id uint) error {
	var body *models.LegislationType
	if err := repo.db.First(&body, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.LegislationType{}, id)
	return result.Error
}

// Activate
func (repo *LegislationTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.LegislationType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.LegislationType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

// Deactivate
func (repo *LegislationTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.LegislationType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.LegislationType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

// ExistsByTitle
func (repo *LegislationTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.LegislationType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

// ExistsByTitleExcludingID
func (repo *LegislationTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.LegislationType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
