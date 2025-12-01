package controllers

import (
	"fmt"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

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

func (tc *TenableSettingsController) Activate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid ID"))
		return
	}

	if err := tc.service.ActivateSettings(uint(id)); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenable settings activated"})
}

func (tc *TenableSettingsController) Deactivate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid ID"))
		return
	}

	if err := tc.service.DeactivateSettings(uint(id)); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenable settings deactivated"})
}

func (tc *TenableSettingsController) TestConnection(c *gin.Context) {
	var input dtos.TestTenableConnectionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	if err := tc.service.TestConnection(input); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection successful"})
}
