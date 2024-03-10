package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SeverityLevelController struct {
	service services.SeverityLevelServiceInterface
	db      *gorm.DB
}

func NewSeverityLevelController(service services.SeverityLevelServiceInterface) *SeverityLevelController {
	return &SeverityLevelController{
		service: service,
	}
}

func (slc *SeverityLevelController) GetAllSeverityLevels(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	SeverityLevels, err := slc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(slc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, SeverityLevels)
}
