package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LifecycleStatusController struct {
	service services.LifecycleStatusServiceInterface
	db      *gorm.DB
}

func NewLifecycleStatusController(service services.LifecycleStatusServiceInterface) *LifecycleStatusController {
	return &LifecycleStatusController{
		service: service,
	}
}

func (lsc *LifecycleStatusController) GetAllLifecycleStatuses(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	lifecycleStatuses, err := lsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(lsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, lifecycleStatuses)
}
