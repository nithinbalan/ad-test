package controllers

import (
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategyProfileStrategicObjectiveController struct {
	service services.StrategyProfileStrategicObjectiveServiceInterface
	db      *gorm.DB
}

func NewStrategyProfileStrategicObjectiveController(service services.StrategyProfileStrategicObjectiveServiceInterface) *StrategyProfileStrategicObjectiveController {
	return &StrategyProfileStrategicObjectiveController{
		service: service,
	}
}

func (snec *StrategyProfileStrategicObjectiveController) GetAllStrategyProfileStrategicObjectives(c *gin.Context) {
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
		"focus_area_ids",
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

	queryParams := strategyDTOs.StrategyProfileStrategicObjectiveListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	objectives, err := snec.service.GetAll(uint(strategyProfileID), &queryParams)
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

	c.JSON(http.StatusOK, objectives)
}

func (snec *StrategyProfileStrategicObjectiveController) GetStrategyProfileStrategicObjectiveDetails(c *gin.Context) {
	// Get strategy profile ID from URL
	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Get objective ID from URL
	objectiveIDStr := c.Param("id")
	objectiveID, err := strconv.ParseUint(objectiveIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Call service to get the objective details
	objective, err := snec.service.GetFullDetailsByID(uint(strategyProfileID), uint(objectiveID))
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

	c.JSON(http.StatusOK, objective)
}

func (snec *StrategyProfileStrategicObjectiveController) CreateStrategyProfileStrategicObjective(c *gin.Context) {
	var FormDTO strategyDTOs.StrategyProfileStrategicObjectiveFormDTO

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

	c.JSON(http.StatusCreated, gin.H{"successful_ids": ids})
}

func (snec *StrategyProfileStrategicObjectiveController) DeleteStrategyProfileStrategicObjectives(c *gin.Context) {
	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	objectiveIDStr := c.Param("id")
	objectiveID, err := strconv.ParseUint(objectiveIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, err.Error())
		return
	}

	err = snec.service.Delete(uint(strategyProfileID), uint(objectiveID), currentUserID)
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
