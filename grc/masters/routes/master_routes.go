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

	masterRoutes := router.Group("/")

	masterRoutes.Use(middlewares.AuthMiddleware())

	//Review Frequency
	reviewFrequencyRepository := repositories.NewReviewFrequencyRepository(db)
	reviewFrequencyService := services.NewReviewFrequencyService(reviewFrequencyRepository)
	reviewFrequencyController := controllers.NewReviewFrequencyController(reviewFrequencyService)

	masterRoutes.GET("/review-frequencies", reviewFrequencyController.GetAllReviewFrequencies)
	masterRoutes.GET("/review-frequencies/:id", reviewFrequencyController.GetReviewFrequencyDetails)
	//Severity Level
	severityLevelRepository := repositories.NewSeverityLevelRepository(db)
	severityLevelService := services.NewSeverityLevelService(severityLevelRepository)
	severityLevelController := controllers.NewSeverityLevelController(severityLevelService)

	masterRoutes.GET("/severity-levels", severityLevelController.GetAllSeverityLevels)
	masterRoutes.GET("/severity-levels/:id", severityLevelController.GetSeverityLevelDetails)

	//Lifecycle Status
	lifecycleStatusRepository := repositories.NewLifecycleStatusRepository(db)
	lifecycleStatusService := services.NewLifecycleStatusService(lifecycleStatusRepository)
	lifecycleStatusController := controllers.NewLifecycleStatusController(lifecycleStatusService)

	masterRoutes.GET("/lifecycle-statuses", lifecycleStatusController.GetAllLifecycleStatuses)
	masterRoutes.GET("/lifecycle-statuses/:id", lifecycleStatusController.GetLifecycleStatusDetails)

}
