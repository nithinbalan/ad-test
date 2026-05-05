package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskUpdatePlanTimelineTypeController struct {
	service services.RiskUpdatePlanTimelineTypeServiceInterface
	db      *gorm.DB
}

func NewRiskUpdatePlanTimelineTypeController(service services.RiskUpdatePlanTimelineTypeServiceInterface) *RiskUpdatePlanTimelineTypeController {
	return &RiskUpdatePlanTimelineTypeController{
		service: service,
	}
}

func (rtc *RiskUpdatePlanTimelineTypeController) GetAllRiskUpdatePlanTimelineTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	timelineTypes, err := rtc.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(rtc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, timelineTypes)
}

func (rtc *RiskUpdatePlanTimelineTypeController) GetRiskUpdatePlanTimelineTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	timelineType, err := rtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, timelineType)
}
