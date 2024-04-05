package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ActionPlanStatusController struct {
	service services.ActionPlanStatusServiceInterface
	db      *gorm.DB
}

func NewActionPlanStatusController(service services.ActionPlanStatusServiceInterface) *ActionPlanStatusController {
	return &ActionPlanStatusController{
		service: service,
	}
}

func (apsc *ActionPlanStatusController) GetAllActionPlanStatuses(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	actionPlanStatuses, err := apsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(apsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, actionPlanStatuses)
}

func (apsc *ActionPlanStatusController) GetActionPlanStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInternalErrResp(apsc.db, c, "Invalid id")
		return
	}

	actionPlanStatus, err := apsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, actionPlanStatus)
}
