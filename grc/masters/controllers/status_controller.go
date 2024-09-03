package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StatusController struct {
	service services.StatusServiceInterface
	db      *gorm.DB
}

func NewStatusController(service services.StatusServiceInterface) *StatusController {
	return &StatusController{
		service: service,
	}
}

func (sc *StatusController) GetAllStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	statuses, err := sc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(sc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, statuses)
}

func (sc *StatusController) GetStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInternalErrResp(sc.db, c, "Invalid id")
		return
	}

	status, err := sc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, status)
}
