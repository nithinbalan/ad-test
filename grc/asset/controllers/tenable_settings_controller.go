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

func (tc *TenableSettingsController) GetSettings(c *gin.Context) {

	data, err := tc.service.GetSettings()
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (tc *TenableSettingsController) UpdateConnection(c *gin.Context) {
	var input dtos.UpdateTenableConnectionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	if err := tc.service.UpdateConnectionSettings(input); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Connection settings updated"})
}

func (tc *TenableSettingsController) UpdateSyncSettings(c *gin.Context) {
	var input dtos.UpdateTenableSyncSettingsDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	if err := tc.service.UpdateSyncSettings(input); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sync settings updated"})
}

func (tc *TenableSettingsController) Activate(c *gin.Context) {
	if err := tc.service.ActivateSettings(); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tenable settings activated"})
}

func (tc *TenableSettingsController) Deactivate(c *gin.Context) {
	if err := tc.service.DeactivateSettings(); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tenable settings deactivated"})
}

func (tc *TenableSettingsController) TestConnection(c *gin.Context) {
	err := tc.service.TestConnection()
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Connection successful"})
}
