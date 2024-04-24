package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlEfficiencyStatusController struct {
	service services.ControlEfficiencyStatusServiceInterface
	db      *gorm.DB
}

func NewControlEfficiencyStatusController(service services.ControlEfficiencyStatusServiceInterface) *ControlEfficiencyStatusController {
	return &ControlEfficiencyStatusController{
		service: service,
	}
}

func (cesc *ControlEfficiencyStatusController) GetAllControlEfficiencyStatuses(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	controlEfficiencyStatuses, err := cesc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(cesc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, controlEfficiencyStatuses)
}

func (cesc *ControlEfficiencyStatusController) GetControlEfficiencyStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlEfficiencyStatus, err := cesc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, controlEfficiencyStatus)
}
