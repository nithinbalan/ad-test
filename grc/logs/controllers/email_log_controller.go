package controllers

import (
	"errors"
	logDTOs "jethur/grc/logs/dtos"
	"jethur/grc/logs/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmailLogController struct {
	service services.EmailLogServiceInterface
	db      *gorm.DB
}

func NewEmailLogController(service services.EmailLogServiceInterface) *EmailLogController {
	return &EmailLogController{
		service: service,
	}
}

func (elc *EmailLogController) GetAllEmailLogs(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	allowedParams := []string{
		"email_status_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	validSortFields := map[string]string{
		"id":                    "emails.id",
		"subject":               "emails.subject",
		"module_id":             "emails.module_id",
		"sub_module_id":         "emails.sub_module_id",
		"email_status_id":       "emails.email_status_id",
		"is_permanent_failed":   "emails.is_permanent_failed",
		"attempt_count":         "emails.attempt_count",
		"max_attempts":          "emails.max_attempts",
		"last_attempt_at":       "emails.last_attempt_at",
		"sent_at":               "emails.sent_at",
		"created_at":            "emails.created_at",
		"email_status_title_en": "email_statuses.title_en",
		"email_status_title_ar": "email_statuses.title_ar",
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(elc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(elc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := logDTOs.EmailLogListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		StartDate:   startDate,
		EndDate:     endDate,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
		Filters:     filters,
	}

	logs, err := elc.service.GetAll(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(elc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, logs)
}

func (elc *EmailLogController) GetEmailLogDetails(c *gin.Context) {
	logIDStr := c.Param("logId")

	logID, err := strconv.ParseUint(logIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	emailLog, err := elc.service.GetFullDetailsByID(uint(logID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(elc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, emailLog)
}

func (elc *EmailLogController) GetSummary(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	queryParams := logDTOs.EmailLogListQueryParamsDTO{
		StartDate: startDate,
		EndDate:   endDate,
	}

	summary, err := elc.service.GetSummary(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(elc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, summary)
}
