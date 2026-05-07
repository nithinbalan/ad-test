package controllers

import (
	"fmt"
	"jethur/grc/logs/dtos"
	"jethur/grc/logs/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	defaultPageSize = 10
	maxPageSize     = 200
)

type APIRequestLogController struct {
	service services.APIRequestLogServiceInterface
	db      *gorm.DB
}

func NewAPIRequestLogController(db *gorm.DB, service services.APIRequestLogServiceInterface) *APIRequestLogController {
	return &APIRequestLogController{
		service: service,
		db:      db,
	}
}

func (arlc *APIRequestLogController) GetAllAPIRequestLogs(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "requested_at")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid page parameter"))
		return
	}
	if page < 1 {
		page = 1
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", defaultPageSize)
	if err != nil {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid page_size parameter"))
		return
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	validSortFields := map[string]string{
		"id":           "api_request_logs.id",
		"method":       "api_request_logs.method",
		"url":          "api_request_logs.url",
		"requested_at": "api_request_logs.requested_at",
		"status_code":  "api_request_logs.status_code",
		"duration_ms":  "api_request_logs.duration_ms",
		"log_status":   "api_request_logs.log_status",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid sort_by"))
		return
	}

	queryParams := dtos.APIRequestLogListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		QuickFilter: quickFilter,
	}

	result, err := arlc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(arlc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (arlc *APIRequestLogController) GetAPIRequestLogsSummary(c *gin.Context) {
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.APIRequestLogListQueryParamsDTO{
		Keyword:     keyword,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		QuickFilter: quickFilter,
	}

	result, err := arlc.service.GetSummary(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(arlc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
