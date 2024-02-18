package routes

import (
	"jethur/grc/masters/controllers"
	"jethur/grc/masters/repositories/implementations"

	"github.com/gin-gonic/gin"
)

func RegisterMasterRoutes(router *gin.Engine) {

	//Review Frequency
	reviewFrequencyRepository := implementations.NewReviewFrequencyRepository()
	reviewFrequencyController := controllers.NewReviewFrequencyController(reviewFrequencyRepository)
	router.GET("/review-frequencies", reviewFrequencyController.GetAllReviewFrequencies)
}
