package controllers

import (
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TimezoneController struct {
	service services.TimezoneServiceInterface
	db      *gorm.DB
}

func NewTimezoneController(service services.TimezoneServiceInterface) *TimezoneController {
	return &TimezoneController{
		service: service,
	}
}

func (tc *TimezoneController) GetAllTimezones(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	timezones, err := tc.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(tc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, timezones)
}

func (tc *TimezoneController) GetTimezoneDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	timezone, err := tc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, timezone)
}
