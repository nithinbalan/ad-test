package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommunicationActionPlanStatusController struct {
	service services.CommunicationActionPlanStatusServiceInterface
	db      *gorm.DB
}

func NewCommunicationActionPlanStatusController(service services.CommunicationActionPlanStatusServiceInterface) *CommunicationActionPlanStatusController {
	return &CommunicationActionPlanStatusController{
		service: service,
	}
}

func (tesc *CommunicationActionPlanStatusController) GetAllCommunicationActionPlanStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	communicationActionPlanStatuses, err := tesc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tesc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, communicationActionPlanStatuses)
}

func (tesc *CommunicationActionPlanStatusController) GetCommunicationActionPlanStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	communicationActionPlanStatus, err := tesc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, communicationActionPlanStatus)
}
