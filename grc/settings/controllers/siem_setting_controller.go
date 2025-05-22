package controllers

import (
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SIEMSettingController struct {
	service services.SIEMSettingServiceInterface
	db      *gorm.DB
}

func NewSIEMSettingController(service services.SIEMSettingServiceInterface) *SIEMSettingController {
	return &SIEMSettingController{
		service: service,
	}
}

func (ssc *SIEMSettingController) GetSIEMSettings(c *gin.Context) {
	siem, err := ssc.service.GetFirst()
	if err != nil {
		utility.SendInternalErrResp(ssc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, siem)
}

func (ssc *SIEMSettingController) UpdateSIEMSettings(c *gin.Context) {
	var req dtos.SIEMSettingDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ssc.db, c, err.Error())
		return
	}

	err = ssc.service.Update(currentUserID, req)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
		default:
			utility.SendInternalErrResp(ssc.db, c, err.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}
