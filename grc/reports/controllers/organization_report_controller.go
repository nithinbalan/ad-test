package controllers

import (
	"jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationReportController struct {
	service services.OrganizationReportServiceInterface
	db      *gorm.DB
}

func NewOrganizationReportController(service services.OrganizationReportServiceInterface) *OrganizationReportController {
	return &OrganizationReportController{
		service: service,
	}
}

func (orc *OrganizationReportController) GetUsersByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	usersByStatusReport, err := orc.service.GetUsersByStatus(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, usersByStatusReport)
}
