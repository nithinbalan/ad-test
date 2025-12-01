package controllers

import (
	"net/http"
	"strconv"

	libraryDTOs "jethur/grc/library/dtos"
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

	result, err := clc.service.GetClausesByVersion(uint(versionID), uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (clc *ComplianceLibraryController) GetAllMappings(c *gin.Context) {
	data, err := clc.service.GetAllControlMapping()
	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (clc *ComplianceLibraryController) GetControlMappings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	export := false
	if v, exists := c.Get("export"); exists {
		export = v.(bool)
	} else if s := c.Query("export"); s != "" {
		export, _ = strconv.ParseBool(s)
	}
	page, _ := utility.ParseQueryParamToInt(c, "page", 1)
	pageSize, _ := utility.ParseQueryParamToInt(c, "page_size", 10)

	params := libraryDTOs.ControlMappingListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Export:   export,
	}

	result, err := clc.service.GetControlMappings(&params)
	if err != nil {
		utility.SendInternalErrResp(clc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, _ := utility.GetCurrentUserID(c)
		file, err := clc.service.ExportControlMappings(result.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(clc.db, c, err.Error())
			return
		}
		defer file.Close()

		c.Header("Content-Disposition", `attachment; filename="Control_Mapping.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(file.Name())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (clc *ComplianceLibraryController) GetControlMappingsExport(c *gin.Context) {
	c.Set("export", true)
	clc.GetControlMappings(c)
}
