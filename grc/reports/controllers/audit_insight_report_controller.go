package controllers

import (
	"errors"
	auditInsightService "jethur/grc/audit-insights/services"
	auditDTOs "jethur/grc/audit/dtos"
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditInsightReportController struct {
	service                           services.AuditInsightReportServiceInterface
	auditInsightsSvc                  auditInsightService.AuditInsightsAuditServiceInterface
	auditInsightsFindingSvc           auditInsightService.AuditInsightsAuditFindingServiceInterface
	auditInsightsFindingActionPlanSvc auditInsightService.AuditInsightsFindingActionPlanServiceInterface
	auditInsightsAuditPlanSvc         auditInsightService.AuditInsightsAuditPlanServiceInterface
	auditInsightsDraftReportSvc       auditInsightService.AuditInsightsDraftReportServiceInterface
	auditInsightsFinalReportSvc       auditInsightService.AuditInsightsFinalReportServiceInterface
	db                                *gorm.DB
}

func NewAuditInsightReportController(service services.AuditInsightReportServiceInterface,
	auditInsightsSvc auditInsightService.AuditInsightsAuditServiceInterface,
	auditInsightsFindingSvc auditInsightService.AuditInsightsAuditFindingServiceInterface,
	auditInsightsFindingActionPlanSvc auditInsightService.AuditInsightsFindingActionPlanServiceInterface,
	auditInsightsAuditPlanSvc auditInsightService.AuditInsightsAuditPlanServiceInterface,
	auditInsightsDraftReportSvc auditInsightService.AuditInsightsDraftReportServiceInterface,
	auditInsightsFinalReportSvc auditInsightService.AuditInsightsFinalReportServiceInterface,
) *AuditInsightReportController {
	return &AuditInsightReportController{
		service:                           service,
		auditInsightsSvc:                  auditInsightsSvc,
		auditInsightsFindingSvc:           auditInsightsFindingSvc,
		auditInsightsFindingActionPlanSvc: auditInsightsFindingActionPlanSvc,
		auditInsightsAuditPlanSvc:         auditInsightsAuditPlanSvc,
		auditInsightsDraftReportSvc:       auditInsightsDraftReportSvc,
		auditInsightsFinalReportSvc:       auditInsightsFinalReportSvc,
	}
}

func (airc *AuditInsightReportController) GetAuditByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditByStatus(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditByYear(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditByYear(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditByMonth(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditByMonth(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditPlanByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditPlanByCategory(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditPlanByTeamLead(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditPlanByTeamLead(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditPlanByAuditableItem(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditPlanByAuditableItem(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditPlanByMonth(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditPlanByMonth(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditPlanVsAudits(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditPlanVsAudits(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditFindingByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditFindingByCategory(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}
func (airc *AuditInsightReportController) GetAuditFindingByAuditCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditFindingByAuditCategory(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditFindingByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditFindingByStatus(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditFindingBySeverityLevel(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditFindingBySeverityLevel(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditFindingByMonth(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	results, err := airc.service.GetAuditFindingByMonth(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditFindingByYear(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	results, err := airc.service.GetAuditFindingByYear(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditActionPlanByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	results, err := airc.service.GetAuditActionPlanByStatus(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditActionPlanByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	results, err := airc.service.GetAuditActionPlanByResponsibleUser(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditDraftReportByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	results, err := airc.service.GetAuditDraftReportByStatus(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditFinalReportByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	results, err := airc.service.GetAuditFinalReportByStatus(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetOpenActionPlanByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetOpenActionPlanByResponsibleUser(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetOverdueActionPlanByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetOverdueActionPlanByResponsibleUser(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAuditInsightOverdueActionPlans(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}

	results, err := airc.service.GetAuditInsightOverdueActionPlans(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (airc *AuditInsightReportController) GetAllAuditInsightsAudits(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_category_ids",
		"audit_team_ids",
		"team_lead_ids",
		"audit_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := auditDTOs.AuditListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		StartDate:   startDate,
		EndDate:     endDate,
		Export:      export,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	audits, err := airc.auditInsightsSvc.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
			return
		}

		excelFile, err := airc.auditInsightsSvc.ExportAuditInsights(audits.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Audit Insights.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, audits)
}

func (airc *AuditInsightReportController) GetAllAuditInsightsAuditsExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	airc.GetAllAuditInsightsAudits(c)
}

func (airc *AuditInsightReportController) GetAllAuditInsightsActionPlans(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	// Check if `findingId` exists in the route

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid export parameter")
		return
	}
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

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid page_size parameter")
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
		Export:      export,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	audits, err := airc.auditInsightsFindingActionPlanSvc.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
			return
		}

		excelFile, err := airc.auditInsightsFindingActionPlanSvc.ExportAuditsightsFindingActionPlans(audits.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Audit Insights Finding Action Plans.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, audits)
}

func (airc *AuditInsightReportController) GetAuditInsightsActionPlanExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	airc.GetAllAuditInsightsActionPlans(c)
}

func (airc *AuditInsightReportController) GetAllAuditInsightsAuditFindings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_ids",
		"auditable_item_ids",
		"audit_finding_category_ids",
		"severity_level_ids",
		"audit_finding_status_ids",
		"audit_category_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid page_size parameter")
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
	queryParams := auditDTOs.AuditFindingListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	audits, err := airc.auditInsightsFindingSvc.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
			return
		}

		excelFile, err := airc.auditInsightsFindingSvc.ExportAuditInsightsAuditFindings(audits.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Audit Insights Audit Findings.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, audits)
}

func (airc *AuditInsightReportController) GetAuditInsightsAuditFindingsExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	airc.GetAllAuditInsightsAuditFindings(c)
}

func (airc *AuditInsightReportController) GetAllAuditInsightsAuditPlans(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid export parameter")
		return
	}
	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_category_ids",
		"audit_team_ids",
		"team_lead_ids",
		"audit_status_ids",
		"auditable_item_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	startDate := c.Query("start_date")
	targetDate := c.Query("target_date")

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := auditDTOs.AuditPlanListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		StartDate:   startDate,
		TargetDate:  targetDate,
		Export:      export,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	auditPlans, err := airc.auditInsightsAuditPlanSvc.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
			return
		}

		excelFile, err := airc.auditInsightsAuditPlanSvc.ExportAuditInsightsAuditPlans(auditPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Audit Insights Audit Plans.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, auditPlans)
}

func (airc *AuditInsightReportController) GetAuditInsightsAuditPlansExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	airc.GetAllAuditInsightsAuditPlans(c)
}

func (airc *AuditInsightReportController) GetAllAuditInsightsDraftReports(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid export parameter")
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
		utility.SendInternalErrResp(airc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid page_size parameter")
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
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	auditDraftReports, err := airc.auditInsightsDraftReportSvc.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
			return
		}

		excelFile, err := airc.auditInsightsDraftReportSvc.ExportAuditInsightsDraftReports(auditDraftReports.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
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

func (airc *AuditInsightReportController) GetAuditInsightsDraftReportsExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	airc.GetAllAuditInsightsDraftReports(c)
}

func (airc *AuditInsightReportController) GetAllAuditInsightsFinalReports(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.DefaultQuery("sort_by", "audit_plan_title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_ids",
		"audit_category_ids",
		"audit_team_ids",
		"team_lead_ids",
		"audit_status_ids",
		"audit_final_report_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	reportDate := c.Query("report_date")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, "Invalid page_size parameter")
		return
	}
	validSortFields := map[string]string{
		"audit_plan_title":                   "audit_plans.title",
		"start_date":                         "audits.start_date",
		"end_date":                           "audits.end_date",
		"audit_category_title_en":            "audit_categories.title_en",
		"audit_category_title_ar":            "audit_categories.title_ar",
		"audit_status_title_en":              "audit_statuses.title_en",
		"audit_status_title_ar":              "audit_statuses.title_ar",
		"audit_final_report_status_title_en": "audit_final_report_statuses.title_en",
		"audit_final_report_status_title_ar": "audit_final_report_statuses.title_ar",
		"team_lead_first_name_en":            "team_lead.first_name_en",
		"team_lead_first_name_ar":            "team_lead.first_name_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := auditDTOs.AuditFinalReportListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
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
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	auditFinalReports, err := airc.auditInsightsFinalReportSvc.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(airc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
			return
		}

		excelFile, err := airc.auditInsightsFinalReportSvc.ExportAuditInsightsFinalReports(auditFinalReports.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(airc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Audit Insights Final Reports.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, auditFinalReports)
}

func (airc *AuditInsightReportController) GetAuditInsightsFinalReportsExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	airc.GetAllAuditInsightsFinalReports(c)
}
