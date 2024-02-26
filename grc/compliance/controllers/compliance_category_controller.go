package controllers

import (
	"jethur/grc/compliance/repositories/interfaces"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceCategoryController struct {
	repo interfaces.ComplianceCategoryRepositoryInterface
	db   *gorm.DB
}

func NewComplianceCategoryController(repo interfaces.ComplianceCategoryRepositoryInterface, db *gorm.DB) *ComplianceCategoryController {
	return &ComplianceCategoryController{
		repo: repo,
		db:   db,
	}
}

func (rfc *ComplianceCategoryController) GetAllComplianceCategories(c *gin.Context) {
	complianceCategories, err := rfc.repo.GetAll(c)
	if err != nil {
		utility.SendErrorResponse(rfc.db, c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, complianceCategories)
}
