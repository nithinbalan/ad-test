package controllers

import (
	"errors"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentService "jethur/grc/incident/services"
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentReportController struct {
	service               services.IncidentReportServiceInterface
	incidentRegisterSvc   incidentService.IncidentRegisterServiceInterface
	incidentActionPlanSvc incidentService.IncidentActionPlanServiceInterface

	db *gorm.DB
}

func NewIncidentReportController(service services.IncidentReportServiceInterface,
	incidentRegisterSvc incidentService.IncidentRegisterServiceInterface,
	incidentActionPlanSvc incidentService.IncidentActionPlanServiceInterface) *IncidentReportController {
	return &IncidentReportController{
		service:               service,
		incidentRegisterSvc:   incidentRegisterSvc,
		incidentActionPlanSvc: incidentActionPlanSvc,
	}
}

func (irc *IncidentReportController) GetAllIncidents(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"incident_category_ids",
		"incident_sub_category_ids",
		"severity_level_ids",
		"incident_status_ids",
		"reported_user_ids",
		"incident_type_ids",
		"severity_level_ids",
		"responsible_unit_ids",
		"responsible_user_ids",
	}
	validSortFields := map[string]string{
		"incident_at":                    "incident_register.incident_at",
		"reported_at":                    "incident_register.reported_at",
		"reference_number":               "incident_register.reference_number",
		"title":                          "incident_register.title",
		"description":                    "incident_register.description",
		"incident_category_en":           "incident_categories.title_en",
		"incident_category_ar":           "incident_categories.title_ar",
		"incident_type_en":               "incident_types.title_en",
		"incident_type_ar":               "incident_types.title_ar",
		"incident_status_en":             "incident_statuses.title_en",
		"incident_status_ar":             "incident_statuses.title_ar",
		"severity_level_en":              "severity_levels.title_en",
		"severity_level_ar":              "severity_levels.title_ar",
		"severity_level_color_code":      "severity_levels.color_code",
		"responsible_unit_title_en":      "organization_units.title_en",
		"responsible_unit_title_ar":      "organization_units.title_ar",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"responsible_user_last_name_en":  "responsible_user.last_name_en",
		"responsible_user_last_name_ar":  "responsible_user.last_name_ar",
		"responsible_user_email":         "responsible_user.email",
		"responsible_user_initial":       "responsible_user.initial",
		"reported_user_first_name_en":    "reported_user.first_name_en",
		"reported_user_first_name_ar":    "reported_user.first_name_ar",
		"reported_user_last_name_en":     "reported_user.last_name_en",
		"reported_user_last_name_ar":     "reported_user.last_name_ar",
		"location":                       "incident_register.location",
		"quick_action":                   "incident_register.quick_action",
		"investigation_initiated_at":     "incident_investigations.initiated_at",
		"investigation_description":      "incident_investigations.description",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	incidentAt := c.Query("incident_at")
	reportedAt := c.Query("reported")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := incidentDTOs.IncidentRegisterListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		IncidentAt:  incidentAt,
		ReportedAt:  reportedAt,
		Export:      export,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}
	incidents, err := irc.incidentRegisterSvc.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(irc.db, c, err.Error())
			return
		}

		excelFile, err := irc.incidentRegisterSvc.ExportIncidentRegister(incidents.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(irc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Incident Register.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, incidents)
}

func (irc *IncidentReportController) GetIncidentExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	irc.GetAllIncidents(c)
}

func (irc *IncidentReportController) GetIncidentByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentByCategory(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentBySeverityLevel(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentBySeverityLevel(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentByReportedUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentByReportedUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentByReportedYear(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentByReportedYear(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}
func (irc *IncidentReportController) GetIncidentByAsset(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentByAsset(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}
func (irc *IncidentReportController) GetIncidentBySLA(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentBySLA(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentInvestigationSummary(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentInvestigationSummary(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentRCA(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentRCA(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentClosure(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentClosure(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetOpenIncidentsWithRisk(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetOpenIncidentsWithRisk(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentPendingActionPlan(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentPendingActionPlan(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentActionPlanByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentActionPlanByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentActionPlanByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentActionPlanByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetIncidentActionPlanByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := irc.service.GetIncidentActionPlanByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (irc *IncidentReportController) GetAllIncidentActionPlans(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	var incidentID uint64
	var err error

	// Check if incidentID is provided in the URL
	if incidentIDStr != "" {
		incidentID, err = strconv.ParseUint(incidentIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		incidentID = 0
	}
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"severity_level_ids",
		"responsible_unit_ids",
		"responsible_user_ids",
		"action_plan_status_ids",
		"incident_register_ids",
	}

	validSortFields := map[string]string{
		"title":               "incident_action_plans.title",
		"description":         "incident_action_plans.description",
		"incident_title":      "incident_register.title",
		"severity_level_en":   "severity_levels.title_en",
		"severity_level_ar":   "severity_levels.title_ar",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"target_date":         "incident_action_plans.target_date",
		"percentage":          "incident_action_plans.percentage",
		"status_en":           "action_plan_statuses.title_en",
		"status_ar":           "action_plan_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := incidentDTOs.IncidentActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         dbField,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
		Export:         export,
		DateFrom:       dateFrom,
		DateTo:         dateTo,
	}
	incidentActionPlans, err := irc.incidentActionPlanSvc.GetAll(uint(incidentID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(irc.db, c, err.Error())
			return
		}

		excelFile, err := irc.incidentActionPlanSvc.ExportIncidentActionPlan(incidentActionPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(irc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Incident action plan.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, incidentActionPlans)
}

func (irc *IncidentReportController) GetIncidentActionPlanExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	irc.GetAllIncidentActionPlans(c)
}
