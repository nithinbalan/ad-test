package controllers

import (
	"jethur/grc/awareness-reporting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessEvaluationRatingController struct {
	service services.AwarenessEvaluationRatingServiceInterface
	db      *gorm.DB
}

func NewAwarenessEvaluationRatingController(service services.AwarenessEvaluationRatingServiceInterface) *AwarenessEvaluationRatingController {
	return &AwarenessEvaluationRatingController{
		service: service,
	}
}

func (aerc *AwarenessEvaluationRatingController) GetAllAwarenessEvaluationRatings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	awarenessEvaluationRatings, err := aerc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(aerc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, awarenessEvaluationRatings)
}

func (aerc *AwarenessEvaluationRatingController) GetAwarenessEvaluationRatingDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	awarenessEvaluationRating, err := aerc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, awarenessEvaluationRating)
}