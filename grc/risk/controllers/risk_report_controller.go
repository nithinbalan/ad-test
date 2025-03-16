package controllers

import (
	"net/http"
	"strconv"

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

func (rrc *RiskReportController) GetAllReports(c *gin.Context) {

	riskReportID, err := strconv.Atoi(c.Param("riskReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reports, err := rrc.service.GetAll(uint(riskReportID))
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

	c.JSON(http.StatusOK, reports)
}

func (rrc *RiskReportController) GetRiskReportDetails(c *gin.Context) {

	riskReportID, err := strconv.Atoi(c.Param("riskReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := rrc.service.GetRiskReportDetails(uint(riskReportID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, details)
}

func (rrc *RiskReportController) CreateRiskReport(c *gin.Context) {
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

	riskReportID, err := rrc.service.Create(currentUserID, &request)
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

func (rrc *RiskReportController) UpdateReport(c *gin.Context) {
	riskReportID, err := strconv.Atoi(c.Param("riskReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.RiskReportSectionFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.Update(uint(riskReportID), uint(sectionID), currentUserID, &formDTO)
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

	c.Status(http.StatusOK)
}
