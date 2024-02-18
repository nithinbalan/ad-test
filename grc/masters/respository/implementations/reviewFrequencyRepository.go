package implementations

import (
	"jethur/grc/masters/models"
	"jethur/utility"

	"github.com/gin-gonic/gin"
)

type ReviewFrequencyRepository struct{}

func NewReviewFrequencyRepository() *ReviewFrequencyRepository {
	return &ReviewFrequencyRepository{}
}

func (repo *ReviewFrequencyRepository) GetAll(c *gin.Context) ([]models.ReviewFrequency, error) {

	gormDB, err := utility.GetDBFromContext(c)
	if err != nil {
		return nil, err
	}

	var reviewFrequencies []models.ReviewFrequency

	result := gormDB.Find(&reviewFrequencies)

	return reviewFrequencies, result.Error
}
