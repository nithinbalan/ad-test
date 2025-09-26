package controllers

import (
	"fmt"
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskServices "jethur/grc/risk/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ERMDashboardController struct {
	service        services.ERMDashboardServiceInterface
	riskHeatMapSvc riskServices.RiskHeatmapServiceInterface
	db             *gorm.DB
}

func NewERMDashboardController(service services.ERMDashboardServiceInterface,
	riskHeatMapSvc riskServices.RiskHeatmapServiceInterface) *ERMDashboardController {
	return &ERMDashboardController{
		service:        service,
		riskHeatMapSvc: riskHeatMapSvc,
	}
}

// ERM Risk Overview
func (edc *ERMDashboardController) GetERMRiskOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	overview, err := edc.service.GetERMRiskOverview(&queryParams)
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

	c.JSON(http.StatusOK, overview)
}

// ERM Risk Distribution
func (edc *ERMDashboardController) GetERMRiskDistribution(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	distribution, err := edc.service.GetERMRiskDistribution(&queryParams)
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

	c.JSON(http.StatusOK, distribution)
}

// ERM Overdue Risks
func (edc *ERMDashboardController) GetERMOverdueRisks(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	overdueRisks, err := edc.service.GetERMOverdueRisks(&queryParams)
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

	c.JSON(http.StatusOK, overdueRisks)
}

// ERM Risks Above Appetite
func (edc *ERMDashboardController) GetERMRisksAboveAppetite(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	risksAboveAppetite, err := edc.service.GetERMRisksAboveAppetite(&queryParams)
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

	c.JSON(http.StatusOK, risksAboveAppetite)
}

// ERM Leadership Attention Metrics
func (edc *ERMDashboardController) GetERMLeadershipAttentionMetrics(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	leadershipMetrics, err := edc.service.GetERMLeadershipAttentionMetrics(&queryParams)
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

	c.JSON(http.StatusOK, leadershipMetrics)
}

// ERM Control Efficiency Summary
func (edc *ERMDashboardController) GetERMControlEfficiencySummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := edc.service.GetERMControlEfficiencySummary(&queryParams)
	fmt.Println(summary)
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

// ERM Control Efficiency Breakdown
func (edc *ERMDashboardController) GetERMControlEfficiencyBreakdown(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := edc.service.GetERMControlEfficiencyBreakdown(&queryParams)
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

// ERM Controls Overview
func (edc *ERMDashboardController) GetERMControlsOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMControlsOverview(&queryParams)
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

// ERM Risk Appetite Overview
func (edc *ERMDashboardController) GetERMRiskAppetiteOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMRiskAppetiteOverview(&queryParams)
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

// ERM Risk Appetite Summary
func (edc *ERMDashboardController) GetERMRiskAppetiteSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMRiskAppetiteSummary(&queryParams)
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

// ERM Risk Aging Status
func (edc *ERMDashboardController) GetERMRiskAgingStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMRiskAgingStatus(&queryParams)
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

// ERM Risk Treatment Plan Progress
func (edc *ERMDashboardController) GetERMRiskTreatmentPlanProgress(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMRiskTreatmentPlanProgress(&queryParams)
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

// ERM Strategic Objectives Risk Overview
func (edc *ERMDashboardController) GetERMStrategicObjectivesRiskOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMStrategicObjectivesRiskOverview(&queryParams)
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

// ERM Strategic Objectives Risk Summary
func (edc *ERMDashboardController) GetERMStrategicObjectivesRiskSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMStrategicObjectivesRiskSummary(&queryParams)
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

// ERM Risk Summary By Department
func (edc *ERMDashboardController) GetERMRiskSummaryByDepartment(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	summary, err := edc.service.GetERMRiskSummaryByDepartment(&queryParams)
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

// ERM Risk Heat Map
func (edc *ERMDashboardController) GetERMRiskHeatMap(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := riskDTOs.RiskHeatmapQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	riskHeatmap, err := edc.riskHeatMapSvc.GetRiskHeatmapByRiskMatrix(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(edc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, riskHeatmap)
}

func (edc *ERMDashboardController) GetERMRiskCategorySummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMRiskCategorySummary(&queryParams)
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

func (edc *ERMDashboardController) GetERMRiskSourceSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMRiskSourceSummary(&queryParams)
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

func (edc *ERMDashboardController) GetERMRiskStatusSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMRiskStatusSummary(&queryParams)
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

func (edc *ERMDashboardController) GetERMRiskDepartmentSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMRiskDepartmentSummary(&queryParams)
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

func (edc *ERMDashboardController) GetERMRisksIdentifiedSummary(c *gin.Context) {

	summary, err := edc.service.GetERMRisksIdentifiedSummary()
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

func (edc *ERMDashboardController) GetERMRiskResponseStrategySummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.ERMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	summary, err := edc.service.GetERMRiskResponseStrategySummary(&queryParams)
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
