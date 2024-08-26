package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditFindingStatusController struct {
	service services.AuditFindingStatusServiceInterface
	db      *gorm.DB
}

func NewAuditFindingStatusController(service services.AuditFindingStatusServiceInterface) *AuditFindingStatusController {
	return &AuditFindingStatusController{
		service: service,
	}
}

func (afsc *AuditFindingStatusController) GetAllAuditFindingStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	auditFindingStatuses, err := afsc.service.GetAll(sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(afsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditFindingStatuses)
}

func (afsc *AuditFindingStatusController) GetAuditFindingStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditFindingStatus, err := afsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, auditFindingStatus)
}
