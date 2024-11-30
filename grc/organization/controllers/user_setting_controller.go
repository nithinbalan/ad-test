package controllers

import (
	organizationDTOs "jethur/grc/organization/dtos"
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

func (usc *UserSettingController) GetSetting(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	userSetting, err := usc.service.GetUserSetting(uint(userID))
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

	c.JSON(http.StatusOK, userSetting)
}

func (usc *UserSettingController) UpdateSetting(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var userSettingFormDTO organizationDTOs.UserSettingFormDTO
	if err := c.ShouldBindJSON(&userSettingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(usc.db, c, err.Error())
		return
	}

	err = usc.service.UpdateUserSetting(uint(userID), currentUserID, &userSettingFormDTO)
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

func (usc *UserSettingController) EnableMainNotification(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = usc.service.EnableMainNotification(uint(userID))

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

func (usc *UserSettingController) DisableMainNotification(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = usc.service.DisableMainNotification(uint(userID))

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

func (usc *UserSettingController) EnableMFA(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(usc.db, c, err.Error())
		return
	}

	err = usc.service.EnableMFA(uint(userID), currentUserID)

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

func (usc *UserSettingController) DisableMFA(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(usc.db, c, err.Error())
		return
	}

	err = usc.service.DisableMFA(uint(userID), currentUserID)

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

func (usc *UserSettingController) NewPassword(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = usc.service.ResetPassword(uint(userID))
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
