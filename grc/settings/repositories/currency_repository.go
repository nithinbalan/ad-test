package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type CurrencyRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.Currency, error)
	GetByID(id uint) (*models.Currency, error)
}

// implementation
type CurrencyRepository struct {
	db *gorm.DB
}

func NewCurrencyRepository(db *gorm.DB) *CurrencyRepository {
	return &CurrencyRepository{
		db: db,
	}
}

func (repo *CurrencyRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.Currency, error) {
	var currencies []*models.Currency

	query := repo.db

	if keyword != "" {
		query = query.Where("(title LIKE ?)", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&currencies)

	return currencies, result.Error
}

func (repo *CurrencyRepository) GetByID(id uint) (*models.Currency, error) {
	var currency *models.Currency
	result := repo.db.First(&currency, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return currency, result.Error
}
