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

type StrategyProfileActionPlanController struct {
	service services.StrategyProfileActionPlanServiceInterface
	db      *gorm.DB
}

func NewStrategyProfileActionPlanController(service services.StrategyProfileActionPlanServiceInterface) *StrategyProfileActionPlanController {
	return &StrategyProfileActionPlanController{
		service: service,
	}
}

func (snec *StrategyProfileActionPlanController) GetAllStrategyProfileActionPlans(c *gin.Context) {
	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"strategy_profile_focus_area_ids",
		"strategy_profile_goal_ids",
		"strategy_profile_objective_ids",
		"strategy_profile_initiative_ids",
		"strategy_profile_project_ids",
		"responsible_user_ids",
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

	queryParams := strategyDTOs.StrategyProfileActionPlanListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
	}

	actionPlan, err := snec.service.GetAll(uint(strategyProfileID), &queryParams)
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

	c.JSON(http.StatusOK, actionPlan)
}

func (snec *StrategyProfileActionPlanController) GetStrategyProfileActionPlanDetails(c *gin.Context) {
	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlan, err := snec.service.GetFullDetailsByID(uint(strategyProfileID), uint(id))
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

	c.JSON(http.StatusOK, actionPlan)
}

func (snec *StrategyProfileActionPlanController) CreateStrategyProfileActionPlan(c *gin.Context) {
	var FormDTO strategyDTOs.StrategyProfileActionPlanFormDTO

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

	id, err := snec.service.Create(uint(strategyProfileID), currentUserID, &FormDTO)

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

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (snec *StrategyProfileActionPlanController) UpdateStrategyProfileActionPlan(c *gin.Context) {
	var FormDTO strategyDTOs.StrategyProfileActionPlanFormDTO

	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
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

	err = snec.service.Update(uint(strategyProfileID), uint(id), currentUserID, &FormDTO)
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

func (snec *StrategyProfileActionPlanController) DeleteStrategyProfileActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(snec.db, c, err.Error())
		return
	}

	failedIDs, err := snec.service.Delete(uint(strategyProfileID), request.IDs, currentUserID)
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

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
