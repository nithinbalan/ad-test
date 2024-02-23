package implementations

import (
	"jethur/grc/masters/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewFrequencyRepository struct {
	db *gorm.DB
}

func NewReviewFrequencyRepository(db *gorm.DB) *ReviewFrequencyRepository {
	return &ReviewFrequencyRepository{
		db: db,
	}
}

func (repo *ReviewFrequencyRepository) GetAll(c *gin.Context) ([]models.ReviewFrequency, error) {

	var reviewFrequencies []models.ReviewFrequency

	result := repo.db.Find(&reviewFrequencies)

	return reviewFrequencies, result.Error
}
