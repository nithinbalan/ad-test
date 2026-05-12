package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanWorkflowController struct {
	service services.RiskRegisterUpdatePlanWorkflowServiceInterface
	db      *gorm.DB
}

func NewRiskRegisterUpdatePlanWorkflowController(service services.RiskRegisterUpdatePlanWorkflowServiceInterface) *RiskRegisterUpdatePlanWorkflowController {
	return &RiskRegisterUpdatePlanWorkflowController{
		service: service,
	}
}

func (rrupwc *RiskRegisterUpdatePlanWorkflowController) GetWorkflowList(c *gin.Context) {
	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := rrupwc.service.GetRiskRegisterUpdatePlanWorkflowList(uint(riskRegisterUpdatePlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
