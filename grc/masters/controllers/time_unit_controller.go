package controllers

import (
	"fmt"
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TimeUnitController struct {
	service services.TimeUnitServiceInterface
	db      *gorm.DB
}

func NewTimeUnitController(service services.TimeUnitServiceInterface) *TimeUnitController {
	return &TimeUnitController{
		service: service,
	}
}

func (icc *TimeUnitController) GetAllTimeUnits(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	fmt.Println("1")
	timeUnits, err := icc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, timeUnits)
}

func (icc *TimeUnitController) GetTimeUnitDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	timeUnit, err := icc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, timeUnit)
}
