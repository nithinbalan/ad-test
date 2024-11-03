package controllers

import (
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetStatusController struct {
	service services.AssetStatusServiceInterface
	db      *gorm.DB
}

func NewAssetStatusController(service services.AssetStatusServiceInterface) *AssetStatusController {
	return &AssetStatusController{
		service: service,
	}
}

func (asc *AssetStatusController) GetAllAssetStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	assetStatuses, err := asc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, assetStatuses)
}

func (asc *AssetStatusController) GetAssetStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetStatus, err := asc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, assetStatus)
}
