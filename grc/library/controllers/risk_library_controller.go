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

type RiskLibraryController struct {
	service services.RiskLibraryServiceInterface
	db      *gorm.DB
}

func NewRiskLibraryController(service services.RiskLibraryServiceInterface) *RiskLibraryController {
	return &RiskLibraryController{
		service: service,
	}
}

func (rlc *RiskLibraryController) GetAllRiskLibrary(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	categories := c.QueryArray("category")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(rlc.db, c, "Invalid export parameter")
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rlc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rlc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := libraryDTOs.RiskLibraryListQueryParamsDTO{
		Sort:       sort,
		SortBy:     sortBy,
		Keyword:    keyword,
		Categories: categories,
		Page:       page,
		PageSize:   pageSize,
		Export:     export,
	}
	RiskLibrary, err := rlc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(rlc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, _ := utility.GetCurrentUserID(c)
		excelFile, err := rlc.service.ExportRiskLibrary(RiskLibrary.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(rlc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Disposition", `attachment; filename="Risk_Library.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, RiskLibrary)
}

func (rlc *RiskLibraryController) RiskCountByCategory(c *gin.Context) {
	result, err := rlc.service.GetRiskCategoryCounts()
	if err != nil {
		utility.SendInternalErrResp(rlc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (rlc *RiskLibraryController) GetRiskLibraryExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"
	rlc.GetAllRiskLibrary(c)
}
