package controllers

import (
	"jethur/grc/models"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskSettingsController struct {
	service services.RiskSettingsServiceInterface
	db      *gorm.DB
}

func NewRiskSettingsController(service services.RiskSettingsServiceInterface) *RiskSettingsController {
	return &RiskSettingsController{
		service: service,
	}
}

// GET: Fetch the Single Risk Setting
func (rsc *RiskSettingsController) GetAllRiskSettings(c *gin.Context) {
	riskSetting, err := rsc.service.Get()
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, riskSetting)
}

// PUT: Update the Single Risk Setting
func (rsc *RiskSettingsController) UpdateRiskSettings(c *gin.Context) {
	var updatedSetting *models.RiskSetting
	if err := c.ShouldBindJSON(&updatedSetting); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	updatedSetting, err = rsc.service.Update(currentUserID, updatedSetting)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, updatedSetting)
}
