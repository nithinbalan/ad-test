package controllers

import (
	"jethur/grc/masters/repositories/interfaces"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewFrequencyController struct {
	repo interfaces.ReviewFrequencyRepositoryInterface
	db   *gorm.DB
}

func NewReviewFrequencyController(repo interfaces.ReviewFrequencyRepositoryInterface, db *gorm.DB) *ReviewFrequencyController {
	return &ReviewFrequencyController{
		repo: repo,
		db:   db,
	}
}

func (rfc *ReviewFrequencyController) GetAllReviewFrequencies(c *gin.Context) {
	articles, err := rfc.repo.GetAll(c)
	if err != nil {
		utility.SendErrorResponse(rfc.db, c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, articles)
}
