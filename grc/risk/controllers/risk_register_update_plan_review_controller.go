package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanReviewController struct {
	service services.RiskRegisterUpdatePlanReviewServiceInterface
	db      *gorm.DB
}

func NewRiskRegisterUpdatePlanReviewController(service services.RiskRegisterUpdatePlanReviewServiceInterface) *RiskRegisterUpdatePlanReviewController {
	return &RiskRegisterUpdatePlanReviewController{
		service: service,
	}
}

func (rruprc *RiskRegisterUpdatePlanReviewController) SubmitRiskRegisterUpdatePlan(c *gin.Context) {
	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rruprc.db, c, err.Error())
		return
	}

	err = rruprc.service.Submit(uint(riskRegisterUpdatePlanID), currentUserID)
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
			utility.SendInternalErrResp(rruprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rruprc *RiskRegisterUpdatePlanReviewController) ApproveRiskRegisterUpdatePlan(c *gin.Context) {
	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rruprc.db, c, err.Error())
		return
	}

	err = rruprc.service.Approve(uint(riskRegisterUpdatePlanID), currentUserID)
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
			utility.SendInternalErrResp(rruprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rruprc *RiskRegisterUpdatePlanReviewController) RejectRiskRegisterUpdatePlan(c *gin.Context) {
	riskRegisterUpdatePlanIDStr := c.Param("id")
	riskRegisterUpdatePlanID, err := strconv.ParseUint(riskRegisterUpdatePlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rruprc.db, c, err.Error())
		return
	}

	err = rruprc.service.Reject(uint(riskRegisterUpdatePlanID), currentUserID)
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
			utility.SendInternalErrResp(rruprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
