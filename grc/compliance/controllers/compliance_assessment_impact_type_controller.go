package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentImpactTypeController struct {
	service services.ComplianceAssessmentImpactTypeServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentImpactTypeController(service services.ComplianceAssessmentImpactTypeServiceInterface) *ComplianceAssessmentImpactTypeController {
	return &ComplianceAssessmentImpactTypeController{
		service: service,
	}
}

func (caitc *ComplianceAssessmentImpactTypeController) GetAllComplianceAssessmentImpactTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	impactTypes, err := caitc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(caitc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, impactTypes)
}

func (caitc *ComplianceAssessmentImpactTypeController) GetComplianceAssessmentImpactTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	impactType, err := caitc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, impactType)
}
