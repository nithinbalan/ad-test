package controllers

import (
	"errors"
	businessFunctionDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BusinessFunctionController struct {
	service services.BusinessFunctionServiceInterface
	db      *gorm.DB
}

func NewBusinessFunctionController(service services.BusinessFunctionServiceInterface) *BusinessFunctionController {
	return &BusinessFunctionController{
		service: service,
	}
}

func (bfc *BusinessFunctionController) GetAllBusinessFunctions(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"severity_level_ids",
	}
	validSortFields := map[string]string{
		"reference_number":               "business_functions.reference_number",
		"title":                          "business_functions.title",
		"description":                    "business_functions.description",
		"severity_level_en":              "severity_levels.title_en",
		"severity_level_ar":              "severity_levels.title_ar",
		"severity_level_color_code":      "severity_levels.color_code",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"responsible_user_last_name_en":  "responsible_user.last_name_en",
		"responsible_user_last_name_ar":  "responsible_user.last_name_ar",
		"responsible_user_email":         "responsible_user.email",
		"responsible_user_initial":       "responsible_user.initial",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := businessFunctionDTOs.BusinessFunctionListQueryParamsDTO{
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
	businesFunctions, err := bfc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, businesFunctions)
}

func (bfc *BusinessFunctionController) GetBusinessFunctionDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	businessFunction, err := bfc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(bfc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, businessFunction)
}

func (bfc *BusinessFunctionController) CreateBusinessFunction(c *gin.Context) {
	var businessFunctionFormDTO businessFunctionDTOs.BusinessFunctionFormDTO

	if err := c.ShouldBindJSON(&businessFunctionFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, err.Error())
		return
	}

	id, err := bfc.service.Create(currentUserID, &businessFunctionFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(bfc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (bfc *BusinessFunctionController) UpdateBusinessFunction(c *gin.Context) {
	var businessFunctionFormDTO businessFunctionDTOs.BusinessFunctionFormDTO

	businesFunctionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&businessFunctionFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, err.Error())
		return
	}

	err = bfc.service.Update(uint(businesFunctionID), currentUserID, &businessFunctionFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(bfc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (bfc *BusinessFunctionController) DeleteBusinessFunctions(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, err.Error())
		return
	}

	failedIDs, err := bfc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
