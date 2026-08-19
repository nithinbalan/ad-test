package controllers

import (
	settingDTOs "jethur/grc/settings/dtos"
	settingSvc "jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentWorkflowLevelController struct {
	workflowSvc settingSvc.WorkflowServiceInterface
	db          *gorm.DB
}

func NewIncidentWorkflowLevelController(
	workflowSvc settingSvc.WorkflowServiceInterface,
) *IncidentWorkflowLevelController {
	return &IncidentWorkflowLevelController{
		workflowSvc: workflowSvc,
	}
}

func (iwc *IncidentWorkflowLevelController) CreateWorkflowLevel(c *gin.Context) {
	var workflowLevelFormDTO settingDTOs.WorkflowLevelFormDTO
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
		utility.SendInternalErrResp(iwc.db, c, err.Error())
		return
	}

	createdId, err := iwc.workflowSvc.CreateWorkFlowLevel(uint(workflowID), currentUserID, &workflowLevelFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iwc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, createdId)
}

func (iwc *IncidentWorkflowLevelController) DeleteWorkflowLevel(c *gin.Context) {
	workflowLevelIDStr := c.Param("levelId")
	workflowLevelID, err := strconv.ParseUint(workflowLevelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow level ID"})
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iwc.db, c, err.Error())
		return
	}

	if err := iwc.workflowSvc.DeleteWorkflowLevel(currentUserID, uint(workflowLevelID)); err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow level deleted successfully"})
}
