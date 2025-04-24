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

type AuditInsightsAuditPlanController struct {
	service services.AuditInsightsAuditPlanServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsAuditPlanController(service services.AuditInsightsAuditPlanServiceInterface) *AuditInsightsAuditPlanController {
	return &AuditInsightsAuditPlanController{
		service: service,
	}
}

func (aiapc *AuditInsightsAuditPlanController) GetAllAuditInsightsAuditPlans(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(aiapc.db, c, "Invalid export parameter")
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
	targetDate := c.Query("target_date")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(aiapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(aiapc.db, c, "Invalid page_size parameter")
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
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aiapc.db, c, err.Error())
		return
	}
	auditPlans, err := aiapc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aiapc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(aiapc.db, c, err.Error())
			return
		}

		excelFile, err := aiapc.service.ExportAuditInsightsAuditPlans(auditPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(aiapc.db, c, err.Error())
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

func (aiapc *AuditInsightsAuditPlanController) GetAuditInsightsAuditPlanDetails(c *gin.Context) {
	idStr := c.Param("auditPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditPlan, err := aiapc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aiapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditPlan)
}

func (aiapc *AuditInsightsAuditPlanController) GetAuditInsightsAuditPlansExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	aiapc.GetAllAuditInsightsAuditPlans(c)
}
