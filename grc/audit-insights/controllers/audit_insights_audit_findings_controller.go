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
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(aiafc.db, c, "Invalid export parameter")
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
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(aiafc.db, c, err.Error())
			return
		}

		excelFile, err := aiafc.service.ExportAuditInsightsAuditFindings(audits.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(aiafc.db, c, err.Error())
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

func (aiafc *AuditInsightsAuditFindingController) GetAuditInsightFindingCountByCategory(c *gin.Context) {
	auditIDStr := c.Param("auditId")

	auditId, err := strconv.ParseUint(auditIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	findingCounts, err := aiafc.service.GetAuditInsightFindingCountByCategory(uint(auditId))
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
	c.JSON(http.StatusOK, findingCounts)
}

func (aiafc *AuditInsightsAuditFindingController) GetAuditInsightsAuditFindingsExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	aiafc.GetAllAuditInsightsAuditFindings(c)
}
