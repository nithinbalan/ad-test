package controllers

import (
	"jethur/grc/awareness-reporting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessCampaignStatusController struct {
	service services.AwarenessCampaignStatusServiceInterface
	db      *gorm.DB
}

func NewAwarenessCampaignStatusController(service services.AwarenessCampaignStatusServiceInterface,) *AwarenessCampaignStatusController {
	return &AwarenessCampaignStatusController{
		service: service,
	}
}

func (acsc *AwarenessCampaignStatusController) GetAllAwarenessCampaignStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	awarenessCampaignStatuses, err := acsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(acsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, awarenessCampaignStatuses)
}

func (acsc *AwarenessCampaignStatusController) GetAwarenessCampaignStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	awarenessCampaignStatus, err := acsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, awarenessCampaignStatus)
}