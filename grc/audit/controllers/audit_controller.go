package controllers

import (
	"errors"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditController struct {
	service services.AuditServiceInterface
	db      *gorm.DB
}

func NewAuditController(service services.AuditServiceInterface) *AuditController {
	return &AuditController{
		service: service,
	}
}

func (apc *AuditController) GetAllAudits(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.DefaultQuery("sort_by", "audit_plan_title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, "Invalid export parameter")
		return
	}
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
		utility.SendInternalErrResp(apc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, "Invalid page_size parameter")
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
		Export:      export,
	}

	audits, err := apc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}
	// If export is true, generate and return the Excel file
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(apc.db, c, err.Error())
			return
		}

		excelFile, err := apc.service.ExportAudits(audits.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(apc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Audits.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, audits)
}

func (apc *AuditController) GetAuditsExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	apc.GetAllAudits(c)
}

func (apc *AuditController) GetAuditDetails(c *gin.Context) {
	idStr := c.Param("auditId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	audit, err := apc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(apc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, audit)
}

func (apc *AuditController) CreateAudit(c *gin.Context) {
	var auditFormDTO auditDTOs.AuditFormDTO

	if err := c.ShouldBindJSON(&auditFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	id, err := apc.service.Create(currentUserID, &auditFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(apc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (apc *AuditController) UpdateAudit(c *gin.Context) {
	var auditFormDTO auditDTOs.AuditFormDTO

	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	err = apc.service.Update(uint(auditID), currentUserID, &auditFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (apc *AuditController) DeleteAudits(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	failedIDs, err := apc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (apc *AuditController) GetAllAuditableItems(c *gin.Context) {
	auditIDStr := c.Param("auditId")

	auditID, err := strconv.ParseUint(auditIDStr, 10, 64)
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
		"audit_category_ids",
		"audit_team_ids",
		"team_lead_ids",
		"audit_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := auditDTOs.AuditAuditableItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditPlans, err := apc.service.GetAllAuditableItems(uint(auditID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditPlans)
}
