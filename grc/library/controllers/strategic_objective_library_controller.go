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

type StrategicObjectiveLibraryController struct {
	service services.StrategicObjectiveLibraryServiceInterface
	db      *gorm.DB
}

func NewStrategicObjectiveLibraryController(service services.StrategicObjectiveLibraryServiceInterface) *StrategicObjectiveLibraryController {
	return &StrategicObjectiveLibraryController{
		service: service,
	}
}

func (solc *StrategicObjectiveLibraryController) GetAll(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	export, _ := strconv.ParseBool(c.DefaultQuery("export", "false"))
	page, _ := utility.ParseQueryParamToInt(c, "page", 1)
	pageSize, _ := utility.ParseQueryParamToInt(c, "page_size", 10)
	categories := c.QueryArray("categories")

	params := libraryDTOs.StrategicObjectiveLibraryListQueryParamsDTO{
		Sort:       sort,
		SortBy:     sortBy,
		Keyword:    keyword,
		Page:       page,
		PageSize:   pageSize,
		Export:     export,
		Categories: categories,
	}

	StrategicObjective, err := solc.service.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(solc.db, c, err.Error())
		return
	}
	if export {
		currentUserID, _ := utility.GetCurrentUserID(c)
		excelFile, err := solc.service.Export(StrategicObjective.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(solc.db, c, err.Error())
			return
		}
		defer excelFile.Close()
		c.Header("Content-Disposition", `attachment; filename="strategic_objective_library.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, StrategicObjective)
}

func (solc *StrategicObjectiveLibraryController) Export(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"
	solc.GetAll(c)
}

func (solc *StrategicObjectiveLibraryController) CountByCategory(c *gin.Context) {
	result, err := solc.service.GetCategoryCounts()
	if err != nil {
		utility.SendInternalErrResp(solc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
