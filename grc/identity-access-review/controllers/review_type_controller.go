package controllers

import (
	"jethur/grc/identity-access-review/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewTypeController struct {
	service services.ReviewTypeServiceInterface
	db      *gorm.DB
}

func NewReviewTypeController(service services.ReviewTypeServiceInterface) *ReviewTypeController {
	return &ReviewTypeController{
		service: service,
	}
}

func (rtc *ReviewTypeController) GetAllReviewTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	reviewTypes, err := rtc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rtc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reviewTypes)
}

func (rtc *ReviewTypeController) GetReviewTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reviewType, err := rtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, reviewType)
}
