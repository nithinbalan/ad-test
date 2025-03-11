package controllers

import (
	"net/http"

	"jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskReportController struct {
	service services.RiskReportServiceInterface
	db      *gorm.DB
}

func NewRiskReportController(service services.RiskReportServiceInterface) *RiskReportController {
	return &RiskReportController{
		service: service,
	}
}

// GenerateRiskReport generates a new risk report for risks in a given date range
func (rrc *RiskReportController) GenerateRiskReport(c *gin.Context) {
	var request dtos.RiskReportFormDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	riskReportID, err := rrc.service.GenerateRiskReport(currentUserID, &request)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Risk report generated successfully",
		"risk_report_id": riskReportID,
	})
}
