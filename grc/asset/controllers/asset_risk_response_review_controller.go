package controllers

import (
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRiskResponseReviewController struct {
	service services.AssetRiskResponseReviewServiceInterface
	db      *gorm.DB
}

func NewAssetRiskResponseReviewController(service services.AssetRiskResponseReviewServiceInterface) *AssetRiskResponseReviewController {
	return &AssetRiskResponseReviewController{
		service: service,
	}
}

func (arrc *AssetRiskResponseReviewController) SubmitAssetRiskResponse(c *gin.Context) {
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

func (arrc *AssetRiskResponseReviewController) ApproveAssetRiskResponse(c *gin.Context) {
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

func (arrc *AssetRiskResponseReviewController) RevertAssetRiskResponse(c *gin.Context) {
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
