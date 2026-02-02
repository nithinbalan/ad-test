package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type SecrecyClassificationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceSecrecyClassification, error)
	GetByID(id uint) (*models.ComplianceSecrecyClassification, error)
	Create(classification *models.ComplianceSecrecyClassification) (*models.ComplianceSecrecyClassification, error)
	Update(id uint, classification *models.ComplianceSecrecyClassification) (*models.ComplianceSecrecyClassification, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
	ExistsByTypeExcludingID(id uint, t string) (bool, error)
	ExistsByType(titleAr string) (bool, error)
}

// Implementation
type SecrecyClassificationRepository struct {
	db *gorm.DB
}

func NewSecrecyClassificationRepository(db *gorm.DB) *SecrecyClassificationRepository {
	return &SecrecyClassificationRepository{db: db}
}

func (repo *SecrecyClassificationRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.ComplianceSecrecyClassification, error) {
	var classifications []*models.ComplianceSecrecyClassification
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

	result := query.Order(sortBy + " " + sort).Find(&classifications)
	return classifications, result.Error
}

func (repo *SecrecyClassificationRepository) GetByID(id uint) (*models.ComplianceSecrecyClassification, error) {
	var classification *models.ComplianceSecrecyClassification
	result := repo.db.First(&classification, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return classification, nil
}

func (repo *SecrecyClassificationRepository) Create(classification *models.ComplianceSecrecyClassification) (*models.ComplianceSecrecyClassification, error) {
	result := repo.db.Create(&classification)
	return classification, result.Error
}

func (repo *SecrecyClassificationRepository) Update(id uint, updatedClassification *models.ComplianceSecrecyClassification) (*models.ComplianceSecrecyClassification, error) {
	var classification *models.ComplianceSecrecyClassification

	if err := repo.db.First(&classification, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&classification).Updates(updatedClassification).Error; err != nil {
		return nil, err
	}

	return classification, nil
}

func (repo *SecrecyClassificationRepository) Delete(id uint) error {
	var classification *models.ComplianceSecrecyClassification
	if err := repo.db.First(&classification, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ComplianceSecrecyClassification{}, id)
	return result.Error
}

func (repo *SecrecyClassificationRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ComplianceSecrecyClassification{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ComplianceSecrecyClassification{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *SecrecyClassificationRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ComplianceSecrecyClassification{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ComplianceSecrecyClassification{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *SecrecyClassificationRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ComplianceSecrecyClassification{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *SecrecyClassificationRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ComplianceSecrecyClassification{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
func (repo *SecrecyClassificationRepository) ExistsByType(t string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ComplianceSecrecyClassification{}).
		Where("type = ?", t).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
func (repo *SecrecyClassificationRepository) ExistsByTypeExcludingID(id uint, t string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ComplianceSecrecyClassification{}).
		Where("id <> ? AND type = ?", id, t).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
