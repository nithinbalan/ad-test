package controllers

import (
	"jethur/grc/compliance/repositories/interfaces"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ComplianceRequirementCategoryController struct {
	repo interfaces.ComplianceRequirementCategoryRepositoryInterface
}

func NewComplianceRequirementCategoryController(repo interfaces.ComplianceRequirementCategoryRepositoryInterface) *ComplianceRequirementCategoryController {
	return &ComplianceRequirementCategoryController{
		repo: repo,
	}
}

func (rfc *ComplianceRequirementCategoryController) GetAllComplianceRequirementCategories(c *gin.Context) {
	complianceRequirementCategories, err := rfc.repo.GetAll(c)
	if err != nil {
		utility.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, complianceRequirementCategories)
}
