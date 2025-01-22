package controllers

import (
	"errors"
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskServices "jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskReportController struct {
	service           services.RiskReportServiceInterface
	kriRegisterSvc    riskServices.KRIRegisterServiceInterface
	riskRegisterSvc   riskServices.RiskRegisterServiceInterface
	riskActionPlanSvc riskServices.RiskActionPlanServiceInterface

	db *gorm.DB
}

func NewRiskReportController(service services.RiskReportServiceInterface,
	kriRegisterSvc riskServices.KRIRegisterServiceInterface,
	riskRegisterSvc riskServices.RiskRegisterServiceInterface,
	riskActionPlanSvc riskServices.RiskActionPlanServiceInterface) *RiskReportController {
	return &RiskReportController{
		service:           service,
		kriRegisterSvc:    kriRegisterSvc,
		riskRegisterSvc:   riskRegisterSvc,
		riskActionPlanSvc: riskActionPlanSvc,
	}
}

func (rrc *RiskReportController) GetAllKRIRegisters(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid export parameter")
		return
	}
	// Define an array of allowed parameter names
	allowedParams := []string{
		"kri_category_ids",
		"risk_level_ids",
		"review_frequency_ids",
		"responsible_user_ids",
		"responsible_unit_ids",
	}

	validSortFields := map[string]string{
		"reference_number":    "kri_register.reference_number",
		"title":               "kri_register.title",
		"description":         "kri_register.description",
		"kri_category_en":     "kri_categories.title_en",
		"kri_category_ar":     "kri_categories.title_ar",
		"calculation_type":    "kri_register.calculation_type",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"low_threshold":       "kri_register.low_threshold",
		"medium_threshold":    "kri_register.medium_threshold",
		"high_threshold":      "kri_register.high_threshold",
		"current_value":       "kri_register.current_value",
		"unit_en":             "units.title_en",
		"unit_ar":             "units.title_ar",
		"risk_level_en":       "risk_levels.title_en",
		"risk_level_ar":       "risk_levels.title_ar",
		"review_frequency_en": "review_frequencies.title_en",
		"review_frequency_ar": "review_frequencies.title_ar",
		"review_at":           "kri_register.last_reviewed_date",
		"next_review_date":    "kri_register.next_review_date",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid page_size parameter")
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

	queryParams := riskDTOs.KRIRegisterListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
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
	kriRegisters, err := rrc.kriRegisterSvc.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(rrc.db, c, err.Error())
			return
		}

		excelFile, err := rrc.kriRegisterSvc.ExportKriRegister(kriRegisters.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(rrc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="KRI register.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, kriRegisters)
}

func (rrc *RiskReportController) GetKriRegisterExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	rrc.GetAllKRIRegisters(c)
}

