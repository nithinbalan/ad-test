package controllers

import (
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentChangeRequestInitiateTypeController struct {
	service services.DocumentChangeRequestInitiateTypeServiceInterface
	db      *gorm.DB
}

func NewDocumentChangeRequestInitiateTypeController(service services.DocumentChangeRequestInitiateTypeServiceInterface) *DocumentChangeRequestInitiateTypeController {
	return &DocumentChangeRequestInitiateTypeController{
		service: service,
	}
}

func (dcritc *DocumentChangeRequestInitiateTypeController) GetAllDocumentChangeRequestInitiateTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	changeRequestInitiateTypes, err := dcritc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(dcritc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, changeRequestInitiateTypes)
}

func (dcritc *DocumentChangeRequestInitiateTypeController) GetDocumentChangeRequestInitiateTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestInitiateType, err := dcritc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, changeRequestInitiateType)
}
