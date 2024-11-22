package controllers

import (
	"jethur/grc/audit-insights/services"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditInsightsAuditFindingController struct {
	service services.AuditInsightsAuditFindingServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsAuditFindingController(service services.AuditInsightsAuditFindingServiceInterface) *AuditInsightsAuditFindingController {
	return &AuditInsightsAuditFindingController{
		service: service,
	}
}

func (aiafc *AuditInsightsAuditFindingController) GetAllAuditInsightsAuditFindings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_ids",
		"auditable_item_ids",
		"audit_finding_category_ids",
		"severity_level_ids",
		"audit_finding_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(aiafc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aiafc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := auditDTOs.AuditFindingListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aiafc.db, c, err.Error())
		return
	}
	audits, err := aiafc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aiafc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, audits)
}

func (aiafc *AuditInsightsAuditFindingController) GetAuditInsightsAuditFindingDetails(c *gin.Context) {
	idStr := c.Param("findingId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	audit, err := aiafc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aiafc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, audit)
}
