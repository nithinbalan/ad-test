package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowStatusController struct {
	service services.WorkflowStatusServiceInterface
	db      *gorm.DB
}

func NewWorkflowStatusController(service services.WorkflowStatusServiceInterface) *WorkflowStatusController {
	return &WorkflowStatusController{
		service: service,
	}
}

func (wsc *WorkflowStatusController) GetAllWorkflowStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	workflowStatuses, err := wsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(wsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, workflowStatuses)
}

func (wsc *WorkflowStatusController) GetWorkflowStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInternalErrResp(wsc.db, c, "Invalid id")
		return
	}

	workflowStatus, err := wsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, workflowStatus)
}
