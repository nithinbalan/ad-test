package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReminderController struct {
	service services.ReminderServiceInterface
	db      *gorm.DB
}

func NewReminderController(service services.ReminderServiceInterface) *ReminderController {
	return &ReminderController{
		service: service,
	}
}

func (rc *ReminderController) GetAllReminderes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	allowedParams := []string{
		"module_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := dtos.ReminderListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	reminders, err := rc.service.GetAll(queryParams)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, reminders)
}

func (rc *ReminderController) GetReminderDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reminder, err := rc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, reminder)
}

func (rc *ReminderController) GetAllEnabledModules(c *gin.Context) {
	modules, err := rc.service.GetAllEnabledModules()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, modules)
}

func (rc *ReminderController) GetReminderByModuleID(c *gin.Context) {
	moduleIDStr := c.Param("moduleId")

	moduleID, err := strconv.ParseUint(moduleIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	stakeholderStatuses, err := rc.service.GetAllRemindersByModuleID(uint(moduleID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, stakeholderStatuses)
}

func (rc *ReminderController) GetAllReminderContexts(c *gin.Context) {
	data, err := rc.service.GetAllReminderContexts()
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}
