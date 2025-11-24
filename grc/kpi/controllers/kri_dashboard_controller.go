package controllers

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KRIDashboardController struct {
	service services.KRIDashboardServiceInterface
	db      *gorm.DB
}

func NewKRIDashboardController(service services.KRIDashboardServiceInterface) *KRIDashboardController {
	return &KRIDashboardController{
		service: service,
	}
}

func (kdc *KRIDashboardController) GetKRIOverview(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.KriDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	response, err := kdc.service.GetKRIOverview(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

func (kdc *KRIDashboardController) GetKRIReviewStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.KriDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	response, err := kdc.service.GetKRIReviewStatus(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

func (kdc *KRIDashboardController) GetKRIStatusCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.KriDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	response, err := kdc.service.GetKRIStatusCategory(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

func (kdc *KRIDashboardController) GetKRISeverityBreakdown(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.KriDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	response, err := kdc.service.GetKRISeverityBreakdown(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, response)
}

func (kdc *KRIDashboardController) GetTopKRIRiskExposure(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.KriDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	response, err := kdc.service.GetTopKRIRiskExposure(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, response)
}

func (kdc *KRIDashboardController) GetKRISeverityByOrganizationUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.KriDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	response, err := kdc.service.GetKRISeverityByOrganizationUnit(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

func (kdc *KRIDashboardController) GetKRISeverityByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := kpiDTOs.KriDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	response, err := kdc.service.GetKRISeverityByCategory(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

func (kdc *KRIDashboardController) GetKRIMostFrequentBreaches(c *gin.Context) {
	// dateFrom := c.Query("date_from")
	// dateTo := c.Query("date_to")

	// queryParams := kpiDTOs.KriDashboardQueryParamsDTO{
	// 	DateFrom: dateFrom,
	// 	DateTo:   dateTo,
	// }

	response, err := kdc.service.GetKRIMostFrequentBreaches()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, response)
}
