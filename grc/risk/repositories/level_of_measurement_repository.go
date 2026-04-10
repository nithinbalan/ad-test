package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type LevelOfMeasurementRepositoryInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.LevelOfMeasurement, error)
	GetByID(id uint) (*models.LevelOfMeasurement, error)
}

type LevelOfMeasurementRepository struct {
	db *gorm.DB
}

func NewLevelOfMeasurementRepository(db *gorm.DB) LevelOfMeasurementRepositoryInterface {
	return &LevelOfMeasurementRepository{db: db}
}

func (repo *LevelOfMeasurementRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.LevelOfMeasurement, error) {
	var data []*models.LevelOfMeasurement
	query := repo.db.Model(&models.LevelOfMeasurement{})

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	err := query.Order(sortBy + " " + sort).Find(&data).Error
	return data, err
}

func (repo *LevelOfMeasurementRepository) GetByID(id uint) (*models.LevelOfMeasurement, error) {
	var data models.LevelOfMeasurement
	err := repo.db.First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}