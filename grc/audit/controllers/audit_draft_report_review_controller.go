package controllers

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditDraftReportReviewController struct {
	service services.AuditDraftReportReviewServiceInterface
	db      *gorm.DB
}

func NewAuditDraftReportReviewController(service services.AuditDraftReportReviewServiceInterface) *AuditDraftReportReviewController {
	return &AuditDraftReportReviewController{
		service: service,
	}
}

func (adrrc *AuditDraftReportReviewController) SubmitAuditDraftReport(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	auditDraftReportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrrc.db, c, err.Error())
		return
	}
	var auditReportWorkflowFormDTO dtos.AuditReportWorkflowSubmitFormDTO

	if err := c.ShouldBindJSON(&auditReportWorkflowFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = adrrc.service.Submit(uint(auditID), uint(auditDraftReportID), currentUserID, &auditReportWorkflowFormDTO)
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
			utility.SendInternalErrResp(adrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (adrrc *AuditDraftReportReviewController) ApproveAuditDraftReport(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	auditDraftReportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrrc.db, c, err.Error())
		return
	}
	var auditReportWorkflowFormDTO dtos.AuditReportWorkflowApproveFormDTO

	if err := c.ShouldBindJSON(&auditReportWorkflowFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = adrrc.service.Approve(uint(auditID), uint(auditDraftReportID), currentUserID, &auditReportWorkflowFormDTO)
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
			utility.SendInternalErrResp(adrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (adrrc *AuditDraftReportReviewController) RevertAuditDraftReport(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	auditDraftReportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrrc.db, c, err.Error())
		return
	}
	var auditReportWorkflowFormDTO dtos.AuditReportWorkflowApproveFormDTO

	if err := c.ShouldBindJSON(&auditReportWorkflowFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = adrrc.service.Revert(uint(auditID), uint(auditDraftReportID), currentUserID, &auditReportWorkflowFormDTO)
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
			utility.SendInternalErrResp(adrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
