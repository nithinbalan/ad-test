package controllers

import (
	"net/http"

	"jethur/grc/settings/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RAGSettingController struct {
	service services.RAGSettingServiceInterface
	db      *gorm.DB
}

func NewRAGSettingController(service services.RAGSettingServiceInterface) *RAGSettingController {
	return &RAGSettingController{
		service: service,
	}
}

func (rsc *RAGSettingController) GetAllRAGSetting(c *gin.Context) {
	data, err := rsc.service.GetAllRAGSetting()
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}
