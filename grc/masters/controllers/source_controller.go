package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SourceController struct {
	service services.SourceServiceInterface
	db      *gorm.DB
}

func NewSourceController(service services.SourceServiceInterface) *SourceController {
	return &SourceController{
		service: service,
	}
}

func (sc *SourceController) GetAllSources(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	sources, err := sc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(sc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, sources)
}

func (sc *SourceController) GetSourceDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInternalErrResp(sc.db, c, "Invalid id")
		return
	}

	source, err := sc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, source)
}
