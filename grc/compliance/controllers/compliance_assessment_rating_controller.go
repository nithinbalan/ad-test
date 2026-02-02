package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentRatingController struct {
	service services.ComplianceAssessmentRatingServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentRatingController(service services.ComplianceAssessmentRatingServiceInterface) *ComplianceAssessmentRatingController {
	return &ComplianceAssessmentRatingController{
		service: service,
	}
}

func (carc *ComplianceAssessmentRatingController) GetAllComplianceAssessmentRatings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	ratings, err := carc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(carc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, ratings)
}

func (carc *ComplianceAssessmentRatingController) GetComplianceAssessmentRatingDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	rating, err := carc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, rating)
}
