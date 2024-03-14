package routes

import (
	activityLogRepository "jethur/grc/logs/repositories"
	activityLogServie "jethur/grc/logs/services"
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

	activityLogRepo := activityLogRepository.NewActivityLogRepository(db)
	activityLogSvc := activityLogServie.NewActivityLogService(activityLogRepo)

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

	//Checklist
	checklistRepository := repositories.NewChecklistRepository(db)
	checklistService := services.NewChecklistService(checklistRepository, activityLogSvc)
	checklistController := controllers.NewChecklistController(checklistService)

	masterRoutes.GET("/checklists", checklistController.GetAllChecklists)
	masterRoutes.POST("/checklists", checklistController.CreateChecklist)
	masterRoutes.GET("/checklists/:id", checklistController.GetChecklistDetails)
	masterRoutes.PUT("/checklists/:id", checklistController.UpdateChecklist)
	masterRoutes.POST("/checklists/delete", checklistController.DeleteChecklists)
	masterRoutes.POST("/checklists/activate", checklistController.ActivateChecklists)
	masterRoutes.POST("/checklists/deactivate", checklistController.DeactivateChecklists)

}
