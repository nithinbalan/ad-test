package controllers

import (
	"jethur/grc/kpi/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IdentityAccessReviewDashboardController struct {
	service services.IdentityAccessReviewDashboardServiceInterface
	db      *gorm.DB
}

func NewIdentityAccessReviewDashboardController(service services.IdentityAccessReviewDashboardServiceInterface, db *gorm.DB) *IdentityAccessReviewDashboardController {
	return &IdentityAccessReviewDashboardController{service: service, db: db}
}

func (controller *IdentityAccessReviewDashboardController) GetIdentityAccessReviewSummary(c *gin.Context) {
	response, err := controller.service.GetIdentityAccessReviewSummary()
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *IdentityAccessReviewDashboardController) GetCampaignsByLifecycle(c *gin.Context) {
	response, err := controller.service.GetCampaignsByLifecycle()
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *IdentityAccessReviewDashboardController) GetCompletionBySystem(c *gin.Context) {
	response, err := controller.service.GetCompletionBySystem()
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *IdentityAccessReviewDashboardController) GetActiveAndRecentCampaigns(c *gin.Context) {
	response, err := controller.service.GetActiveAndRecentCampaigns()
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
