package controllers

import (
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskReviewController struct {
	service services.RiskReviewServiceInterface
	db      *gorm.DB
}

func NewRiskReviewController(service services.RiskReviewServiceInterface) *RiskReviewController {
	return &RiskReviewController{
		service: service,
	}
}

func (rrc *RiskReviewController) SubmitRisk(c *gin.Context) {
	riskIdStr := c.Param("id")
	riskId, err := strconv.ParseUint(riskIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.Submit(uint(riskId), currentUserID)
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

func (rrc *RiskReviewController) ApproveRisk(c *gin.Context) {
	riskIdStr := c.Param("id")
	riskId, err := strconv.ParseUint(riskIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.Approve(uint(riskId), currentUserID)
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

func (rrc *RiskReviewController) RejectRisk(c *gin.Context) {
	riskIdStr := c.Param("id")
	riskId, err := strconv.ParseUint(riskIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.Reject(uint(riskId), currentUserID)
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
