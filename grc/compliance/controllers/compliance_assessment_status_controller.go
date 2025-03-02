package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentStatusController struct {
	service services.ComplianceAssessmentStatusServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentStatusController(service services.ComplianceAssessmentStatusServiceInterface) *ComplianceAssessmentStatusController {
	return &ComplianceAssessmentStatusController{
		service: service,
	}
}

func (casc *ComplianceAssessmentStatusController) GetAllComplianceAssessmentStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	assessmentStatuses, err := casc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(casc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, assessmentStatuses)
}

func (casc *ComplianceAssessmentStatusController) GetComplianceAssessmentStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assessmentStatus, err := casc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, assessmentStatus)
}
