package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskResponseController struct {
	service services.RiskResponseServiceInterface
	db      *gorm.DB
}

func NewRiskResponseController(service services.RiskResponseServiceInterface) *RiskResponseController {
	return &RiskResponseController{
		service: service,
	}
}

func (rrc *RiskResponseController) GetRiskResponseDetails(c *gin.Context) {
	riskIDStr := c.Param("riskId")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskResponse, err := rrc.service.GetFullDetailsByRiskID(uint(riskID))
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

	c.JSON(http.StatusOK, riskResponse)
}
