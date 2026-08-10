package controllers

import (
	"jethur/grc/identity-access-review/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewCampaignStatusController struct {
	service services.ReviewCampaignStatusServiceInterface
}

func NewReviewCampaignStatusController(service services.ReviewCampaignStatusServiceInterface) *ReviewCampaignStatusController {
	return &ReviewCampaignStatusController{
		service: service,
	}
}

func (rcsc *ReviewCampaignStatusController) GetAllReviewCampaignStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	reviewCampaignStatuses, err := rcsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reviewCampaignStatuses)
}

func (rcsc *ReviewCampaignStatusController) GetReviewCampaignStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reviewCampaignStatus, err := rcsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, reviewCampaignStatus)
}
