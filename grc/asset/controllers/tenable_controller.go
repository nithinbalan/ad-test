package controllers

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenableController struct {
	service services.TenableServiceInterface
	db      *gorm.DB
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

	response, err := tc.service.GetAllAssets(&params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (tc *TenableController) ValidateTenableAssets(c *gin.Context) {
	var input []map[string]interface{}

	// Client sends edited assets from Tenable
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, err.Error())
		return
	}

	// Validate client-sent assets (not fetching)
	assets, err := tc.service.ValidateTenableAssets(input, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, assets)
}

func (tc *TenableController) ImportAssets(c *gin.Context) {
	var input []map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, err.Error())
		return
	}

	createdIDs, err := tc.service.ImportAssets(input, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"created_ids": createdIDs,
	})
}
