package controllers

import (
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TriggerConditionController struct {
	service services.TriggerConditionServiceInterface
	db      *gorm.DB
}

func NewTriggerConditionController(service services.TriggerConditionServiceInterface) *TriggerConditionController {
	return &TriggerConditionController{
		service: service,
	}
}

func (tcc *TriggerConditionController) GetAllTriggerConditiones(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	triggerConditions, err := tcc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, triggerConditions)
}

func (tcc *TriggerConditionController) GetTriggerConditionDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	triggerCondition, err := tcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, triggerCondition)
}
