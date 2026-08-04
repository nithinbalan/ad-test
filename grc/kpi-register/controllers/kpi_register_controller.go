package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	kpiDTOs "jethur/grc/kpi-register/dtos"
	"jethur/grc/kpi-register/services"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPIRegisterController struct {
	service services.KPIRegisterServiceInterface
	db      *gorm.DB
}

func NewKPIRegisterController(service services.KPIRegisterServiceInterface) *KPIRegisterController {
	return &KPIRegisterController{
		service: service,
	}
}

func (krc *KPIRegisterController) GetAllKPIRegister(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	var export bool
	if val, exists := c.Get("export"); exists {
		export = val.(bool)
	} else {
		export, _ = strconv.ParseBool(c.DefaultQuery("export", "false"))
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"kpi_category_ids",
		"kpi_status_ids",
		"kpi_review_status_ids",
		"review_frequency_ids",
		"calculation_type_ids",
		"unit_ids",
		"kpi_ids",
		"responsible_user_ids",
		"kpi_achievement_logic_ids",

		// ✅ add these
		"types",
		"kpi_status_types",
		"kpi_review_status_types",
	}
	validSortFields := map[string]string{
		"id":                             "kpi_register.id",
		"created_at":                     "kpi_register.created_at",
		"reference_number":               "kpi_register.reference_number",
		"title":                          "kpi_register.title",
		"type":                           "kpi_register.type",
		"description":                    "kpi_register.description",
		"actual_value":                   "kpi_register.actual_value",
		"achievement":                    "kpi_register.achievement",
		"target_value":                   "kpi_register.target_value",
		"start_date":                     "kpi_register.start_date",
		"end_date":                       "kpi_register.end_date",
		"kpi_category_title_en":          "kpi_categories.title_en",
		"kpi_category_title_ar":          "kpi_categories.title_ar",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"kpi_status_title_en":            "kpi_statuses.title_en",
		"kpi_status_title_ar":            "kpi_statuses.title_ar",
		"kpi_review_status_title_en":     "kpi_review_statuses.title_en",
		"kpi_review_status_title_ar":     "kpi_review_statuses.title_ar",
		"review_frequency_title_en":      "review_frequencies.title_en",
		"review_frequency_title_ar":      "review_frequencies.title_ar",
		"calculation_type_title_en":      "calculation_types.title_en",
		"calculation_type_title_ar":      "calculation_types.title_ar",
		"unit_title_en":                  "units.title_en",
		"unit_title_ar":                  "units.title_ar",
		"kpi_achievement_logic_title_en": "kpi_achievement_logics.title_en",
		"kpi_achievement_logic_title_ar": "kpi_achievement_logics.title_ar",
		"target_description":             "kpi_register.target_description",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}
	queryParams := kpiDTOs.KPIRegisterListQueryParamsDTO{
		Sort:          sort,
		SortBy:        sortField,
		Keyword:       keyword,
		Page:          page,
		PageSize:      pageSize,
		Filters:       filters,
		Export:        export,
		QuickFilter:   quickFilter,
		CurrentUserID: int(currentUserID),
	}
	kpiRegister, err := krc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(krc.db, c, err.Error())
			return
		}

		excelFile, err := krc.service.ExportKPIRegister(kpiRegister.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(krc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="KPI register.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, kpiRegister)
}

func (krc *KPIRegisterController) GetKPIRegisterDetails(c *gin.Context) {
	idStr := c.Param("kpiRegisterId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiRegister, err := krc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(krc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, kpiRegister)
}

func (krc *KPIRegisterController) CreateKPIRegister(c *gin.Context) {
	var kpiRegisterFormDTO kpiDTOs.KPIRegisterFormDTO

	if err := c.ShouldBindJSON(&kpiRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	id, err := krc.service.Create(currentUserID, &kpiRegisterFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(krc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (krc *KPIRegisterController) UpdateKPIRegister(c *gin.Context) {
	var kpiRegisterFormDTO kpiDTOs.KPIRegisterFormDTO

	kpiRegisterID, err := strconv.Atoi(c.Param("kpiRegisterId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&kpiRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	err = krc.service.Update(uint(kpiRegisterID), currentUserID, &kpiRegisterFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(krc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (krc *KPIRegisterController) DeleteKPIRegister(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	failedIDs, err := krc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (krc *KPIRegisterController) GetKPITrendAnalysis(c *gin.Context) {
	idStr := c.Param("kpiRegisterId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlTrend, err := krc.service.GetTrendAnalysis(uint(id), 10)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(krc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controlTrend)
}

// In kpi_register_controller.go
func (krc *KPIRegisterController) GetKPILatestValueUpdate(c *gin.Context) {
	idStr := c.Param("kpiRegisterId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	latestValueUpdate, err := krc.service.GetLatestValueUpdateByKPIRegisterID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(krc.db, c, e.Error())
			return
		}
	}

	// If no value update is found, return an empty object
	if latestValueUpdate == nil {
		c.JSON(http.StatusOK, struct{}{})
		return
	}

	c.JSON(http.StatusOK, latestValueUpdate)
}

func (krc *KPIRegisterController) GenerateKPIRegisterTemplate(c *gin.Context) {
	file, err := krc.service.GenerateKPIRegisterTemplate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate KPI Register template",
		})
		return
	}
	defer file.Close()

	c.Header(
		"Content-Disposition",
		"attachment; filename=KPIRegisterTemplate.xlsx",
	)
	c.Header(
		"Content-Type",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	)

	c.File(file.Name())
}

func (krc *KPIRegisterController) ValidateKPIRegisterTemplate(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File is required",
		})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open file",
		})
		return
	}
	defer openedFile.Close()

	validatedData, err := krc.service.ValidateKPIRegisterTemplate(openedFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, validatedData)
}

