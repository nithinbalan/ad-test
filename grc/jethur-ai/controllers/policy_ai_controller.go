package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/jethur-ai/services"
	"jethur/utility"
	"net/http"
)

type PolicyAIController struct {
	service services.PolicyAIServiceInterface
	db      *gorm.DB
}

func NewPolicyAIController(service services.PolicyAIServiceInterface) *PolicyAIController {
	return &PolicyAIController{
		service: service,
	}
}

func (paic *PolicyAIController) GetPolicyAutofill(c *gin.Context) {
	var req dtos.PolicyAutofillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := paic.service.GetPolicyAutofill(req)
	if err != nil {
		utility.SendInternalErrResp(paic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}