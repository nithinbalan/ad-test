package controllers

import (
	"jethur/grc/audit-insights/dtos"
	"jethur/grc/audit-insights/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditInsightsDraftFindingController struct {
	service services.AuditInsightsDraftFindingServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsDraftFindingController(service services.AuditInsightsDraftFindingServiceInterface) *AuditInsightsDraftFindingController {
	return &AuditInsightsDraftFindingController{
		service: service,
	}
}

func (aidfc *AuditInsightsDraftFindingController) AcceptAuditDraftFinding(c *gin.Context) {
	auditDraftFindingID, err := strconv.Atoi(c.Param("findingId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aidfc.db, c, err.Error())
		return
	}
	var auditInsightsDraftFindingApproveFormDTO dtos.AuditInsightsDraftFindingApproveFormDTO

	if err := c.ShouldBindJSON(&auditInsightsDraftFindingApproveFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = aidfc.service.Accept(uint(auditDraftFindingID), currentUserID, &auditInsightsDraftFindingApproveFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(aidfc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (aidfc *AuditInsightsDraftFindingController) ChallengeAuditDraftFinding(c *gin.Context) {
	auditDraftFindingID, err := strconv.Atoi(c.Param("findingId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aidfc.db, c, err.Error())
		return
	}
	var auditInsightsDraftFindingChallengeFormDTO dtos.AuditInsightsDraftFindingChallengeFormDTO

	if err := c.ShouldBindJSON(&auditInsightsDraftFindingChallengeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = aidfc.service.Challenge(uint(auditDraftFindingID), currentUserID, &auditInsightsDraftFindingChallengeFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(aidfc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
