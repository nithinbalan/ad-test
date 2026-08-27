package controllers

import (
	"errors"
	auditDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditDashboardController struct {
	service services.AuditDashboardServiceInterface
	db      *gorm.DB
}

func NewAuditDashboardController(service services.AuditDashboardServiceInterface, db *gorm.DB) *AuditDashboardController {
	return &AuditDashboardController{service: service, db: db}
}

func (adc *AuditDashboardController) buildParams(c *gin.Context) *auditDTOs.AuditDashboardQueryParamsDTO {
	year, _ := strconv.Atoi(c.Query("year"))
	return &auditDTOs.AuditDashboardQueryParamsDTO{
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Year:     year,
	}
}

func (adc *AuditDashboardController) GetAuditStatusBreakdown(c *gin.Context) {
	result, err := adc.service.GetAuditStatusBreakdown(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AuditDashboardController) GetAuditSummary(c *gin.Context) {
	result, err := adc.service.GetAuditSummary(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AuditDashboardController) GetAuditsByStatus(c *gin.Context) {
	statusIDStr := c.Param("statusId")
	statusID, err := strconv.ParseUint(statusIDStr, 10, 64)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid statusId parameter"))
		return
	}
	sortOrder := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	exportFromContext, exists := c.Get("export")
	var export bool
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(adc.db, c, "Invalid export parameter")
			return
		}
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"start_date":       "audits.start_date",
		"end_date":         "audits.end_date",
		"audit_plan_title": "audit_plans.title",
		"created_at":       "audits.created_at",
		"id":               "audits.id",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	queryParams := auditDTOs.AuditDashboardListQueryParamsDTO{
		Sort:     sortOrder,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Export:   export,
	}

	audits, err := adc.service.GetAuditsByStatus(uint(statusID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(adc.db, c, err.Error())
			return
		}
		excelFile, err := adc.service.ExportAuditsByStatus(audits.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(adc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		fileName := adc.service.GetAuditsByStatusExportFilename(uint(statusID))
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="`+fileName+`"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}
	c.JSON(http.StatusOK, audits)
}

func (adc *AuditDashboardController) GetAuditsByStatusExport(c *gin.Context) {
	c.Set("export", true)
	adc.GetAuditsByStatus(c)
}

func (adc *AuditDashboardController) GetAuditFindingsSummary(c *gin.Context) {
	result, err := adc.service.GetAuditFindingsSummary(adc.buildParams(c))
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AuditDashboardController) GetAuditFindingsSeverityBreakdown(c *gin.Context) {
	result, err := adc.service.GetAuditFindingsSeverityBreakdown()
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AuditDashboardController) GetAuditFindingsListBySeverityLevel(c *gin.Context) {
	sortOrder := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")

	severityIDStr := c.Param("severityId")
	severityLevelID, err := strconv.ParseUint(severityIDStr, 10, 64)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid severityId parameter"))
		return
	}

	exportFromContext, exists := c.Get("export")
	var export bool
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(adc.db, c, "Invalid export parameter")
			return
		}
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"reference_number": "audit_findings.reference_number",
		"title":            "audit_findings.title",
		"created_at":       "audit_findings.created_at",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	queryParams := auditDTOs.AuditDashboardListQueryParamsDTO{
		Sort:     sortOrder,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Export:   export,
	}

	audits, err := adc.service.GetAuditFindingsBySeverityLevel(uint(severityLevelID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(adc.db, c, err.Error())
			return
		}
		excelFile, err := adc.service.ExportAuditFindingsBySeverityLevel(audits.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(adc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		fileName := adc.service.GetAuditFindingsBySeverityLevelExportFilename(uint(severityLevelID))
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="`+fileName+`"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}
	c.JSON(http.StatusOK, audits)
}

func (adc *AuditDashboardController) GetAuditFindingsListBySeverityLevelExport(c *gin.Context) {
	c.Set("export", true)
	adc.GetAuditFindingsListBySeverityLevel(c)
}

func (adc *AuditDashboardController) GetAuditFindingsTrend(c *gin.Context) {
	result, err := adc.service.GetAuditFindingsTrend()
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (adc *AuditDashboardController) GetOverdueAuditActionPlans(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(adc.db, c, "Invalid export parameter")
			return
		}
	}

	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"title":            "audit_finding_action_plans.title",
		"responsible_user": "responsible_user.first_name_en",
		"department":       "responsible_user_department_title_en",
		"severity_level":   "severity_levels.weight",
		"target_date":      "audit_finding_action_plans.target_date",
		"overdue_days":     "overdue_days",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := auditDTOs.AuditDashboardListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Export:   export,
	}

	audits, err := adc.service.GetOverdueAuditActionPlans(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(adc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(adc.db, c, err.Error())
			return
		}
		excelFile, err := adc.service.ExportOverdueAuditActionPlans(audits.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(adc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Overdue Action Plans.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}
	c.JSON(http.StatusOK, audits)
}

func (adc *AuditDashboardController) GetOverdueAuditActionPlansExport(c *gin.Context) {
	c.Set("export", true)
	adc.GetOverdueAuditActionPlans(c)
}
