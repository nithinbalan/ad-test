package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementStatusController struct {
	service services.ComplianceAssessmentRequirementStatusServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentRequirementStatusController(service services.ComplianceAssessmentRequirementStatusServiceInterface) *ComplianceAssessmentRequirementStatusController {
	return &ComplianceAssessmentRequirementStatusController{
		service: service,
	}
}

func (carsc *ComplianceAssessmentRequirementStatusController) GetAllComplianceAssessmentRequirementStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	requirementStatuses, err := carsc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(carsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, requirementStatuses)
}

func (carsc *ComplianceAssessmentRequirementStatusController) GetComplianceAssessmentRequirementStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementStatus, err := carsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, requirementStatus)
}
