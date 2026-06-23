package controllers

import (
	"jethur/grc/awareness-reporting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessHandlingStatusController struct {
	service services.AwarenessHandlingStatusServiceInterface
	db      *gorm.DB
}

func NewAwarenessHandlingStatusController(service services.AwarenessHandlingStatusServiceInterface) *AwarenessHandlingStatusController {
	return &AwarenessHandlingStatusController{
		service: service,
	}
}

func (ahsc *AwarenessHandlingStatusController) GetAllAwarenessHandlingStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	awarenessHandlingStatuses, err := ahsc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(ahsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, awarenessHandlingStatuses)
}

func (ahsc *AwarenessHandlingStatusController) GetAwarenessHandlingStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	awarenessHandlingStatus, err := ahsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, awarenessHandlingStatus)
}