package controllers

import (
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlReviewStatusController struct {
	service services.ControlReviewStatusServiceInterface
	db      *gorm.DB
}

func NewControlReviewStatusController(service services.ControlReviewStatusServiceInterface) *ControlReviewStatusController {
	return &ControlReviewStatusController{
		service: service,
	}
}

func (crsc *ControlReviewStatusController) GetAllControlReviewStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	controlReviewStatuses, err := crsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(crsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, controlReviewStatuses)
}

func (crsc *ControlReviewStatusController) GetControlReviewStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlReviewStatus, err := crsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, controlReviewStatus)
}
