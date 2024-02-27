package repositories

import (
	"jethur/grc/masters/models"

	"gorm.io/gorm"
)

// Repository Interfaces
type ReviewFrequencyRepositoryInterface interface {
	GetAll() ([]models.ReviewFrequency, error)
	// GetByID(c *gin.Context, id uint) (models.ReviewFrequency, error)
	// Create(c *gin.Context, article models.ReviewFrequency) (models.ReviewFrequency, error)
	// Update(c *gin.Context, id uint, article models.ReviewFrequency) (models.ReviewFrequency, error)
	// Delete(c *gin.Context, id uint) error
}

// Repository & Impelmentation
type ReviewFrequencyRepository struct {
	db *gorm.DB
}

func NewReviewFrequencyRepository(db *gorm.DB) *ReviewFrequencyRepository {
	return &ReviewFrequencyRepository{
		db: db,
	}
}

func (repo *ReviewFrequencyRepository) GetAll() ([]models.ReviewFrequency, error) {

	var reviewFrequencies []models.ReviewFrequency

	result := repo.db.Find(&reviewFrequencies)

	return reviewFrequencies, result.Error
}
