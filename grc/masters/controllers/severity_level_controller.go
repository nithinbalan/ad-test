package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

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

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	severityLevels, err := slc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(slc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, severityLevels)
}

func (slc *SeverityLevelController) GetSeverityLevelDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	severityLevel, err := slc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, severityLevel)
}
