package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceDetailDashboardController struct {
	service services.ComplianceDetailDashboardServiceInterface
	db      *gorm.DB
}

func NewComplianceDetailDashboardController(service services.ComplianceDetailDashboardServiceInterface) *ComplianceDetailDashboardController {
	return &ComplianceDetailDashboardController{
		service: service,
	}
}

func (cddc *ComplianceDetailDashboardController) GetComplianceDetailSummary(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetComplianceDetailSummary(uint(complianceID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cddc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}
