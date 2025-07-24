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

type StrategyProfileProjectMilestoneController struct {
	service services.StrategyProfileProjectMilestoneServiceInterface
	db      *gorm.DB
}

func NewStrategyProfileProjectMilestoneController(service services.StrategyProfileProjectMilestoneServiceInterface) *StrategyProfileProjectMilestoneController {
	return &StrategyProfileProjectMilestoneController{
		service: service,
	}
}

func (sppmc *StrategyProfileProjectMilestoneController) GetAllStrategyProfileProjectMilestones(c *gin.Context) {
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
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(sppmc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(sppmc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := strategyDTOs.StrategyProfileProjectMilestoneListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	strategicIntiatives, err := sppmc.service.GetAll(uint(strategyProfileID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sppmc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, strategicIntiatives)
}

func (sppmc *StrategyProfileProjectMilestoneController) CreateStrategyProfileProjectMilestone(c *gin.Context) {
	var FormDTO strategyDTOs.StrategyProfileProjectMilestoneFormDTO

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
		utility.SendInternalErrResp(sppmc.db, c, err.Error())
		return
	}

	ids, err := sppmc.service.Create(uint(strategyProfileID), currentUserID, &FormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sppmc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.SuccessfulIDsResponse{SuccessfulIds: ids})
}

func (sppmc *StrategyProfileProjectMilestoneController) DeleteStrategyProfileProjectMilestone(c *gin.Context) {
	strategyProfileIDStr := c.Param("strategyProfileId")
	strategyProfileID, err := strconv.ParseUint(strategyProfileIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategyProfileProjectMilestoneIDStr := c.Param("id")
	strategyProfileProjectMilestoneID, err := strconv.ParseUint(strategyProfileProjectMilestoneIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sppmc.db, c, err.Error())
		return
	}

	err = sppmc.service.Delete(uint(strategyProfileID), uint(strategyProfileProjectMilestoneID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sppmc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
