package controllers

import (
	"net/http"

	"jethur/grc/library/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceLibraryController struct {
	service services.ComplianceLibraryServiceInterface
	db      *gorm.DB
}

func NewComplianceLibraryController(service services.ComplianceLibraryServiceInterface, db *gorm.DB) *ComplianceLibraryController {
	return &ComplianceLibraryController{
		service: service,
		db:      db,
	}
}

func (clc *ComplianceLibraryController) GetAll(c *gin.Context) {
	data, err := clc.service.GetAll()
	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}
