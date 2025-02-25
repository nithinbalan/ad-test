package controllers

import (
	"errors"
	documentDTOs "jethur/grc/documents/dtos"
	documentServices "jethur/grc/documents/services"
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentReportController struct {
	service                  services.DocumentReportServiceInterface
	documentSvc              documentServices.DocumentServiceInterface
	documentChangeRequestSvc documentServices.DocumentChangeRequestServiceInterface
	db                       *gorm.DB
}

func NewDocumentReportController(
	service services.DocumentReportServiceInterface,
	documentSvc documentServices.DocumentServiceInterface,
	documentChangeRequestSvc documentServices.DocumentChangeRequestServiceInterface,
) *DocumentReportController {
	return &DocumentReportController{
		service:                  service,
		documentSvc:              documentSvc,
		documentChangeRequestSvc: documentChangeRequestSvc,
	}
}

func (drc *DocumentReportController) GetAllDocuments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	enableAcknowledgedStatus := c.Query("enable_acknowledge_status")
	documentTypes := utility.ParseCommaSeparatedStringToArray(c.Query("document_types"))
	documentStatus := utility.ParseCommaSeparatedStringToArray(c.Query("document_status"))

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"document_category_ids",
		"document_type_ids",
		"document_access_type_ids",
		"responsible_user_ids",
		"responsible_unit_ids",
		"document_status_ids",
	}

	validSortFields := map[string]string{
		"type_ar":             "document_types.title_ar",
		"type_en":             "document_types.title_en",
		"title_ar":            "documents.title_ar",
		"title_en":            "documents.title_en",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"category_en":         "document_categories.title_en",
		"category_ar":         "document_categories.title_ar",
		"status_en":           "document_statuses.title_ar",
		"status_ar":           "document_statuses.title_en",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := documentDTOs.DocumentListQueryParamsDTO{
		Sort:                     sort,
		SortBy:                   dbField,
		Keyword:                  keyword,
		Status:                   status,
		Page:                     page,
		PageSize:                 pageSize,
		Filters:                  filters,
		QuickFilter:              quickFilter,
		EnableAcknowledgedStatus: enableAcknowledgedStatus == "true",
		DocumentTypes:            documentTypes,
		DocumentStatus:           documentStatus,
		Export:                   export,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	documents, err := drc.documentSvc.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(drc.db, c, err.Error())
			return
		}

		excelFile, err := drc.documentSvc.ExportDocuments(documents.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(drc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Document.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, documents)
}

func (drc *DocumentReportController) GetDocumentExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	drc.GetAllDocuments(c)
}

func (drc *DocumentReportController) GetDocumentByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsByCategory(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentByType(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsByType(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentByAccess(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsByAccess(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentByPolicyAndProcedure(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsByPolicyAndProcedure(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentByPublishedAndArchived(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsByPublishedAndArchived(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentPublishedByYear(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsPublishedByYear(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentPublishedByMonth(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsPublishedByMonth(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetAllDocumentChangeRequests(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "reference_number")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"document_change_request_type_ids",
		"document_change_request_status_ids",
		"document_ids",
	}

	validSortFields := map[string]string{
		"document_change_reuqest_type_ar": "change_request_types.title_ar",
		"document_change_reuqest_type_en": "change_request_types.title_en",
		"reason_for_change":               "document_change_requests.reason_for_change",
		"reference_number":                "document_change_requests.reason_for_change",
		"status_en":                       "change_request_statuses.title_ar",
		"status_ar":                       "change_request_statuses.title_en",
		"document_ar":                     "documents.title_en",
		"document_en":                     "documents.title_en",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page parameter"))
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page_size parameter"))
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	params := documentDTOs.DocumentChangeRequestQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
	}

	documentChangeRequests, err := drc.documentChangeRequestSvc.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(drc.db, c, err.Error())
			return
		}

		excelFile, err := drc.documentChangeRequestSvc.ExportDocumentChangeRequest(documentChangeRequests.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(drc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="DocumentChangeRequests.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, documentChangeRequests)
}

func (drc *DocumentReportController) GetDocumentChangeRequestExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	drc.GetAllDocumentChangeRequests(c)
}

func (drc *DocumentReportController) GetDocumentChangeRequestByType(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentChangeRequestsByType(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentChangeRequestByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentChangeRequestsByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentsUpdatedByChangeType(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsUpdatedByChangeType(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (drc *DocumentReportController) GetDocumentsArchivedByQuarter(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := drc.service.GetDocumentsArchivedByQuarter(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
