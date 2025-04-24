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

type AuditInsightsAuditController struct {
	service services.AuditInsightsAuditServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsAuditController(service services.AuditInsightsAuditServiceInterface) *AuditInsightsAuditController {
	return &AuditInsightsAuditController{
		service: service,
	}
}

func (aiac *AuditInsightsAuditController) GetAllAuditInsightsAudits(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(aiac.db, c, "Invalid export parameter")
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

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(aiac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aiac.db, c, "Invalid page_size parameter")
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
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aiac.db, c, err.Error())
		return
	}
	audits, err := aiac.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aiac.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(aiac.db, c, err.Error())
			return
		}

		excelFile, err := aiac.service.ExportAuditInsights(audits.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(aiac.db, c, err.Error())
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

func (aiac *AuditInsightsAuditController) GetAuditInsightsAuditDetails(c *gin.Context) {
	idStr := c.Param("auditId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	audit, err := aiac.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aiac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, audit)
}

func (aiac *AuditInsightsAuditController) GetAuditInsightsAuditsExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	aiac.GetAllAuditInsightsAudits(c)
}
