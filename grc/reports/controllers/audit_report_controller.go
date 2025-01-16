package controllers

import (
	"errors"
	auditDTOs "jethur/grc/audit/dtos"
	auditService "jethur/grc/audit/services"
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditReportController struct {
	service                   services.AuditReportServiceInterface
	auditPlanSvc              auditService.AuditPlanServiceInterface
	auditSvc                  auditService.AuditServiceInterface
	auditFindingSvc           auditService.AuditFindingServiceInterface
	auditFindingActionPlanSvc auditService.AuditFindingActionPlanServiceInterface

	db *gorm.DB
}

func NewAuditReportController(service services.AuditReportServiceInterface,
	auditPlanSvc auditService.AuditPlanServiceInterface,
	auditSvc auditService.AuditServiceInterface,
	auditFindingSvc auditService.AuditFindingServiceInterface,
	auditFindingActionPlanSvc auditService.AuditFindingActionPlanServiceInterface) *AuditReportController {
	return &AuditReportController{
		service:                   service,
		auditPlanSvc:              auditPlanSvc,
		auditSvc:                  auditSvc,
		auditFindingSvc:           auditFindingSvc,
		auditFindingActionPlanSvc: auditFindingActionPlanSvc,
	}
}

func (arc *AuditReportController) GetAllAuditPlans(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_category_ids",
		"audit_team_ids",
		"team_lead_ids",
		"audit_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	startDate := c.Query("start_date")
	targetDate := c.Query("target_date")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, "Invalid page_size parameter")
		return
	}
	validSortFields := map[string]string{
		"reference_number":        "audit_plans.reference_number",
		"title":                   "audit_plans.title",
		"start_date":              "audit_plans.start_date",
		"target_date":             "audit_plans.target_date",
		"audit_category_title_en": "audit_categories.title_en",
		"audit_category_title_ar": "audit_categories.title_ar",
		"audit_status_en":         "audit_statuses.title_en",
		"audit_status_ar":         "audit_statuses.title_ar",
		"team_lead_first_name_en": "team_lead.first_name_en",
		"team_lead_first_name_ar": "team_lead.first_name_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := auditDTOs.AuditPlanListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		StartDate:   startDate,
		TargetDate:  targetDate,
	}

	auditPlans, err := arc.auditPlanSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditPlans)
}

func (arc *AuditReportController) GetAuditPlanExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	arc.GetAllAuditPlans(c)
}

func (arc *AuditReportController) GetAuditPlanByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditPlanByCategory(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditPlanByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditPlanByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditPlanByTeam(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditPlanByTeam(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditPlanByTeamLead(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditPlanByTeamLead(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditPlanByYear(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditPlanByYear(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditByCategory(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditByOpenVsClosed(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditByOpenVsClosed(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditByMonth(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditByMonth(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditByTeamLead(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditByTeamLead(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditFindingByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditFindingByCategory(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditFindingByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditFindingByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditFindingBySeverityLevel(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditFindingBySeverityLevel(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditFindingByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditFindingByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditFindingByMonth(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditFindingByMonth(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditFindingByAuditor(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditFindingByAuditor(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditActionPlanByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditActionPlanByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditActionPlanBySeverityLevel(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditActionPlanBySeverityLevel(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditActionPlanByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditActionPlanByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditActionPlanByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditActionPlanByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAuditActionPlanOnTimeVsDelayed(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := arc.service.GetAuditActionPlanOnTimeVsDelayed(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (arc *AuditReportController) GetAllAudits(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.DefaultQuery("sort_by", "audit_plan_title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_plan_ids",
		"audit_category_ids",
		"audit_team_ids",
		"team_lead_ids",
		"audit_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, "Invalid page_size parameter")
		return
	}
	validSortFields := map[string]string{
		"audit_plan_title":        "audit_plans.title",
		"start_date":              "audits.start_date",
		"end_date":                "audits.end_date",
		"audit_category_title_en": "audit_categories.title_en",
		"audit_category_title_ar": "audit_categories.title_ar",
		"audit_status_title_en":   "audit_statuses.title_en",
		"audit_status_title_ar":   "audit_statuses.title_ar",
		"team_lead_first_name_en": "team_lead.first_name_en",
		"team_lead_first_name_ar": "team_lead.first_name_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := auditDTOs.AuditListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	audits, err := arc.auditSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, audits)
}

func (arc *AuditReportController) GetAuditExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	arc.GetAllAudits(c)
}

func (arc *AuditReportController) GetAllAuditFindings(c *gin.Context) {

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
		utility.SendInternalErrResp(arc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, "Invalid page_size parameter")
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

	auditFindings, err := arc.auditFindingSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditFindings)
}
func (arc *AuditReportController) GetAuditFindingExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	arc.GetAllAuditFindings(c)
}

func (arc *AuditReportController) GetAllAuditActionPlans(c *gin.Context) {

	var findingID uint64
	var err error

	findingIDStr := c.Param("findingId")

	findingID, err = strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"severity_level_ids",
		"responsible_user_ids",
		"action_plan_status_ids",
		"audit_finding_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := auditDTOs.AuditFindingActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortBy,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
	}
	auditFindingActionPlans, err := arc.auditFindingActionPlanSvc.GetAll(uint(findingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditFindingActionPlans)
}

func (arc *AuditReportController) GetAuditActionPlanExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	arc.GetAllAuditActionPlans(c)
}
