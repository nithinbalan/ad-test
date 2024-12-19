package controllers

import (
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentChangeRequestWorkflowStatusController struct {
	service services.DocumentChangeRequestWorkflowStatusServiceInterface
	db      *gorm.DB
}

func NewDocumentChangeRequestWorkflowStatusController(service services.DocumentChangeRequestWorkflowStatusServiceInterface) *DocumentChangeRequestWorkflowStatusController {
	return &DocumentChangeRequestWorkflowStatusController{
		service: service,
	}
}

func (dcrwsc *DocumentChangeRequestWorkflowStatusController) GetAllDocumentChangeRequestWorkflowStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	statuses, err := dcrwsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(dcrwsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, statuses)
}

func (dcrwsc *DocumentChangeRequestWorkflowStatusController) GetDocumentChangeRequestWorkflowStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	status, err := dcrwsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, status)
}
