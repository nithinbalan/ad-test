package controllers

import (
	"jethur/grc/settings/services" // Ensure the service path correctly points to where WorkflowLevelTypeService is defined.
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowLevelTypeController struct {
	service services.WorkflowLevelTypeServiceInterface
	db      *gorm.DB
}

func NewWorkflowLevelTypeController(service services.WorkflowLevelTypeServiceInterface) *WorkflowLevelTypeController {
	return &WorkflowLevelTypeController{
		service: service,
	}
}

func (wltc *WorkflowLevelTypeController) GetAllWorkflowLevelTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	workflowLevelTypes, err := wltc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(wltc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, workflowLevelTypes)
}

func (wltc *WorkflowLevelTypeController) GetWorkflowLevelTypeDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflowLevelType, err := wltc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, workflowLevelType)
}
