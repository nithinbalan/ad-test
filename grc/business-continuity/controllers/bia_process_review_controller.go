package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BiaProcessReviewController struct {
	service services.BiaProcessReviewServiceInterface
	db      *gorm.DB
}

func NewBiaProcessReviewController(service services.BiaProcessReviewServiceInterface) *BiaProcessReviewController {
	return &BiaProcessReviewController{
		service: service,
	}
}

func (bprc *BiaProcessReviewController) SubmitBiaProcess(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bprc.db, c, err.Error())
		return
	}

	if err := bprc.service.Submit(uint(biaProcessID), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(bprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (bprc *BiaProcessReviewController) ApproveBiaProcess(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bprc.db, c, err.Error())
		return
	}

	if err := bprc.service.Approve(uint(biaProcessID), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(bprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (bprc *BiaProcessReviewController) RejectBiaProcess(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bprc.db, c, err.Error())
		return
	}

	if err := bprc.service.Reject(uint(biaProcessID), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(bprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (bprc *BiaProcessReviewController) ReopenBiaProcess(c *gin.Context) {
	biaProcessIDStr := c.Param("id")
	biaProcessID, err := strconv.ParseUint(biaProcessIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bprc.db, c, err.Error())
		return
	}

	if err := bprc.service.Reopen(uint(biaProcessID), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(bprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
