package controllers

import (
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRiskReviewController struct {
	service services.AssetRiskReviewServiceInterface
	db      *gorm.DB
}

func NewAssetRiskReviewController(service services.AssetRiskReviewServiceInterface) *AssetRiskReviewController {
	return &AssetRiskReviewController{
		service: service,
	}
}

func (arrc *AssetRiskReviewController) SubmitAssetRiskRegister(c *gin.Context) {
	riskIdStr := c.Param("id")
	riskId, err := strconv.ParseUint(riskIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	err = arrc.service.Submit(uint(riskId), currentUserID)
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
			utility.SendInternalErrResp(arrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arrc *AssetRiskReviewController) ApproveAssetRiskRegister(c *gin.Context) {
	riskIdStr := c.Param("id")
	riskId, err := strconv.ParseUint(riskIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	err = arrc.service.Approve(uint(riskId), currentUserID)
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
			utility.SendInternalErrResp(arrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arrc *AssetRiskReviewController) RejectAssetRiskRegister(c *gin.Context) {
	riskIdStr := c.Param("id")
	riskId, err := strconv.ParseUint(riskIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	err = arrc.service.Reject(uint(riskId), currentUserID)
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
			utility.SendInternalErrResp(arrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arrc *AssetRiskReviewController) ResubmitAssetRiskRegister(c *gin.Context) {
	riskIdStr := c.Param("id")
	riskId, err := strconv.ParseUint(riskIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	err = arrc.service.Resubmit(uint(riskId), currentUserID)
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
			utility.SendInternalErrResp(arrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arrc *AssetRiskReviewController) RevertAssetRiskRegister(c *gin.Context) {
	riskIdStr := c.Param("id")
	riskId, err := strconv.ParseUint(riskIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arrc.db, c, err.Error())
		return
	}

	err = arrc.service.Revert(uint(riskId), currentUserID)
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
			utility.SendInternalErrResp(arrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
