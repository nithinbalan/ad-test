package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditReportWorkflowController struct {
	service services.AuditReportWorkflowServiceInterface
	db      *gorm.DB
}

func NewAuditReportWorkflowController(service services.AuditReportWorkflowServiceInterface) *AuditReportWorkflowController {
	return &AuditReportWorkflowController{
		service: service,
	}
}

func (arwc *AuditReportWorkflowController) GetWorkflowList(c *gin.Context) {
	auditIDStr := c.Param("auditId")
	auditID, err := strconv.ParseUint(auditIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	workflows, err := arwc.service.GetAuditReportWorkflowList(uint(auditID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, workflows)
}
