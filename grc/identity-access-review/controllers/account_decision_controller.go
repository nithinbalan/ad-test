package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/identity-access-review/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
)

type AccountDecisionController struct {
	service services.AccountDecisionServiceInterface
}

func NewAccountDecisionController(service services.AccountDecisionServiceInterface) *AccountDecisionController {
	return &AccountDecisionController{service: service}
}

func (adc *AccountDecisionController) GetAllAccountDecisions(c *gin.Context) {
	accountDecisions, err := adc.service.GetAll(
		c.Query("sort_order"),
		c.Query("sort_by"),
		c.Query("keyword"),
		c.Query("status"),
	)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, accountDecisions)
}

func (adc *AccountDecisionController) GetAccountDecisionDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	accountDecision, err := adc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	c.JSON(http.StatusOK, accountDecision)
}
