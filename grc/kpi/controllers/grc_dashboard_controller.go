package controllers

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GRCDashboardController struct {
	service services.GRCDashboardServiceInterface
	db      *gorm.DB
}

func NewGRCDashboardController(service services.GRCDashboardServiceInterface) *GRCDashboardController {
	return &GRCDashboardController{
		service: service,
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
