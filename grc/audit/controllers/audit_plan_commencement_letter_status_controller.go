package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/audit/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditPlanCommencementLetterStatusController struct {
	service services.AuditPlanCommencementLetterStatusServiceInterface
	db      *gorm.DB
}

func NewAuditPlanCommencementLetterStatusController(service services.AuditPlanCommencementLetterStatusServiceInterface) *AuditPlanCommencementLetterStatusController {
	return &AuditPlanCommencementLetterStatusController{
		service: service,
	}
}

// GetAllStatuses handles fetching all audit plan commencement letter statuses
func (aps *AuditPlanCommencementLetterStatusController) GetAllStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	statuses, err := aps.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(aps.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, statuses)
}

// GetStatusDetails handles fetching the details of a specific audit plan commencement letter status by ID
func (aps *AuditPlanCommencementLetterStatusController) GetStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	status, err := aps.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, status)
}
