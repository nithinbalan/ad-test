package implementations

import (
	"jethur/grc/compliance/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceCategoryRepository struct {
	db *gorm.DB
}

func NewComplianceCategoryRepository(db *gorm.DB) *ComplianceCategoryRepository {
	return &ComplianceCategoryRepository{
		db: db,
	}
}

func (repo *ComplianceCategoryRepository) GetAll(c *gin.Context) ([]models.ComplianceCategory, error) {

	var complianceCategories []models.ComplianceCategory

	// Get parameters from the request query string
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")
	if sortOrder == "" {
		sortOrder = "asc"
	}
	search := c.Query("search")
	status := c.Query("status")

	// Build the query based on sorting, search and filter parameters
	query := repo.db.Model(&models.ComplianceCategory{})
	if sortBy != "" {
		query = query.Order(sortBy + " " + sortOrder)
	}
	if search != "" {
		query = query.Where("title_en LIKE ?", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	result := query.Find(&complianceCategories)

	return complianceCategories, result.Error
}
