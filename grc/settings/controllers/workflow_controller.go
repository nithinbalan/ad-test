package controllers

import (
	"fmt"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowController struct {
	service services.WorkflowServiceInterface
	db      *gorm.DB
}

func NewWorkflowController(service services.WorkflowServiceInterface) *WorkflowController {
	return &WorkflowController{
		service: service,
	}
}

func (wfc *WorkflowController) GetAllWorkflows(c *gin.Context) {
	workflow_type := c.Query("type")
	module_identifier := c.Query("module")
	workflows, err := wfc.service.GetAllWorkflowsByType(workflow_type, module_identifier)
	if err != nil {
		utility.SendInternalErrResp(wfc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, workflows)
}

func (wfc *WorkflowController) GetWorkflowlevels(c *gin.Context) {
	idStr := c.Param("workflowId")
	workflow_id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := wfc.service.GetWorkflowLevels(uint(workflow_id))
	if err != nil {
		utility.SendInternalErrResp(wfc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, workflows)
}

func (wfc *WorkflowController) CreateWorkflowLevel(c *gin.Context) {
	var workflowLevelFormDTO settingsDTOs.WorkflowLevelFormDTO
	if err := c.ShouldBindJSON(&workflowLevelFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	workflowIDStr := c.Param("workflowId")
	workflowID, err := strconv.ParseUint(workflowIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow level ID"})
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(wfc.db, c, err.Error())
		return
	}

	createdId, err := wfc.service.CreateWorkFlowLevel(uint(workflowID), currentUserID, &workflowLevelFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(wfc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, createdId)
}

func (wfc *WorkflowController) DeleteWorkflowLevel(c *gin.Context) {
	fmt.Println("hii")
	workflowLevelIDStr := c.Param("levelId")
	workflowLevelID, err := strconv.ParseUint(workflowLevelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow level ID"})
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(wfc.db, c, err.Error())
		return
	}

	if err := wfc.service.DeleteWorkflowLevel(currentUserID, uint(workflowLevelID)); err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(wfc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow level deleted successfully"})
}
