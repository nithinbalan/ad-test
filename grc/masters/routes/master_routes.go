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
	reviewFrequencyRepo := repositories.NewReviewFrequencyRepository(db)
	reviewFrequencySvc := services.NewReviewFrequencyService(reviewFrequencyRepo)
	reviewFrequencyController := controllers.NewReviewFrequencyController(reviewFrequencySvc)

	masterRoutes.GET("/review-frequencies", reviewFrequencyController.GetAllReviewFrequencies)
	masterRoutes.GET("/review-frequencies/:id", reviewFrequencyController.GetReviewFrequencyDetails)
	//Severity Level
	severityLevelRepo := repositories.NewSeverityLevelRepository(db)
	severityLevelSvc := services.NewSeverityLevelService(severityLevelRepo)
	severityLevelController := controllers.NewSeverityLevelController(severityLevelSvc)

	masterRoutes.GET("/severity-levels", severityLevelController.GetAllSeverityLevels)
	masterRoutes.GET("/severity-levels/:id", severityLevelController.GetSeverityLevelDetails)

	//Lifecycle Status
	lifecycleStatusRepo := repositories.NewLifecycleStatusRepository(db)
	lifecycleStatusSvc := services.NewLifecycleStatusService(lifecycleStatusRepo)
	lifecycleStatusController := controllers.NewLifecycleStatusController(lifecycleStatusSvc)

	masterRoutes.GET("/lifecycle-statuses", lifecycleStatusController.GetAllLifecycleStatuses)
	masterRoutes.GET("/lifecycle-statuses/:id", lifecycleStatusController.GetLifecycleStatusDetails)

	//Checklist
	checklistRepo := repositories.NewChecklistRepository(db)
	checklistSvc := services.NewChecklistService(checklistRepo, activityLogSvc)
	checklistController := controllers.NewChecklistController(checklistSvc)

	masterRoutes.GET("/checklists", checklistController.GetAllChecklists)
	masterRoutes.POST("/checklists", checklistController.CreateChecklist)
	masterRoutes.GET("/checklists/:id", checklistController.GetChecklistDetails)
	masterRoutes.PUT("/checklists/:id", checklistController.UpdateChecklist)
	masterRoutes.POST("/checklists/delete", checklistController.DeleteChecklists)
	masterRoutes.POST("/checklists/activate", checklistController.ActivateChecklists)
	masterRoutes.POST("/checklists/deactivate", checklistController.DeactivateChecklists)

}
