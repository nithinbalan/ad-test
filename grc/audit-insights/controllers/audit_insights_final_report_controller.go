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

type AuditInsightsFinalReportController struct {
	service services.AuditInsightsFinalReportServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsFinalReportController(service services.AuditInsightsFinalReportServiceInterface) *AuditInsightsFinalReportController {
	return &AuditInsightsFinalReportController{
		service: service,
	}
}

func (aidrc *AuditInsightsFinalReportController) GetAllAuditInsightsFinalReports(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_ids",
		"audit_category_ids",
		"audit_team_ids",
		"team_lead_ids",
		"audit_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	reportDate := c.Query("report_date")

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

	queryParams := auditDTOs.AuditFinalReportListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		ReportDate:  reportDate,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aidrc.db, c, err.Error())
		return
	}
	auditFinalReports, err := aidrc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aidrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditFinalReports)
}

func (aidrc *AuditInsightsFinalReportController) GetAuditInsightsFinalReportDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aidrc.db, c, err.Error())
		return
	}
	audit, err := aidrc.service.GetFullDetailsByID(uint(id), currentUserID)
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
	c.JSON(http.StatusOK, audit)
}
