package controllers

import (
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicInitiativeStatusController struct {
	service services.StrategicInitiativeStatusServiceInterface
	db      *gorm.DB
}

func NewStrategicInitiativeStatusController(service services.StrategicInitiativeStatusServiceInterface) *StrategicInitiativeStatusController {
	return &StrategicInitiativeStatusController{
		service: service,
	}
}

func (c *StrategicInitiativeStatusController) GetAllStrategicInitiativeStatuses(ctx *gin.Context) {
	sort := ctx.Query("sort_order")
	sortBy := ctx.Query("sort_by")
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")

	items, err := c.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(c.db, ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (c *StrategicInitiativeStatusController) GetStrategicInitiativeStatusDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	item, err := c.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	ctx.JSON(http.StatusOK, item)
}
