package controllers

import (
	"jethur/grc/settings/services" // Ensure the service path correctly points to where WorkflowRoleTypeService is defined.
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowRoleTypeController struct {
	service services.WorkflowRoleTypeServiceInterface
	db      *gorm.DB
}

func NewWorkflowRoleTypeController(service services.WorkflowRoleTypeServiceInterface) *WorkflowRoleTypeController {
	return &WorkflowRoleTypeController{
		service: service,
	}
}

func (wrtc *WorkflowRoleTypeController) GetAllWorkflowRoleTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	workflowRoleTypes, err := wrtc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(wrtc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, workflowRoleTypes)
}

func (wrtc *WorkflowRoleTypeController) GetWorkflowRoleTypeDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflowRoleType, err := wrtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, workflowRoleType)
}
