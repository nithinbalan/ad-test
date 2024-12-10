package controllers

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	riskServices "jethur/grc/risk/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GRCDashboardController struct {
	service        services.GRCDashboardServiceInterface
	riskHeatMapSvc riskServices.RiskHeatmapServiceInterface
	db             *gorm.DB
}

func NewGRCDashboardController(service services.GRCDashboardServiceInterface,
	riskHeatMapSvc riskServices.RiskHeatmapServiceInterface) *GRCDashboardController {
	return &GRCDashboardController{
		service:        service,
		riskHeatMapSvc: riskHeatMapSvc,
	}
}

// GRC Compliance Summary
func (gdc *GRCDashboardController) GetGRCComplianceSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCComplianceSummary(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

// GRC Risk Summary
func (gdc *GRCDashboardController) GetGRCRiskSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCRiskSummary(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCAuditSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCAuditSummary(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCDocumentSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCDocumentSummary(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCOrganizationSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCOrganizationSummary(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCIncidentSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCIncidentSummary(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCComplianceOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCComplianceOverview(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCRiskOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCRiskOverview(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCDocumentOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCDocumentOverview(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCTopCompliances(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCTopCompliances(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCTopRisks(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCTopRisks(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCTopPublishedDocuments(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCTopPublishedDocuments(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCDocumentTypeChart(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := gdc.service.GetGRCDocumentTypeChart(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}

func (gdc *GRCDashboardController) GetGRCRiskHeatMap(c *gin.Context) {
	riskHeatmap, err := gdc.riskHeatMapSvc.GetRiskHeatmapByRiskMatrix()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(gdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, riskHeatmap)
}
