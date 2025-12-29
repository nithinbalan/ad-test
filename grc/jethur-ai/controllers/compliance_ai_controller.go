package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/jethur-ai/services"
	"jethur/utility"
	"net/http"
	"strconv"
)

type ComplianceAIController struct {
	service services.ComplianceAIServiceInterface
	db      *gorm.DB
}

func NewComplianceAIController(service services.ComplianceAIServiceInterface) *ComplianceAIController {
	return &ComplianceAIController{
		service: service,
	}
}

func (caic *ComplianceAIController) GetComplianceAutofill(c *gin.Context) {
	var req dtos.ComplianceAutofillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := caic.service.GetComplianceAutofill(req)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (caic *ComplianceAIController) GetComplianceSummary(c *gin.Context) {
	complianceID, err := strconv.ParseUint(c.Param("complianceId"), 10 ,64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid compliance ID",
		})
		return
	}

	result, err := caic.service.GetComplianceSummary(uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (caic *ComplianceAIController) GetComplianceClauseGuidance(c *gin.Context) {
	clauseID, err := strconv.ParseUint(c.Param("clauseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid clause ID",
		})
		return
	}

	result, err := caic.service.GetComplianceClauseGuidance(uint(clauseID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}




