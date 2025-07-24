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

type StrategyProfileProjectActionPlanController struct {
	service services.StrategyProfileProjectActionPlanServiceInterface
	db      *gorm.DB
}

func NewStrategyProfileProjectActionPlanController(service services.StrategyProfileProjectActionPlanServiceInterface) *StrategyProfileProjectActionPlanController {
	return &StrategyProfileProjectActionPlanController{
		service: service,
	}
}

func (sppac *StrategyProfileProjectActionPlanController) GetAllStrategyProfileProjectActionPlans(c *gin.Context) {
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
		"project_ids",
		"project_milestone_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(sppac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(sppac.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := strategyDTOs.StrategyProfileProjectActionPlanListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	strategicIntiatives, err := sppac.service.GetAll(uint(strategyProfileID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sppac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, strategicIntiatives)
}

func (sppac *StrategyProfileProjectActionPlanController) CreateStrategyProfileProjectActionPlan(c *gin.Context) {
	var FormDTO strategyDTOs.StrategyProfileProjectActionPlanFormDTO

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
		utility.SendInternalErrResp(sppac.db, c, err.Error())
		return
	}

	ids, err := sppac.service.Create(uint(strategyProfileID), currentUserID, &FormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sppac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.SuccessfulIDsResponse{SuccessfulIds: ids})
}

func (sppac *StrategyProfileProjectActionPlanController) DeleteStrategyProfileProjectActionPlan(c *gin.Context) {
	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategyProfileProjectActionPlanIDStr := c.Param("id")
	strategyProfileProjectActionPlanID, err := strconv.ParseUint(strategyProfileProjectActionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sppac.db, c, err.Error())
		return
	}

	err = sppac.service.Delete(uint(strategyProfileID), uint(strategyProfileProjectActionPlanID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sppac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
