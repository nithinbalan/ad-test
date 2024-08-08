package controllers

import (
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

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

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
	queryParams := documentDTOs.DocumentAuditLogListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	if err != nil {
		utility.SendInternalErrResp(dalc.db, c, "Invalid page_size parameter")
		return
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
