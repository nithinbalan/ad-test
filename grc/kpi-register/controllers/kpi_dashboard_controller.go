package controllers

import (
	"errors"
	"jethur/grc/kpi-register/dtos"
	"jethur/grc/kpi-register/services"
	"jethur/utility"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPIDashboardController struct {
	service services.KPIDashboardServiceInterface
	db      *gorm.DB
}

func NewKPIDashboardController(
	service services.KPIDashboardServiceInterface,
) *KPIDashboardController {
	return &KPIDashboardController{
		service: service,
	}
}

func (kdc *KPIDashboardController) GetKPISummary(c *gin.Context) {
	var domainID *uint

	if domainParam := c.Query("domain_id"); domainParam != "" {
		id, err := strconv.ParseUint(domainParam, 10, 64)
		if err != nil {
			utility.SendInvalidErrResp(c, utility.NewValidationError(
				err,
				"invalid",
				"domain_id",
			))
			return
		}
		temp := uint(id)
		domainID = &temp
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"responsible_user_ids",
		"kpi_status_ids",
		"kpi_category_ids",
	}

	filters := utility.ParseQueryFilters(
		allowedParams,
		c.Request.URL.Query(),
	)

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		DomainID: domainID,
		Filters:  filters,
		IsKRI:    isKRI,
	}

	result, err := kdc.service.GetKPISummary(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetKPIStatusSummary(c *gin.Context) {
	var domainID *uint

	if domainParam := c.Query("domain_id"); domainParam != "" {
		id, err := strconv.ParseUint(domainParam, 10, 64)
		if err != nil {
			utility.SendInvalidErrResp(c, utility.NewValidationError(
				err,
				"invalid",
				"domain_id",
			))
			return
		}
		temp := uint(id)
		domainID = &temp
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"responsible_user_ids",
		"kpi_status_ids",
		"kpi_category_ids",
	}

	filters := utility.ParseQueryFilters(
		allowedParams,
		c.Request.URL.Query(),
	)

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		DomainID: domainID,
		Filters:  filters,
		IsKRI:    isKRI,
	}

	result, err := kdc.service.GetKPIStatusSummary(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetKPIPassFailSummary(c *gin.Context) {
	var domainID *uint

	if domainParam := c.Query("domain_id"); domainParam != "" {
		id, err := strconv.ParseUint(domainParam, 10, 64)
		if err != nil {
			utility.SendInvalidErrResp(c, utility.NewValidationError(
				err,
				"invalid",
				"domain_id",
			))
			return
		}
		temp := uint(id)
		domainID = &temp
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"responsible_user_ids",
		"kpi_status_ids",
		"kpi_category_ids",
	}

	filters := utility.ParseQueryFilters(
		allowedParams,
		c.Request.URL.Query(),
	)

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		DomainID: domainID,
		Filters:  filters,
		IsKRI:    isKRI,
	}

	result, err := kdc.service.GetKPIPassFailSummary(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetTopPerformingKPIs(c *gin.Context) {

	startDateParam := c.Query("start_date")
	endDateParam := c.Query("end_date")
	quickFilter := c.QueryArray("quick_filter")
	periodType := c.Query("period_type")

	var startDate *time.Time
	if startDateParam != "" {
		parsed, err := time.Parse("2006-01-02", startDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if endDateParam != "" {
		parsed, err := time.Parse("2006-01-02", endDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		endDate = &parsed
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"kpi_category_ids",
		"kpi_status_ids",
		"kpi_review_status_ids",
		"review_frequency_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		StartDate:   startDate,
		EndDate:     endDate,
		Filters:     filters,
		QuickFilter: quickFilter,
		PeriodType:  &periodType,
		IsKRI:       isKRI,
	}

	result, err := kdc.service.GetTopPerformingKPIs(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetBottomPerformingKPIs(c *gin.Context) {

	startDateParam := c.Query("start_date")
	endDateParam := c.Query("end_date")
	quickFilter := c.QueryArray("quick_filter")
	periodType := c.Query("period_type")

	var startDate *time.Time
	if startDateParam != "" {
		parsed, err := time.Parse("2006-01-02", startDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if endDateParam != "" {
		parsed, err := time.Parse("2006-01-02", endDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		endDate = &parsed
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"kpi_category_ids",
		"kpi_status_ids",
		"kpi_review_status_ids",
		"review_frequency_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		StartDate:   startDate,
		EndDate:     endDate,
		Filters:     filters,
		QuickFilter: quickFilter,
		PeriodType:  &periodType,
		IsKRI:       isKRI,
	}

	result, err := kdc.service.GetBottomPerformingKPIs(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetCategorywiseKPIStats(c *gin.Context) {
	var domainID *uint

	if domainParam := c.Query("domain_id"); domainParam != "" {
		id, err := strconv.ParseUint(domainParam, 10, 64)
		if err != nil {
			utility.SendInvalidErrResp(c, utility.NewValidationError(
				err,
				"invalid",
				"domain_id",
			))
			return
		}
		temp := uint(id)
		domainID = &temp
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"responsible_user_ids",
		"kpi_status_ids",
		"kpi_category_ids",
	}

	filters := utility.ParseQueryFilters(
		allowedParams,
		c.Request.URL.Query(),
	)

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		DomainID: domainID,
		Filters:  filters,
		IsKRI:    isKRI,
	}

	result, err := kdc.service.GetCategorywiseKPIStats(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetDepartmentwiseKPIStats(c *gin.Context) {
	var domainID *uint

	if domainParam := c.Query("domain_id"); domainParam != "" {
		id, err := strconv.ParseUint(domainParam, 10, 64)
		if err != nil {
			utility.SendInvalidErrResp(c, utility.NewValidationError(
				err,
				"invalid",
				"domain_id",
			))
			return
		}
		temp := uint(id)
		domainID = &temp
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"responsible_user_ids",
		"kpi_status_ids",
		"kpi_category_ids",
	}

	filters := utility.ParseQueryFilters(
		allowedParams,
		c.Request.URL.Query(),
	)

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		DomainID: domainID,
		Filters:  filters,
		IsKRI:    isKRI,
	}

	result, err := kdc.service.GetDepartmentwiseKPIStats(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetOverdueForAchievementDetails(c *gin.Context) {
	var domainID *uint

	if domainParam := c.Query("domain_id"); domainParam != "" {
		id, err := strconv.ParseUint(domainParam, 10, 64)
		if err != nil {
			utility.SendInvalidErrResp(c, utility.NewValidationError(
				err,
				"invalid",
				"domain_id",
			))
			return
		}
		temp := uint(id)
		domainID = &temp
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"responsible_user_ids",
		"kpi_status_ids",
		"kpi_category_ids",
	}

	filters := utility.ParseQueryFilters(
		allowedParams,
		c.Request.URL.Query(),
	)

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		DomainID: domainID,
		Filters:  filters,
		IsKRI:    isKRI,
	}

	result, err := kdc.service.GetOverdueForAchievementDetails(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetOverdueReviewDetails(c *gin.Context) {
	var domainID *uint

	if domainParam := c.Query("domain_id"); domainParam != "" {
		id, err := strconv.ParseUint(domainParam, 10, 64)
		if err != nil {
			utility.SendInvalidErrResp(c, utility.NewValidationError(
				err,
				"invalid",
				"domain_id",
			))
			return
		}
		temp := uint(id)
		domainID = &temp
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"responsible_user_ids",
		"kpi_status_ids",
		"kpi_category_ids",
	}

	filters := utility.ParseQueryFilters(
		allowedParams,
		c.Request.URL.Query(),
	)

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		DomainID: domainID,
		Filters:  filters,
		IsKRI:    isKRI,
	}

	result, err := kdc.service.GetOverdueForReviewDetails(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetKPIStatusOverview(c *gin.Context) {

	startDateParam := c.Query("start_date")
	endDateParam := c.Query("end_date")
	quickFilter := c.QueryArray("quick_filter")
	periodType := c.Query("period_type")

	var startDate *time.Time
	if startDateParam != "" {
		parsed, err := time.Parse("2006-01-02", startDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if endDateParam != "" {
		parsed, err := time.Parse("2006-01-02", endDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		endDate = &parsed
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"kpi_category_ids",
		"kpi_status_ids",
		"kpi_review_status_ids",
		"review_frequency_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		StartDate:   startDate,
		EndDate:     endDate,
		Filters:     filters,
		QuickFilter: quickFilter,
		PeriodType:  &periodType,
		IsKRI:       isKRI,
	}

	result, err := kdc.service.GetKPIStatusOverview(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetAchievementSummary(c *gin.Context) {

	startDateParam := c.Query("start_date")
	endDateParam := c.Query("end_date")
	quickFilter := c.QueryArray("quick_filter")
	periodType := c.Query("period_type")

	var startDate *time.Time
	if startDateParam != "" {
		parsed, err := time.Parse("2006-01-02", startDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if endDateParam != "" {
		parsed, err := time.Parse("2006-01-02", endDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		endDate = &parsed
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"kpi_category_ids",
		"kpi_status_ids",
		"kpi_review_status_ids",
		"review_frequency_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		StartDate:   startDate,
		EndDate:     endDate,
		Filters:     filters,
		QuickFilter: quickFilter,
		PeriodType:  &periodType,
		IsKRI:       isKRI,
	}

	result, err := kdc.service.GetAchievementSummary(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetKPIByReviewStatus(c *gin.Context) {

	startDateParam := c.Query("start_date")
	endDateParam := c.Query("end_date")
	quickFilter := c.QueryArray("quick_filter")
	periodType := c.Query("period_type")

	var startDate *time.Time
	if startDateParam != "" {
		parsed, err := time.Parse("2006-01-02", startDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if endDateParam != "" {
		parsed, err := time.Parse("2006-01-02", endDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		endDate = &parsed
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"kpi_category_ids",
		"kpi_status_ids",
		"kpi_review_status_ids",
		"review_frequency_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		StartDate:   startDate,
		EndDate:     endDate,
		Filters:     filters,
		QuickFilter: quickFilter,
		PeriodType:  &periodType,
		IsKRI:       isKRI,
	}

	result, err := kdc.service.GetKPIByReviewStatus(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetKPIAchievementByDomain(c *gin.Context) {

	startDateParam := c.Query("start_date")
	endDateParam := c.Query("end_date")
	quickFilter := c.QueryArray("quick_filter")
	periodType := c.Query("period_type")

	var startDate *time.Time
	if startDateParam != "" {
		parsed, err := time.Parse("2006-01-02", startDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if endDateParam != "" {
		parsed, err := time.Parse("2006-01-02", endDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		endDate = &parsed
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"kpi_category_ids",
		"kpi_status_ids",
		"kpi_review_status_ids",
		"review_frequency_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		StartDate:   startDate,
		EndDate:     endDate,
		Filters:     filters,
		QuickFilter: quickFilter,
		PeriodType:  &periodType,
		IsKRI:       isKRI,
	}

	result, err := kdc.service.GetKPIAchievementByDomain(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetAllKPIs(c *gin.Context) {
	quickFilter := c.QueryArray("quick_filter")
	startDateParam := c.Query("start_date")
	endDateParam := c.Query("end_date")
	periodType := c.Query("period_type")

	var startDate *time.Time
	if startDateParam != "" {
		parsed, err := time.Parse("2006-01-02", startDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if endDateParam != "" {
		parsed, err := time.Parse("2006-01-02", endDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		endDate = &parsed
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"kpi_category_ids",
		"kpi_status_ids",
		"kpi_review_status_ids",
		"review_frequency_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(kdc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kdc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := dtos.KPIDashboardQueryParamsDTO{
		StartDate:   startDate,
		EndDate:     endDate,
		PeriodType:  &periodType,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		IsKRI:       isKRI,
	}

	kpis, err := kdc.service.GetAllKPIs(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(kdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, kpis)
}

func (kdc *KPIDashboardController) GetTopKPIs(c *gin.Context) {

	startDateParam := c.Query("start_date")
	endDateParam := c.Query("end_date")
	quickFilter := c.QueryArray("quick_filter")
	periodType := c.Query("period_type")

	var startDate *time.Time
	if startDateParam != "" {
		parsed, err := time.Parse("2006-01-02", startDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if endDateParam != "" {
		parsed, err := time.Parse("2006-01-02", endDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		endDate = &parsed
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"kpi_category_ids",
		"kpi_status_ids",
		"kpi_review_status_ids",
		"review_frequency_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		StartDate:   startDate,
		EndDate:     endDate,
		Filters:     filters,
		QuickFilter: quickFilter,
		PeriodType:  &periodType,
		IsKRI:       isKRI,
	}

	result, err := kdc.service.GetTopKPIs(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (kdc *KPIDashboardController) GetBottomKPIs(c *gin.Context) {

	startDateParam := c.Query("start_date")
	endDateParam := c.Query("end_date")
	quickFilter := c.QueryArray("quick_filter")
	periodType := c.Query("period_type")

	var startDate *time.Time
	if startDateParam != "" {
		parsed, err := time.Parse("2006-01-02", startDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if endDateParam != "" {
		parsed, err := time.Parse("2006-01-02", endDateParam)
		if err != nil {
			utility.SendInternalErrResp(kdc.db, c, "Invalid date parameter")
			return
		}
		endDate = &parsed
	}

	var isKRI *bool
	if isKRIStr := c.Query("is_kri"); isKRIStr != "" {
		val, err := strconv.ParseBool(isKRIStr)
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid is_kri parameter"))
			return
		}
		isKRI = &val
	}

	allowedParams := []string{
		"kpi_category_ids",
		"kpi_status_ids",
		"kpi_review_status_ids",
		"review_frequency_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := &dtos.KPIDashboardQueryParamsDTO{
		StartDate:   startDate,
		EndDate:     endDate,
		Filters:     filters,
		QuickFilter: quickFilter,
		PeriodType:  &periodType,
		IsKRI:       isKRI,
	}

	result, err := kdc.service.GetBottomKPIs(queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(kdc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}
