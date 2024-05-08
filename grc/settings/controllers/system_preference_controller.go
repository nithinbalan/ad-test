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

type SystemPreferenceController struct {
	service services.SystemPreferenceServiceInterface
	db      *gorm.DB
}

func NewSystemPreferenceController(service services.SystemPreferenceServiceInterface) *SystemPreferenceController {
	return &SystemPreferenceController{
		service: service,
	}
}

func (spc *SystemPreferenceController) GetSystemPreference(c *gin.Context) {
	systemPreference, err := spc.service.GetSystemPreference()
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, systemPreference)
}

func (spc *SystemPreferenceController) UpdateSystemPreference(c *gin.Context) {
	systemPreferenceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var systemPreferenceFormDTO settingDTOs.SystemPreferenceFormDTO
	if err := c.ShouldBindJSON(&systemPreferenceFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(spc.db, c, err.Error())
		return
	}

	err = spc.service.Update(uint(systemPreferenceID), currentUserID, &systemPreferenceFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(spc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
