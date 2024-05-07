package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskStatusController struct {
	service services.RiskStatusServiceInterface
	db      *gorm.DB
}

func NewRiskStatusController(service services.RiskStatusServiceInterface) *RiskStatusController {
	return &RiskStatusController{
		service: service,
	}
}

func (rsc *RiskStatusController) GetAllRiskStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	riskStatuses, err := rsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskStatuses)
}

func (rsc *RiskStatusController) GetRiskStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskStatus, err := rsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskStatus)
}
