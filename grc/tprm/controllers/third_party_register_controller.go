package controllers

import (
	"errors"
	"jethur/grc/tprm/dtos"
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyRegisterController struct {
	service services.ThirdPartyRegisterServiceInterface
	db      *gorm.DB
}

func NewThirdPartyRegisterController(service services.ThirdPartyRegisterServiceInterface) *ThirdPartyRegisterController {
	return &ThirdPartyRegisterController{
		service: service,
	}
}

func (tprc *ThirdPartyRegisterController) GetAllThirdPartyRegister(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, "Invalid export parameter")
		return
	}
	allowedParams := []string{
		"third_party_type_ids", "country_ids", "organization_unit_ids",
		"severity_level_ids", "responsible_user_ids",
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

	query := dtos.ThirdPartyRegisterListQueryParamsDTO{
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

	thirdParties, err := tprc.service.GetAll(&query)
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

		excelFile, err := tprc.service.ExportThirdPartyRegister(thirdParties.Items, currentUserID)
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

func (tprc *ThirdPartyRegisterController) GetThirdPartyRegisterExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	tprc.GetAllThirdPartyRegister(c)
}

func (tprc *ThirdPartyRegisterController) GetThirdPartyRegisterDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	result, err := tprc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (tprc *ThirdPartyRegisterController) CreateThirdPartyRegister(c *gin.Context) {
	var form dtos.ThirdPartyRegisterFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	id, err := tprc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (tprc *ThirdPartyRegisterController) UpdateThirdPartyRegister(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ThirdPartyRegisterFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	err = tprc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tprc *ThirdPartyRegisterController) DeleteThirdPartyRegister(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	failedIDs, err := tprc.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
