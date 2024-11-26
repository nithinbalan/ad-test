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
	quickFilter := c.QueryArray("quick_filter")

	queryParams := dtos.ModuleListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		QuickFilter: quickFilter,
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

func (mc *ModuleController) ReorderModule(c *gin.Context) {

	var modules []dtos.ModuleReorderDTO
	if err := c.ShouldBindJSON(&modules); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	err := mc.service.ReorderModule(modules)
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

func (mc *ModuleController) ReorderSubModule(c *gin.Context) {
	var modules []dtos.ModuleReorderDTO
	if err := c.ShouldBindJSON(&modules); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	err := mc.service.ReorderSubModule(modules)
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

func (mc *ModuleController) UpdateModule(c *gin.Context) {

	var moduleFormDTO dtos.ModuleFormDTO

	moduleId, err := strconv.Atoi(c.Param("moduleId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&moduleFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = mc.service.UpdateModule(uint(moduleId), moduleFormDTO)
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

func (mc *ModuleController) UpdateSubModule(c *gin.Context) {

	var subModuleFormDTO dtos.SubModuleFormDTO

	moduleId, err := strconv.Atoi(c.Param("subModuleId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&subModuleFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = mc.service.UpdateSubModule(uint(moduleId), subModuleFormDTO)
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
