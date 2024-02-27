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
	articles, err := rfc.service.GetAll()
	if err != nil {
		utility.SendErrorResponse(rfc.db, c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, articles)
}
