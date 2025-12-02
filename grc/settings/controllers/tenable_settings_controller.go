package controllers

import (
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
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

func (tc *TenableSettingsController) SaveTenableSettings(c *gin.Context) {
	var dto dtos.UpdateTenableSettingsDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	if err := tc.service.SaveTenableSettings(currentUserID, dto); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenable settings saved successfully"})
}

func (tc *TenableSettingsController) DeleteTenableSettings(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	if err := tc.service.DeleteTenableSettings(currentUserID); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenable settings deleted successfully"})
}

func (tc *TenableSettingsController) TestConnection(c *gin.Context) {
	var input dtos.TestTenableConnectionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	if err := tc.service.TestConnection(input); err != nil {
		switch err.(type) {
		case *utility.TenableConnectionError:
			utility.SendCustomErrResp(c, http.StatusUnprocessableEntity, err.Error())
			return
		default:
			utility.SendInternalErrResp(nil, c, err.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
