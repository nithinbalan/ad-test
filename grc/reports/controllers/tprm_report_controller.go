package controllers

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	tprmDTOs "jethur/grc/tprm/dtos"
	tprmServices "jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TPRMReportController struct {
	db                      *gorm.DB
	service                 services.TPRMReportServiceInterface
	thirdPartyRegisterSvc   tprmServices.ThirdPartyRegisterServiceInterface
	thirdPartyAssessmentSvc tprmServices.ThirdPartyAssessmentServiceInterface
}

func NewTPRMReportController(service services.TPRMReportServiceInterface,
	thirdPartyRegisterSvc tprmServices.ThirdPartyRegisterServiceInterface,
	thirdPartyAssessmentSvc tprmServices.ThirdPartyAssessmentServiceInterface) *TPRMReportController {
	return &TPRMReportController{
		service:                 service,
		thirdPartyRegisterSvc:   thirdPartyRegisterSvc,
		thirdPartyAssessmentSvc: thirdPartyAssessmentSvc,
	}
}

func (tprc *TPRMReportController) GetAllThirdPartyRegister(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(tprc.db, c, "Invalid export parameter")
			return
		}
	}
	allowedParams := []string{
		"third_party_type_ids",
		"country_ids",
		"organization_unit_ids",
		"severity_level_ids",
		"responsible_user_ids",
		"latest_status_ids",
		"latest_risk_level_ids",
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

	validSortFields := map[string]string{
		"title":                          "third_party_register.title",
		"reference_number":               "third_party_register.reference_number",
		"website":                        "third_party_register.website",
		"email":                          "third_party_register.email",
		"phone_number":                   "third_party_register.phone_number",
		"onboarding_date":                "third_party_register.onboarding_date",
		"service_start_date":             "third_party_register.service_start_date",
		"service_end_date":               "third_party_register.service_end_date",
		"country_title_en":               "countries.title_en",
		"country_title_ar":               "countries.title_ar",
		"third_party_type_en":            "third_party_types.title_en",
		"third_party_type_ar":            "third_party_types.title_ar",
		"organization_unit_title_en":     "organization_units.title_en",
		"organization_unit_title_ar":     "organization_units.title_ar",
		"severity_level_title_en":        "severity_levels.title_en",
		"severity_level_title_ar":        "severity_levels.title_ar",
		"severity_level_color_code":      "severity_levels.color_code",
		"severity_level_type":            "severity_levels.type",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_last_name_en":  "responsible_user.last_name_en",
		"created_user_first_name_en":     "created_user.first_name_en",
		"created_user_last_name_en":      "created_user.last_name_en",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid sort_by parameter"))
		return
	}

	query := tprmDTOs.ThirdPartyRegisterListQueryParamsDTO{
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

	thirdParties, err := tprc.thirdPartyRegisterSvc.GetAll(&query)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(tprc.db, c, err.Error())
			return
		}

		excelFile, err := tprc.thirdPartyRegisterSvc.ExportThirdPartyRegister(thirdParties.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(tprc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Third Party Register.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, thirdParties)
}

func (tprc *TPRMReportController) GetThirdPartyRegisterExport(c *gin.Context) {
	c.Set("export", true)

	tprc.GetAllThirdPartyRegister(c)
}

func (tprc *TPRMReportController) GetThirdPartyByType(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetThirdPartyByType(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tprc *TPRMReportController) GetThirdPartyByCountry(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetThirdPartyByCountry(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tprc *TPRMReportController) GetThirdPartyByDepartment(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetThirdPartyByDepartment(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tprc *TPRMReportController) GetThirdPartyByCriticality(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetThirdPartyByCriticality(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tprc *TPRMReportController) GetThirdPartyByLatestAssessmentStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetThirdPartyByLatestAssessmentStatus(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tprc *TPRMReportController) GetThirdPartyByRiskRating(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetThirdPartyByRiskRating(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tprc *TPRMReportController) GetThirdPartyServiceExpiry(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetThirdPartyServiceExpiry(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tprc *TPRMReportController) GetAllThirdPartyAssessment(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(tprc.db, c, "Invalid export parameter")
			return
		}
	}
	allowedParams := []string{
		"third_party_register_ids", 
		"question_set_ids",
		"third_party_assessment_status_ids",
		"risk_level_ids",
		"third_party_type_ids",
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

	validSortFields := map[string]string{
		"third_party_register_title": "third_party_register.title",
		"title":                      "third_party_assessments.title",
		"reference_number":           "third_party_assessments.reference_number",
		"due_date":                   "third_party_assessments.due_date",
		"question_set_title_en":      "question_sets.title_en",
		"question_set_title_ar":      "question_sets.title_ar",
		"created_user_first_name_en": "created_user.first_name_en",
		"created_user_last_name_en":  "created_user.last_name_en",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid sort_by parameter"))
		return
	}

	query := tprmDTOs.ThirdPartyAssessmentListQueryParamsDTO{
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

	thirdParties, err := tprc.thirdPartyAssessmentSvc.GetAll(&query)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(tprc.db, c, err.Error())
			return
		}

		excelFile, err := tprc.thirdPartyAssessmentSvc.ExportThirdPartyAssessments(thirdParties.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(tprc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="ThirdPartyAssessments.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, thirdParties)
}

func (tprc *TPRMReportController) GetThirdPartyAssessmentExport(c *gin.Context) {
	c.Set("export", true)

	tprc.GetAllThirdPartyAssessment(c)
}

func (tprc *TPRMReportController) GetAssessmentsByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetAssessmentsByStatus(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tprc *TPRMReportController) GetAssessmentsByRiskRating(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetAssessmentsByRiskRating(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tprc *TPRMReportController) GetAssessmentsByThirdPartyCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetAssessmentsByThirdPartyCategory(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tprc *TPRMReportController) GetAssessmentsByTemplate(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetAssessmentsByTemplate(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tprc *TPRMReportController) GetOverdueAssessments(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	params := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	result, err := tprc.service.GetOverdueAssessments(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
