package controllers

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TenableController struct {
	service services.TenableServiceInterface
}

func NewTenableController(service services.TenableServiceInterface) *TenableController {
	return &TenableController{service: service}
}

func (tc *TenableController) GetAllAssets(c *gin.Context) {
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	params := dtos.TenableAssetListQueryParamsDTO{
		Page:     page,
		PageSize: pageSize,
	}

	response, err := tc.service.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
