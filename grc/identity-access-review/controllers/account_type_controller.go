package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/identity-access-review/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
)

type AccountTypeController struct {
	service services.AccountTypeServiceInterface
}

func NewAccountTypeController(service services.AccountTypeServiceInterface) *AccountTypeController {
	return &AccountTypeController{service: service}
}

func (atc *AccountTypeController) GetAllAccountTypes(c *gin.Context) {
	accountTypes, err := atc.service.GetAll(
		c.Query("sort_order"),
		c.Query("sort_by"),
		c.Query("keyword"),
		c.Query("status"),
	)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, accountTypes)
}

func (atc *AccountTypeController) GetAccountTypeDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	accountType, err := atc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	c.JSON(http.StatusOK, accountType)
}
