package controllers

import (
	"jethur/grc/masters/repositories/interfaces"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReviewFrequencyController struct {
	repo interfaces.ReviewFrequencyRepositoryInterface
}

func NewReviewFrequencyController(repo interfaces.ReviewFrequencyRepositoryInterface) *ReviewFrequencyController {
	return &ReviewFrequencyController{
		repo: repo,
	}
}

func (rfc *ReviewFrequencyController) GetAllReviewFrequencies(c *gin.Context) {
	articles, err := rfc.repo.GetAll(c)
	if err != nil {
		utility.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, articles)
}
