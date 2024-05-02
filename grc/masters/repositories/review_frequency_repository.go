package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ReviewFrequencyRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewFrequency, error)
	GetByID(id uint) (*models.ReviewFrequency, error)
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

	var reviewFrequencies []*models.ReviewFrequency

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

func (repo *ReviewFrequencyRepository) GetByID(id uint) (*models.ReviewFrequency, error) {
	var reviewFrequency *models.ReviewFrequency
	result := repo.db.First(&reviewFrequency, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return reviewFrequency, result.Error
}
