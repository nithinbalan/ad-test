package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	unauthenticatedRoutes := router.Group("/")
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
	threatIntelligenceLikelihoodController := app.ThreatIntelligenceModule.ThreatIntelligenceLikelihoodController
	threatIntelligenceSeverityController := app.ThreatIntelligenceModule.ThreatIntelligenceSeverityController
	threatIntelligenceRiskLevelController := app.ThreatIntelligenceModule.ThreatIntelligenceRiskLevelController
	threatIntelligenceRiskLevelValueController := app.ThreatIntelligenceModule.ThreatIntelligenceRiskLevelValueController
	threatIntelligencePlannedAdvisoryController := app.ThreatIntelligenceModule.ThreatIntelligencePlannedAdvisoryController

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

	// Threat Intelligence Likelihoods
	threatIntelligenceRoutes.GET("/threat-intelligence-likelihoods", threatIntelligenceLikelihoodController.GetAllThreatIntelligenceLikelihoods)
	threatIntelligenceRoutes.POST("/threat-intelligence-likelihoods", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceLikelihoodController.CreateThreatIntelligenceLikelihood)
	threatIntelligenceRoutes.GET("/threat-intelligence-likelihoods/:id", threatIntelligenceLikelihoodController.GetThreatIntelligenceLikelihoodDetails)
	threatIntelligenceRoutes.PUT("/threat-intelligence-likelihoods/:id", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceLikelihoodController.UpdateThreatIntelligenceLikelihood)
	threatIntelligenceRoutes.POST("/threat-intelligence-likelihoods/delete", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceLikelihoodController.DeleteThreatIntelligenceLikelihoods)
	threatIntelligenceRoutes.PUT("/threat-intelligence-likelihoods/activate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceLikelihoodController.ActivateThreatIntelligenceLikelihood)
	threatIntelligenceRoutes.PUT("/threat-intelligence-likelihoods/deactivate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceLikelihoodController.DeactivateThreatIntelligenceLikelihood)

	// Threat Intelligence Severities
	threatIntelligenceRoutes.GET("/threat-intelligence-severities", threatIntelligenceSeverityController.GetAllThreatIntelligenceSeverities)
	threatIntelligenceRoutes.POST("/threat-intelligence-severities", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceSeverityController.CreateThreatIntelligenceSeverity)
	threatIntelligenceRoutes.GET("/threat-intelligence-severities/:id", threatIntelligenceSeverityController.GetThreatIntelligenceSeverityDetails)
	threatIntelligenceRoutes.PUT("/threat-intelligence-severities/:id", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceSeverityController.UpdateThreatIntelligenceSeverity)
	threatIntelligenceRoutes.POST("/threat-intelligence-severities/delete", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceSeverityController.DeleteThreatIntelligenceSeverities)
	threatIntelligenceRoutes.PUT("/threat-intelligence-severities/activate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceSeverityController.ActivateThreatIntelligenceSeverity)
	threatIntelligenceRoutes.PUT("/threat-intelligence-severities/deactivate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceSeverityController.DeactivateThreatIntelligenceSeverity)

	// Threat Intelligence Risk Levels
	threatIntelligenceRoutes.GET("/threat-intelligence-risk-levels", threatIntelligenceRiskLevelController.GetAllThreatIntelligenceRiskLevels)
	threatIntelligenceRoutes.POST("/threat-intelligence-risk-levels", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceRiskLevelController.CreateThreatIntelligenceRiskLevel)
	threatIntelligenceRoutes.GET("/threat-intelligence-risk-levels/:id", threatIntelligenceRiskLevelController.GetThreatIntelligenceRiskLevelDetails)
	threatIntelligenceRoutes.PUT("/threat-intelligence-risk-levels/:id", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceRiskLevelController.UpdateThreatIntelligenceRiskLevel)
	threatIntelligenceRoutes.POST("/threat-intelligence-risk-levels/delete", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceRiskLevelController.DeleteThreatIntelligenceRiskLevels)
	threatIntelligenceRoutes.PUT("/threat-intelligence-risk-levels/activate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceRiskLevelController.ActivateThreatIntelligenceRiskLevel)
	threatIntelligenceRoutes.PUT("/threat-intelligence-risk-levels/deactivate", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceRiskLevelController.DeactivateThreatIntelligenceRiskLevel)

	// Threat Intelligence Risk Level Values
	threatIntelligenceRoutes.GET("/threat-intelligence-risk-level-values", threatIntelligenceRiskLevelValueController.GetAllThreatIntelligenceRiskLevelValues)
	threatIntelligenceRoutes.POST("/threat-intelligence-risk-level-values", middlewares.RolePermissionMiddleware(db, "threat_intelligence_setting", "modify"), threatIntelligenceRiskLevelValueController.CreateThreatIntelligenceRiskLevelValue)
	threatIntelligenceRoutes.GET("/threat-intelligence-matrix-configuration", threatIntelligenceRiskLevelValueController.GetThreatIntelligenceRiskLevelMatrix)

	// Threat Intelligence Planned Advisory
	unauthenticatedRoutes.GET("/threat-intelligence-planned-advisories/:id/document-preview", threatIntelligencePlannedAdvisoryController.GetThreatIntelligencePlannedAdvisoryDocumentPreview)
	unauthenticatedRoutes.GET("/threat-intelligence-planned-advisories/:id/document-download", threatIntelligencePlannedAdvisoryController.DownloadThreatIntelligencePlannedAdvisoryDocument)
	threatIntelligenceRoutes.GET("/threat-intelligence-planned-advisories/export", middlewares.RolePermissionMiddleware(db, "planned_advisory", "view"), threatIntelligencePlannedAdvisoryController.ExportThreatIntelligencePlannedAdvisories)
	threatIntelligenceRoutes.GET("/threat-intelligence-planned-advisories", middlewares.RolePermissionMiddleware(db, "planned_advisory", "view"), threatIntelligencePlannedAdvisoryController.GetAllThreatIntelligencePlannedAdvisories)
	threatIntelligenceRoutes.GET("/threat-intelligence-planned-advisories/:id", middlewares.RolePermissionMiddleware(db, "planned_advisory", "view"), threatIntelligencePlannedAdvisoryController.GetThreatIntelligencePlannedAdvisoryDetails)
	threatIntelligenceRoutes.POST("/threat-intelligence-planned-advisories", middlewares.RolePermissionMiddleware(db, "planned_advisory", "modify"), threatIntelligencePlannedAdvisoryController.CreateThreatIntelligencePlannedAdvisory)
	threatIntelligenceRoutes.PUT("/threat-intelligence-planned-advisories/:id", middlewares.RolePermissionMiddleware(db, "planned_advisory", "modify"), threatIntelligencePlannedAdvisoryController.UpdateThreatIntelligencePlannedAdvisory)
	threatIntelligenceRoutes.POST("/threat-intelligence-planned-advisories/:id/status", middlewares.RolePermissionMiddleware(db, "planned_advisory", "modify"), threatIntelligencePlannedAdvisoryController.UpdateThreatIntelligencePlannedAdvisoryStatus)
	threatIntelligenceRoutes.POST("/threat-intelligence-planned-advisories/delete", middlewares.RolePermissionMiddleware(db, "planned_advisory", "modify"), threatIntelligencePlannedAdvisoryController.DeleteThreatIntelligencePlannedAdvisories)
}
