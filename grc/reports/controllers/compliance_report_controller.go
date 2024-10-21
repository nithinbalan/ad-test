package controllers

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceReportController struct {
	service services.ComplianceReportServiceInterface
	db      *gorm.DB
}

func NewComplianceReportController(service services.ComplianceReportServiceInterface) *ComplianceReportController {
	return &ComplianceReportController{
		service: service,
	}
}

func (crc *ComplianceReportController) GetComplianceByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	complianceReport, err := crc.service.GetComplianceByCategory(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceReport)
}
