package controllers

import (
	"errors"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditFindingController struct {
	service services.AuditFindingServiceInterface
	db      *gorm.DB
}

func NewAuditFindingController(service services.AuditFindingServiceInterface) *AuditFindingController {
	return &AuditFindingController{
		service: service,
	}
}

func (afc *AuditFindingController) GetAllAuditFindings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(afc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
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
		utility.SendInternalErrResp(afc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afc.db, c, "Invalid page_size parameter")
		return
	}
	validSortFields := map[string]string{
		"reference_number":                "audit_findings.reference_number",
		"title":                           "audit_findings.title",
		"description":                     "audit_findings.description",
		"recommendations":                 "audit_findings.recommendations",
		"audit_id":                        "audits.id",
		"audit_title":                     "audit_plans.title",
		"auditable_item_title_en":         "auditable_items.title_en",
		"auditable_item_title_ar":         "auditable_items.title_ar",
		"audit_finding_category_title_en": "audit_finding_categories.title_en",
		"audit_finding_category_title_ar": "audit_finding_categories.title_ar",
		"severity_level_title_en":         "severity_levels.title_en",
		"severity_level_title_ar":         "severity_levels.title_ar",
		"audit_finding_status_title_en":   "audit_finding_statuses.title_en",
		"audit_finding_status_title_ar":   "audit_finding_statuses.title_ar",
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

	auditFindings, err := afc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(afc.db, c, err.Error())
		return
	}
	// If export is true, generate and return the Excel file
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(afc.db, c, err.Error())
			return
		}

		excelFile, err := afc.service.ExportAuditFindings(auditFindings.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(afc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Audit Findings.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, auditFindings)
}

func (afc *AuditFindingController) GetAuditFindingExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	afc.GetAllAuditFindings(c)
}

func (afc *AuditFindingController) GetAuditFindingDetails(c *gin.Context) {

	idStr := c.Param("findingId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditFinding, err := afc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditFinding)
}

func (afc *AuditFindingController) CreateAuditFinding(c *gin.Context) {
	var auditFindingFormDTO auditDTOs.AuditFindingFormDTO

	if err := c.ShouldBindJSON(&auditFindingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afc.db, c, err.Error())
		return
	}

	id, err := afc.service.Create(currentUserID, &auditFindingFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (afc *AuditFindingController) UpdateAuditFinding(c *gin.Context) {
	var auditFindingFormDTO auditDTOs.AuditFindingFormDTO

	id, err := strconv.Atoi(c.Param("findingId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditFindingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afc.db, c, err.Error())
		return
	}

	err = afc.service.Update(uint(id), currentUserID, &auditFindingFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (afc *AuditFindingController) DeleteAuditFindings(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afc.db, c, err.Error())
		return
	}

	failedIDs, err := afc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(afc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (afc *AuditFindingController) GetAuditFindingCountByCategory(c *gin.Context) {
	auditIDStr := c.Param("auditId")

	auditId, err := strconv.ParseUint(auditIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	findingCounts, err := afc.service.GetAuditFindingCountByCategory(uint(auditId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, findingCounts)
}

func (afc *AuditFindingController) GetAuditFindingDocumentPreview(c *gin.Context) {

	idStr := c.Param("findingId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditFindingDocumentPath, err := afc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(afc.db, c, err.Error())
		}
		return
	}
	c.File(auditFindingDocumentPath)
}

func (afc *AuditFindingController) DownloadAuditFindingDocument(c *gin.Context) {

	idStr := c.Param("findingId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditFindingDocumentPath, err := afc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(afc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(auditFindingDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(auditFindingDocumentPath)
}
