package controllers

import (
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentChangeRequestStatusController struct {
	service services.DocumentChangeRequestStatusServiceInterface
	db      *gorm.DB
}

func NewDocumentChangeRequestStatusController(service services.DocumentChangeRequestStatusServiceInterface) *DocumentChangeRequestStatusController {
	return &DocumentChangeRequestStatusController{
		service: service,
	}
}

func (dcrsc *DocumentChangeRequestStatusController) GetAllDocumentChangeRequestStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	documentChangeRequestStatuses, err := dcrsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(dcrsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentChangeRequestStatuses)
}

func (dcrsc *DocumentChangeRequestStatusController) GetDocumentChangeRequestStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentChangeRequestStatus, err := dcrsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, documentChangeRequestStatus)
}
