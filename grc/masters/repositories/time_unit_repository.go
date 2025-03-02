package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type TimeUnitRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TimeUnit, error)
	GetByID(id uint) (*models.TimeUnit, error)
}

// implementation
type TimeUnitRepository struct {
	db *gorm.DB
}

func NewTimeUnitRepository(db *gorm.DB) *TimeUnitRepository {
	return &TimeUnitRepository{
		db: db,
	}
}

func (repo *TimeUnitRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TimeUnit, error) {
	var timeUnits []*models.TimeUnit
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
	result := query.Find(&timeUnits)

	return timeUnits, result.Error
}

func (repo *TimeUnitRepository) GetByID(id uint) (*models.TimeUnit, error) {
	var category *models.TimeUnit
	result := repo.db.First(&category, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}
