package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type IntegrationCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IntegrationCategory, error)
}

// implementation
type IntegrationCategoryRepository struct {
	db *gorm.DB
}

func NewIntegrationCategoryRepository(db *gorm.DB) *IntegrationCategoryRepository {
	return &IntegrationCategoryRepository{
		db: db,
	}
}

func (repo *IntegrationCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IntegrationCategory, error) {

	var reviewFrequencies []*models.IntegrationCategory

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
	result := query.Find(&reviewFrequencies)

	return reviewFrequencies, result.Error
}
