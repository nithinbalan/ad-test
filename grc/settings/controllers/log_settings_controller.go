package controllers

import (
	logDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogSettingsController struct {
	service services.LogSettingsServiceInterface
	db      *gorm.DB
}

func NewLogSettingsController(service services.LogSettingsServiceInterface) *LogSettingsController {
	return &LogSettingsController{
		service: service,
	}
}

func (lsc *LogSettingsController) GetLogSettings(c *gin.Context) {
	logSettings, err := lsc.service.GetLogSettings()
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, logSettings)
}

func (lsc *LogSettingsController) UpdateLogSettings(c *gin.Context) {
	logSettingsID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var logSettingsFormDTO logDTOs.LogSettingsFormDTO
	if err := c.ShouldBindJSON(&logSettingsFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(lsc.db, c, err.Error())
		return
	}

	err = lsc.service.Update(uint(logSettingsID), currentUserID, &logSettingsFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(lsc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
