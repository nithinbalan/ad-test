package controllers

import (
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPIStatusController struct {
	service services.KPIStatusServiceInterface
	db      *gorm.DB
}

func NewKPIStatusController(service services.KPIStatusServiceInterface) *KPIStatusController {
	return &KPIStatusController{
		service: service,
	}
}

func (lsc *KPIStatusController) GetAllKPIStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	kpiStatuses, err := lsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(lsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, kpiStatuses)
}

func (lsc *KPIStatusController) GetKPIStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiStatus, err := lsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, kpiStatus)
}
