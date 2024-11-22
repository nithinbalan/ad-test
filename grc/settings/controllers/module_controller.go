package controllers

import (
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"strconv"

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

func (mc *ModuleController) EnableModule(c *gin.Context) {
	moduleId, err := strconv.Atoi(c.Param("moduleId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	_, err = mc.service.EnableModule(uint(moduleId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mc.db, c, e.Error())
			return
		}
	}
}

func (mc *ModuleController) DisableModule(c *gin.Context) {
	moduleId, err := strconv.Atoi(c.Param("moduleId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	_, err = mc.service.DisableModule(uint(moduleId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mc.db, c, e.Error())
			return
		}
	}
}
