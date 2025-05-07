package controllers

import (
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
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

func (adsc *ActiveDirectorySettingController) ADLogin(c *gin.Context) {
	var req dtos.ADLoginRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if err := adsc.service.Authenticate(req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Authenticated successfully"})
}
