package controllers

import (
	"net/http"

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
	data, err := rcc.service.GetAllReminderContexts()
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}
