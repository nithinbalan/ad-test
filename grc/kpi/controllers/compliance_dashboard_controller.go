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
