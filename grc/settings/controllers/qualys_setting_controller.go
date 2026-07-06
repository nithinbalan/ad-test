package controllers

import (
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QualysSettingsController struct {
	service services.QualysSettingsServiceInterface
}

func NewQualysSettingsController(service services.QualysSettingsServiceInterface) *QualysSettingsController {
	return &QualysSettingsController{service: service}
}

func (tc *QualysSettingsController) GetQualysSettings(c *gin.Context) {

	data, err := tc.service.GetQualysSettings()
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (tc *QualysSettingsController) SaveQualysSettings(c *gin.Context) {
	var dto dtos.UpdateQualysSettingsDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	if err := tc.service.SaveQualysSettings(currentUserID, dto); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (tc *QualysSettingsController) DeleteQualysSettings(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	if err := tc.service.DeleteQualysSettings(currentUserID); err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (tc *QualysSettingsController) TestConnection(c *gin.Context) {
	var input dtos.TestQualysConnectionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	if err := tc.service.TestConnection(input); err != nil {
		switch e := err.(type) {
		case *utility.QualysConnectionError:
			utility.SendCustomErrResp(c, http.StatusUnprocessableEntity, err.Error())
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return

		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return

		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return

		default:
			utility.SendInternalErrResp(nil, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tc *QualysSettingsController) CheckSavedConnection(c *gin.Context) {

	err := tc.service.CheckSavedConnection()
	if err != nil {
		switch e := err.(type) {
		case *utility.QualysConnectionError:
			utility.SendCustomErrResp(c, http.StatusUnprocessableEntity, err.Error())
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return

		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return

		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return

		default:
			utility.SendInternalErrResp(nil, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Qualys connection (from saved config) successful",
	})
}