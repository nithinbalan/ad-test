package controllers

import (
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentWorkflowStatusController struct {
	service services.DocumentWorkflowStatusServiceInterface
	db      *gorm.DB
}

func NewDocumentWorkflowStatusController(service services.DocumentWorkflowStatusServiceInterface) *DocumentWorkflowStatusController {
	return &DocumentWorkflowStatusController{
		service: service,
	}
}

func (dwsc *DocumentWorkflowStatusController) GetAllDocumentWorkflowStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	statuses, err := dwsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(dwsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, statuses)
}

func (dwsc *DocumentWorkflowStatusController) GetDocumentWorkflowStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	status, err := dwsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, status)
}
