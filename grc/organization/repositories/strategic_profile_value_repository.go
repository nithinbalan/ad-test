package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type StrategicProfileValueRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.StrategicProfileValue, error)
	GetByID(id uint) (*models.StrategicProfileValue, error)
	GetFullDetailsByID(id uint) (*models.StrategicProfileValue, error)
	Create(strategicProfileValue *models.StrategicProfileValue) (uint, error)
	Update(strategicProfileValue *models.StrategicProfileValue) error
	Delete(id uint) error
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type StrategicProfileValueRepository struct {
	db *gorm.DB
}

func NewStrategicProfileValueRepository(db *gorm.DB) *StrategicProfileValueRepository {
	return &StrategicProfileValueRepository{
		db: db,
	}
}

func (repo *StrategicProfileValueRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.StrategicProfileValue, error) {
	var strategicProfileValues []*models.StrategicProfileValue

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
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&strategicProfileValues)

	return strategicProfileValues, result.Error
}

func (repo *StrategicProfileValueRepository) GetByID(id uint) (*models.StrategicProfileValue, error) {
	var strategicProfileValue *models.StrategicProfileValue
	result := repo.db.First(&strategicProfileValue, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategicProfileValue, nil
}

func (repo *StrategicProfileValueRepository) GetFullDetailsByID(id uint) (*models.StrategicProfileValue, error) {
	var strategicProfileValue *models.StrategicProfileValue

	result := repo.db.
		Preload("StrategicProfile").
		First(&strategicProfileValue, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategicProfileValue, nil
}

func (repo *StrategicProfileValueRepository) Create(strategicProfileValue *models.StrategicProfileValue) (uint, error) {
	result := repo.db.Create(&strategicProfileValue)

	return strategicProfileValue.ID, result.Error
}

func (repo *StrategicProfileValueRepository) Update(updatedStrategicProfileValue *models.StrategicProfileValue) error {
	if err := repo.db.Save(updatedStrategicProfileValue).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategicProfileValueRepository) Delete(id uint) error {
	var strategicProfileValue *models.StrategicProfileValue
	if err := repo.db.First(&strategicProfileValue, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.StrategicProfileValue{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *StrategicProfileValueRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.StrategicProfileValue{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}

	return count > 0, nil
}

func (repo *StrategicProfileValueRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.StrategicProfileValue{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
