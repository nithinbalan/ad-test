package controllers

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditReportReviewController struct {
	service services.AuditReportReviewServiceInterface
	db      *gorm.DB
}

func NewAuditReportReviewController(service services.AuditReportReviewServiceInterface) *AuditReportReviewController {
	return &AuditReportReviewController{
		service: service,
	}
}

func (arrc *AuditReportReviewController) SubmitAuditFinalReport(c *gin.Context) {
	auditIdStr := c.Param("auditId")
	auditId, err := strconv.ParseUint(auditIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	err = arrc.service.Submit(uint(auditId), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arrc *AuditReportReviewController) GetComments(c *gin.Context) {
	auditIdStr := c.Param("auditId")
	auditId, err := strconv.ParseUint(auditIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := arrc.service.GetComments(uint(auditId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (arrc *AuditReportReviewController) AddComment(c *gin.Context) {
	var auditReportReviewCommentFormDTO auditDTOs.AuditReportReviewCommentFormDTO
	if err := c.ShouldBindJSON(&auditReportReviewCommentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	auditIdStr := c.Param("auditId")
	auditId, err := strconv.ParseUint(auditIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	err = arrc.service.AddComment(uint(auditId), currentUserID, &auditReportReviewCommentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(arrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
