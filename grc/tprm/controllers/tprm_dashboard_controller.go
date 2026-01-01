package controllers

import (
	tprmDTOs "jethur/grc/tprm/dtos"
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TPRMDashboardController struct {
	service services.TPRMDashboardServiceInterface
	db      *gorm.DB
}

func NewTPRMDashboardController(service services.TPRMDashboardServiceInterface) *TPRMDashboardController {
	return &TPRMDashboardController{
		service: service,
	}
}

func (tdc *TPRMDashboardController) GetThirdPartyInsights(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	allowedParams := []string{
		"organization_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := tprmDTOs.TPRMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Filters:  filters,
	}

	response, err := tdc.service.GetThirdPartyInsights(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, response)
}

func (tdc *TPRMDashboardController) GetAssessmentInsights(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	allowedParams := []string{
		"organization_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := tprmDTOs.TPRMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Filters:  filters,
	}

	response, err := tdc.service.GetAssessmentInsights(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, response)
}

func (tdc *TPRMDashboardController) GetTPRMSystemStatistics(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	allowedParams := []string{
		"organization_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := tprmDTOs.TPRMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Filters:  filters,
	}

	response, err := tdc.service.GetSystemStatistics(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, response)
}

func (tdc *TPRMDashboardController) GetOngoingAssessments(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	allowedParams := []string{
		"organization_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := tprmDTOs.TPRMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Filters:  filters,
	}

	response, err := tdc.service.GetOngoingAssessments(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, response)
}

func (tdc *TPRMDashboardController) GetThirdPartyServiceExpiry(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	allowedParams := []string{
		"organization_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := tprmDTOs.TPRMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Filters:  filters,
	}

	response, err := tdc.service.GetThirdPartyServiceExpiry(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, response)
}

func (tdc *TPRMDashboardController) GetAssessmentByTemplateType(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	allowedParams := []string{
		"organization_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := tprmDTOs.TPRMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Filters:  filters,
	}

	response, err := tdc.service.GetAssessmentByTemplateType(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, response)
}

func (tdc *TPRMDashboardController) GetThirdPartyByType(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	allowedParams := []string{
		"organization_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := tprmDTOs.TPRMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Filters:  filters,
	}

	response, err := tdc.service.GetThirdPartyByType(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, response)
}

func (tdc *TPRMDashboardController) GetAssessmentRiskMonthlyTrend(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	allowedParams := []string{
		"organization_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := tprmDTOs.TPRMDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Filters:  filters,
	}

	response, err := tdc.service.GetAssessmentRiskMonthlyTrend(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, response)
}
