package controllers

import (
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"

	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModuleController struct {
	service services.ModuleServiceInterface
	db      *gorm.DB
}

func NewModuleController(service services.ModuleServiceInterface) *ModuleController {
	return &ModuleController{
		service: service,
	}
}

func (mc *ModuleController) GetAllModules(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")

	queryParams := dtos.ModuleListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}

	modules, err := mc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(mc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, modules)
}
