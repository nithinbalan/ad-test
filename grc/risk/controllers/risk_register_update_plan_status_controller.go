package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanStatusController struct {
	service services.RiskRegisterUpdatePlanStatusServiceInterface
	db      *gorm.DB
}

func NewRiskRegisterUpdatePlanStatusController(service services.RiskRegisterUpdatePlanStatusServiceInterface) *RiskRegisterUpdatePlanStatusController {
	return &RiskRegisterUpdatePlanStatusController{
		service: service,
	}
}

func (rrupsc *RiskRegisterUpdatePlanStatusController) GetAllRiskRegisterUpdatePlanStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	riskRegisterUpdatePlanStatuses, err := rrupsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rrupsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskRegisterUpdatePlanStatuses)
}

func (rrupsc *RiskRegisterUpdatePlanStatusController) GetRiskRegisterUpdatePlanStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskRegisterUpdatePlanStatus, err := rrupsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskRegisterUpdatePlanStatus)
}
