package controllers

import (
	"jethur/grc/models"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GovernanceSettingsController struct {
	service services.GovernanceSettingsServiceInterface
	db      *gorm.DB
}

func NewGovernanceSettingsController(service services.GovernanceSettingsServiceInterface) *GovernanceSettingsController {
	return &GovernanceSettingsController{
		service: service,
	}
}

// GET: Fetch Governance Settings
func (gsc *GovernanceSettingsController) GetAllGovernanceSettings(c *gin.Context) {
	governanceSetting, err := gsc.service.Get()
	if err != nil {
		utility.SendInternalErrResp(gsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, governanceSetting)
}

// PUT: Update Governance Settings
func (gsc *GovernanceSettingsController) UpdateGovernanceSettings(c *gin.Context) {
	var updatedSetting *models.GovernanceSetting

	if err := c.ShouldBindJSON(&updatedSetting); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(gsc.db, c, err.Error())
		return
	}

	updatedSetting, err = gsc.service.Update(currentUserID, updatedSetting)
	if err != nil {
		utility.SendInternalErrResp(gsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, updatedSetting)
}
