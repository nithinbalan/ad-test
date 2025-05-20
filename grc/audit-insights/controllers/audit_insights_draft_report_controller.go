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

type AuditInsightsDraftReportController struct {
	service services.AuditInsightsDraftReportServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsDraftReportController(service services.AuditInsightsDraftReportServiceInterface) *AuditInsightsDraftReportController {
	return &AuditInsightsDraftReportController{
		service: service,
	}
}

func (aidrc *AuditInsightsDraftReportController) GetAllAuditInsightsDraftReports(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(aidrc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_ids",
		"audit_category_ids",
		"audit_team_ids",
		"team_lead_ids",
		"audit_status_ids",
		"audit_draft_report_status_ids",
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

	queryParams := auditDTOs.AuditDraftReportListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		ReportDate:  reportDate,
		Export:      export,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aidrc.db, c, err.Error())
		return
	}
	auditDraftReports, err := aidrc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aidrc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(aidrc.db, c, err.Error())
			return
		}

		excelFile, err := aidrc.service.ExportAuditInsightsDraftReports(auditDraftReports.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(aidrc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Audit Insights Draf Reports.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, auditDraftReports)
}

func (aidrc *AuditInsightsDraftReportController) GetAuditInsightsDraftReportExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	aidrc.GetAllAuditInsightsDraftReports(c)
}

func (aidrc *AuditInsightsDraftReportController) GetAuditInsightsDraftReportDetails(c *gin.Context) {
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