func (rrc *RiskReportController) GetKRIByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetKRIByCategory(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetKRIByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetKRIByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetKRIByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetKRIByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetKRIByRiskLevel(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetKRIByRiskLevel(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetKRIByReviewFrequency(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetKRIByReviewFrequency(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetAllRiskRegisters(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	valueStr := c.Query("value")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"risk_category_ids",
		"risk_source_ids",
		"responsible_unit_ids",
		"responsible_user_ids",
		"risk_status_ids",
		"review_frequency_ids",
		"inherent_risk_rating_ids",
		"residual_risk_rating_ids",
		"current_likelihood_ids",
		"current_impact_ids",
		"current_risk_rating_ids",
		"risk_response_strategy_ids",
	}

	validSortFields := map[string]string{
		"reference_number":          "risk_register.reference_number",
		"title":                     "risk_register.title",
		"description":               "risk_register.description",
		"risk_category_en":          "risk_categories.title_en",
		"risk_category_ar":          "risk_categories.title_ar",
		"organization_unit_en":      "organization_units.title_en",
		"organization_unit_ar":      "organization_units.title_ar",
		"responsible_user_en":       "users.first_name_en",
		"responsible_user_ar":       "users.first_name_ar",
		"risk_source_en":            "risk_sources.title_en",
		"risk_source_ar":            "risk_sources.title_ar",
		"risk_status_en":            "risk_statuses.title_en",
		"risk_status_ar":            "risk_statuses.title_ar",
		"inherent_likelihood_en":    "inherent_likelihood.title_en",
		"inherent_likelihood_ar":    "inherent_likelihood.title_ar",
		"inherent_impact_en":        "inherent_impact.title_en",
		"inherent_impact_ar":        "inherent_impact.title_ar",
		"inherent_risk_rating_en":   "inherent_risk_rating.title_en",
		"inherent_risk_rating_ar":   "inherent_risk_rating.title_ar",
		"risk_response_strategy_en": "risk_response_strategies.title_en",
		"risk_response_strategy_ar": "risk_response_strategies.title_ar",
		"residual_likelihood_en":    "residual_likelihood.title_en",
		"residual_likelihood_ar":    "residual_likelihood.title_ar",
		"residual_impact_en":        "residual_impact.title_en",
		"residual_impact_ar":        "residual_impact.title_ar",
		"residual_risk_rating_en":   "residual_risk_rating.title_en",
		"residual_risk_rating_ar":   "residual_risk_rating.title_ar",
		"reported_at":               "risk_register.reported_at",
		"review_frequency_en":       "review_frequencies.title_en",
		"review_frequency_ar":       "review_frequencies.title_ar",
		"next_review_date":          "risk_register.next_review_date",
	}

	// Retrieve RiskAppetiteToleranceLevelID from query
	riskAppetiteToleranceLevelIDStr := c.Query("risk_appetite_tolerance_level_id")
	var riskAppetiteToleranceLevelID uint
	if riskAppetiteToleranceLevelIDStr != "" {
		parsedID, err := strconv.ParseUint(riskAppetiteToleranceLevelIDStr, 10, 32)
		if err != nil {
			utility.SendInternalErrResp(rrc.db, c, "Invalid RiskAppetiteToleranceLevelID parameter")
			return
		}
		riskAppetiteToleranceLevelID = uint(parsedID) // Explicit casting from uint64 to uint
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	reportedAtFrom := c.Query("reported_at")
	reportedAtTo := c.Query("reported_at")
	nextReviewDateFrom := c.Query("next_review_date_from")
	nextReviewDateTo := c.Query("next_review_date_to")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(rrc.db, c, "Invalid value parameter")
			return
		}
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := riskDTOs.RiskRegisterListQueryParamsDTO{
		Sort:                         sort,
		SortBy:                       dbField,
		Keyword:                      keyword,
		Status:                       status,
		Page:                         page,
		PageSize:                     pageSize,
		Value:                        value,
		Filters:                      filters,
		QuickFilter:                  quickFilter,
		ReportedAtFrom:               reportedAtFrom,
		ReportedAtTo:                 reportedAtTo,
		NextReviewDateFrom:           nextReviewDateFrom,
		NextReviewDateTo:             nextReviewDateTo,
		Export:                       export,
		DateFrom:                     dateFrom,
		DateTo:                       dateTo,
		RiskAppetiteToleranceLevelID: riskAppetiteToleranceLevelID,
	}
	risks, err := rrc.riskRegisterSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(rrc.db, c, err.Error())
			return
		}

		excelFile, err := rrc.riskRegisterSvc.ExportRiskRegister(risks.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(rrc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Risk register.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, risks)
}

func (rrc *RiskReportController) GetRiskRegisterExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	rrc.GetAllRiskRegisters(c)
}

func (rrc *RiskReportController) GetRiskByCategory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskByCategory(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetRiskByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetRiskByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetRiskByRiskSource(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskByRiskSource(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetRiskByInherentRiskRating(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskByInherentRiskRating(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetRiskByResidualRiskRating(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskByResidualRiskRating(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetRiskByRiskStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskByRiskStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetRiskByRiskResponseStrategy(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskByRiskResponseStrategy(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetRiskByRiskAppetiteToleranceLevel(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskByRiskAppetiteToleranceLevel(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetRiskByKRIEnabledStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskByKRIEnabledStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetAllRiskActionPlans(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"severity_level_ids",
		"responsible_user_ids",
		"action_plan_status_ids",
		"risk_register_ids",
	}

	validSortFields := map[string]string{
		"title":               "risk_action_plans.title",
		"risk_title":          "risk_register.title",
		"severity_level_en":   "severity_levels.title_en",
		"severity_level_ar":   "severity_levels.title_ar",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"target_date":         "risk_action_plans.target_date",
		"percentage":          "risk_action_plans.percentage",
		"status_en":           "action_plan_statuses.title_en",
		"status_ar":           "action_plan_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := riskDTOs.RiskActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         dbField,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
		Export:         export,
		DateFrom:       dateFrom,
		DateTo:         dateTo,
	}
	riskActionPlans, err := rrc.riskActionPlanSvc.GetAll(uint(riskID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(rrc.db, c, err.Error())
			return
		}

		excelFile, err := rrc.riskActionPlanSvc.ExportRiskActionPlan(riskActionPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(rrc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Risk action plan.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, riskActionPlans)
}

func (rrc *RiskReportController) GetRiskActionPlanExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	rrc.GetAllRiskActionPlans(c)
}

func (rrc *RiskReportController) GetRiskActionPlansByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskActionPlansByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetRiskActionPlansByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskActionPlansByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (rrc *RiskReportController) GetRiskActionPlansBySeverityLevel(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := rrc.service.GetRiskActionPlansBySeverityLevel(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}
