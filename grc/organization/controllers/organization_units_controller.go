package controllers

import (
	"errors"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"jethur/grc/organization/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationUnitsController struct {
	service services.OrganizationUnitsServiceInterface
	db      *gorm.DB
}

func NewOrganizationUnitsController(
	service services.OrganizationUnitsServiceInterface,
) *OrganizationUnitsController {
	return &OrganizationUnitsController{
		service: service,
	}
}

func (osc *OrganizationUnitsController) GetOrganizationStructure(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")

	organizationStructure, err := osc.service.GetOrganizationStructure(sort, sortBy)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, organizationStructure)
}

func (osc *OrganizationUnitsController) GetUsersByOrganizationUnitID(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "first_name_en")
	Level := c.DefaultQuery("level", "1")
	// Get organization unit ID from URL parameter
	unitID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	level, err := strconv.ParseUint(Level, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"first_name_en":  "users.first_name_en",
		"first_name_ar":  "users.first_name_ar",
		"designation_en": "designations.title_en",
		"designation_ar": "designations.title_ar",
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Page:     page,
		PageSize: pageSize,
		Level:    uint(level),
	}

	// Call service to get users by organization unit ID
	users, err := osc.service.GetUsersByOrganizationUnit(&queryParams, uint(unitID), uint(level))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, users)
}

func (osc *OrganizationUnitsController) GetOrganizationUnitSummary(c *gin.Context) {
	unitID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := osc.service.GetOrganizationUnitSummary(uint(unitID))
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (osc *OrganizationUnitsController) UpdatePrimaryOrganizationUnitHead(c *gin.Context) {
	var form organizationDTOs.DepartmentHeadUpdateDTO

	unitID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, err.Error())
		return
	}

	err = osc.service.UpdatePrimaryOrganizationUnitHead(uint(unitID), uint(userID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
		}
		return
	}
	c.Status(http.StatusOK)
}

func (osc *OrganizationUnitsController) GetAllPoliciesByOrganizationUnit(c *gin.Context) {
	unitIDStr := c.Param("id")

	unitID, err := strconv.ParseUint(unitIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	policies, err := osc.service.GetAllPoliciesByOrganizationUnit(uint(unitID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, policies)
}

func (osc *OrganizationUnitsController) GetAllRiskByOrganizationUnit(c *gin.Context) {

	unitIDStr := c.Param("id")

	unitID, err := strconv.ParseUint(unitIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	valueStr := c.Query("value")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(osc.db, c, "Invalid value parameter")
			return
		}
	}
	allowedParams := []string{
		"current_likelihood_ids",
		"current_impact_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Value:    value,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
	risks, err := osc.service.GetAllRiskByOrganizationUnit(uint(unitID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (osc *OrganizationUnitsController) GetParentOrganizationUnitWithHead(c *gin.Context) {
	unitIDStr := c.Param("id")

	unitID, err := strconv.ParseUint(unitIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := osc.service.GetParentOrganizationUnitWithHead(uint(unitID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, result)
}
