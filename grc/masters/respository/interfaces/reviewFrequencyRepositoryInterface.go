package interfaces

import (
	"jethur/grc/masters/models"

	"github.com/gin-gonic/gin"
)

type ReviewFrequencyRepositoryInterface interface {
	GetAll(c *gin.Context) ([]models.ReviewFrequency, error)
	// GetByID(c *gin.Context, id uint) (models.ReviewFrequency, error)
	// Create(c *gin.Context, article models.ReviewFrequency) (models.ReviewFrequency, error)
	// Update(c *gin.Context, id uint, article models.ReviewFrequency) (models.ReviewFrequency, error)
	// Delete(c *gin.Context, id uint) error
}
