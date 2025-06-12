package controllers

import (
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/jethur-ai/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonAIController struct {
	service services.CommonAIServiceInterface
}

func NewCommonAIController(service services.CommonAIServiceInterface) *CommonAIController {
	return &CommonAIController{service: service}
}

func (cc *CommonAIController) RewriteDescription(c *gin.Context) {
	var req dtos.DescriptionRewriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := cc.service.RewriteDescription(req)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (cc *CommonAIController) TranslateLanguage(c *gin.Context) {
	var req dtos.LanguageTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := cc.service.TranslateLanguage(req)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
