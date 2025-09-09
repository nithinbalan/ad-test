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

type ArticleLibraryController struct {
	service services.ArticleLibraryServiceInterface
	db      *gorm.DB
}

func NewArticleLibraryController(service services.ArticleLibraryServiceInterface) *ArticleLibraryController {
	return &ArticleLibraryController{
		service: service,
	}
}

func (alc *ArticleLibraryController) GetAll(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"article_category_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, _ := utility.ParseQueryParamToInt(c, "page", 1)
	pageSize, _ := utility.ParseQueryParamToInt(c, "page_size", 10)

	params := libraryDTOs.ArticleLibraryListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Filters:     filters,
		QuickFilter: quickFilter,
		Page:        page,
		PageSize:    pageSize,
	}

	StrategicObjective, err := alc.service.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(alc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, StrategicObjective)
}

func (alc *ArticleLibraryController) GetArticleLibraryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	Article, err := alc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(alc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, Article)
}

func (alc *ArticleLibraryController) ArticleCategories(c *gin.Context) {
	result, err := alc.service.GetArticleCategories()
	if err != nil {
		utility.SendInternalErrResp(alc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
