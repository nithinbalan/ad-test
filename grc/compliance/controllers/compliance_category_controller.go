package controllers

import (
	"jethur/grc/compliance/repositories/interfaces"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ComplianceCategoryController struct {
	repo interfaces.ComplianceCategoryRepositoryInterface
}

func NewComplianceCategoryController(repo interfaces.ComplianceCategoryRepositoryInterface) *ComplianceCategoryController {
	return &ComplianceCategoryController{
		repo: repo,
	}
}

func (rfc *ComplianceCategoryController) GetAllComplianceCategories(c *gin.Context) {
	complianceCategories, err := rfc.repo.GetAll(c)
	if err != nil {
		utility.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, complianceCategories)
}
