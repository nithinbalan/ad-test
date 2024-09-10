package controllers

import (
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentStatusController struct {
	service services.DocumentStatusServiceInterface
	db      *gorm.DB
}

func NewDocumentStatusController(service services.DocumentStatusServiceInterface) *DocumentStatusController {
	return &DocumentStatusController{
		service: service,
	}
}

func (dsc *DocumentStatusController) GetAllDocumentStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	documentStatuses, err := dsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(dsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentStatuses)
}

func (dsc *DocumentStatusController) GetDocumentStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentStatus, err := dsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, documentStatus)
}
