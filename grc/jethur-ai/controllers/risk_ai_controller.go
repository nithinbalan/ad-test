package controllers

import (
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/jethur-ai/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskAIController struct {
	service services.RiskAIServiceInterface
	db      *gorm.DB
}

func NewRiskAIController(service services.RiskAIServiceInterface) *RiskAIController {
	return &RiskAIController{
		service: service,
	}
}

func (rc *RiskAIController) GetRiskAutofill(c *gin.Context) {
	var req dtos.RiskAutofillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := rc.service.GetRiskAutofill(req)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (rc *RiskAIController) GetRiskRCA(c *gin.Context) {
	var req dtos.RiskRCARequest
	if err := c.ShouldBind(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := rc.service.GetRiskRCA(req)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (rc *RiskAIController) GetRiskAnalysis(c *gin.Context) {
	var req dtos.RiskAIRequest
	if err := c.ShouldBind(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := rc.service.GetRiskAnalysis(req)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (rc *RiskAIController) GetRiskResponseStrategyAnalysis(c *gin.Context) {
	var req dtos.RiskAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := rc.service.GetRiskResponseStrategyAnalysis(req)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (rc *RiskAIController) GetRiskControlAnalysis(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Risk ID"})
		return
	}

	result, err := rc.service.GetRiskControlAnalysis(uint(id))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (rc *RiskAIController) GetRiskKriAnalysis(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Risk ID"})
		return
	}

	result, err := rc.service.GetRiskKriAnalysis(uint(id))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (rc *RiskAIController) GetRiskResidualControlAnalysis(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Risk ID"})
		return
	}

	result, err := rc.service.GetRiskResidualControlAnalysis(uint(id))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (rc *RiskAIController) GetRiskResidualRiskAnalysis(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("riskId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Risk ID"})
		return
	}

	result, err := rc.service.GetRiskResidualRiskAnalysis(uint(id))
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (rc *RiskAIController) ChatWithRisk(c *gin.Context) {
	var req dtos.RiskAIChatRequest
	if err := c.ShouldBind(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	result, err := rc.service.ChatWithRisk(req, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
