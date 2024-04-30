package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlTypeController struct {
	service services.ControlTypeServiceInterface
	db      *gorm.DB
}

func NewControlTypeController(service services.ControlTypeServiceInterface) *ControlTypeController {
	return &ControlTypeController{
		service: service,
	}
}

func (ctc *ControlTypeController) GetAllControlTypes(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	controlTypes, err := ctc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ctc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, controlTypes)
}
