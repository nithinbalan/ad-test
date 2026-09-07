package controllers

import (
	"jethur/grc/threat-intelligence/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlannedAdvisoryStatusController struct {
	service services.PlannedAdvisoryStatusServiceInterface
	db      *gorm.DB
}

func NewPlannedAdvisoryStatusController(service services.PlannedAdvisoryStatusServiceInterface) *PlannedAdvisoryStatusController {
	return &PlannedAdvisoryStatusController{
		service: service,
	}
}

func (pasc *PlannedAdvisoryStatusController) GetAllPlannedAdvisoryStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	plannedAdvisoryStatuses, err := pasc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(pasc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, plannedAdvisoryStatuses)
}

func (pasc *PlannedAdvisoryStatusController) GetPlannedAdvisoryStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	plannedAdvisoryStatus, err := pasc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, plannedAdvisoryStatus)
}
