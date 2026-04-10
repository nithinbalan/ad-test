package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LevelOfMeasurementController struct {
	service services.LevelOfMeasurementServiceInterface
	db      *gorm.DB
}

func NewLevelOfMeasurementController(service services.LevelOfMeasurementServiceInterface) *LevelOfMeasurementController {
	return &LevelOfMeasurementController{
		service: service,
	}
}

func (c *LevelOfMeasurementController) GetAllLevelOfMeasurements(ctx *gin.Context) {
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

func (c *LevelOfMeasurementController) GetLevelOfMeasurementDetails(ctx *gin.Context) {
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
