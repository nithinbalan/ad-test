package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConformanceController struct {
	service services.ConformanceServiceInterface
	db      *gorm.DB
}

func NewConformanceController(service services.ConformanceServiceInterface) *ConformanceController {
	return &ConformanceController{
		service: service,
	}
}

func (cc *ConformanceController) GetAllConformances(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	conformances, err := cc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, conformances)
}

func (cc *ConformanceController) GetConformanceDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid id")
		return
	}

	conformance, err := cc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, conformance)
}
