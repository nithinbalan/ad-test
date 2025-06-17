package controllers

import (
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EntryTypeController struct {
	service services.EntryTypeServiceInterface
	db      *gorm.DB
}

func NewEntryTypeController(service services.EntryTypeServiceInterface) *EntryTypeController {
	return &EntryTypeController{
		service: service,
	}
}

func (dsc *EntryTypeController) GetAllEntryTypees(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	entryTypees, err := dsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(dsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, entryTypees)
}

func (dsc *EntryTypeController) GetEntryTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	entryType, err := dsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, entryType)
}
