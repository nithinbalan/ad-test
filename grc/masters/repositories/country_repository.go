package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type CountryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.Country, error)
	GetByID(id uint) (*models.Country, error)
}

// implementation
type CountryRepository struct {
	db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) *CountryRepository {
	return &CountryRepository{
		db: db,
	}
}

func (repo *CountryRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.Country, error) {
	var countries []*models.Country

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ?)", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&countries)

	return countries, result.Error
}

func (repo *CountryRepository) GetByID(id uint) (*models.Country, error) {
	var country *models.Country
	result := repo.db.First(&country, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return country, result.Error
}
