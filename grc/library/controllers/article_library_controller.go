package controllers

import (
	"net/http"

	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleLibraryController struct {
	service services.ArticleLibraryServiceInterface
	db      *gorm.DB
}

func NewArticleLibraryController(service services.ArticleLibraryServiceInterface) *ArticleLibraryController {
	return &ArticleLibraryController{
		service: service,
	}
}

func (vlc *ArticleLibraryController) GetAll(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, _ := utility.ParseQueryParamToInt(c, "page", 1)
	pageSize, _ := utility.ParseQueryParamToInt(c, "page_size", 10)

	params := libraryDTOs.ArticleLibraryListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	StrategicObjective, err := vlc.service.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(vlc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, StrategicObjective)
}
