package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditStatusController struct {
	service services.AuditStatusServiceInterface
	db      *gorm.DB
}

func NewAuditStatusController(service services.AuditStatusServiceInterface) *AuditStatusController {
	return &AuditStatusController{
		service: service,
	}
}

func (asc *AuditStatusController) GetAllAuditStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	auditStatuses, err := asc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditStatuses)
}

func (asc *AuditStatusController) GetAuditStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditStatus, err := asc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, auditStatus)
}
