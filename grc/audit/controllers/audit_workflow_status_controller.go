package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditWorkflowStatusController struct {
	service services.AuditWorkflowStatusServiceInterface
	db      *gorm.DB
}

func NewAuditWorkflowStatusController(service services.AuditWorkflowStatusServiceInterface) *AuditWorkflowStatusController {
	return &AuditWorkflowStatusController{
		service: service,
	}
}

func (awsc *AuditWorkflowStatusController) GetAllAuditWorkflowStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	auditDraftReportStatuses, err := awsc.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(awsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditDraftReportStatuses)
}

func (awsc *AuditWorkflowStatusController) GetAuditWorkflowStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditDraftReportStatus, err := awsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, auditDraftReportStatus)
}
