package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskResponseReviewController struct {
	service services.RiskResponseReviewServiceInterface
	db      *gorm.DB
}

func NewRiskResponseReviewController(service services.RiskResponseReviewServiceInterface) *RiskResponseReviewController {
	return &RiskResponseReviewController{
		service: service,
	}
}

func (rrrc *RiskResponseReviewController) SubmitRiskResponse(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrrc.db, c, err.Error())
		return
	}

	err = rrrc.service.Submit(uint(riskResponseID), currentUserID)
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
			utility.SendInternalErrResp(rrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrrc *RiskResponseReviewController) ApproveRiskResponse(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrrc.db, c, err.Error())
		return
	}

	err = rrrc.service.Approve(uint(riskResponseID), currentUserID)
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
			utility.SendInternalErrResp(rrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrrc *RiskResponseReviewController) RejectRiskResponse(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrrc.db, c, err.Error())
		return
	}

	err = rrrc.service.Reject(uint(riskResponseID), currentUserID)
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
			utility.SendInternalErrResp(rrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrrc *RiskResponseReviewController) ReopenRiskResponse(c *gin.Context) {
	riskResponseIDStr := c.Param("id")
	riskResponseID, err := strconv.ParseUint(riskResponseIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrrc.db, c, err.Error())
		return
	}

	err = rrrc.service.Reopen(uint(riskResponseID), currentUserID)
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
			utility.SendInternalErrResp(rrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
