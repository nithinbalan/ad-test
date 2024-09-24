package controllers

import (
	"jethur/grc/audit-insights/services"
	"jethur/grc/audit/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditInsightsDocumentRequestController struct {
	service services.AuditInsightsDocumentRequestServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsDocumentRequestController(service services.AuditInsightsDocumentRequestServiceInterface) *AuditInsightsDocumentRequestController {
	return &AuditInsightsDocumentRequestController{
		service: service,
	}
}

func (aidrc *AuditInsightsDocumentRequestController) GetAllAuditInsightsDocumentRequests(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	requestedAt := c.Query("requested_at")

	allowedParams := []string{
		"audit_ids",
		"audit_document_request_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(aidrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aidrc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := dtos.AuditDocumentRequestListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		RequestedAt: requestedAt,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aidrc.db, c, err.Error())
		return
	}
	auditDocumentRequests, err := aidrc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aidrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditDocumentRequests)
}

func (aidrc *AuditInsightsDocumentRequestController) GetAuditInsightsDocumentRequestDetails(c *gin.Context) {
	idStr := c.Param("documentRequestId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditDocumentRequest, err := aidrc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aidrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditDocumentRequest)
}
