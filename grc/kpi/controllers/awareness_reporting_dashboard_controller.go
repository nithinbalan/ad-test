package controllers

import (
	awarenessDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessDashboardController struct {
	service services.AwarenessDashboardServiceInterface
	db      *gorm.DB
}

func NewAwarenessDashboardController(service services.AwarenessDashboardServiceInterface) *AwarenessDashboardController {
	return &AwarenessDashboardController{service: service}
}

func (adc *AwarenessDashboardController) buildParams(c *gin.Context) *awarenessDTOs.AwarenessDashboardQueryParamsDTO {
	year, _ := strconv.Atoi(c.Query("year"))
	return &awarenessDTOs.AwarenessDashboardQueryParamsDTO{
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Year:     year,
	}
}

func (adc *AwarenessDashboardController) GetAwarenessPlanSummary(c *gin.Context) {
	result, err := adc.service.GetAwarenessPlanSummary(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AwarenessDashboardController) GetCampaignsByType(c *gin.Context) {
	result, err := adc.service.GetCampaignsByType(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AwarenessDashboardController) GetPriorityDistribution(c *gin.Context) {
	result, err := adc.service.GetPriorityDistribution(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AwarenessDashboardController) GetCampaignsPerMonth(c *gin.Context) {
	result, err := adc.service.GetCampaignsPerMonth(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AwarenessDashboardController) GetChannelMixPerMonth(c *gin.Context) {
	result, err := adc.service.GetChannelMixPerMonth(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AwarenessDashboardController) GetCoverageByAudience(c *gin.Context) {
	result, err := adc.service.GetCoverageByAudience(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AwarenessDashboardController) GetAwarenessProgressSummary(c *gin.Context) {
	result, err := adc.service.GetAwarenessProgressSummary(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AwarenessDashboardController) GetCompletionByChannel(c *gin.Context) {
	result, err := adc.service.GetCompletionByChannel(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AwarenessDashboardController) GetCompletionPerCampaign(c *gin.Context) {
	result, err := adc.service.GetCompletionPerCampaign(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AwarenessDashboardController) GetSurveyCompletionStatus(c *gin.Context) {
	result, err := adc.service.GetSurveyCompletionStatus(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AwarenessDashboardController) GetCompletionByDepartment(c *gin.Context) {
	result, err := adc.service.GetCompletionByDepartment(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AwarenessDashboardController) GetPhishingResults(c *gin.Context) {
	result, err := adc.service.GetPhishingResults(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}