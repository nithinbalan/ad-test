package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/jethur-ai/services"
	"jethur/utility"
	"net/http"
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
