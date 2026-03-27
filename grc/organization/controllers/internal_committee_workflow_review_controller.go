package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InternalCommitteeWorkflowReviewController struct {
	service services.InternalCommitteeWorkflowReviewServiceInterface
	db      *gorm.DB
}

func NewInternalCommitteeWorkflowReviewController(service services.InternalCommitteeWorkflowReviewServiceInterface) *InternalCommitteeWorkflowReviewController {
	return &InternalCommitteeWorkflowReviewController{
		service: service,
	}
}

func (icwc *InternalCommitteeWorkflowReviewController) SubmitInternalCommittee(c *gin.Context) {
	committeeIdStr := c.Param("committeeId")
	committeeId, err := strconv.ParseUint(committeeIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icwc.db, c, err.Error())
		return
	}

	err = icwc.service.Submit(uint(committeeId), currentUserID)
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
			utility.SendInternalErrResp(icwc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (icwc *InternalCommitteeWorkflowReviewController) ApproveInternalCommittee(c *gin.Context) {
	committeeIdStr := c.Param("committeeId")
	committeeId, err := strconv.ParseUint(committeeIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icwc.db, c, err.Error())
		return
	}

	err = icwc.service.Approve(uint(committeeId), currentUserID)
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
			utility.SendInternalErrResp(icwc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (icwc *InternalCommitteeWorkflowReviewController) RevertInternalCommittee(c *gin.Context) {
	committeeIdStr := c.Param("committeeId")
	committeeId, err := strconv.ParseUint(committeeIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icwc.db, c, err.Error())
		return
	}

	err = icwc.service.Revert(uint(committeeId), currentUserID)
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
			utility.SendInternalErrResp(icwc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
