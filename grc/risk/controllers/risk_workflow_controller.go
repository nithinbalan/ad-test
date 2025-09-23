package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskWorkflowController struct {
	service services.RiskWorkflowServiceInterface
	db      *gorm.DB
}

func NewRiskWorkflowController(service services.RiskWorkflowServiceInterface) *RiskWorkflowController {
	return &RiskWorkflowController{
		service: service,
	}
}

func (arwc *RiskWorkflowController) GetWorkflowList(c *gin.Context) {
	riskIDStr := c.Param("id")
	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := arwc.service.GetRiskWorkflowList(uint(riskID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
