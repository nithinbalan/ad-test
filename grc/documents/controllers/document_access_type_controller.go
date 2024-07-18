package controllers

import (
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentAccessTypeController struct {
	service services.DocumentAccessTypeServiceInterface
	db      *gorm.DB
}

func NewDocumentAccessTypeController(service services.DocumentAccessTypeServiceInterface) *DocumentAccessTypeController {
	return &DocumentAccessTypeController{
		service: service,
	}
}

func (dac *DocumentAccessTypeController) GetAllDocumentAccessTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	accessTypes, err := dac.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(dac.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, accessTypes)
}

func (dac *DocumentAccessTypeController) GetDocumentAccessTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	accessType, err := dac.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, accessType)
}
