package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompetencyLevelController struct {
	service services.CompetencyLevelServiceInterface
	db      *gorm.DB
}

func NewCompetencyLevelController(service services.CompetencyLevelServiceInterface) *CompetencyLevelController {
	return &CompetencyLevelController{
		service: service,
	}
}

func (clc *CompetencyLevelController) GetAllCompetencyLevels(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	competencyLevels, err := clc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, competencyLevels)
}

func (clc *CompetencyLevelController) GetCompetencyLevelDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlEfficiencyStatus, err := clc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, controlEfficiencyStatus)
}
