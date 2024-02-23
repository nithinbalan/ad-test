package interfaces

import (
	"jethur/grc/compliance/models"

	"github.com/gin-gonic/gin"
)

type ComplianceCategoryRepositoryInterface interface {
	GetAll(c *gin.Context) ([]models.ComplianceCategory, error)
	// GetByID(c *gin.Context, id uint) (models.ComplianceCategory, error)
	// Create(c *gin.Context, article models.ComplianceCategory) (models.ComplianceCategory, error)
	// Update(c *gin.Context, id uint, article models.ComplianceCategory) (models.ComplianceCategory, error)
	// Delete(c *gin.Context, id uint) error
}
