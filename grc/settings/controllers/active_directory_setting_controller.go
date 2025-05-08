package controllers

import (
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ActiveDirectorySettingController struct {
	service services.ActiveDirectorySettingServiceInterface
	db      *gorm.DB
}

func NewActiveDirectorySettingController(service services.ActiveDirectorySettingServiceInterface) *ActiveDirectorySettingController {
	return &ActiveDirectorySettingController{
		service: service,
	}
}
func (adsc *ActiveDirectorySettingController) GetActiveDirectorySettings(c *gin.Context) {
	ad, err := adsc.service.GetFirst()
	if err != nil {
		utility.SendInternalErrResp(adsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, ad)
}

func (adsc *ActiveDirectorySettingController) CreateActiveDirectorySetting(c *gin.Context) {
	var req dtos.ActiveDirectorySettingDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adsc.db, c, err.Error())
		return
	}
	err = adsc.service.Create(currentUserID, req)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adsc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)

}
