package controllers

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceServices "jethur/grc/compliance/services"
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceReportController struct {
	service               services.ComplianceReportServiceInterface
	complianceRegisterSvc complianceServices.ComplianceRegisterServiceInterface
	db                    *gorm.DB
}

func NewComplianceReportController(service services.ComplianceReportServiceInterface,
	complianceRegisterSvc complianceServices.ComplianceRegisterServiceInterface) *ComplianceReportController {
	return &ComplianceReportController{
		service:               service,
		complianceRegisterSvc: complianceRegisterSvc,
	}
}

func (crc *ComplianceReportController) GetAllComplianceRegisters(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_category_ids",
		"severity_level_ids",
		"compliance_status_ids",
		"responsible_user_ids",
		"responsible_unit_ids",
		"review_frequency_ids",
		"review_responsible_user_ids",
		"lifecycle_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse deadline_from, deadline_to, issue_date_from, issue_date_to, expiry_date_from, expiry_date_to, next_review_date_from, next_review_date_to parameters
	deadlineFrom := c.Query("deadline_from")
	deadlineTo := c.Query("deadline_to")
	issueDateFrom := c.Query("issue_date_from")
	issueDateTo := c.Query("issue_date_to")
	expiryDateFrom := c.Query("expiry_date_from")
	expiryDateTo := c.Query("expiry_date_to")
	nextReviewDateFrom := c.Query("next_review_date_from")
	nextReviewDateTo := c.Query("next_review_date_to")

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page_size parameter")
		return
	}
	validSortFields := map[string]string{
		"reference_number":           "compliance_register.reference_number",
		"standard_code":              "compliance_register.standard_code",
		"title_en":                   "compliance_register.title_en",
		"title_ar":                   "compliance_register.title_ar",
		"description_en":             "compliance_register.description_en",
		"description_ar":             "compliance_register.description_ar",
		"category_title_en":          "compliance_categories.title_en",
		"category_title_ar":          "compliance_categories.title_ar",
		"severity_level_en":          "severity_levels.title_en",
		"severity_level_ar":          "severity_levels.title_ar",
		"responsible_unit_en":        "organization_units.title_en",
		"responsible_unit_ar":        "organization_units.title_ar",
		"responsible_user_en":        "responsible_user.first_name_en",
		"responsible_user_ar":        "responsible_user.first_name_ar",
		"review_frequency_title_en":  "review_frequencies.title_en",
		"review_frequency_title_ar":  "review_frequencies.title_ar",
		"review_responsible_user_en": "review_responsible_user.first_name_en",
		"review_responsible_user_ar": "review_responsible_user.first_name_ar",
		"status_en":                  "lifecycle_statuses.title_en",
		"status_ar":                  "lifecycle_statuses.title_ar",
		"compliance_status_en":       "compliance_statuses.title_en",
		"compliance_status_ar":       "compliance_statuses.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := complianceDTOs.ComplianceRegisterListQueryParamsDTO{
		Sort:               sort,
		SortBy:             sortField,
		Keyword:            keyword,
		Status:             status,
		Page:               page,
		PageSize:           pageSize,
		Filters:            filters,
		QuickFilter:        quickFilter,
		DeadlineFrom:       deadlineFrom,
		DeadlineTo:         deadlineTo,
		IssueDateFrom:      issueDateFrom,
		IssueDateTo:        issueDateTo,
		ExpiryDateFrom:     expiryDateFrom,
		ExpiryDateTo:       expiryDateTo,
		NextReviewDateFrom: nextReviewDateFrom,
		NextReviewDateTo:   nextReviewDateTo,
		DateFrom:           dateFrom,
		DateTo:             dateTo,
		Export:             export,
	}

	compliances, err := crc.complianceRegisterSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	// If export is true, generate and return the Excel file
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(crc.db, c, err.Error())
			return
		}

		excelFile, err := crc.complianceRegisterSvc.ExportComplianceRegisters(compliances.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(crc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Compliance Register.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	// Otherwise, return the data as JSON
	c.JSON(http.StatusOK, compliances)
}

func (crc *ComplianceReportController) ComplianceRegisterReportExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	crc.GetAllComplianceRegisters(c)
}

func (crc *ComplianceReportController) GetComplianceByComplianceCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	complianceReport, err := crc.service.GetComplianceByComplianceCategory(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceReport)
}

func (crc *ComplianceReportController) GetComplianceBySeverityLevel(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := crc.service.GetComplianceBySeverityLevel(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (crc *ComplianceReportController) GetComplianceByComplianceStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := crc.service.GetComplianceByComplianceStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (crc *ComplianceReportController) GetComplianceByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := crc.service.GetComplianceByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (crc *ComplianceReportController) GetComplianceByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	result, err := crc.service.GetComplianceByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (crc *ComplianceReportController) GetComplianceRequirementByComplianceStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	results, err := crc.service.GetComplianceRequirementByComplianceStatus(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, results)
}

func (crc *ComplianceReportController) GetComplianceRequirementBySeverityLevel(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	results, err := crc.service.GetComplianceRequirementBySeverity(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, results)
}

func (crc *ComplianceReportController) GetComplianceRequirementByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	results, err := crc.service.GetComplianceRequirementByCategory(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, results)
}

func (crc *ComplianceReportController) GetComplianceRequirementByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	results, err := crc.service.GetComplianceRequirementByResponsibleUser(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, results)
}
