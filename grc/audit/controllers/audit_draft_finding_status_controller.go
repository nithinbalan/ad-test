package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditDraftFindingStatusController struct {
	service services.AuditDraftFindingStatusServiceInterface
	db      *gorm.DB
}

func NewAuditDraftFindingStatusController(service services.AuditDraftFindingStatusServiceInterface) *AuditDraftFindingStatusController {
	return &AuditDraftFindingStatusController{
		service: service,
	}
}

func (adfsc *AuditDraftFindingStatusController) GetAllAuditDraftFindingStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	auditDraftFindingStatuses, err := adfsc.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(adfsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditDraftFindingStatuses)
}

func (adfsc *AuditDraftFindingStatusController) GetAuditDraftFindingStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditDraftFindingStatus, err := adfsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, auditDraftFindingStatus)
}
