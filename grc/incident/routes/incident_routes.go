package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	incidentRoutes := router.Group("/")

	incidentRoutes.Use(middlewares.AuthMiddleware())

	incidentCategoryController := app.IncidentModule.IncidentCategoryController
	incidentSubCategoryController := app.IncidentModule.IncidentSubCategoryController
	incidentTypeController := app.IncidentModule.IncidentTypeController
	incidentStatusController := app.IncidentModule.IncidentStatusController
	incidentRegisterController := app.IncidentModule.IncidentRegisterController
	incidentInvestigationStatusController := app.IncidentModule.IncidentInvestigationStatusController
	incidentInvestigationController := app.IncidentModule.IncidentInvestigationController

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

	//Incident Status
	incidentRoutes.GET("/incident-statuses", incidentStatusController.GetAllIncidentStatuses)
	incidentRoutes.GET("/incident-statuses/:id", incidentStatusController.GetIncidentStatusDetails)

	//Incident Register
	incidentRoutes.GET("/incident-register", incidentRegisterController.GetAllIncidentRegisters)
	incidentRoutes.GET("/incident-register/:id", incidentRegisterController.GetIncidentRegisterDetails)
	incidentRoutes.POST("/incident-register", incidentRegisterController.CreateIncidentRegister)
	incidentRoutes.PUT("/incident-register/:id", incidentRegisterController.UpdateIncidentRegister)
	incidentRoutes.POST("/incident-register/delete", incidentRegisterController.DeleteIncidentRegisters)

	//Incident Investigation Status
	incidentRoutes.GET("/incident-investigation-statuses", incidentInvestigationStatusController.GetAllIncidentInvestigationStatuses)
	incidentRoutes.GET("/incident-investigation-statuses/:id", incidentInvestigationStatusController.GetIncidentInvestigationStatusDetails)

	//Incident Investigation
	incidentRoutes.GET("/incident/:incidentId/investigations", incidentInvestigationController.GetAllIncidentInvestigations)
	incidentRoutes.GET("/incident/:incidentId/investigations/:id", incidentInvestigationController.GetIncidentInvestigationDetails)
	incidentRoutes.POST("/incident/:incidentId/investigations", incidentInvestigationController.CreateIncidentInvestigation)
	incidentRoutes.PUT("/incident/:incidentId/investigations/:id", incidentInvestigationController.UpdateIncidentInvestigation)
	incidentRoutes.POST("/incident/:incidentId/investigations/delete", incidentInvestigationController.DeleteIncidentInvestigations)

}
