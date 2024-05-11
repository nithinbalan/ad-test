package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type TimezoneRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.Timezone, error)
	GetByID(id uint) (*models.Timezone, error)
}

// Implementation
type TimezoneRepository struct {
	db *gorm.DB
}

func NewTimezoneRepository(db *gorm.DB) *TimezoneRepository {
	return &TimezoneRepository{
		db: db,
	}
}

func (repo *TimezoneRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.Timezone, error) {
	var timezones []*models.Timezone

	query := repo.db

	if keyword != "" {
		query = query.Where("(name LIKE ?)", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "name"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&timezones)

	return timezones, result.Error
}

func (repo *TimezoneRepository) GetByID(id uint) (*models.Timezone, error) {
	var timezone *models.Timezone
	result := repo.db.First(&timezone, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return timezone, result.Error
}
