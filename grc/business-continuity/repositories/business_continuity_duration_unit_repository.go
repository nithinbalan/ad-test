package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type BusinessContinuityDurationUnitRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BusinessContinuityDurationUnit, error)
	GetByID(id uint) (*models.BusinessContinuityDurationUnit, error)
	GetByType(durationType string) (*models.BusinessContinuityDurationUnit, error)
}

// implementation
type BusinessContinuityDurationUnitRepository struct {
	db *gorm.DB
}

func NewBusinessContinuityDurationUnitRepository(db *gorm.DB) *BusinessContinuityDurationUnitRepository {
	return &BusinessContinuityDurationUnitRepository{
		db: db,
	}
}

func (bcdur *BusinessContinuityDurationUnitRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BusinessContinuityDurationUnit, error) {
	var businessContinuityDurationUnits []*models.BusinessContinuityDurationUnit
	query := bcdur.db.Model(&models.BusinessContinuityDurationUnit{})

	if keyword != "" {
		query = query.Where("title_en LIKE ? OR title_ar LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if sort != "" && sortBy != "" {
		query = query.Order(sortBy + " " + sort)
	} else {
		query = query.Order("id ASC")
	}

	err := query.Find(&businessContinuityDurationUnits).Error
	return businessContinuityDurationUnits, err
}

func (bcdur *BusinessContinuityDurationUnitRepository) GetByID(id uint) (*models.BusinessContinuityDurationUnit, error) {
	var businessContinuityDurationUnit models.BusinessContinuityDurationUnit
	err := bcdur.db.Where("id = ?", id).First(&businessContinuityDurationUnit).Error
	return &businessContinuityDurationUnit, err
}

func (bcdur *BusinessContinuityDurationUnitRepository) GetByType(durationType string) (*models.BusinessContinuityDurationUnit, error) {
	var businessContinuityDurationUnit models.BusinessContinuityDurationUnit
	err := bcdur.db.Where("type = ?", durationType).First(&businessContinuityDurationUnit).Error
	return &businessContinuityDurationUnit, err
}
