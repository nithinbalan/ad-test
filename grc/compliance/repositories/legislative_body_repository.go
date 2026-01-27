package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type LegislativeBodyRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LegislativeBody, error)
	GetByID(id uint) (*models.LegislativeBody, error)
	Create(body *models.LegislativeBody) (*models.LegislativeBody, error)
	Update(id uint, body *models.LegislativeBody) (*models.LegislativeBody, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type LegislativeBodyRepository struct {
	db *gorm.DB
}

func NewLegislativeBodyRepository(db *gorm.DB) *LegislativeBodyRepository {
	return &LegislativeBodyRepository{db: db}
}

// GetAll
func (repo *LegislativeBodyRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LegislativeBody, error) {
	var bodies []*models.LegislativeBody
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
	result := query.Find(&bodies)

	return bodies, result.Error
}

// GetByID
func (repo *LegislativeBodyRepository) GetByID(id uint) (*models.LegislativeBody, error) {
	var body *models.LegislativeBody
	result := repo.db.First(&body, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return body, nil
}

// Create
func (repo *LegislativeBodyRepository) Create(body *models.LegislativeBody) (*models.LegislativeBody, error) {
	result := repo.db.Create(&body)
	return body, result.Error
}

// Update
func (repo *LegislativeBodyRepository) Update(id uint, updatedBody *models.LegislativeBody) (*models.LegislativeBody, error) {
	var body *models.LegislativeBody

	if err := repo.db.First(&body, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&body).Updates(updatedBody).Error; err != nil {
		return nil, err
	}

	return body, nil
}

// Delete
func (repo *LegislativeBodyRepository) Delete(id uint) error {
	var body *models.LegislativeBody
	if err := repo.db.First(&body, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.LegislativeBody{}, id)
	return result.Error
}

// Activate
func (repo *LegislativeBodyRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.LegislativeBody{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.LegislativeBody{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

// Deactivate
func (repo *LegislativeBodyRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.LegislativeBody{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.LegislativeBody{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

// ExistsByTitle
func (repo *LegislativeBodyRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.LegislativeBody{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

// ExistsByTitleExcludingID
func (repo *LegislativeBodyRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.LegislativeBody{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
