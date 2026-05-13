package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BusinessContinuityDurationUnitController struct {
	service services.BusinessContinuityDurationUnitServiceInterface
	db      *gorm.DB
}

func NewBusinessContinuityDurationUnitController(service services.BusinessContinuityDurationUnitServiceInterface) *BusinessContinuityDurationUnitController {
	return &BusinessContinuityDurationUnitController{
		service: service,
	}
}

func (bcdusc *BusinessContinuityDurationUnitController) GetAllBusinessContinuityDurationUnits(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	businessContinuityDurationUnits, err := bcdusc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(bcdusc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, businessContinuityDurationUnits)
}

func (bcdusc *BusinessContinuityDurationUnitController) GetBusinessContinuityDurationUnitDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	businessContinuityDurationUnit, err := bcdusc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, businessContinuityDurationUnit)
}
