package controllers

import (
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationUnitController struct {
	service services.OrganizationUnitServiceInterface
	db      *gorm.DB
}

func NewOrganizationUnitController(service services.OrganizationUnitServiceInterface) *OrganizationUnitController {
	return &OrganizationUnitController{
		service: service,
	}
}

func (ouc *OrganizationUnitController) GetAllOrganizationUnits(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	primaryStr := c.Query("primary")
	levelStr := c.Query("level")
	parentUnitIDsStr := c.Query("parent_unit_ids")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"organization_level_ids",
		"organization_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse primary as a boolean
	var primary bool
	if primaryStr != "" {
		// Parse primary as a boolean if it's not empty
		var err error
		primary, err = strconv.ParseBool(primaryStr)
		if err != nil {
			utility.SendInternalErrResp(ouc.db, c, "Invalid primary parameter")
			return
		}
	}

	var level uint
	var err error
	if levelStr != "" {
		levelParsed, err := strconv.ParseUint(levelStr, 10, 32) // Parse as uint32 to avoid overflow on 32-bit systems
		if err != nil {
			utility.SendInternalErrResp(ouc.db, c, "Invalid level parameter")
			return
		}
		level = uint(levelParsed)
	}

	var parentUnitIDs []uint
	if parentUnitIDsStr != "" {
		for _, idStr := range strings.Split(parentUnitIDsStr, ",") {
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				utility.SendInternalErrResp(ouc.db, c, "Invalid parent unit ID")
				return
			}
			parentUnitIDs = append(parentUnitIDs, uint(id))
		}
	}

	queryParams := settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Sort:          sort,
		SortBy:        sortBy,
		Keyword:       keyword,
		Status:        status,
		Filters:       filters,
		Primary:       primary,
		Level:         level,
		ParentUnitIDs: parentUnitIDs,
	}

	organizationUnits, err := ouc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(ouc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, organizationUnits)
}

func (ouc *OrganizationUnitController) GetOrganizationUnitDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	organizationUnit, err := ouc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, organizationUnit)
}

func (ouc *OrganizationUnitController) CreateOrganizationUnit(c *gin.Context) {
	var organizationUnitFormDTO settingsDTOs.OrganizationUnitFormDTO

	if err := c.ShouldBindJSON(&organizationUnitFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ouc.db, c, err.Error())
		return
	}

	createdOrganizationUnit, err := ouc.service.Create(currentUserID, &organizationUnitFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ouc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdOrganizationUnit)
}

func (ouc *OrganizationUnitController) UpdateOrganizationUnit(c *gin.Context) {
	var organizationUnitFormDTO settingsDTOs.OrganizationUnitFormDTO

	organizationUnitID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&organizationUnitFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ouc.db, c, err.Error())
		return
	}

	err = ouc.service.Update(uint(organizationUnitID), currentUserID, &organizationUnitFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ouc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ouc *OrganizationUnitController) DeleteOrganizationUnits(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ouc.db, c, err.Error())
		return
	}

	failedIDs, err := ouc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ouc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ouc *OrganizationUnitController) ActivateOrganizationUnits(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ouc.db, c, err.Error())
		return
	}

	err = ouc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ouc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ouc *OrganizationUnitController) DeactivateOrganizationUnits(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ouc.db, c, err.Error())
		return
	}

	err = ouc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ouc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
