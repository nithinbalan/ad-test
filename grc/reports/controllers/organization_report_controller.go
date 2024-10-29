package controllers

import (
	"errors"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationReportController struct {
	service    services.OrganizationReportServiceInterface
	userSvc    organizationServices.UserServiceInterface
	processSvc organizationServices.ProcessServiceInterface
	controlSvc organizationServices.ControlServiceInterface
	db         *gorm.DB
}

func NewOrganizationReportController(service services.OrganizationReportServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	processSvc organizationServices.ProcessServiceInterface,
	controlSvc organizationServices.ControlServiceInterface) *OrganizationReportController {
	return &OrganizationReportController{
		service:    service,
		userSvc:    userSvc,
		processSvc: processSvc,
		controlSvc: controlSvc,
	}
}

func (orc *OrganizationReportController) GetAllUsers(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "first_name_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	userIdsStr := c.Query("user_ids")
	roleIdsStr := c.Query("role_ids")

	userIds, err := utility.ParseCommaSeparatedStringToIntArray(userIdsStr)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid 'user_ids' parameter"))
		return
	}

	roleIds, err := utility.ParseCommaSeparatedStringToIntArray(roleIdsStr)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid 'role_ids' parameter"))
		return
	}

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"designation_ids",
		"organization_unit_1_ids",
		"organization_unit_2_ids",
		"organization_unit_3_ids",
		"organization_unit_4_ids",
		"organization_unit_5_ids",
		"user_status_ids",
	}

	validSortFields := map[string]string{
		"first_name_en":  "users.first_name_en",
		"first_name_ar":  "users.first_name_ar",
		"email":          "users.email",
		"mobile_number":  "users.mobile_number",
		"designation_en": "designations.title_en",
		"designation_ar": "designations.title_ar",
		"role_en":        "roles_agg.roles_en",
		"role_ar":        "roles_agg.roles_ar",
		"user_status_en": "user_statuses.title_en",
		"user_status_ar": "user_statuses.title_ar",
		"name_en":        "CONCAT(users.first_name_en, ' ', users.last_name_en)",
		"name_ar":        "CONCAT(users.first_name_ar, ' ', users.last_name_ar)",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := organizationDTOs.UserListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		UserIDs:     userIds,
		RoleIDs:     roleIds,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	users, err := orc.userSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	if export {

		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(orc.db, c, err.Error())
			return
		}

		excelFile, err := orc.userSvc.ExportUsers(users.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(orc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Users.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, users)
}

func (orc *OrganizationReportController) GetUsersExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	orc.GetAllUsers(c)
}

func (orc *OrganizationReportController) GetUsersByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	usersByStatusReport, err := orc.service.GetUsersByStatus(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, usersByStatusReport)
}

func (orc *OrganizationReportController) GetUsersByRole(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	usersByRoleReport, err := orc.service.GetUsersByRole(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, usersByRoleReport)
}

func (orc *OrganizationReportController) GetUsersByDesignation(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	usersByDesignationReport, err := orc.service.GetUsersByDesignation(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, usersByDesignationReport)
}

func (orc *OrganizationReportController) GetAllProcesses(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"process_group_ids",
		"owner_unit_ids",
		"owner_ids",
		"lifecycle_status_ids",
		"risk_matrix_risk_rating_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page_size parameter")
		return
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	validSortFields := map[string]string{
		"reference_number":     "processes.reference_number",
		"title_en":             "processes.title_en",
		"title_ar":             "processes.title_ar",
		"description":          "processes.description",
		"process_group_en":     "process_groups.title_en",
		"process_group_ar":     "process_groups.title_ar",
		"organization_unit_en": "organization_units.title_en",
		"organization_unit_ar": "organization_units.title_ar",
		"owner_en":             "owner.first_name_en",
		"owner_ar":             "owner.first_name_ar",
		"status_en":            "lifecycle_statuses.title_en",
		"status_ar":            "lifecycle_statuses.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationDTOs.ProcessListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	processes, err := orc.processSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	if export {

		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(orc.db, c, err.Error())
			return
		}

		excelFile, err := orc.processSvc.ExportProcesses(processes.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(orc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Processes.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, processes)
}

func (orc *OrganizationReportController) GetProcessesExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	orc.GetAllProcesses(c)
}

func (orc *OrganizationReportController) GetProcessesByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	processByResponsibleUnitReport, err := orc.service.GetProcessesByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, processByResponsibleUnitReport)
}

func (orc *OrganizationReportController) GetProcessesByRiskRating(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	processesByRiskRating, err := orc.service.GetProcessesByRiskRating(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, processesByRiskRating)
}

func (orc *OrganizationReportController) GetProcessesByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	processesByResponsibleUser, err := orc.service.GetProcessesByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, processesByResponsibleUser)
}

func (orc *OrganizationReportController) GetProcessesByProcessGroup(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	processesByProcessGroup, err := orc.service.GetProcessesByProcessGroup(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, processesByProcessGroup)
}

