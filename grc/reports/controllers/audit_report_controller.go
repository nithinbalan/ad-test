package controllers

import (
	"errors"
	auditDTOs "jethur/grc/audit/dtos"
	auditService "jethur/grc/audit/services"
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditReportController struct {
	service      services.AuditReportServiceInterface
	auditPlanSvc auditService.AuditPlanServiceInterface
	auditSvc     auditService.AuditServiceInterface

	db *gorm.DB
}

func NewAuditReportController(service services.AuditReportServiceInterface,
	auditPlanSvc auditService.AuditPlanServiceInterface,
	auditSvc auditService.AuditServiceInterface) *AuditReportController {
	return &AuditReportController{
		service:      service,
		auditPlanSvc: auditPlanSvc,
		auditSvc:     auditSvc,
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
