package controllers

import (
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentChangeRequestTypeController struct {
	service services.DocumentChangeRequestTypeServiceInterface
	db      *gorm.DB
}

func NewDocumentChangeRequestTypeController(service services.DocumentChangeRequestTypeServiceInterface) *DocumentChangeRequestTypeController {
	return &DocumentChangeRequestTypeController{
		service: service,
	}
}

func (dcrtc *DocumentChangeRequestTypeController) GetAllDocumentChangeRequestTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	changeRequestTypes, err := dcrtc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(dcrtc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, changeRequestTypes)
}

func (dcrtc *DocumentChangeRequestTypeController) GetDocumentChangeRequestTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestType, err := dcrtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, changeRequestType)
}
