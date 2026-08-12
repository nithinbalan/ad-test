package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type IncidentSensitivityClassificationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentSensitivityClassification, error)
	GetByID(id uint) (*models.IncidentSensitivityClassification, error)
	Create(classification *models.IncidentSensitivityClassification) (*models.IncidentSensitivityClassification, error)
	Update(id uint, classification *models.IncidentSensitivityClassification) (*models.IncidentSensitivityClassification, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type IncidentSensitivityClassificationRepository struct {
	db *gorm.DB
}

func NewIncidentSensitivityClassificationRepository(db *gorm.DB) *IncidentSensitivityClassificationRepository {
	return &IncidentSensitivityClassificationRepository{
		db: db,
	}
}

func (repo *IncidentSensitivityClassificationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentSensitivityClassification, error) {
	var incidentSensitivityClassifications []*models.IncidentSensitivityClassification
	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&incidentSensitivityClassifications)

	return incidentSensitivityClassifications, result.Error
}

func (repo *IncidentSensitivityClassificationRepository) GetByID(id uint) (*models.IncidentSensitivityClassification, error) {
	var classification *models.IncidentSensitivityClassification
	result := repo.db.First(&classification, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return classification, nil
}

func (repo *IncidentSensitivityClassificationRepository) Create(classification *models.IncidentSensitivityClassification) (*models.IncidentSensitivityClassification, error) {
	result := repo.db.Create(&classification)

	return classification, result.Error
}

func (repo *IncidentSensitivityClassificationRepository) Update(id uint, updatedClassification *models.IncidentSensitivityClassification) (*models.IncidentSensitivityClassification, error) {
	var classification *models.IncidentSensitivityClassification
	if err := repo.db.First(&classification, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&classification).Updates(updatedClassification).Error; err != nil {
		return nil, err
	}

	return classification, nil
}

func (repo *IncidentSensitivityClassificationRepository) Delete(id uint) error {
	var classification *models.IncidentSensitivityClassification
	if err := repo.db.First(&classification, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.IncidentSensitivityClassification{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentSensitivityClassificationRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.IncidentSensitivityClassification{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentSensitivityClassification{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *IncidentSensitivityClassificationRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.IncidentSensitivityClassification{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentSensitivityClassification{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *IncidentSensitivityClassificationRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.IncidentSensitivityClassification{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *IncidentSensitivityClassificationRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.IncidentSensitivityClassification{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
