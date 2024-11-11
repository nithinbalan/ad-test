package controllers

import (
	logDTOs "jethur/grc/logs/dtos"
	"jethur/grc/logs/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ActivityLogController struct {
	service services.ActivityLogServiceInterface
	db      *gorm.DB
}

func NewActivityLogController(service services.ActivityLogServiceInterface) *ActivityLogController {
	return &ActivityLogController{
		service: service,
	}
}

func (alc *ActivityLogController) GetAllActivityLogs(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"user_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(alc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(alc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := logDTOs.ActivityLogListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	logs, err := alc.service.GetAll(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(alc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, logs)
}

func (alc *ActivityLogController) GetActivityLogDetails(c *gin.Context) {
	logIDStr := c.Param("logId")

	logID, err := strconv.ParseUint(logIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	activityLog, err := alc.service.GetFullDetailsByID(uint(logID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(alc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, activityLog)
}
