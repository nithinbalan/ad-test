package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KRIRegisterController struct {
	service services.KRIRegisterServiceInterface
	db      *gorm.DB
}

func NewKRIRegisterController(service services.KRIRegisterServiceInterface) *KRIRegisterController {
	return &KRIRegisterController{
		service: service,
	}
}

func (krc *KRIRegisterController) GetAllKRIRegisters(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, "Invalid export parameter")
		return
	}
	// Define an array of allowed parameter names
	allowedParams := []string{
		"kri_category_ids",
		"risk_level_ids",
		"review_frequency_ids",
		"owner_ids",
		"owner_unit_ids",
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
		utility.SendInternalErrResp(krc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

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
	}
	kriRegisters, err := krc.service.GetAll(&queryParams)

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

		excelFile, err := krc.service.ExportKriRegister(kriRegisters.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(krc.db, c, err.Error())
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

func (krc *KRIRegisterController) GetKRIRegisterDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kriRegister, err := krc.service.GetFullDetailsByID(uint(id))
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
	c.JSON(http.StatusOK, kriRegister)
}

func (krc *KRIRegisterController) CreateKRIRegister(c *gin.Context) {

	var kriRegisterFormDTO riskDTOs.KRIRegisterFormDTO
	if err := c.ShouldBindJSON(&kriRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}
	id, err := krc.service.Create(currentUserID, &kriRegisterFormDTO)
	if err != nil {
		switch e := err.(type) {
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

func (krc *KRIRegisterController) UpdateKRIRegister(c *gin.Context) {

	kriRegisterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var kriRegisterFormDTO riskDTOs.KRIRegisterFormDTO
	if err := c.ShouldBindJSON(&kriRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(krc.db, c, err.Error())
		return
	}

	err = krc.service.Update(uint(kriRegisterID), currentUserID, &kriRegisterFormDTO)
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

func (krc *KRIRegisterController) DeleteKRIRegisters(c *gin.Context) {

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

func (krc *KRIRegisterController) GetKriRegisterExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	krc.GetAllKRIRegisters(c)
}
