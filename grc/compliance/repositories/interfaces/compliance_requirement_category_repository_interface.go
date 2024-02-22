package interfaces

import (
	"jethur/grc/compliance/models"

	"github.com/gin-gonic/gin"
)

type ComplianceRequirementCategoryRepositoryInterface interface {
	GetAll(c *gin.Context) ([]models.ComplianceRequirementCategory, error)
	// GetByID(c *gin.Context, id uint) (models.ComplianceRequirementCategory, error)
	// Create(c *gin.Context, article models.ComplianceRequirementCategory) (models.ComplianceRequirementCategory, error)
	// Update(c *gin.Context, id uint, article models.ComplianceRequirementCategory) (models.ComplianceRequirementCategory, error)
	// Delete(c *gin.Context, id uint) error
}
