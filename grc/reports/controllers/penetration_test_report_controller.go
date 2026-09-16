package controllers

import (
	"errors"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	penetrationTestServices "jethur/grc/penetration-testing/services"
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestReportController struct {
	service                   services.PenetrationTestReportServiceInterface
	penetrationTestSvc        penetrationTestServices.PenetrationTestServiceInterface
	penetrationTestFindingSvc penetrationTestServices.PenetrationTestFindingServiceInterface

	db *gorm.DB
}

func NewPenetrationTestReportController(service services.PenetrationTestReportServiceInterface,
	penetrationTestSvc penetrationTestServices.PenetrationTestServiceInterface,
	penetrationTestFindingSvc penetrationTestServices.PenetrationTestFindingServiceInterface) *PenetrationTestReportController {
	return &PenetrationTestReportController{
		service:                   service,
		penetrationTestSvc:        penetrationTestSvc,
		penetrationTestFindingSvc: penetrationTestFindingSvc,
	}
}

func (ptrc *PenetrationTestReportController) GetAllPenetrationTests(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "system_application_name")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(ptrc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"penetration_test_status_ids",
		"cs_responsible_employee_ids",
		"responsible_unit_ids",
		"created_by_ids",
		"years",
		"periods",
		"penetration_test_type_ids",
		"penetration_test_plan_ids",
		"penetration_test_classification_ids",
		"penetration_test_executor_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	validSortFields := map[string]string{
		"id":                             "penetration_tests.id",
		"penetration_test_status_ar":     "penetration_test_statuses.title_ar",
		"responsible_user_first_name_ar": "cs_responsible_employee.first_name_ar",
		"responsible_user_last_name_en":  "cs_responsible_employee.last_name_en",
		"responsible_user_last_name_ar":  "cs_responsible_employee.last_name_ar",
		"responsible_user_email":         "cs_responsible_employee.email",
		"created_user_first_name_ar":     "created_user.first_name_ar",
		"created_user_last_name_en":      "created_user.last_name_en",
		"created_user_last_name_ar":      "created_user.last_name_ar",
		"penetration_test_status_en":     "penetration_test_statuses.title_en",
		"responsible_user_first_name_en": "cs_responsible_employee.first_name_en",
		"responsible_unit_title_ar":      "responsible_unit.title_ar",
		"responsible_unit_title_en":      "responsible_unit.title_en",
		"created_user_first_name_en":     "created_user.first_name_en",
		"reference_number":               "penetration_tests.reference_number",
		"system_application_name":        "penetration_tests.system_application_name",
		"year":                           "penetration_tests.year",
		"period":                         "penetration_tests.period",
		"created_at":                     "penetration_tests.created_at",
		"test_type_en":                   "penetration_test_types.title_en",
		"test_type_ar":                   "penetration_test_types.title_ar",
		"classification_en":              "penetration_test_classifications.title_en",
		"classification_ar":              "penetration_test_classifications.title_ar",
		"test_plan_en":                   "penetration_test_plans.title_en",
		"test_plan_ar":                   "penetration_test_plans.title_ar",
		"executor_en":                    "penetration_test_executors.title_en",
		"executor_ar":                    "penetration_test_executors.title_ar",
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := penetrationTestDTOs.PenetrationTestListQueryParamsDTO{
		Export:      export,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	penetrationTests, err := ptrc.penetrationTestSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(ptrc.db, c, err.Error())
			return
		}
		excelFile, err := ptrc.penetrationTestSvc.ExportPenetrationTest(penetrationTests.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(ptrc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Penetration Tests.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, penetrationTests)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestExport(c *gin.Context) {
	c.Set("export", true)

	ptrc.GetAllPenetrationTests(c)
}

func (ptrc *PenetrationTestReportController) GetAllPenetrationTestFindings(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(ptrc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"responsible_unit_ids",
		"penetration_test_ids",
		"action_plan_status_ids",
		"severity_level_ids",
		"owner_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	validSortFields := map[string]string{
		"id":                        "penetration_test_findings.id",
		"system_application_name":   "penetration_tests.system_application_name",
		"responsible_user_ar":       "owner_user.first_name_ar",
		"severity_level_ar":         "severity_levels.title_ar",
		"description":               "penetration_test_findings.description",
		"responsible_user_en":       "owner_user.first_name_en",
		"responsible_unit_title_ar": "responsible_unit.title_ar",
		"responsible_unit_title_en": "responsible_unit.title_en",
		"reference_number":          "penetration_test_findings.reference_number",
		"title":                     "penetration_test_findings.title",
		"severity_level_en":         "severity_levels.title_en",
		"status_en":                 "action_plan_statuses.title_en",
		"status_ar":                 "action_plan_statuses.title_ar",
		"percentage":                "penetration_test_findings.percentage",
		"target_remediation_date":   "penetration_test_findings.target_remediation_date",
		"created_at":                "penetration_test_findings.created_at",
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, "Invalid page parameter")
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	penetrationTestID, err := utility.ParseQueryParamToInt(c, "penetration_test_id", 0)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid penetration_test_id parameter"))
		return
	}

	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	queryParams := penetrationTestDTOs.PenetrationTestFindingListQueryParamsDTO{
		TargetDateFrom:    targetDateFrom,
		TargetDateTo:      targetDateTo,
		Export:            export,
		Sort:              sort,
		SortBy:            sortField,
		Keyword:           keyword,
		Page:              page,
		PageSize:          pageSize,
		PenetrationTestID: uint(penetrationTestID),
		Filters:           filters,
		QuickFilter:       quickFilter,
		DateFrom:          dateFrom,
		DateTo:            dateTo,
	}
	findings, err := ptrc.penetrationTestFindingSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(ptrc.db, c, err.Error())
			return
		}
		excelFile, err := ptrc.penetrationTestFindingSvc.ExportPenetrationTestFinding(findings.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(ptrc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Penetration Test Findings.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, findings)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestFindingExport(c *gin.Context) {
	c.Set("export", true)

	ptrc.GetAllPenetrationTestFindings(c)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestByType(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestByType(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestByClassification(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestByClassification(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestByPlan(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestByPlan(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestByExecutor(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestByExecutor(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestByCreatedUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestByCreatedUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestByYear(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestByYear(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestByPeriod(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestByPeriod(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestFindingByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestFindingByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestFindingByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestFindingByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestFindingByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestFindingByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestFindingBySeverityLevel(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := ptrc.service.GetPenetrationTestFindingBySeverityLevel(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}
