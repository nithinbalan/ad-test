package controllers

import (
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/jethur-ai/services"
	"jethur/utility"
	"net/http"

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
	var req dtos.RiskAnalysisRequest
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
