package implementations

import (
	"jethur/grc/compliance/models"
	"jethur/utility"

	"github.com/gin-gonic/gin"
)

type ComplianceCategoryRepository struct{}

func NewComplianceCategoryRepository() *ComplianceCategoryRepository {
	return &ComplianceCategoryRepository{}
}

func (repo *ComplianceCategoryRepository) GetAll(c *gin.Context) ([]models.ComplianceCategory, error) {

	gormDB, err := utility.GetDBFromContext(c)
	if err != nil {
		return nil, err
	}

	var complianceCategories []models.ComplianceCategory

	// Get parameters from the request query string
    sortBy := c.Query("sortBy")
    sortOrder := c.Query("sortOrder")
    search := c.Query("search")

    // Build the query based on sorting and search parameters
    query := gormDB.Model(&models.ComplianceCategory{})
    if sortBy != "" && sortOrder != "" {
        query = query.Order(sortBy + " " + sortOrder)
    }
    if search != "" {
        query = query.Where("title_en LIKE ?", "%"+search+"%")
    }

	result := query.Find(&complianceCategories)

	return complianceCategories, result.Error
}
