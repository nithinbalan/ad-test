package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/grc/organizational-context/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StakeholderController struct {
	service services.StakeholderServiceInterface
	db      *gorm.DB
}

func NewStakeholderController(service services.StakeholderServiceInterface) *StakeholderController {
	return &StakeholderController{
		service: service,
	}
}

func (sc *StakeholderController) GetAllStakeholders(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"title_en",
		"stakeholder_category_ids",
		"type_ids",
		"influence_level_ids",
		"impact_level_ids",
		"organization_unit_ids",
		"third_party_register_ids",
		"user_ids",
		"client_ids",
	}
	validSortFields := map[string]string{
		"title_en":                   "stakeholders.title_en",
		"influence_level_en":         "influence_levels.title_en",
		"influence_level_ar":         "influence_levels.title_ar",
		"influence_level_color_code": "influence_levels.color_code",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(sc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(sc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := organizationalContextDTOs.StakeholderListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	stakeholders, err := sc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(sc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, stakeholders)
}

func (sc *StakeholderController) GetStakeholderDetails(c *gin.Context) {
	idStr := c.Param("stakeholderId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	stakeholder, err := sc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, stakeholder)
}

func (sc *StakeholderController) CreateStakeholder(c *gin.Context) {

	var stakeholderFormDTO organizationalContextDTOs.StakeholderFormDTO

	if err := c.ShouldBindJSON(&stakeholderFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sc.db, c, err.Error())
		return
	}

	id, err := sc.service.Create(currentUserID, &stakeholderFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (sc *StakeholderController) UpdateStakeholder(c *gin.Context) {
	var stakeholderFormDTO organizationalContextDTOs.StakeholderFormDTO

	stakeholderID, err := strconv.Atoi(c.Param("stakeholderId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&stakeholderFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sc.db, c, err.Error())
		return
	}

	err = sc.service.Update(uint(stakeholderID), currentUserID, &stakeholderFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (sc *StakeholderController) DeleteStakeholders(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sc.db, c, err.Error())
		return
	}

	failedIDs, err := sc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
