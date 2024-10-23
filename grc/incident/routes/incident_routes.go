package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	incidentRoutes := router.Group("/")

	incidentRoutes.Use(middlewares.AuthMiddleware())

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	incidentCategoryController := app.IncidentModule.IncidentCategoryController
	incidentSubCategoryController := app.IncidentModule.IncidentSubCategoryController
	incidentTypeController := app.IncidentModule.IncidentTypeController
	incidentStatusController := app.IncidentModule.IncidentStatusController
	incidentRegisterController := app.IncidentModule.IncidentRegisterController
	incidentInvestigationStatusController := app.IncidentModule.IncidentInvestigationStatusController
	incidentInvestigationController := app.IncidentModule.IncidentInvestigationController
	incidentActionPlanController := app.IncidentModule.IncidentActionPlanController
	incidentActionPlanStatusUpdateController := app.IncidentModule.IncidentActionPlanStatusUpdateController
	incidentRootCauseAnalysisController := app.IncidentModule.IncidentRootCauseAnalysisController
	incidentCommunicationLogController := app.IncidentModule.IncidentCommunicationLogController
	incidentClosureController := app.IncidentModule.IncidentClosureController
	incidentPostReviewEffectivenessController := app.IncidentModule.IncidentPostReviewEffectivenessController
	incidentPostReviewController := app.IncidentModule.IncidentPostReviewController
	incidentAssociatedItemController := app.IncidentModule.IncidentAssociatedItemController
	incidentTimelineController := app.IncidentModule.IncidentTimelineController

	//Incident Category
	incidentRoutes.GET("/incident-categories", incidentCategoryController.GetAllIncidentCategories)
	incidentRoutes.GET("/incident-categories/:id", incidentCategoryController.GetIncidentCategoryDetails)
	incidentRoutes.POST("/incident-categories", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentCategoryController.CreateIncidentCategory)
	incidentRoutes.PUT("/incident-categories/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentCategoryController.UpdateIncidentCategory)
	incidentRoutes.POST("/incident-categories/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentCategoryController.DeleteIncidentCategories)
	incidentRoutes.PUT("/incident-categories/activate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentCategoryController.ActivateIncidentCategories)
	incidentRoutes.PUT("/incident-categories/deactivate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentCategoryController.DeactivateIncidentCategories)

	//Incident Sub Category
	incidentRoutes.GET("/incident-sub-categories", incidentSubCategoryController.GetAllIncidentSubCategories)
	incidentRoutes.GET("/incident-sub-categories/:id", incidentSubCategoryController.GetIncidentSubCategoryDetails)
	incidentRoutes.POST("/incident-sub-categories", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSubCategoryController.CreateIncidentSubCategory)
	incidentRoutes.PUT("/incident-sub-categories/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSubCategoryController.UpdateIncidentSubCategory)
	incidentRoutes.POST("/incident-sub-categories/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSubCategoryController.DeleteIncidentSubCategories)
	incidentRoutes.PUT("/incident-sub-categories/activate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSubCategoryController.ActivateIncidentSubCategories)
	incidentRoutes.PUT("/incident-sub-categories/deactivate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSubCategoryController.DeactivateIncidentSubCategories)

	//Incident Type
	incidentRoutes.GET("/incident-types", incidentTypeController.GetAllIncidentTypes)
	incidentRoutes.GET("/incident-types/:id", incidentTypeController.GetIncidentTypeDetails)
	incidentRoutes.POST("/incident-types", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentTypeController.CreateIncidentType)
	incidentRoutes.PUT("/incident-types/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentTypeController.UpdateIncidentType)
	incidentRoutes.POST("/incident-types/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentTypeController.DeleteIncidentTypes)
	incidentRoutes.PUT("/incident-types/activate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentTypeController.ActivateIncidentTypes)
	incidentRoutes.PUT("/incident-types/deactivate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentTypeController.DeactivateIncidentTypes)

	//Incident Status
	incidentRoutes.GET("/incident-statuses", incidentStatusController.GetAllIncidentStatuses)
	incidentRoutes.GET("/incident-statuses/:id", incidentStatusController.GetIncidentStatusDetails)

	//Incident Register
	unauthenticatedRoutes.GET("/incident-register/:id/document-preview", incidentRegisterController.GetIncidentRegisterDocumentPreview)
	unauthenticatedRoutes.GET("/incident-register/:id/document-download", incidentRegisterController.DownloadIncidentRegisterDocument)
	incidentRoutes.GET("/incident-register/export", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentRegisterController.GetIncidentRegisterExport)
	incidentRoutes.GET("/incident-register", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentRegisterController.GetAllIncidentRegisters)
	incidentRoutes.GET("/incident-register/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentRegisterController.GetIncidentRegisterDetails)
	incidentRoutes.POST("/incident-register", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRegisterController.CreateIncidentRegister)
	incidentRoutes.PUT("/incident-register/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRegisterController.UpdateIncidentRegister)
	incidentRoutes.POST("/incident-register/delete", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRegisterController.DeleteIncidentRegisters)

	//Incident Investigation Status
	incidentRoutes.GET("/incident-investigation-statuses", incidentInvestigationStatusController.GetAllIncidentInvestigationStatuses)
	incidentRoutes.GET("/incident-investigation-statuses/:id", incidentInvestigationStatusController.GetIncidentInvestigationStatusDetails)

	//Incident Investigation
	incidentRoutes.GET("/incidents/:incidentID/investigation", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentInvestigationController.GetIncidentInvestigationDetails)
	incidentRoutes.POST("/incidents/:incidentID/investigation", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentInvestigationController.CreateIncidentInvestigation)
	incidentRoutes.PUT("/incidents/:incidentID/investigation", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentInvestigationController.UpdateIncidentInvestigation)
	incidentRoutes.DELETE("/incidents/:incidentID/investigation", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentInvestigationController.DeleteIncidentInvestigations)

	//Incident Action Plan
	incidentRoutes.GET("/incident-action-plans/export", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanController.GetIncidentActionPlanExport)
	incidentRoutes.GET("/incident-action-plans", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanController.GetAllIncidentActionPlans)
	incidentRoutes.GET("/incidents/:incidentID/action-plans", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanController.GetAllIncidentActionPlans)
	incidentRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanController.GetIncidentActionPlanDetails)
	incidentRoutes.POST("/incidents/:incidentID/action-plans", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "modify"), incidentActionPlanController.CreateIncidentActionPlan)
	incidentRoutes.PUT("/incidents/:incidentID/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "modify"), incidentActionPlanController.UpdateIncidentActionPlan)
	incidentRoutes.POST("/incidents/:incidentID/action-plans/delete", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "modify"), incidentActionPlanController.DeleteIncidentActionPlans)

	//Incident Action Plan Status Updates
	incidentRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanStatusUpdateController.GetAllIncidentActionPlanStatusUpdates)
	incidentRoutes.POST("/incidents/:incidentID/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "modify"), incidentActionPlanStatusUpdateController.CreateIncidentActionPlanStatusUpdate)
	incidentRoutes.PUT("/incidents/:incidentID/action-plans/:actionPlanId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "modify"), incidentActionPlanStatusUpdateController.UpdateIncidentActionPlanStatusUpdate)
	incidentRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanStatusUpdateController.GetLatestIncidentActionPlanStatusUpdate)

	//Incident Root Cause Analyses
	incidentRoutes.GET("/incidents/:incidentID/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentRootCauseAnalysisController.GetAllIncidentRootCauseAnalyses)
	incidentRoutes.POST("/incidents/:incidentID/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRootCauseAnalysisController.CreateIncidentRootCauseAnalysis)
	incidentRoutes.PUT("/incidents/:incidentID/root-cause-analysis/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRootCauseAnalysisController.UpdateIncidentRootCauseAnalysis)
	incidentRoutes.DELETE("/incidents/:incidentID/root-cause-analysis/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRootCauseAnalysisController.DeleteIncidentRootCauseAnalyses)

	//Incident Communication logs
	incidentRoutes.GET("/incidents/:incidentID/communication-logs", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentCommunicationLogController.GetAllIncidentCommunicationLogs)
	incidentRoutes.GET("/incidents/:incidentID/communication-logs/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentCommunicationLogController.GetIncidentCommunicationLogDetails)
	incidentRoutes.POST("/incidents/:incidentID/communication-logs", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentCommunicationLogController.CreateIncidentCommunicationLog)
	incidentRoutes.PUT("/incidents/:incidentID/communication-logs/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentCommunicationLogController.UpdateIncidentCommunicationLog)
	incidentRoutes.POST("/incidents/:incidentID/communication-logs/delete", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentCommunicationLogController.DeleteIncidentCommunicationLogs)

	//Incident Closure
	incidentRoutes.GET("/incidents/:incidentID/closure", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentClosureController.GetIncidentClosureDetails)
	incidentRoutes.POST("/incidents/:incidentID/closure", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureController.CreateIncidentClosure)
	incidentRoutes.PUT("/incidents/:incidentID/closure", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureController.UpdateIncidentClosure)
	incidentRoutes.DELETE("/incidents/:incidentID/closure", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureController.DeleteIncidentClosures)

	//Incident Post Review Effectiveness
	incidentRoutes.GET("/incident-post-review-effectivenesses", incidentPostReviewEffectivenessController.GetAllIncidentPostReviewEffectivenesses)
	incidentRoutes.GET("/incident-post-review-effectivenesses/:id", incidentPostReviewEffectivenessController.GetIncidentPostReviewEffectivenessDetails)

	//Incident PostReview
	incidentRoutes.GET("/incidents/:incidentID/post-review", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentPostReviewController.GetIncidentPostReviewDetails)
	incidentRoutes.POST("/incidents/:incidentID/post-review", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentPostReviewController.CreateIncidentPostReview)
	incidentRoutes.PUT("/incidents/:incidentID/post-review", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentPostReviewController.UpdateIncidentPostReview)
	incidentRoutes.DELETE("/incidents/:incidentID/post-review", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentPostReviewController.DeleteIncidentPostReviews)

	//Incident Associated Items
	incidentRoutes.GET("/incidents/:incidentID/controls", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentControls)
	incidentRoutes.GET("/incidents/:incidentID/compliance-requirements", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentComplianceRequirements)
	incidentRoutes.GET("/incidents/:incidentID/risks", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentRisks)
	incidentRoutes.GET("/incidents/:incidentID/processes", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentProcesses)
	incidentRoutes.GET("/incidents/:incidentID/slas", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentSLAs)
	incidentRoutes.GET("/incidents/:incidentID/documents", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentDocuments)

	//Incident timelines
	incidentRoutes.GET("/incidents/:incidentID/timelines", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentTimelineController.GetAllIncidentTimelines)

}
