package routes

import (
	"jethur/grc/masters/controllers"
	"jethur/grc/masters/repositories/implementations"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	//Review Frequency
	reviewFrequencyRepository := implementations.NewReviewFrequencyRepository(db)
	reviewFrequencyController := controllers.NewReviewFrequencyController(reviewFrequencyRepository, db)
	router.GET("/review-frequencies", reviewFrequencyController.GetAllReviewFrequencies)
}
