package controllers

import (
	"net/http"

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
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	RiskLibrary, err := rlc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(rlc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, RiskLibrary)
}

func (rlc *RiskLibraryController) GetRiskCategoryCounts(c *gin.Context) {
	result, err := rlc.service.GetRiskCategoryCounts()
	if err != nil {
		utility.SendInternalErrResp(rlc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
