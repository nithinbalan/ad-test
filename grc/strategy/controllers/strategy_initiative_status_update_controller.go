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

type StrategyInitiativeStatusUpdateController struct {
	service services.StrategyInitiativeStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewStrategyInitiativeStatusUpdateController(service services.StrategyInitiativeStatusUpdateServiceInterface) *StrategyInitiativeStatusUpdateController {
	return &StrategyInitiativeStatusUpdateController{
		service: service,
	}
}

func (sisuc *StrategyInitiativeStatusUpdateController) GetAllStrategyInitiativeStatusUpdates(c *gin.Context) {
	strategyInitiativeIDStr := c.Param("id")
	strategyInitiativeID, err := strconv.ParseUint(strategyInitiativeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	list, err := sisuc.service.GetAll(uint(strategyInitiativeID), sort, sortBy, keyword)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sisuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, list)
}

func (sisuc *StrategyInitiativeStatusUpdateController) CreateStrategyInitiativeStatusUpdate(c *gin.Context) {
	var dto strategyDTOs.StrategyInitiativeStatusUpdateFormDTO

	strategyInitiativeIDStr := c.Param("id")
	strategyInitiativeID, err := strconv.ParseUint(strategyInitiativeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sisuc.db, c, err.Error())
		return
	}

	id, err := sisuc.service.Create(uint(strategyInitiativeID), currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sisuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (sisuc *StrategyInitiativeStatusUpdateController) UpdateStrategyInitiativeStatusUpdate(c *gin.Context) {
	var dto strategyDTOs.StrategyInitiativeStatusUpdateFormDTO

	strategyInitiativeIDStr := c.Param("id")
	strategyInitiativeID, err := strconv.ParseUint(strategyInitiativeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	statusUpdateIDStr := c.Param("statusUpdateId")
	statusUpdateID, err := strconv.ParseUint(statusUpdateIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sisuc.db, c, err.Error())
		return
	}

	err = sisuc.service.Update(uint(strategyInitiativeID), uint(statusUpdateID), currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sisuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (sisuc *StrategyInitiativeStatusUpdateController) GetLatestStrategyInitiativeStatusUpdate(c *gin.Context) {
	strategyInitiativeIDStr := c.Param("id")
	strategyInitiativeID, err := strconv.ParseUint(strategyInitiativeIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	item, err := sisuc.service.GetLatestStrategyInitiativeStatusUpdate(uint(strategyInitiativeID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sisuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, item)
}
