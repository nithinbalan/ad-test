package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProcessCriticalityLevelController struct {
	service services.ProcessCriticalityLevelServiceInterface
	db      *gorm.DB
}

func NewProcessCriticalityLevelController(service services.ProcessCriticalityLevelServiceInterface) *ProcessCriticalityLevelController {
	return &ProcessCriticalityLevelController{
		service: service,
	}
}

func (pclc *ProcessCriticalityLevelController) GetAllProcessCriticalityLevels(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	processCriticalityLevels, err := pclc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(pclc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, processCriticalityLevels)
}

func (pclc *ProcessCriticalityLevelController) GetProcessCriticalityLevelDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	processCriticalityLevel, err := pclc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, processCriticalityLevel)
}
