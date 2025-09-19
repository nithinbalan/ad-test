package controllers

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ERMDashboardController struct {
	service services.ERMDashboardServiceInterface
	db      *gorm.DB
}

func NewERMDashboardController(service services.ERMDashboardServiceInterface) *ERMDashboardController {
	return &ERMDashboardController{
		service: service,
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
