package controllers

import (
	incidentDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentDashboardController struct {
	service services.IncidentDashboardServiceInterface
	db      *gorm.DB
}

func NewIncidentDashboardController(service services.IncidentDashboardServiceInterface, db *gorm.DB) *IncidentDashboardController {
	return &IncidentDashboardController{service: service, db: db}
}

func (idc *IncidentDashboardController) buildParams(c *gin.Context) *incidentDTOs.IncidentDashboardQueryParamsDTO {
	year, _ := strconv.Atoi(c.Query("year"))
	return &incidentDTOs.IncidentDashboardQueryParamsDTO{
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Year:     year,
	}
}

func (controller *IncidentDashboardController) GetIncidentStatusSummary(c *gin.Context) {
	params := controller.buildParams(c)
	response, err := controller.service.GetIncidentStatusSummary(params)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *IncidentDashboardController) GetIncidentSeverityBreakdown(c *gin.Context) {
	params := controller.buildParams(c)
	response, err := controller.service.GetIncidentSeverityBreakdown(params)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *IncidentDashboardController) GetIncidentAcknowledgementSummary(c *gin.Context) {
	params := controller.buildParams(c)
	response, err := controller.service.GetIncidentAcknowledgementSummary(params)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *IncidentDashboardController) GetIncidentSLASummary(c *gin.Context) {
	params := controller.buildParams(c)
	response, err := controller.service.GetIncidentSLASummary(params)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *IncidentDashboardController) GetIncidentCategorySummary(c *gin.Context) {
	params := controller.buildParams(c)
	response, err := controller.service.GetIncidentCategorySummary(params)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *IncidentDashboardController) GetIncidentTypeSummary(c *gin.Context) {
	params := controller.buildParams(c)
	response, err := controller.service.GetIncidentTypeSummary(params)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *IncidentDashboardController) GetIncidentTrend(c *gin.Context) {
	params := controller.buildParams(c)
	response, err := controller.service.GetIncidentTrend(params)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *IncidentDashboardController) GetIncidentActionPlanSummary(c *gin.Context) {
	params := controller.buildParams(c)
	response, err := controller.service.GetIncidentActionPlanSummary(params)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *IncidentDashboardController) GetOverdueActionPlans(c *gin.Context) {
	params := controller.buildParams(c)
	response, err := controller.service.GetOverdueActionPlans(params)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
