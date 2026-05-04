package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskResponseWorkflowController struct {
	service services.RiskResponseWorkflowServiceInterface
	db      *gorm.DB
}

func NewRiskResponseWorkflowController(service services.RiskResponseWorkflowServiceInterface) *RiskResponseWorkflowController {
	return &RiskResponseWorkflowController{
		service: service,
	}
}

func (rrwc *RiskResponseWorkflowController) GetWorkflowList(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := rrwc.service.GetRiskResponseWorkflowList(uint(riskResponseID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
