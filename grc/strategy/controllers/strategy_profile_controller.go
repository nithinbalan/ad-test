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

type StrategyProfileController struct {
	service services.StrategyProfileServiceInterface
	db      *gorm.DB
}

func NewStrategyProfileController(service services.StrategyProfileServiceInterface) *StrategyProfileController {
	return &StrategyProfileController{
		service: service,
	}
}

func (spc *StrategyProfileController) GetAllStrategyProfiles(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"title",
	}
	validSortFields := map[string]string{
		"title": "strategy_profiles.title",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := strategyDTOs.StrategyProfileListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
	strategyProfiles, err := spc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, strategyProfiles)
}

func (spc *StrategyProfileController) GetStrategyProfileDetails(c *gin.Context) {
	idStr := c.Param("strategyProfileId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategyProfile, err := spc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, strategyProfile)
}

func (spc *StrategyProfileController) CreateStrategyProfile(c *gin.Context) {

	var strategyProfileFormDTO strategyDTOs.StrategyProfileFormDTO

	if err := c.ShouldBindJSON(&strategyProfileFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	id, err := spc.service.Create(currentUserID, &strategyProfileFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (spc *StrategyProfileController) UpdateStrategyProfile(c *gin.Context) {
	var strategyProfileFormDTO strategyDTOs.StrategyProfileFormDTO

	strategyProfileID, err := strconv.Atoi(c.Param("strategyProfileId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&strategyProfileFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	err = spc.service.Update(uint(strategyProfileID), currentUserID, &strategyProfileFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (spc *StrategyProfileController) DeleteStrategyProfiles(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	failedIDs, err := spc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
