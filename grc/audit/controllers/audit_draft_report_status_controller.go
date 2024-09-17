package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditDraftReportStatusController struct {
	service services.AuditDraftReportStatusServiceInterface
	db      *gorm.DB
}

func NewAuditDraftReportStatusController(service services.AuditDraftReportStatusServiceInterface) *AuditDraftReportStatusController {
	return &AuditDraftReportStatusController{
		service: service,
	}
}

func (adrsc *AuditDraftReportStatusController) GetAllAuditDraftReportStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	auditDraftReportStatuses, err := adrsc.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(adrsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditDraftReportStatuses)
}

func (adrsc *AuditDraftReportStatusController) GetAuditDraftReportStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditDraftReportStatus, err := adrsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, auditDraftReportStatus)
}
