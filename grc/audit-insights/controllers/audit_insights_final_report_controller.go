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
	sortBy := c.DefaultQuery("sort_by", "audit_plan_title")
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
		"audit_final_report_status_ids",
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
		utility.SendInternalErrResp(aidrc.db, c, err.Error())
		return
	}
	auditFinalReports, err := aidrc.service.GetAll(&queryParams, currentUserID)
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

		excelFile, err := aidrc.service.ExportAuditInsightsFinalReports(auditFinalReports.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(aidrc.db, c, err.Error())
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

func (aidrc *AuditInsightsFinalReportController) GetAuditInsightsFinalReportExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	aidrc.GetAllAuditInsightsFinalReports(c)
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
