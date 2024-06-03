package controllers

import (
	"jethur/grc/task/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskStatusController struct {
	service services.TaskStatusServiceInterface
	db      *gorm.DB
}

func NewTaskStatusController(service services.TaskStatusServiceInterface) *TaskStatusController {
	return &TaskStatusController{
		service: service,
	}
}

func (tsc *TaskStatusController) GetAllTaskStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	taskStatuses, err := tsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, taskStatuses)
}

func (tsc *TaskStatusController) GetTaskStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	taskStatus, err := tsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, taskStatus)
}
