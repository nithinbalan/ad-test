package controllers

import (
	"jethur/grc/tprm/dtos"
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TPRMSettingController struct {
	service services.TPRMSettingServiceInterface
	db      *gorm.DB
}

func NewTPRMSettingController(service services.TPRMSettingServiceInterface) *TPRMSettingController {
	return &TPRMSettingController{
		service: service,
	}
}

func (tsc *TPRMSettingController) GetTPRMSetting(c *gin.Context) {
	result, err := tsc.service.GetTPRMSetting()
	if err != nil {
		utility.SendInternalErrResp(tsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (tsc *TPRMSettingController) UpdateTPRMSetting(c *gin.Context) {

	var settingFormDTO dtos.TprmSettingFormDTO

	if err := c.ShouldBindJSON(&settingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tsc.db, c, err.Error())
		return
	}

	err = tsc.service.UpdateTPRMSetting(currentUserID, settingFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tsc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
