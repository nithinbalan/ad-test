package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategyProfileStrategicIntiativeController struct {
	service services.StrategyProfileStrategicIntiativeServiceInterface
	db      *gorm.DB
}

func NewStrategyProfileStrategicIntiativeController(service services.StrategyProfileStrategicIntiativeServiceInterface) *StrategyProfileStrategicIntiativeController {
	return &StrategyProfileStrategicIntiativeController{
		service: service,
	}
}

func (snec *StrategyProfileStrategicIntiativeController) GetAllStrategyProfileStrategicIntiatives(c *gin.Context) {
	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "id")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"strategy_profile_ids",
		"strategic_objective_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := strategyDTOs.StrategyProfileStrategicIntiativeListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	strategicIntiatives, err := snec.service.GetAll(uint(strategyProfileID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(snec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, strategicIntiatives)
}

func (snec *StrategyProfileStrategicIntiativeController) GetStrategyProfileStrategicIntiativeDetails(c *gin.Context) {
	// Get strategy profile ID from URL
	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Get strategic intiative ID from URL
	strategicIntiativeIDStr := c.Param("id")
	strategicIntiativeID, err := strconv.ParseUint(strategicIntiativeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Call service to get the strategic intiative details
	strategicIntiative, err := snec.service.GetFullDetailsByID(uint(strategyProfileID), uint(strategicIntiativeID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(snec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, strategicIntiative)
}

func (snec *StrategyProfileStrategicIntiativeController) CreateStrategyProfileStrategicIntiative(c *gin.Context) {
	var FormDTO strategyDTOs.StrategyProfileStrategicIntiativeFormDTO

	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&FormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, err.Error())
		return
	}

	ids, err := snec.service.Create(uint(strategyProfileID), currentUserID, &FormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(snec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.SuccessfulIDsResponse{SuccessfulIds: ids})
}

func (snec *StrategyProfileStrategicIntiativeController) DeleteStrategyProfileStrategicIntiative(c *gin.Context) {
	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicIntiativeIDStr := c.Param("id")
	strategicIntiativeID, err := strconv.ParseUint(strategicIntiativeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, err.Error())
		return
	}

	err = snec.service.Delete(uint(strategyProfileID), uint(strategicIntiativeID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(snec.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
