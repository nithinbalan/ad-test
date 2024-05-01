package controllers

import (
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CurrencyController struct {
	service services.CurrencyServiceInterface
	db      *gorm.DB
}

func NewCurrencyController(service services.CurrencyServiceInterface) *CurrencyController {
	return &CurrencyController{
		service: service,
	}
}

func (cc *CurrencyController) GetAllCurrencies(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	currencies, err := cc.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, currencies)
}

func (cc *CurrencyController) GetCurrencyDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currency, err := cc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, currency)
}
