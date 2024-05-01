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

	incidentSubCategoryRepo := incidentRepositories.NewIncidentSubCategoryRepository(db)
	incidentSubCategorySvc := incidentServices.NewIncidentSubCategoryService(incidentSubCategoryRepo, activityLogSvc, incidentCategorySvc)
	incidentSubCategoryController := controllers.NewIncidentSubCategoryController(incidentSubCategorySvc)

	incidentTypeRepo := incidentRepositories.NewIncidentTypeRepository(db)
	incidentTypeSvc := incidentServices.NewIncidentTypeService(incidentTypeRepo, activityLogSvc)
	incidentTypeController := controllers.NewIncidentTypeController(incidentTypeSvc)

	//Incident Category
	incidentRoutes.GET("/incident-categories", incidentCategoryController.GetAllIncidentCategories)
	incidentRoutes.GET("/incident-categories/:id", incidentCategoryController.GetIncidentCategoryDetails)
	incidentRoutes.POST("/incident-categories", incidentCategoryController.CreateIncidentCategory)
	incidentRoutes.PUT("/incident-categories/:id", incidentCategoryController.UpdateIncidentCategory)
	incidentRoutes.POST("/incident-categories/delete", incidentCategoryController.DeleteIncidentCategories)
	incidentRoutes.PUT("/incident-categories/activate", incidentCategoryController.ActivateIncidentCategories)
	incidentRoutes.PUT("/incident-categories/deactivate", incidentCategoryController.DeactivateIncidentCategories)

	//Incident Sub Category
	incidentRoutes.GET("/incident-sub-categories", incidentSubCategoryController.GetAllIncidentSubCategories)
	incidentRoutes.GET("/incident-sub-categories/:id", incidentSubCategoryController.GetIncidentSubCategoryDetails)
	incidentRoutes.POST("/incident-sub-categories", incidentSubCategoryController.CreateIncidentSubCategory)
	incidentRoutes.PUT("/incident-sub-categories/:id", incidentSubCategoryController.UpdateIncidentSubCategory)
	incidentRoutes.POST("/incident-sub-categories/delete", incidentSubCategoryController.DeleteIncidentSubCategories)
	incidentRoutes.PUT("/incident-sub-categories/activate", incidentSubCategoryController.ActivateIncidentSubCategories)
	incidentRoutes.PUT("/incident-sub-categories/deactivate", incidentSubCategoryController.DeactivateIncidentSubCategories)
	
	//Incident Type
	incidentRoutes.GET("/incident-types", incidentTypeController.GetAllIncidentTypes)
	incidentRoutes.GET("/incident-types/:id", incidentTypeController.GetIncidentTypeDetails)
	incidentRoutes.POST("/incident-types", incidentTypeController.CreateIncidentType)
	incidentRoutes.PUT("/incident-types/:id", incidentTypeController.UpdateIncidentType)
	incidentRoutes.POST("/incident-types/delete", incidentTypeController.DeleteIncidentTypes)
	incidentRoutes.PUT("/incident-types/activate", incidentTypeController.ActivateIncidentTypes)
	incidentRoutes.PUT("/incident-types/deactivate", incidentTypeController.DeactivateIncidentTypes)

}
