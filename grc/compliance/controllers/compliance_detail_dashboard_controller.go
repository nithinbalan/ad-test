package controllers

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
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
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "domain_title_en")
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cddc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cddc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"domain_title_en": "domain_title_en",
		"domain_title_ar": "domain_title_ar",
		"compliant":       "compliant",
		"partial":         "partial",
		"non_compliant":   "non_compliant",
		"pending":         "pending",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := complianceDTOs.ComplianceDashboardQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Page:     page,
		PageSize: pageSize,
	}

	summary, err := cddc.service.GetComplianceDomainWiseOverview(&queryParams, uint(complianceID))
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
