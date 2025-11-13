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

func (cddc *ComplianceDetailDashboardController) GetComplianceClauseSummary(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetComplianceClauseSummary(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetTop5ClausePerformance(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetTop5ClausePerformance(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetBottom5ClausePerformance(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetBottom5ClausePerformance(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetComplianceDomainWiseOverview(c *gin.Context) {

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetComplianceDomainWiseOverview(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetPolicyByOrganizationUnit(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetPolicyByOrganizationUnit(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetControlEffectivenessSummary(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetControlEffectivenessSummary(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetControlImplementationSummary(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetControlImplementationSummary(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetControlSummary(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetControlSummary(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetControlEfficiencySummary(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetControlEfficiencySummary(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetControlCoverageByDomain(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetControlCoverageByDomain(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetCompliancePolicySummary(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetCompliancePolicySummary(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetPolicyAcknowledgementByOrganizationUnit(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetPolicyAcknowledgementByOrganizationUnit(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetActionPlanOverview(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetActionPlanOverview(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetActionPlanStatusSummary(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetActionPlanStatusSummary(uint(complianceID))
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

func (cddc *ComplianceDetailDashboardController) GetActionPlansByOrganizationUnit(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := cddc.service.GetActionPlansByOrganizationUnit(uint(complianceID))
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