func (orc *OrganizationReportController) GetAllControls(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	complianceRegisterIDsStr := c.Query("compliance_register_ids")

	complianceRegisterIDs, err := utility.ParseCommaSeparatedStringToIntArray(complianceRegisterIDsStr)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid 'compliance_register_ids' parameter"))
		return
	}

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"control_type_ids",
		"control_category_ids",
		"owner_unit_ids",
		"owner_ids",
		"control_implementation_status_ids",
		"control_efficiency_status_ids",
		"test_responsible_user_ids",
		"lifecycle_status_ids",
		"risk_matrix_risk_rating_ids",
	}

	validSortFields := map[string]string{
		"reference_number":         "controls.reference_number",
		"code":                     "controls.code",
		"title_en":                 "controls.title_en",
		"title_ar":                 "controls.title_ar",
		"description":              "controls.description",
		"objective":                "controls.objective",
		"control_category_en":      "control_categories.title_en",
		"control_category_ar":      "control_categories.title_ar",
		"control_type_en":          "control_types.title_en",
		"control_type_ar":          "control_types.title_ar",
		"owner_en":                 "owner.first_name_en",
		"owner_ar":                 "owner.first_name_ar",
		"review_frequency_en":      "review_frequencies.title_en",
		"review_frequency_ar":      "review_frequencies.title_ar",
		"implementation_status_en": "control_implementation_statuses.title_en",
		"implementation_status_ar": "control_implementation_statuses.title_ar",
		"implementation_date":      "controls.implemented_at",
		"efficiency_status_en":     "control_efficiency_statuses.title_en",
		"efficiency_status_ar":     "control_efficiency_statuses.title_ar",
		"last_test_date":           "controls.last_test_date",
		"next_test_date":           "controls.next_test_date",
		"test_frequency_en":        "test_frequencies.title_en",
		"test_frequency_ar":        "test_frequencies.title_ar",
		"last_review_note":         "controls.last_review_note",
		"next_review_date":         "controls.next_review_date",
		"status_en":                "lifecycle_statuses.title_en",
		"status_ar":                "lifecycle_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	//Parse last_review_date_from, last_review_date_to, next_review_date_from, next_review_date_to, last_test_date_from, last_test_date_to, next_test_date_from, next_test_date_to, implemented_at_from, implemented_at_to parameters
	lastReviewDateFrom := c.Query("last_review_date_from")
	lastReviewDateTo := c.Query("last_review_date_to")
	nextReviewDateFrom := c.Query("next_review_date_from")
	nextReviewDateTo := c.Query("next_review_date_to")
	lastTestDateFrom := c.Query("last_test_date_from")
	lastTestDateTo := c.Query("last_test_date_to")
	nextTestDateFrom := c.Query("next_test_date_from")
	nextTestDateTo := c.Query("next_test_date_to")
	implementedAtFrom := c.Query("implemented_at_from")
	implementedAtTo := c.Query("implemented_at_to")

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationDTOs.ControlListQueryParamsDTO{
		Sort:                  sort,
		SortBy:                dbField,
		Keyword:               keyword,
		Status:                status,
		Page:                  page,
		PageSize:              pageSize,
		Filters:               filters,
		QuickFilter:           quickFilter,
		LastReviewDateFrom:    lastReviewDateFrom,
		LastReviewDateTo:      lastReviewDateTo,
		NextReviewDateFrom:    nextReviewDateFrom,
		NextReviewDateTo:      nextReviewDateTo,
		LastTestDateFrom:      lastTestDateFrom,
		LastTestDateTo:        lastTestDateTo,
		NextTestDateFrom:      nextTestDateFrom,
		NextTestDateTo:        nextTestDateTo,
		ImplementedAtFrom:     implementedAtFrom,
		ImplementedAtTo:       implementedAtTo,
		Export:                export,
		DateFrom:              dateFrom,
		DateTo:                dateTo,
		ComplianceRegisterIDs: complianceRegisterIDs,
	}

	controls, err := orc.controlSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(orc.db, c, err.Error())
			return
		}

		excelFile, err := orc.controlSvc.ExportControls(controls.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(orc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Controls.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, controls)
}

func (orc *OrganizationReportController) GetControlsExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	orc.GetAllControls(c)
}

func (orc *OrganizationReportController) GetControlsByControlType(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	controlsByControlTypeReport, err := orc.service.GetControlsByControlType(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, controlsByControlTypeReport)
}

func (orc *OrganizationReportController) GetControlsByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := orc.service.GetControlsByCategory(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (orc *OrganizationReportController) GetControlsByComplianceRegister(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := orc.service.GetControlsByComplianceRegister(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (orc *OrganizationReportController) GetControlsByRiskRating(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := orc.service.GetControlsByRiskRating(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (orc *OrganizationReportController) GetControlsByImplementationStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := orc.service.GetControlsByImplementationStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (orc *OrganizationReportController) GetControlsByEfficiencyStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := orc.service.GetControlsByEfficiencyStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

// controllers/OrganizationReportController.go
func (orc *OrganizationReportController) GetControlsByApplicability(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	controlsByApplicabilityReport, err := orc.service.GetControlsByApplicability(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, controlsByApplicabilityReport)
}
