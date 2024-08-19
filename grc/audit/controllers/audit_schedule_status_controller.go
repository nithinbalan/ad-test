package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditScheduleStatusController struct {
	service services.AuditScheduleStatusServiceInterface
	db      *gorm.DB
}

func NewAuditScheduleStatusController(service services.AuditScheduleStatusServiceInterface) *AuditScheduleStatusController {
	return &AuditScheduleStatusController{
		service: service,
	}
}

func (assc *AuditScheduleStatusController) GetAllAuditScheduleStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	auditScheduleStatuses, err := assc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(assc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditScheduleStatuses)
}

func (assc *AuditScheduleStatusController) GetAuditScheduleStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditScheduleStatus, err := assc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, auditScheduleStatus)
}
