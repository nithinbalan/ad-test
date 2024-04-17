package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserSettingController struct {
	service services.UserSettingServiceInterface
	db      *gorm.DB
}

func NewUserSettingController(service services.UserSettingServiceInterface) *UserSettingController {
	return &UserSettingController{
		service: service,
	}
}

func (usc *UserSettingController) GetNotificationSettings(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	notifications, err := usc.service.GetNotificationSettings(uint(userID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(usc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, notifications)
}

func (usc *UserSettingController) EnableNotification(c *gin.Context) {
	usc.toggleNotification(c, true)
}

func (usc *UserSettingController) DisableNotification(c *gin.Context) {
	usc.toggleNotification(c, false)
}

func (usc *UserSettingController) toggleNotification(c *gin.Context, enable bool) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	notificationMasterID, err := strconv.Atoi(c.Param("notificationMasterId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = usc.service.UpdateNotificationSetting(uint(userID), uint(notificationMasterID), enable)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(usc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
