package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlImplementationStatusController struct {
	service services.ControlImplementationStatusServiceInterface
	db      *gorm.DB
}

func NewControlIMplementationStatusController(service services.ControlImplementationStatusServiceInterface) *ControlImplementationStatusController {
	return &ControlImplementationStatusController{
		service: service,
	}
}

func (cisc *ControlImplementationStatusController) GetAllControlImplementationStatuses(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	controlImplementationStatuses, err := cisc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(cisc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, controlImplementationStatuses)
}

func (cisc *ControlImplementationStatusController) GetControlImplementationStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlImplementationStatus, err := cisc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, controlImplementationStatus)
}
