package controllers

import (
	"errors"
	"net/http"
	"strconv"

	organizationDTOs "jethur/grc/organization/dtos"

	"jethur/grc/organization/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationStructureController struct {
	service services.OrganizationStructureServiceInterface
	db      *gorm.DB
}

func NewOrganizationStructureController(
	service services.OrganizationStructureServiceInterface,
) *OrganizationStructureController {
	return &OrganizationStructureController{
		service: service,
	}
}

func (osc *OrganizationStructureController) GetOrganizationStructure(c *gin.Context) {
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

func (osc *OrganizationStructureController) GetUsersByOrganizationUnitID(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "first_name_en")
	// Get organization unit ID from URL parameter
	unitID, err := strconv.ParseUint(c.Param("unitId"), 10, 32)
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
	}

	// Call service to get users by organization unit ID
	users, err := osc.service.GetUsersByOrganizationUnit(&queryParams, uint(unitID))
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

func (osc *OrganizationStructureController) GetOrganizationUnitSummary(c *gin.Context) {
	unitID, err := strconv.ParseUint(c.Param("unitId"), 10, 32)
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
