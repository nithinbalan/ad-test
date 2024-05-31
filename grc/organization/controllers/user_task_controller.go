package controllers

import (
	"jethur/grc/organization/services"
	taskDTOs "jethur/grc/task/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserTaskController struct {
	service services.UserTaskServiceInterface
	db      *gorm.DB
}

func NewUserTaskController(service services.UserTaskServiceInterface) *UserTaskController {
	return &UserTaskController{
		service: service,
	}
}
func (utc *UserTaskController) GetAllUserTasks(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_action_plan_ids",
		"sla_action_plan_ids",
		"incident_action_plan_ids",
		"control_action_plan_ids",
		"severity_level_ids",
		"module_ids",
		"task_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())
	// Filter for responsible_user_id=userID
	filters["responsible_user_id"] = []string{strconv.Itoa(userID)}

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(utc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(utc.db, c, "Invalid page_size parameter")
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
	userTasks, err := utc.service.GetAll(uint(userID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(utc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, userTasks)
}

// func (utc *UserTaskController) GetUserTaskSummary(c *gin.Context) {
// 	userID, err := strconv.Atoi(c.Param("userId"))
// 	if err != nil {
// 		utility.SendNotFoundErrResp(c)
// 		return
// 	}

// 	userTaskSummary, err := utc.service.GetUserTaskSummary(uint(userID))
// 	if err != nil {
// 		switch e := err.(type) {
// 		case *utility.ResourceNotFoundError:
// 			utility.SendNotFoundErrResp(c)
// 			return
// 		default:
// 			utility.SendInternalErrResp(utc.db, c, e.Error())
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusOK, userTaskSummary)
// }

