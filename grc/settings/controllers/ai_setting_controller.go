package controllers

import (
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AISettingController struct {
	moduleService services.AISettingServiceInterface
	db            *gorm.DB
}

func NewAISettingController(moduleSvc services.AISettingServiceInterface) *AISettingController {
	return &AISettingController{
		moduleService: moduleSvc,
	}
}

func (asc *AISettingController) EnableAISetting(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}

	err = asc.moduleService.EnableAISetting(currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(asc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (asc *AISettingController) DisableAISetting(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}

	err = asc.moduleService.DisableAISetting(currentUserID)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (asc *AISettingController) GetAllAIModules(c *gin.Context) {
	result, err := asc.moduleService.GetAll()
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (asc *AISettingController) GetFeaturesByModuleID(c *gin.Context) {
	moduleIDStr := c.Param("id")
	moduleID, err := strconv.Atoi(moduleIDStr)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := asc.moduleService.GetByModuleID(moduleID)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
