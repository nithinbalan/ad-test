package controllers

import (
	"jethur/grc/models"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskResponseStrategyController struct {
	service services.RiskResponseStrategyServiceInterface
	db      *gorm.DB
}

func NewRiskResponseStrategyController(service services.RiskResponseStrategyServiceInterface) *RiskResponseStrategyController {
	return &RiskResponseStrategyController{
		service: service,
	}
}

func (rrsc *RiskResponseStrategyController) GetAllRiskResponseStrategies(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	riskResponseStrategies, err := rrsc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(rrsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskResponseStrategies)
}

func (rrsc *RiskResponseStrategyController) GetRiskResponseStrategyDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskResponseStrategy, err := rrsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskResponseStrategy)
}

func (rrsc *RiskResponseStrategyController) CreateRiskResponseStrategy(c *gin.Context) {
	var strategy *models.RiskResponseStrategy

	if err := c.ShouldBindJSON(&strategy); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrsc.db, c, err.Error())
		return
	}

	createdStrategy, err := rrsc.service.Create(currentUserID, strategy)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdStrategy)
}

func (rrsc *RiskResponseStrategyController) UpdateRiskResponseStrategy(c *gin.Context) {
	strategyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedStrategy *models.RiskResponseStrategy
	if err := c.ShouldBindJSON(&updatedStrategy); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrsc.db, c, err.Error())
		return
	}

	updatedStrategy, err = rrsc.service.Update(uint(strategyID), currentUserID, updatedStrategy)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedStrategy)
}

func (rrsc *RiskResponseStrategyController) DeleteRiskResponseStrategies(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrsc.db, c, err.Error())
		return
	}

	failedIDs, err := rrsc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rrsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rrsc *RiskResponseStrategyController) ActivateRiskResponseStrategies(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrsc.db, c, err.Error())
		return
	}

	err = rrsc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rrsc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rrsc *RiskResponseStrategyController) DeactivateRiskResponseStrategies(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrsc.db, c, err.Error())
		return
	}

	err = rrsc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rrsc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
