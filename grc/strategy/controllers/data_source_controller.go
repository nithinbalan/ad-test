package controllers

import (
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DataSourceController struct {
	service services.DataSourceServiceInterface
	db      *gorm.DB
}

func NewDataSourceController(service services.DataSourceServiceInterface) *DataSourceController {
	return &DataSourceController{
		service: service,
	}
}

func (lsc *DataSourceController) GetAllDataSourcees(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	dataSourcees, err := lsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(lsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, dataSourcees)
}

func (lsc *DataSourceController) GetDataSourceDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	dataSource, err := lsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, dataSource)
}
