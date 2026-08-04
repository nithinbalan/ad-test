package controllers

import (
	"jethur/grc/kpi-register/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPIAchievementLogicController struct {
	service services.KPIAchievementLogicServiceInterface
	db      *gorm.DB
}

func NewKPIAchievementLogicController(service services.KPIAchievementLogicServiceInterface) *KPIAchievementLogicController {
	return &KPIAchievementLogicController{
		service: service,
	}
}

func (kalc *KPIAchievementLogicController) GetAllKPIAchievementLogics(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	kpiAchievementLogics, err := kalc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(kalc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, kpiAchievementLogics)
}

func (kalc *KPIAchievementLogicController) GetKPIAchievementLogicDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiAchievementLogic, err := kalc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, kpiAchievementLogic)
}
