package repositories

import (
	"jethur/grc/masters/models"

	"gorm.io/gorm"
)

type ReviewFrequencyRepositoryInterface interface {
	GetAll() ([]*models.ReviewFrequency, error)
}

type ReviewFrequencyRepository struct {
	db *gorm.DB
}

func NewReviewFrequencyRepository(db *gorm.DB) *ReviewFrequencyRepository {
	return &ReviewFrequencyRepository{
		db: db,
	}
}

func (repo *ReviewFrequencyRepository) GetAll() ([]*models.ReviewFrequency, error) {

	var reviewFrequencies []*models.ReviewFrequency

	result := repo.db.Find(&reviewFrequencies)

	return reviewFrequencies, result.Error
}
