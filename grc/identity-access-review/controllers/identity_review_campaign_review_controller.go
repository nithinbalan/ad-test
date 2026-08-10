package controllers

import (
	"jethur/grc/identity-access-review/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IdentityReviewCampaignReviewController struct {
	service services.IdentityReviewCampaignReviewServiceInterface
	db      *gorm.DB
}

func NewIdentityReviewCampaignReviewController(service services.IdentityReviewCampaignReviewServiceInterface) *IdentityReviewCampaignReviewController {
	return &IdentityReviewCampaignReviewController{
		service: service,
	}
}

func (ircrc *IdentityReviewCampaignReviewController) SubmitIdentityReviewCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircrc.db, c, err.Error())
		return
	}

	err = ircrc.service.Submit(uint(campaignID), currentUserID)
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
			utility.SendInternalErrResp(ircrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ircrc *IdentityReviewCampaignReviewController) ApproveIdentityReviewCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircrc.db, c, err.Error())
		return
	}

	err = ircrc.service.Approve(uint(campaignID), currentUserID)
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
			utility.SendInternalErrResp(ircrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ircrc *IdentityReviewCampaignReviewController) RejectIdentityReviewCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircrc.db, c, err.Error())
		return
	}

	err = ircrc.service.Reject(uint(campaignID), currentUserID)
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
			utility.SendInternalErrResp(ircrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ircrc *IdentityReviewCampaignReviewController) ResubmitIdentityReviewCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircrc.db, c, err.Error())
		return
	}

	err = ircrc.service.Resubmit(uint(campaignID), currentUserID)
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
			utility.SendInternalErrResp(ircrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ircrc *IdentityReviewCampaignReviewController) RevertIdentityReviewCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircrc.db, c, err.Error())
		return
	}

	err = ircrc.service.Revert(uint(campaignID), currentUserID)
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
			utility.SendInternalErrResp(ircrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
