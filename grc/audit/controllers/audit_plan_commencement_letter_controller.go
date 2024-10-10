package controllers

import (
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditPlanCommencementLetterController struct {
	service services.AuditPlanCommencementLetterServiceInterface
	db      *gorm.DB
}

func NewAuditPlanCommencementLetterController(service services.AuditPlanCommencementLetterServiceInterface) *AuditPlanCommencementLetterController {
	return &AuditPlanCommencementLetterController{
		service: service,
	}
}

func (apc *AuditPlanCommencementLetterController) CreateAuditPlanCommencementLetters(c *gin.Context) {

	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	err = apc.service.Create(uint(auditPlanID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
