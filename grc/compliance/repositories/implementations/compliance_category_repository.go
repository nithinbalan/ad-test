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

	result := gormDB.Find(&complianceCategories)

	return complianceCategories, result.Error
}
