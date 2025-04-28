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

func (ctrl *ThirdPartyRegisterController) GetAllThirdPartyRegister(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"third_party_type_ids", "country_ids", "organization_unit_ids",
		"severity_level_ids", "responsible_user_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("Invalid page parameter"))
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("Invalid page_size parameter"))
		return
	}

	validSortFields := map[string]string{
		"title":             "third_party_register.title",
		"reference_number":  "third_party_register.reference_number",
		"country_title_en":  "countries.title_en",
		"third_party_type":  "third_party_types.title_en",
		"responsible_user":  "responsible_user.first_name_en",
		"organization_unit": "organization_units.title_en",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("Invalid sort_by parameter"))
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
	}

	thirdParties, err := ctrl.service.GetAll(&query)
	if err != nil {
		utility.SendInternalErrResp(ctrl.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, thirdParties)
}

func (ctrl *ThirdPartyRegisterController) GetThirdPartyRegisterDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	result, err := ctrl.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendInternalErrResp(ctrl.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ctrl *ThirdPartyRegisterController) CreateThirdPartyRegister(c *gin.Context) {
	var form dtos.ThirdPartyRegisterFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ctrl.db, c, err.Error())
		return
	}

	id, err := ctrl.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ctrl.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (ctrl *ThirdPartyRegisterController) UpdateThirdPartyRegister(c *gin.Context) {
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
		utility.SendInternalErrResp(ctrl.db, c, err.Error())
		return
	}

	err = ctrl.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ctrl.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ctrl *ThirdPartyRegisterController) DeleteThirdPartyRegister(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ctrl.db, c, err.Error())
		return
	}

	failedIDs, err := ctrl.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ctrl.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
