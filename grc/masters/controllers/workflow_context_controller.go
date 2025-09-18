package controllers

import (
	"errors"
	"jethur/grc/masters/dtos"
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowContextController struct {
	service services.WorkflowContextServiceInterface
	db      *gorm.DB
}

func NewWorkflowContextController(service services.WorkflowContextServiceInterface) *WorkflowContextController {
	return &WorkflowContextController{
		service: service,
	}
}

func (wsc *WorkflowContextController) GetAllWorkflowContext(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"module_ids",
	}

	validSortFields := map[string]string{
		"title_en": "workflow_contexts.title_en",
		"title_ar": "workflow_contexts.title_ar",
		"type":     "workflow_contexts.type",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := dtos.WorkflowContextListQueryParamsDTO{
		Sort:    sort,
		SortBy:  dbField,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	workflowContexts, err := wsc.service.GetAll(queryParams)

	if err != nil {
		utility.SendInternalErrResp(wsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, workflowContexts)
}

func (wsc *WorkflowContextController) GetWorkflowContextDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInternalErrResp(wsc.db, c, "Invalid id")
		return
	}

	workflowContext, err := wsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, workflowContext)
}
