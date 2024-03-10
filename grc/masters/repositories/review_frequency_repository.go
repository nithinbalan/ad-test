package repositories

import (
	"jethur/grc/masters/models"

	"gorm.io/gorm"
)

// Interfaces
type ReviewFrequencyRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewFrequency, error)
}

// implementation
type ReviewFrequencyRepository struct {
	db *gorm.DB
}

func NewReviewFrequencyRepository(db *gorm.DB) *ReviewFrequencyRepository {
	return &ReviewFrequencyRepository{
		db: db,
	}
}

func (repo *ReviewFrequencyRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewFrequency, error) {

	var severityLevels []*models.ReviewFrequency

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
	result := query.Find(&severityLevels)

	return severityLevels, result.Error
}
