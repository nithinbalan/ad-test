package routes

import (
	"jethur/grc/incident/controllers"
	incidentRepositories "jethur/grc/incident/repositories"
	incidentServices "jethur/grc/incident/services"
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	incidentRoutes := router.Group("/")

	incidentRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	incidentCategoryRepo := incidentRepositories.NewIncidentCategoryRepository(db)
	incidentCategorySvc := incidentServices.NewIncidentCategoryService(incidentCategoryRepo, activityLogSvc)
	incidentCategoryController := controllers.NewIncidentCategoryController(incidentCategorySvc)

	//Incident Category
	incidentRoutes.GET("/incident-categories", incidentCategoryController.GetAllIncidentCategories)
	incidentRoutes.GET("/incident-categories/:id", incidentCategoryController.GetIncidentCategoryDetails)
	incidentRoutes.POST("/incident-categories", incidentCategoryController.CreateIncidentCategory)
	incidentRoutes.PUT("/incident-categories/:id", incidentCategoryController.UpdateIncidentCategory)
	incidentRoutes.POST("/incident-categories/delete", incidentCategoryController.DeleteIncidentCategories)
	incidentRoutes.PUT("/incident-categories/activate", incidentCategoryController.ActivateIncidentCategories)
	incidentRoutes.PUT("/incident-categories/deactivate", incidentCategoryController.DeactivateIncidentCategories)

}
