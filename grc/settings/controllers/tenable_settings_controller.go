package controllers

import (
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TenableSettingsController struct {
	service services.TenableSettingsServiceInterface
}

func NewTenableSettingsController(service services.TenableSettingsServiceInterface) *TenableSettingsController {
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

	c.Status(http.StatusOK)
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

	c.Status(http.StatusOK)
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

func (tc *TenableSettingsController) CheckSavedConnection(c *gin.Context) {

	err := tc.service.CheckSavedConnection()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return

		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return

		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return

		case *utility.TenableConnectionError:
			utility.SendCustomErrResp(c, http.StatusUnprocessableEntity, e.Error())
			return

		default:
			utility.SendInternalErrResp(nil, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tenable connection (from saved config) successful",
	})
}
