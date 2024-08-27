package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditDocumentRequestStatusController struct {
	service services.AuditDocumentRequestStatusServiceInterface
	db      *gorm.DB
}

func NewAuditDocumentRequestStatusController(service services.AuditDocumentRequestStatusServiceInterface) *AuditDocumentRequestStatusController {
	return &AuditDocumentRequestStatusController{
		service: service,
	}
}

func (adrs *AuditDocumentRequestStatusController) GetAllAuditDocumentRequestStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	auditDocumentRequestStatuses, err := adrs.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(adrs.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditDocumentRequestStatuses)
}

func (adrs *AuditDocumentRequestStatusController) GetAuditDocumentRequestStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditStatus, err := adrs.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, auditStatus)
}
