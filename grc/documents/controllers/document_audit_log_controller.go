package controllers

import (
	"errors"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentAuditLogController struct {
	service services.DocumentAuditLogServiceInterface
	db      *gorm.DB
}

func NewDocumentAuditLogController(service services.DocumentAuditLogServiceInterface) *DocumentAuditLogController {
	return &DocumentAuditLogController{
		service: service,
	}
}

func (dalc *DocumentAuditLogController) GetAllDocumentAuditLogs(c *gin.Context) {
	documentIDStr := c.Param("id")

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "action_date")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	validSortFields := map[string]string{
		"action_date":     "document_audit_logs.created_at",
		"action_en":       "document_audit_log_types.title_en",
		"action_ar":       "document_audit_log_types.title_ar",
		"performed_by_en": "users.first_name_en",
		"performed_by_ar": "users.first_name_ar",
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(dalc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(dalc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := documentDTOs.DocumentAuditLogListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}

	documentTimelines, err := dalc.service.GetAll(uint(documentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dalc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, documentTimelines)
}
