package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlStatusController struct {
	service services.ControlStatusServiceInterface
	db      *gorm.DB
}

func NewControlStatusController(service services.ControlStatusServiceInterface) *ControlStatusController {
	return &ControlStatusController{
		service: service,
	}
}

func (csc *ControlStatusController) GetAllControlStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	typeFilter := c.Query("type")

	controlStatuses, err := csc.service.GetAll(sort, sortBy, keyword, status, typeFilter)

	if err != nil {
		utility.SendInternalErrResp(csc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, controlStatuses)
}

func (csc *ControlStatusController) GetControlStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlStatus, err := csc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, controlStatus)
}
