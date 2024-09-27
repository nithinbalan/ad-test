package controllers

import (
	settingDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmailTemplateSettingsController struct {
	service services.EmailTemplateSettingsServiceInterface
	db      *gorm.DB
}

func NewEmailTemplateSettingsController(service services.EmailTemplateSettingsServiceInterface) *EmailTemplateSettingsController {
	return &EmailTemplateSettingsController{
		service: service,
	}
}

func (etsc *EmailTemplateSettingsController) GetEmailTemplateSettings(c *gin.Context) {
	emailTemplateSettings, err := etsc.service.GetEmailTemplateSettings()
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, emailTemplateSettings)
}

func (etsc *EmailTemplateSettingsController) UpdateEmailTemplateSettings(c *gin.Context) {
	var emailTemplateSettingsFormDTO settingDTOs.EmailTemplateSettingsFormDTO
	if err := c.ShouldBindJSON(&emailTemplateSettingsFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	emailTemplateSettingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(etsc.db, c, err.Error())
		return
	}

	err = etsc.service.Update(uint(emailTemplateSettingID), currentUserID, &emailTemplateSettingsFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(etsc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
