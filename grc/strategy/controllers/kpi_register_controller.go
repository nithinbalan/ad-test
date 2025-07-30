package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
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

	// Define an array of allowed parameter names
	allowedParams := []string{
		"title",
		"kpi_category_ids",
		"kpi_status_ids",
		"review_frequency_ids",
		"calculation_type_ids",
		"unit_ids",
	}
	validSortFields := map[string]string{
		"reference_number":               "kpi_register.reference_number",
		"title":                          "kpi_register.title",
		"description":                    "kpi_register.description",
		"start_date":                     "kpi_register.start_date",
		"end_date":                       "kpi_register.end_date",
		"kpi_category_title_en":          "kpi_categories.title_en",
		"kpi_category_title_ar":          "kpi_categories.title_ar",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"kpi_status_title_en":            "kpi_statuses.title_en",
		"kpi_status_title_ar":            "kpi_statuses.title_ar",
		"review_frequency_title_en":      "review_frequencies.title_en",
		"review_frequency_title_ar":      "review_frequencies.title_ar",
		"calculation_type_title_en":      "calculation_types.title_en",
		"calculation_type_title_ar":      "calculation_types.title_ar",
		"unit_title_en":                  "units.title_en",
		"unit_title_ar":                  "units.title_ar",
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
	queryParams := strategyDTOs.KPIRegisterListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
	kpiRegister, err := krc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
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

	var kpiRegisterFormDTO strategyDTOs.KPIRegisterFormDTO

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
	var kpiRegisterFormDTO strategyDTOs.KPIRegisterFormDTO

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
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	// This will return 200 OK with null in the response body if latestValueUpdate is nil
	c.JSON(http.StatusOK, latestValueUpdate)
}
