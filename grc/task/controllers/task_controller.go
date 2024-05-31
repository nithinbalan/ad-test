package controllers

import (
	taskDTOs "jethur/grc/task/dtos"
	"jethur/grc/task/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	service services.TaskServiceInterface
	db      *gorm.DB
}

func NewTaskController(service services.TaskServiceInterface) *TaskController {
	return &TaskController{
		service: service,
	}
}

func (tc *TaskController) GetAllTasks(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"responsible_user_ids",
		"compliance_action_plan_ids",
		"sla_action_plan_ids",
		"incident_action_plan_ids",
		"control_action_plan_ids",
		"severity_level_ids",
		"module_ids",
		"task_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := taskDTOs.TaskListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortBy,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
	}
	tasks, err := tc.service.GetAll(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, tasks)
}
