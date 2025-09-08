package controllers

import (
	"net/http"

	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReminderContextController struct {
	service services.ReminderContextServiceInterface
	db      *gorm.DB
}

func NewReminderContextController(service services.ReminderContextServiceInterface) *ReminderContextController {
	return &ReminderContextController{
		service: service,
	}
}

func (rcc *ReminderContextController) GetAllReminderContexts(c *gin.Context) {
	allowedParams := []string{
		"module_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := dtos.ReminderContextListQueryParamsDTO{
		Filters: filters,
	}

	data, err := rcc.service.GetAllReminderContexts(queryParams)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}
