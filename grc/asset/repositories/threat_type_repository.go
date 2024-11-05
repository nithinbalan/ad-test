package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ThreatTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThreatType, error)
	GetByID(id uint) (*models.ThreatType, error)
	Create(threatType *models.ThreatType) (*models.ThreatType, error)
	Update(id uint, threatType *models.ThreatType) (*models.ThreatType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type ThreatTypeRepository struct {
	db *gorm.DB
}

func NewThreatTypeRepository(db *gorm.DB) *ThreatTypeRepository {
	return &ThreatTypeRepository{
		db: db,
	}
}

func (repo *ThreatTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThreatType, error) {
	var threatTypes []*models.ThreatType
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
	result := query.Find(&threatTypes)

	return threatTypes, result.Error
}

func (repo *ThreatTypeRepository) GetByID(id uint) (*models.ThreatType, error) {
	var threatType *models.ThreatType
	result := repo.db.First(&threatType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return threatType, nil
}

func (repo *ThreatTypeRepository) Create(threatType *models.ThreatType) (*models.ThreatType, error) {
	result := repo.db.Create(&threatType)

	return threatType, result.Error
}

func (repo *ThreatTypeRepository) Update(id uint, updatedType *models.ThreatType) (*models.ThreatType, error) {
	var threatType *models.ThreatType
	if err := repo.db.First(&threatType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&threatType).Updates(updatedType).Error; err != nil {
		return nil, err
	}

	return threatType, nil
}

func (repo *ThreatTypeRepository) Delete(id uint) error {
	var threatType *models.ThreatType
	if err := repo.db.First(&threatType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ThreatType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ThreatTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.ThreatType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThreatType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *ThreatTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.ThreatType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThreatType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *ThreatTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ThreatType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *ThreatTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ThreatType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
