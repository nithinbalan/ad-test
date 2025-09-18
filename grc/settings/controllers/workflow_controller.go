package controllers

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	settingDTOs "jethur/grc/settings/dtos"
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
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	moduleIDs := c.QueryArray("module_ids")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"workflow_context_ids",
	}

	validSortFields := map[string]string{
		"title_en": "workflows.title_en",
		"title_ar": "workflows.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := settingDTOs.WorkflowListQueryParamsDTO{
		Sort:    sort,
		SortBy:  dbField,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	workflows, err := wfc.service.GetAll(&queryParams, moduleIDs)
	if err != nil {
		utility.SendInternalErrResp(wfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, workflows)
}

func (wfc *WorkflowController) GetWorkflowDetails(c *gin.Context) {
	idStr := c.Param("workflowId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflow, err := wfc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(wfc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, workflow)
}

func (wfc *WorkflowController) CreateWorkflow(c *gin.Context) {
	var workflowFormDTO settingDTOs.WorkflowFormDTO

	if err := c.ShouldBindJSON(&workflowFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(wfc.db, c, err.Error())
		return
	}

	id, err := wfc.service.Create(currentUserID, &workflowFormDTO)

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

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (wfc *WorkflowController) UpdateWorkflow(c *gin.Context) {
	workflowID, err := strconv.Atoi(c.Param("workflowId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var workflowFormDTO settingDTOs.WorkflowFormDTO
	if err := c.ShouldBindJSON(&workflowFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(wfc.db, c, err.Error())
		return
	}

	err = wfc.service.Update(uint(workflowID), currentUserID, &workflowFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(wfc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)

}

func (wfc *WorkflowController) DeleteWorkflows(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(wfc.db, c, err.Error())
		return
	}

	failedIDs, err := wfc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(wfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (wfc *WorkflowController) GetAllWorkflowsByType(c *gin.Context) {
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
