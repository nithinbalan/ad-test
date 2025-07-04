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

func (clc *ComplianceLibraryController) GetNCADetails(c *gin.Context) {
	data, err := clc.service.GetNCADetails()
	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (clc *ComplianceLibraryController) GetSAMADetails(c *gin.Context) {
	data, err := clc.service.GetSAMADetails()
	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (clc *ComplianceLibraryController) GetClauseControls(c *gin.Context) {
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

func (clc *ComplianceLibraryController) GetAllMappings(c *gin.Context) {
	data, err := clc.service.GetGroupedControlMappings()
	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (clc *ComplianceLibraryController) GetGenericControlMappings(c *gin.Context) {
	data, err := clc.service.GetGenericControlMappings()
	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}
