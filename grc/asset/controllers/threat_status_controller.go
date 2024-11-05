package controllers

import (
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThreatStatusController struct {
	service services.ThreatStatusServiceInterface
	db      *gorm.DB
}

func NewThreatStatusController(service services.ThreatStatusServiceInterface) *ThreatStatusController {
	return &ThreatStatusController{
		service: service,
	}
}

func (tsc *ThreatStatusController) GetAllThreatStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	threatStatuses, err := tsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, threatStatuses)
}

func (tsc *ThreatStatusController) GetThreatStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	threatStatus, err := tsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, threatStatus)
}
