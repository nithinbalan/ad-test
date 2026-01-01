package controllers

import (
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyAssessmentReviewController struct {
	service services.ThirdPartyAssessmentReviewServiceInterface
	db      *gorm.DB
}

func NewThirdPartyAssessmentReviewController(service services.ThirdPartyAssessmentReviewServiceInterface) *ThirdPartyAssessmentReviewController {
	return &ThirdPartyAssessmentReviewController{
		service: service,
	}
}

func (trc *ThirdPartyAssessmentReviewController) SubmitAssessment(c *gin.Context) {
	assessmentIdStr := c.Param("assessmentID")
	assessmentId, err := strconv.ParseUint(assessmentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(trc.db, c, err.Error())
		return
	}

	err = trc.service.Submit(uint(assessmentId), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(trc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrc *ThirdPartyAssessmentReviewController) ApproveAssessment(c *gin.Context) {
	assessmentIdStr := c.Param("assessmentID")
	assessmentId, err := strconv.ParseUint(assessmentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.Approve(uint(assessmentId), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrc *ThirdPartyAssessmentReviewController) RejectAssessment(c *gin.Context) {
	assessmentIdStr := c.Param("assessmentID")
	assessmentId, err := strconv.ParseUint(assessmentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.Reject(uint(assessmentId), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrc *ThirdPartyAssessmentReviewController) ReopenAssessment(c *gin.Context) {
	assessmentIdStr := c.Param("assessmentID")
	assessmentId, err := strconv.ParseUint(assessmentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.Reopen(uint(assessmentId), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
