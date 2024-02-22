package implementations

import (
	"jethur/grc/compliance/models"
	"jethur/utility"

	"github.com/gin-gonic/gin"
)

type ComplianceRequirementCategoryRepository struct{}

func NewComplianceRequirementCategoryRepository() *ComplianceRequirementCategoryRepository {
	return &ComplianceRequirementCategoryRepository{}
}

func (repo *ComplianceRequirementCategoryRepository) GetAll(c *gin.Context) ([]models.ComplianceRequirementCategory, error) {

	gormDB, err := utility.GetDBFromContext(c)
	if err != nil {
		return nil, err
	}

	var complianceRequirementCategories []models.ComplianceRequirementCategory

	// Get parameters from the request query string
    sortBy := c.Query("sortBy")
    sortOrder := c.Query("sortOrder") 
    if sortOrder == "" {
    	sortOrder = "asc"
	}
    search := c.Query("search")
    status := c.Query("status")

    // Build the query based on sorting, search and filter parameters
    query := gormDB.Model(&models.ComplianceRequirementCategory{})
    if sortBy != "" {
        query = query.Order(sortBy + " " + sortOrder)
    }
    if search != "" {
        query = query.Where("title_en LIKE ?", "%"+search+"%").
                Or("title_ar LIKE ?", "%"+search+"%")
    }
    if status != "" {
        query = query.Where("status = ?", status)
    }

	result := query.Find(&complianceRequirementCategories)

	return complianceRequirementCategories, result.Error
}
