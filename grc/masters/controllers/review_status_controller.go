package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewStatusController struct {
	service services.ReviewStatusServiceInterface
	db      *gorm.DB
}

func NewReviewStatusController(service services.ReviewStatusServiceInterface) *ReviewStatusController {
	return &ReviewStatusController{
		service: service,
	}
}

func (rsc *ReviewStatusController) GetAllReviewStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	reviewStatuses, err := rsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, reviewStatuses)
}

func (rsc *ReviewStatusController) GetReviewStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, "Invalid id")
		return
	}

	reviewStatus, err := rsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, reviewStatus)
}
