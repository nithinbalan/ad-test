package controllers

import (
	"errors"
	"jethur/grc/audit-insights/services"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditInsightsFindingActionPlanController struct {
	service services.AuditInsightsFindingActionPlanServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsFindingActionPlanController(service services.AuditInsightsFindingActionPlanServiceInterface) *AuditInsightsFindingActionPlanController {
	return &AuditInsightsFindingActionPlanController{
		service: service,
	}
}

func (aifapc *AuditInsightsFindingActionPlanController) GetAllAuditInsightsFindingActionPlans(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_finding_ids",
		"audit_ids",
		"auditable_item_ids",
		"audit_finding_category_ids",
		"severity_level_ids",
		"audit_finding_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())
	
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(aifapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aifapc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"reference_number": "audit_findings.reference_number",
		"title":            "audit_findings.title",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := auditDTOs.AuditFindingActionPlanListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aifapc.db, c, err.Error())
		return
	}
	audits, err := aifapc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aifapc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, audits)
}

func (aifapc *AuditInsightsFindingActionPlanController) GetAuditInsightsFindingActionPlanDetails(c *gin.Context) {
	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	audit, err := aifapc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aifapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, audit)
}
