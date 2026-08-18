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

type DocumentWorkflowLevelController struct {
	workflowSvc settingSvc.WorkflowServiceInterface
	db          *gorm.DB
}

func NewDocumentWorkflowLevelController(
	workflowSvc settingSvc.WorkflowServiceInterface,
) *DocumentWorkflowLevelController {
	return &DocumentWorkflowLevelController{
		workflowSvc: workflowSvc,
	}
}

func (dwc *DocumentWorkflowLevelController) CreateWorkflowLevel(c *gin.Context) {
	var workflowLevelFormDTO settingDTOs.WorkflowLevelFormDTO
	if err := c.ShouldBindJSON(&workflowLevelFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	workflowIDStr := c.Param("id")
	workflowID, err := strconv.ParseUint(workflowIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow level ID"})
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dwc.db, c, err.Error())
		return
	}

	createdId, err := dwc.workflowSvc.CreateWorkFlowLevel(uint(workflowID), currentUserID, &workflowLevelFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dwc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, createdId)
}

func (dwc *DocumentWorkflowLevelController) DeleteWorkflowLevel(c *gin.Context) {
	workflowLevelIDStr := c.Param("levelId")
	workflowLevelID, err := strconv.ParseUint(workflowLevelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow level ID"})
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dwc.db, c, err.Error())
		return
	}

	if err := dwc.workflowSvc.DeleteWorkflowLevel(currentUserID, uint(workflowLevelID)); err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow level deleted successfully"})
}
