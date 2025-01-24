package controllers

import (
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FeedbackOptionController struct {
	service services.FeedbackOptionServiceInterface
	db      *gorm.DB
}

func NewFeedbackOptionController(service services.FeedbackOptionServiceInterface) *FeedbackOptionController {
	return &FeedbackOptionController{
		service: service,
	}
}

func (tesc *FeedbackOptionController) GetAllFeedbackOptions(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	feedbackOptions, err := tesc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tesc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, feedbackOptions)
}

func (tesc *FeedbackOptionController) GetFeedbackOptionDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	feedbackOption, err := tesc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, feedbackOption)
}
