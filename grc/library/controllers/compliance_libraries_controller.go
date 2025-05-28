package controllers

import (
	"net/http"
	"strconv"

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

func (clc *ComplianceLibraryController) GetComplianceLibraryDetails(c *gin.Context) {
	libraryIDStr := c.Param("libraryId")
	libraryID, err := strconv.Atoi(libraryIDStr)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	data, err := clc.service.GetComplianceLibraryDetailsByID(uint(libraryID))
	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}
