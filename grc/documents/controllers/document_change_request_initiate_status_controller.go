package controllers

import (
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentChangeRequestInitiateStatusController struct {
	service services.DocumentChangeRequestInitiateStatusServiceInterface
	db      *gorm.DB
}

func NewDocumentChangeRequestInitiateStatusController(service services.DocumentChangeRequestInitiateStatusServiceInterface) *DocumentChangeRequestInitiateStatusController {
	return &DocumentChangeRequestInitiateStatusController{
		service: service,
	}
}

func (dcris *DocumentChangeRequestInitiateStatusController) GetAllDocumentChangeRequestInitiateStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	documentChangeRequestInitiateStatuses, err := dcris.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(dcris.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentChangeRequestInitiateStatuses)
}

func (dcris *DocumentChangeRequestInitiateStatusController) GetDocumentChangeRequestInitiateStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentChangeRequestInitiateStatus, err := dcris.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, documentChangeRequestInitiateStatus)
}
