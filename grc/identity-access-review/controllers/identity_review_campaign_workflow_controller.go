package controllers

import (
	"jethur/grc/identity-access-review/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IdentityReviewCampaignWorkflowController struct {
	service services.IdentityReviewCampaignWorkflowServiceInterface
}

func NewIdentityReviewCampaignWorkflowController(service services.IdentityReviewCampaignWorkflowServiceInterface) *IdentityReviewCampaignWorkflowController {
	return &IdentityReviewCampaignWorkflowController{service: service}
}

func (ircwc *IdentityReviewCampaignWorkflowController) GetWorkflowList(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	items, err := ircwc.service.GetWorkflowList(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
		}
		return
	}
	c.JSON(http.StatusOK, items)
}