func (krc *KPIRegisterController) ValidateImportKPIRegisters(c *gin.Context) {
	var input []map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	validatedData, err := krc.service.ValidateImportKPIRegisters(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, validatedData)
}

func (krc *KPIRegisterController) ImportKPIRegisters(c *gin.Context) {
	var input []map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(
			krc.db,
			c,
			err.Error(),
		)
		return
	}

	createdIDs, err := krc.service.ImportKPIRegisters(
		input,
		uint(currentUserID),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"created_ids": createdIDs,
	})
}

func (krc *KPIRegisterController) GetKPICategorySummary(c *gin.Context) {

	allowedParams := []string{
		"kpi_status_ids",
		"kpi_category_ids",
		"kpi_review_status_ids",
		"responsible_user_ids",
		"review_frequency_ids",
		"calculation_type_ids",
		"unit_ids",
	}

	filters := utility.ParseQueryFilters(
		allowedParams,
		c.Request.URL.Query(),
	)

	// ✅ Quick Filters
	quickFilters := c.QueryArray("quick_filter")

	// ✅ Current User
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	queryParams := &kpiDTOs.KPIDashboardQueryParamsDTO{
		Filters:     filters,
		QuickFilter: quickFilters,
	}

	result, err := krc.service.GetKPICategorySummary(queryParams, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(krc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}
func (krc *KPIRegisterController) GetKPIStatusSummary(c *gin.Context) {

	allowedParams := []string{
		"kpi_status_ids",
		"kpi_category_ids",
		"kpi_review_status_ids",
		"responsible_user_ids",
		"review_frequency_ids",
		"calculation_type_ids",
		"unit_ids",
	}

	filters := utility.ParseQueryFilters(
		allowedParams,
		c.Request.URL.Query(),
	)

	// ✅ Quick Filters
	quickFilters := c.QueryArray("quick_filter")

	// ✅ Current User
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	queryParams := &kpiDTOs.KPIDashboardQueryParamsDTO{
		Filters:     filters,
		QuickFilter: quickFilters,
	}

	result, err := krc.service.GetKPIStatusSummary(queryParams, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(krc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}
func (krc *KPIRegisterController) GetKPIOverallSummary(c *gin.Context) {

	allowedParams := []string{
		"kpi_status_ids",
		"kpi_category_ids",
		"kpi_review_status_ids",
		"responsible_user_ids",
		"review_frequency_ids",
		"calculation_type_ids",
		"unit_ids",
	}

	filters := utility.ParseQueryFilters(
		allowedParams,
		c.Request.URL.Query(),
	)

	// ✅ Quick Filters (same as list API)
	quickFilters := c.QueryArray("quick_filter")

	// ✅ Current User (for My KPIs)
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	queryParams := &kpiDTOs.KPIDashboardQueryParamsDTO{
		Filters:     filters,
		QuickFilter: quickFilters,
	}

	result, err := krc.service.GetKPIOverallSummary(queryParams, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(krc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (krc *KPIRegisterController) GetKPIByReviewStatus(c *gin.Context) {

	allowedParams := []string{
		"kpi_status_ids",
		"kpi_category_ids",
		"kpi_review_status_ids",
		"responsible_user_ids",
		"review_frequency_ids",
		"calculation_type_ids",
		"unit_ids",
	}

	filters := utility.ParseQueryFilters(
		allowedParams,
		c.Request.URL.Query(),
	)

	// ✅ Quick Filters
	quickFilters := c.QueryArray("quick_filter")

	// ✅ Current User
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	queryParams := &kpiDTOs.KPIDashboardQueryParamsDTO{
		Filters:     filters,
		QuickFilter: quickFilters,
	}

	result, err := krc.service.GetKPIByReviewStatus(queryParams, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(krc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (krc *KPIRegisterController) GetKPIExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"
	c.Set("export", true)
	krc.GetAllKPIRegister(c)
}

func (krc *KPIRegisterController) ActivateKPIRegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	err = krc.service.ActivateKPIRegisters(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (krc *KPIRegisterController) ArchiveKPIRegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	err = krc.service.ArchiveKPIRegisters(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (krc *KPIRegisterController) ForceDeleteKPIRegister(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	failedIDs, err := krc.service.ForceDelete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (krc *KPIRegisterController) GetAllKPIResponsibleUsers(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "first_name_en")
	keyword := c.Query("keyword")
	status := c.Query("status")

	allowedParams := []string{
		"designation_ids",
		"organization_unit_1_ids",
		"organization_unit_2_ids",
		"organization_unit_3_ids",
		"organization_unit_4_ids",
		"organization_unit_5_ids",
		"user_status_ids",
		"license_type_ids",
		"organization_unit_ids",
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
		utility.SendInternalErrResp(krc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, "Invalid page_size parameter")
		return
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationDTOs.UserListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	users, err := krc.service.GetAllResponsibleUsers(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}
