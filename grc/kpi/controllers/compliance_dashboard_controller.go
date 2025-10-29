package controllers

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceDashboardController struct {
	service services.ComplianceDashboardServiceInterface
	db      *gorm.DB
}

func NewComplianceDashboardController(service services.ComplianceDashboardServiceInterface) *ComplianceDashboardController {
	return &ComplianceDashboardController{
		service: service,
	}
}

func (cdc *ComplianceDashboardController) GetComplianceOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	overview, err := cdc.service.GetComplianceOverview(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, overview)
}

func (cdc *ComplianceDashboardController) GetLeadershipAttentionSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := cdc.service.GetLeadershipAttentionSummary(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (cdc *ComplianceDashboardController) GetComplianceBreakdown(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	breakdown, err := cdc.service.GetComplianceBreakdown(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, breakdown)
}

func (cdc *ComplianceDashboardController) GetHighPerformanceCompliances(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := cdc.service.GetHighPerformanceCompliances(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (cdc *ComplianceDashboardController) GetLeastPerformanceCompliances(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := cdc.service.GetLeastPerformanceCompliances(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (cdc *ComplianceDashboardController) GetHighPerformanceComplianceClauses(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := cdc.service.GetHighPerformanceComplianceClauses(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (cdc *ComplianceDashboardController) GetPolicyOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	overview, err := cdc.service.GetPolicyOverview(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, overview)
}

func (cdc *ComplianceDashboardController) GetLeastPerformanceComplianceClauses(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := cdc.service.GetLeastPerformanceComplianceClauses(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (cdc *ComplianceDashboardController) GetPolicyByOrganizationUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := cdc.service.GetPolicyByOrganizationUnit(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (cdc *ComplianceDashboardController) GetPolicyStatusBreakdown(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := cdc.service.GetPolicyStatusBreakdown(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (cdc *ComplianceDashboardController) GetControlSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := cdc.service.GetControlSummary(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (cdc *ComplianceDashboardController) GetComplianceEfficiencySummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := cdc.service.GetComplianceEfficiencySummary(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (cdc *ComplianceDashboardController) GetControlOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	overview, err := cdc.service.GetControlOverview(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, overview)
}

func (cdc *ComplianceDashboardController) GetOrganizationUnitWiseAssessmentCoverage(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	overview, err := cdc.service.GetOrganizationUnitWiseAssessmentCoverage(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, overview)
}

func (cdc *ComplianceDashboardController) GetActionPlanOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	overview, err := cdc.service.GetActionPlanOverview(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, overview)
}

func (cdc *ComplianceDashboardController) GetActionPlanStatusSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := cdc.service.GetActionPlanStatusSummary(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (cdc *ComplianceDashboardController) GetActionPlansByOrganizationUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := cdc.service.GetActionPlansByOrganizationUnit(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (edc *ComplianceDashboardController) GetTopCompliances(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ComplianceDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetTopCompliances(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(edc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}
