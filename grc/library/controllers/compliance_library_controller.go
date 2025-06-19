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

func (clc *ComplianceLibraryController) GetComplianceLibraryDetailsByTitle(c *gin.Context) {
	key := c.Param("title")

	var titleEn string
	switch key {
	case "sama":
		titleEn = "Cyber Security Framework | Saudi Arabian Monetary Authority"
	case "nca":
		titleEn = "Essential Cybersecurity Controls (ECC – 2: 2024)"
	default:
		utility.SendNotFoundErrResp(c)
		return
	}

	data, err := clc.service.GetComplianceLibraryDetailsByTitle(titleEn)
	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (clc *ComplianceLibraryController) GetClauseControlTree(c *gin.Context) {
	versionIDStr := c.Param("versionID")
	complianceIDStr := c.Param("complianceLibraryID")

	versionID, err := strconv.Atoi(versionIDStr)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	complianceID, err := strconv.Atoi(complianceIDStr)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := clc.service.GetClauseTreeByVersion(uint(versionID), uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
