package routes

import (
	"jethur/grc/masters/controllers"
	"jethur/grc/masters/repositories"
	"jethur/grc/masters/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	//Review Frequency
	reviewFrequencyRepository := repositories.NewReviewFrequencyRepository(db)
	reviewFrequencyService := services.NewReviewFrequencyService(reviewFrequencyRepository)

	reviewFrequencyController := controllers.NewReviewFrequencyController(reviewFrequencyService)

	masterRoutes := router.Group("/")

	masterRoutes.Use(middlewares.AuthMiddleware())

	masterRoutes.GET("/review-frequencies", reviewFrequencyController.GetAllReviewFrequencies)
}
