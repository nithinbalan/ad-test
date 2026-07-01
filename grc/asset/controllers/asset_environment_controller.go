package controllers

import (
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetEnvironmentController struct {
	service services.AssetEnvironmentServiceInterface
	db      *gorm.DB
}

func NewAssetEnvironmentController(service services.AssetEnvironmentServiceInterface) *AssetEnvironmentController {
	return &AssetEnvironmentController{
		service: service,
	}
}

func (aec *AssetEnvironmentController) GetAllAssetEnvironments(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	assetEnvironments, err := aec.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(aec.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, assetEnvironments)
}

func (aec *AssetEnvironmentController) GetAssetEnvironmentDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetEnvironment, err := aec.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, assetEnvironment)
}