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

type KpiLibraryController struct {
	service services.KpiLibraryServiceInterface
	db      *gorm.DB
}

func NewKpiLibraryController(service services.KpiLibraryServiceInterface) *KpiLibraryController {
	return &KpiLibraryController{
		service: service,
	}
}

func (klc *KpiLibraryController) GetAllKpi(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	export, _ := strconv.ParseBool(c.DefaultQuery("export", "false"))
	page, _ := utility.ParseQueryParamToInt(c, "page", 1)
	pageSize, _ := utility.ParseQueryParamToInt(c, "page_size", 10)
	categories := c.QueryArray("categories")

	params := libraryDTOs.KPILibraryListQueryParamsDTO{
		Sort:       sort,
		SortBy:     sortBy,
		Keyword:    keyword,
		Page:       page,
		PageSize:   pageSize,
		Export:     export,
		Categories: categories,
	}

	kpis, err := klc.service.GetAllKpi(&params)
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, _ := utility.GetCurrentUserID(c)
		excelFile, err := klc.service.Export(kpis.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(klc.db, c, err.Error())
			return
		}
		defer excelFile.Close()
		c.Header("Content-Disposition", `attachment; filename="KPI_Library.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, kpis)
}

func (klc *KpiLibraryController) GetKpiExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"
	klc.GetAllKpi(c)
}

func (klc *KpiLibraryController) KPICountByCategory(c *gin.Context) {
	result, err := klc.service.KPIGetCategoryCounts()
	if err != nil {
		utility.SendInternalErrResp(klc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
