package controllers

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TenableSettingsController struct {
	service services.TenableSettingsService
}

func NewTenableSettingsController(service services.TenableSettingsService) *TenableSettingsController {
	return &TenableSettingsController{service: service}
}

func (tc *TenableSettingsController) GetTenableSettings(c *gin.Context) {

	data, err := tc.service.GetTenableSettings()
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (tc *TenableSettingsController) UpdateTenableSettings(c *gin.Context) {
	var input dtos.UpdateTenableSettingsDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	if err := tc.service.UpdateTenableSettings(input); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenable settings updated"})
}

func (tc *TenableSettingsController) DeleteTenableSettings(c *gin.Context) {
	if err := tc.service.DeleteTenableSettings(); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenable settings deleted"})
}

func (tc *TenableSettingsController) TestConnection(c *gin.Context) {
	var input dtos.TestTenableConnectionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := tc.service.TestConnection(input)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)

}
