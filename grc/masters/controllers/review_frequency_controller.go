package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewFrequencyController struct {
	service services.ReviewFrequencyServiceInterface
	db      *gorm.DB
}

func NewReviewFrequencyController(service services.ReviewFrequencyServiceInterface) *ReviewFrequencyController {
	return &ReviewFrequencyController{
		service: service,
	}
}

func (rfc *ReviewFrequencyController) GetAllReviewFrequencies(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	reviewFrequencies, err := rfc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rfc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, reviewFrequencies)
}
