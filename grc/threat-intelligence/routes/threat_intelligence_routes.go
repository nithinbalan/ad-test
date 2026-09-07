package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	threatIntelligenceRoutes := router.Group("/")

	threatIntelligenceRoutes.Use(middlewares.AuthMiddleware(app.DB))

	db := app.DB

	plannedAdvisoryStatusController := app.ThreatIntelligenceModule.PlannedAdvisoryStatusController
	plannedAdvisorySourceController := app.ThreatIntelligenceModule.PlannedAdvisorySourceController
	plannedAdvisoryClassificationController := app.ThreatIntelligenceModule.PlannedAdvisoryClassificationController
	plannedAdvisoryTypeController := app.ThreatIntelligenceModule.PlannedAdvisoryTypeController
	plannedAdvisoryImportanceController := app.ThreatIntelligenceModule.PlannedAdvisoryImportanceController
	plannedAdvisorySecurityClassificationController := app.ThreatIntelligenceModule.PlannedAdvisorySecurityClassificationController
	plannedAdvisoryInternalStatusController := app.ThreatIntelligenceModule.PlannedAdvisoryInternalStatusController

	// Planned Advisory Status
	threatIntelligenceRoutes.GET("/planned-advisory-statuses", plannedAdvisoryStatusController.GetAllPlannedAdvisoryStatuses)
	threatIntelligenceRoutes.GET("/planned-advisory-statuses/:id", plannedAdvisoryStatusController.GetPlannedAdvisoryStatusDetails)

	// Planned Advisory Sources
	threatIntelligenceRoutes.GET("/planned-advisory-sources", plannedAdvisorySourceController.GetAllPlannedAdvisorySources)
	threatIntelligenceRoutes.POST("/planned-advisory-sources", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisorySourceController.CreatePlannedAdvisorySource)
	threatIntelligenceRoutes.GET("/planned-advisory-sources/:id", plannedAdvisorySourceController.GetPlannedAdvisorySourceDetails)
	threatIntelligenceRoutes.PUT("/planned-advisory-sources/:id", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisorySourceController.UpdatePlannedAdvisorySource)
	threatIntelligenceRoutes.POST("/planned-advisory-sources/delete", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisorySourceController.DeletePlannedAdvisorySources)
	threatIntelligenceRoutes.PUT("/planned-advisory-sources/activate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisorySourceController.ActivatePlannedAdvisorySources)
	threatIntelligenceRoutes.PUT("/planned-advisory-sources/deactivate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisorySourceController.DeactivatePlannedAdvisorySources)

	// Planned Advisory Classifications
	threatIntelligenceRoutes.GET("/planned-advisory-classifications", plannedAdvisoryClassificationController.GetAllPlannedAdvisoryClassifications)
	threatIntelligenceRoutes.POST("/planned-advisory-classifications", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryClassificationController.CreatePlannedAdvisoryClassification)
	threatIntelligenceRoutes.GET("/planned-advisory-classifications/:id", plannedAdvisoryClassificationController.GetPlannedAdvisoryClassificationDetails)
	threatIntelligenceRoutes.PUT("/planned-advisory-classifications/:id", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryClassificationController.UpdatePlannedAdvisoryClassification)
	threatIntelligenceRoutes.POST("/planned-advisory-classifications/delete", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryClassificationController.DeletePlannedAdvisoryClassifications)
	threatIntelligenceRoutes.PUT("/planned-advisory-classifications/activate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryClassificationController.ActivatePlannedAdvisoryClassifications)
	threatIntelligenceRoutes.PUT("/planned-advisory-classifications/deactivate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryClassificationController.DeactivatePlannedAdvisoryClassifications)

	// Planned Advisory Types
	threatIntelligenceRoutes.GET("/planned-advisory-types", plannedAdvisoryTypeController.GetAllPlannedAdvisoryTypes)
	threatIntelligenceRoutes.POST("/planned-advisory-types", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryTypeController.CreatePlannedAdvisoryType)
	threatIntelligenceRoutes.GET("/planned-advisory-types/:id", plannedAdvisoryTypeController.GetPlannedAdvisoryTypeDetails)
	threatIntelligenceRoutes.PUT("/planned-advisory-types/:id", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryTypeController.UpdatePlannedAdvisoryType)
	threatIntelligenceRoutes.POST("/planned-advisory-types/delete", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryTypeController.DeletePlannedAdvisoryTypes)
	threatIntelligenceRoutes.PUT("/planned-advisory-types/activate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryTypeController.ActivatePlannedAdvisoryTypes)
	threatIntelligenceRoutes.PUT("/planned-advisory-types/deactivate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryTypeController.DeactivatePlannedAdvisoryTypes)

	// Planned Advisory Importances
	threatIntelligenceRoutes.GET("/planned-advisory-importances", plannedAdvisoryImportanceController.GetAllPlannedAdvisoryImportances)
	threatIntelligenceRoutes.POST("/planned-advisory-importances", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryImportanceController.CreatePlannedAdvisoryImportance)
	threatIntelligenceRoutes.GET("/planned-advisory-importances/:id", plannedAdvisoryImportanceController.GetPlannedAdvisoryImportanceDetails)
	threatIntelligenceRoutes.PUT("/planned-advisory-importances/:id", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryImportanceController.UpdatePlannedAdvisoryImportance)
	threatIntelligenceRoutes.POST("/planned-advisory-importances/delete", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryImportanceController.DeletePlannedAdvisoryImportances)
	threatIntelligenceRoutes.PUT("/planned-advisory-importances/activate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryImportanceController.ActivatePlannedAdvisoryImportances)
	threatIntelligenceRoutes.PUT("/planned-advisory-importances/deactivate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryImportanceController.DeactivatePlannedAdvisoryImportances)

	// Planned Advisory Security Classifications
	threatIntelligenceRoutes.GET("/planned-advisory-security-classifications", plannedAdvisorySecurityClassificationController.GetAllPlannedAdvisorySecurityClassifications)
	threatIntelligenceRoutes.POST("/planned-advisory-security-classifications", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisorySecurityClassificationController.CreatePlannedAdvisorySecurityClassification)
	threatIntelligenceRoutes.GET("/planned-advisory-security-classifications/:id", plannedAdvisorySecurityClassificationController.GetPlannedAdvisorySecurityClassificationDetails)
	threatIntelligenceRoutes.PUT("/planned-advisory-security-classifications/:id", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisorySecurityClassificationController.UpdatePlannedAdvisorySecurityClassification)
	threatIntelligenceRoutes.POST("/planned-advisory-security-classifications/delete", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisorySecurityClassificationController.DeletePlannedAdvisorySecurityClassifications)
	threatIntelligenceRoutes.PUT("/planned-advisory-security-classifications/activate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisorySecurityClassificationController.ActivatePlannedAdvisorySecurityClassifications)
	threatIntelligenceRoutes.PUT("/planned-advisory-security-classifications/deactivate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisorySecurityClassificationController.DeactivatePlannedAdvisorySecurityClassifications)

	// Planned Advisory Internal Statuses
	threatIntelligenceRoutes.GET("/planned-advisory-internal-statuses", plannedAdvisoryInternalStatusController.GetAllPlannedAdvisoryInternalStatuses)
	threatIntelligenceRoutes.POST("/planned-advisory-internal-statuses", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryInternalStatusController.CreatePlannedAdvisoryInternalStatus)
	threatIntelligenceRoutes.GET("/planned-advisory-internal-statuses/:id", plannedAdvisoryInternalStatusController.GetPlannedAdvisoryInternalStatusDetails)
	threatIntelligenceRoutes.PUT("/planned-advisory-internal-statuses/:id", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryInternalStatusController.UpdatePlannedAdvisoryInternalStatus)
	threatIntelligenceRoutes.POST("/planned-advisory-internal-statuses/delete", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryInternalStatusController.DeletePlannedAdvisoryInternalStatuses)
	threatIntelligenceRoutes.PUT("/planned-advisory-internal-statuses/activate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryInternalStatusController.ActivatePlannedAdvisoryInternalStatuses)
	threatIntelligenceRoutes.PUT("/planned-advisory-internal-statuses/deactivate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), plannedAdvisoryInternalStatusController.DeactivatePlannedAdvisoryInternalStatuses)
}
