package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CountryController struct {
	service services.CountryServiceInterface
	db      *gorm.DB
}

func NewCountryController(service services.CountryServiceInterface) *CountryController {
	return &CountryController{
		service: service,
	}
}

func (cc *CountryController) GetAllCountries(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	countries, err := cc.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, countries)
}

func (cc *CountryController) GetCountryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	country, err := cc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, country)
}
