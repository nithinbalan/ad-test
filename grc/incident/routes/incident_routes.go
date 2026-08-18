package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	incidentRoutes := router.Group("/")

	incidentRoutes.Use(middlewares.AuthMiddleware(app.DB))

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
	incidentClosureWorkflowController := app.IncidentModule.IncidentClosureWorkflowController
	incidentClosureReviewController := app.IncidentModule.IncidentClosureReviewController
	incidentClosureReviewCommentController := app.IncidentModule.IncidentClosureReviewCommentController
	incidentPostReviewEffectivenessController := app.IncidentModule.IncidentPostReviewEffectivenessController
	incidentPostReviewController := app.IncidentModule.IncidentPostReviewController
	incidentAssociatedItemController := app.IncidentModule.IncidentAssociatedItemController
	incidentTimelineController := app.IncidentModule.IncidentTimelineController
	incidentActionPlanCommentController := app.IncidentModule.IncidentActionPlanCommentController
	IncidentReportController := app.IncidentModule.IncidentReportController
	incidentResponseTeamController := app.IncidentModule.IncidentResponseTeamController
	incidentImpactController := app.IncidentModule.IncidentImpactController
	incidentScopeController := app.IncidentModule.IncidentScopeController
	incidentSensitivityClassificationController := app.IncidentModule.IncidentSensitivityClassificationController
	incidentEvidenceTypeController := app.IncidentModule.IncidentEvidenceTypeController
	incidentEvidenceSourceController := app.IncidentModule.IncidentEvidenceSourceController
	incidentEvidenceController := app.IncidentModule.IncidentEvidenceController
	incidentSLAEscalationMatrixController := app.IncidentModule.IncidentSLAEscalationMatrixController

	// Incident Category
	incidentRoutes.GET("/incident-categories", incidentCategoryController.GetAllIncidentCategories)
	incidentRoutes.GET("/incident-categories/:id", incidentCategoryController.GetIncidentCategoryDetails)
	incidentRoutes.POST("/incident-categories", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentCategoryController.CreateIncidentCategory)
	incidentRoutes.PUT("/incident-categories/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentCategoryController.UpdateIncidentCategory)
	incidentRoutes.POST("/incident-categories/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentCategoryController.DeleteIncidentCategories)
	incidentRoutes.PUT("/incident-categories/activate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentCategoryController.ActivateIncidentCategories)
	incidentRoutes.PUT("/incident-categories/deactivate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentCategoryController.DeactivateIncidentCategories)

	// Incident Sub Category
	incidentRoutes.GET("/incident-sub-categories", incidentSubCategoryController.GetAllIncidentSubCategories)
	incidentRoutes.GET("/incident-sub-categories/:id", incidentSubCategoryController.GetIncidentSubCategoryDetails)
	incidentRoutes.POST("/incident-sub-categories", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSubCategoryController.CreateIncidentSubCategory)
	incidentRoutes.PUT("/incident-sub-categories/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSubCategoryController.UpdateIncidentSubCategory)
	incidentRoutes.POST("/incident-sub-categories/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSubCategoryController.DeleteIncidentSubCategories)
	incidentRoutes.PUT("/incident-sub-categories/activate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSubCategoryController.ActivateIncidentSubCategories)
	incidentRoutes.PUT("/incident-sub-categories/deactivate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSubCategoryController.DeactivateIncidentSubCategories)

	// Incident Type
	incidentRoutes.GET("/incident-types", incidentTypeController.GetAllIncidentTypes)
	incidentRoutes.GET("/incident-types/:id", incidentTypeController.GetIncidentTypeDetails)
	incidentRoutes.POST("/incident-types", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentTypeController.CreateIncidentType)
	incidentRoutes.PUT("/incident-types/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentTypeController.UpdateIncidentType)
	incidentRoutes.POST("/incident-types/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentTypeController.DeleteIncidentTypes)
	incidentRoutes.PUT("/incident-types/activate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentTypeController.ActivateIncidentTypes)
	incidentRoutes.PUT("/incident-types/deactivate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentTypeController.DeactivateIncidentTypes)

	// Incident Status
	incidentRoutes.GET("/incident-statuses", incidentStatusController.GetAllIncidentStatuses)
	incidentRoutes.GET("/incident-statuses/:id", incidentStatusController.GetIncidentStatusDetails)

	// Incident Register
	unauthenticatedRoutes.GET("/incident-register/:id/document-preview", incidentRegisterController.GetIncidentRegisterDocumentPreview)
	unauthenticatedRoutes.GET("/incident-register/:id/document-download", incidentRegisterController.DownloadIncidentRegisterDocument)
	incidentRoutes.GET("/incident-register/export", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentRegisterController.GetIncidentRegisterExport)
	incidentRoutes.GET("/incident-register", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentRegisterController.GetAllIncidentRegisters)
	incidentRoutes.GET("/incident-register/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentRegisterController.GetIncidentRegisterDetails)
	incidentRoutes.POST("/incident-register", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRegisterController.CreateIncidentRegister)
	incidentRoutes.PUT("/incident-register/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRegisterController.UpdateIncidentRegister)
	incidentRoutes.PUT("/incident-register/:id/status", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRegisterController.UpdateIncidentStatus)
	incidentRoutes.PUT("/incident-register/:id/acknowledge", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRegisterController.AcknowledgeIncidentRegister)
	incidentRoutes.POST("/incident-register/delete", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRegisterController.DeleteIncidentRegisters)
	incidentRoutes.POST("/incident-register/force-delete", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRegisterController.DeleteForceIncidentRegisters)

	// Incident Investigation Status
	incidentRoutes.GET("/incident-investigation-statuses", incidentInvestigationStatusController.GetAllIncidentInvestigationStatuses)
	incidentRoutes.GET("/incident-investigation-statuses/:id", incidentInvestigationStatusController.GetIncidentInvestigationStatusDetails)

	// Incident Investigation
	unauthenticatedRoutes.GET("/incidents/:incidentID/investigation/document-preview", incidentInvestigationController.GetIncidentInvestigationDocumentPreview)
	unauthenticatedRoutes.GET("/incidents/:incidentID/investigation/document-download", incidentInvestigationController.DownloadIncidentInvestigationDocument)
	incidentRoutes.GET("/incidents/:incidentID/investigation", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentInvestigationController.GetIncidentInvestigationDetails)
	incidentRoutes.POST("/incidents/:incidentID/investigation", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentInvestigationController.CreateIncidentInvestigation)
	incidentRoutes.PUT("/incidents/:incidentID/investigation", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentInvestigationController.UpdateIncidentInvestigation)
	incidentRoutes.DELETE("/incidents/:incidentID/investigation", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentInvestigationController.DeleteIncidentInvestigations)

	// Incident Action Plan
	incidentRoutes.POST("/incident-action-plans/force-delete", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanController.DeleteForceIncidentActionPlans)
	incidentRoutes.GET("/incident-action-plans/export", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanController.GetIncidentActionPlanExport)
	incidentRoutes.GET("/incident-action-plans", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanController.GetAllIncidentActionPlans)
	incidentRoutes.GET("/incidents/:incidentID/action-plans", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanController.GetAllIncidentActionPlans)
	incidentRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanController.GetIncidentActionPlanDetails)
	incidentRoutes.GET("/incidents/:incidentID/used-action-plan-master-ids", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanController.GetUsedActionPlanMasterIDs)
	incidentRoutes.POST("/incidents/:incidentID/action-plans", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "modify"), incidentActionPlanController.CreateIncidentActionPlan)
	incidentRoutes.PUT("/incidents/:incidentID/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "modify"), incidentActionPlanController.UpdateIncidentActionPlan)
	incidentRoutes.POST("/incidents/:incidentID/action-plans/delete", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "modify"), incidentActionPlanController.DeleteIncidentActionPlans)

	// Incident Action Plan Status Updates
	unauthenticatedRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId/status-updates/:id/document-preview", incidentActionPlanStatusUpdateController.GetIncidentActionPlanStatusUpdateDocumentPreview)
	unauthenticatedRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId/status-updates/:id/document-download", incidentActionPlanStatusUpdateController.DownloadIncidentActionPlanStatusUpdateDocument)
	incidentRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanStatusUpdateController.GetAllIncidentActionPlanStatusUpdates)
	incidentRoutes.POST("/incidents/:incidentID/action-plans/:actionPlanId/status-updates", incidentActionPlanStatusUpdateController.CreateIncidentActionPlanStatusUpdate)
	incidentRoutes.PUT("/incidents/:incidentID/action-plans/:actionPlanId/status-updates/:id", incidentActionPlanStatusUpdateController.UpdateIncidentActionPlanStatusUpdate)
	incidentRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "incident_action_plan", "view"), incidentActionPlanStatusUpdateController.GetLatestIncidentActionPlanStatusUpdate)
	incidentRoutes.PUT("/incidents/:incidentID/action-plans/:actionPlanId/approve", incidentActionPlanStatusUpdateController.ApproveIncidentActionPlan)
	incidentRoutes.PUT("/incidents/:incidentID/action-plans/:actionPlanId/revert", incidentActionPlanStatusUpdateController.RevertIncidentActionPlan)

	// Incident Action Plan Comments
	unauthenticatedRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId/comments/:commentId/document-download", incidentActionPlanCommentController.DownloadIncidentActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId/comments/:commentId/document-preview", incidentActionPlanCommentController.GetIncidentActionPlanCommentDocumentPreview)
	incidentRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentActionPlanCommentController.GetAllComments)
	incidentRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentActionPlanCommentController.GetCommentDetails)
	incidentRoutes.POST("/incidents/:incidentID/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentActionPlanCommentController.CreateComment)
	incidentRoutes.PUT("/incidents/:incidentID/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentActionPlanCommentController.UpdateComment)
	incidentRoutes.POST("/incidents/:incidentID/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentActionPlanCommentController.DeleteComment)

	// Incident Root Cause Analyses
	incidentRoutes.GET("/incidents/:incidentID/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentRootCauseAnalysisController.GetAllIncidentRootCauseAnalyses)
	incidentRoutes.POST("/incidents/:incidentID/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRootCauseAnalysisController.CreateIncidentRootCauseAnalysis)
	incidentRoutes.PUT("/incidents/:incidentID/root-cause-analysis/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRootCauseAnalysisController.UpdateIncidentRootCauseAnalysis)
	incidentRoutes.DELETE("/incidents/:incidentID/root-cause-analysis/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentRootCauseAnalysisController.DeleteIncidentRootCauseAnalyses)

	// Incident Communication logs
	incidentRoutes.GET("/incidents/:incidentID/communication-logs", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentCommunicationLogController.GetAllIncidentCommunicationLogs)
	incidentRoutes.GET("/incidents/:incidentID/communication-logs/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentCommunicationLogController.GetIncidentCommunicationLogDetails)
	incidentRoutes.POST("/incidents/:incidentID/communication-logs", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentCommunicationLogController.CreateIncidentCommunicationLog)
	incidentRoutes.PUT("/incidents/:incidentID/communication-logs/:id", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentCommunicationLogController.UpdateIncidentCommunicationLog)
	incidentRoutes.POST("/incidents/:incidentID/communication-logs/delete", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentCommunicationLogController.DeleteIncidentCommunicationLogs)

	// Incident Closure
	incidentRoutes.GET("/incidents/:incidentID/closure", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentClosureController.GetIncidentClosureDetails)
	incidentRoutes.POST("/incidents/:incidentID/closure", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureController.CreateIncidentClosure)
	incidentRoutes.PUT("/incidents/:incidentID/closure", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureController.UpdateIncidentClosure)
	incidentRoutes.DELETE("/incidents/:incidentID/closure", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureController.DeleteIncidentClosures)

	// Incident Post Review Effectiveness
	incidentRoutes.GET("/incident-post-review-effectivenesses", incidentPostReviewEffectivenessController.GetAllIncidentPostReviewEffectivenesses)
	incidentRoutes.GET("/incident-post-review-effectivenesses/:id", incidentPostReviewEffectivenessController.GetIncidentPostReviewEffectivenessDetails)

	// Incident PostReview
	incidentRoutes.GET("/incidents/:incidentID/post-review", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentPostReviewController.GetIncidentPostReviewDetails)
	incidentRoutes.POST("/incidents/:incidentID/post-review", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentPostReviewController.CreateIncidentPostReview)
	incidentRoutes.PUT("/incidents/:incidentID/post-review", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentPostReviewController.UpdateIncidentPostReview)
	incidentRoutes.DELETE("/incidents/:incidentID/post-review", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentPostReviewController.DeleteIncidentPostReviews)

	// Incident Associated Items
	incidentRoutes.GET("/incidents/:incidentID/controls", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentControls)
	incidentRoutes.GET("/incidents/:incidentID/compliance-requirements", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentComplianceRequirements)
	incidentRoutes.GET("/incidents/:incidentID/risks", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentRisks)
	incidentRoutes.GET("/incidents/:incidentID/processes", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentProcesses)
	incidentRoutes.GET("/incidents/:incidentID/slas", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentSLAs)
	incidentRoutes.GET("/incidents/:incidentID/documents", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentDocuments)
	incidentRoutes.GET("/incidents/:incidentID/assets", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentAssets)
	incidentRoutes.POST("/incidents/:incidentID/assets", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.CreateIncidentAsset)
	incidentRoutes.DELETE("/incidents/:incidentID/assets/:assetId", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.DeleteIncidentAsset)
	incidentRoutes.GET("/incidents/:incidentID/clauses", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentClauses)
	incidentRoutes.PUT("/incidents/:incidentID/clauses", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.UpdateClauses)

	incidentRoutes.POST("/incidents/:incidentID/risks", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.CreateIncidentRisk)
	incidentRoutes.DELETE("/incidents/:incidentID/risks/:riskId", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.DeleteIncidentRisk)

	incidentRoutes.POST("/incidents/:incidentID/documents", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.CreateIncidentDocument)
	incidentRoutes.DELETE("/incidents/:incidentID/documents/:documentId", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.DeleteIncidentDocument)

	incidentRoutes.POST("/incidents/:incidentID/controls", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.CreateIncidentControl)
	incidentRoutes.DELETE("/incidents/:incidentID/controls/:controlId", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.DeleteIncidentControl)

	incidentRoutes.POST("/incidents/:incidentID/processes", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.CreateIncidentProcess)
	incidentRoutes.DELETE("/incidents/:incidentID/processes/:processId", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.DeleteIncidentProcess)

	incidentRoutes.GET("/incidents/:incidentID/strategic-objectives", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.GetAllIncidentStrategicObjectives)
	incidentRoutes.POST("/incidents/:incidentID/strategic-objectives", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.CreateIncidentStrategicObjective)
	incidentRoutes.DELETE("/incidents/:incidentID/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.DeleteIncidentStrategicObjective)

	incidentRoutes.GET("/incidents/:incidentID/kpis", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.GetAllIncidentKPI)
	incidentRoutes.POST("/incidents/:incidentID/kpis", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.CreateIncidentKPI)
	incidentRoutes.DELETE("/incidents/:incidentID/kpis/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.DeleteIncidentKPI)

	incidentRoutes.GET("/incidents/:incidentID/vulnerabilities", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentAssociatedItemController.GetAllIncidentVulnerabilities)
	incidentRoutes.POST("/incidents/:incidentID/vulnerabilities", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.CreateIncidentVulnerability)
	incidentRoutes.DELETE("/incidents/:incidentID/vulnerabilities/:vulnerabilityId", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentAssociatedItemController.DeleteIncidentVulnerability)

	// Incident timelines
	incidentRoutes.GET("/incidents/:incidentID/timelines", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentTimelineController.GetAllIncidentTimelines)

	incidentRoutes.GET("/incidents/:incidentID/export-pdf", IncidentReportController.ExportIncidentReportToPDF)

	// Incident reports
	incidentRoutes.GET("/incidents-reports", IncidentReportController.GetAllIncidentReports)
	incidentRoutes.GET("/incidents-reports/:incidentReportId", IncidentReportController.GetIncidentReportDetails)
	incidentRoutes.POST("/incidents-reports", IncidentReportController.CreateIncidentReport)
	incidentRoutes.DELETE("/incidents-reports/:incidentReportId", IncidentReportController.DeleteIncidentReports)
	incidentRoutes.GET("/incidents-reports/:incidentReportId/sections/:sectionId", IncidentReportController.GetIncidentReportSection)
	incidentRoutes.PUT("/incidents-reports/:incidentReportId/sections/:sectionId", IncidentReportController.UpdateIncidentReportSection)
	incidentRoutes.DELETE("/incidents-reports/:incidentReportId/sections/:sectionId", IncidentReportController.DeleteIncidentReportSection)
	incidentRoutes.GET("/incidents-reports/:incidentReportId/export-pdf", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), IncidentReportController.GetIncidentReportExportZip)

	// Incident Response Team
	incidentRoutes.GET("/incident-response-teams", incidentResponseTeamController.GetAllIncidentResponseTeams)
	incidentRoutes.GET("/incident-response-teams/:id", incidentResponseTeamController.GetIncidentResponseTeamDetails)
	incidentRoutes.POST("/incident-response-teams", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentResponseTeamController.CreateIncidentResponseTeam)
	incidentRoutes.PUT("/incident-response-teams/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentResponseTeamController.UpdateIncidentResponseTeam)
	incidentRoutes.POST("/incident-response-teams/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentResponseTeamController.DeleteIncidentResponseTeams)

	// Incident Impact
	incidentRoutes.GET("/incident-impacts", incidentImpactController.GetAllIncidentImpacts)
	incidentRoutes.GET("/incident-impacts/:id", incidentImpactController.GetIncidentImpactDetails)
	incidentRoutes.POST("/incident-impacts", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentImpactController.CreateIncidentImpact)
	incidentRoutes.PUT("/incident-impacts/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentImpactController.UpdateIncidentImpact)
	incidentRoutes.POST("/incident-impacts/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentImpactController.DeleteIncidentImpacts)
	incidentRoutes.PUT("/incident-impacts/activate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentImpactController.ActivateIncidentImpacts)
	incidentRoutes.PUT("/incident-impacts/deactivate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentImpactController.DeactivateIncidentImpacts)

	// Incident Scope
	incidentRoutes.GET("/incident-scopes", incidentScopeController.GetAllIncidentScopes)
	incidentRoutes.GET("/incident-scopes/:id", incidentScopeController.GetIncidentScopeDetails)
	incidentRoutes.POST("/incident-scopes", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentScopeController.CreateIncidentScope)
	incidentRoutes.PUT("/incident-scopes/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentScopeController.UpdateIncidentScope)
	incidentRoutes.POST("/incident-scopes/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentScopeController.DeleteIncidentScopes)
	incidentRoutes.PUT("/incident-scopes/activate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentScopeController.ActivateIncidentScopes)
	incidentRoutes.PUT("/incident-scopes/deactivate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentScopeController.DeactivateIncidentScopes)

	// Incident Sensitivity Classification
	incidentRoutes.GET("/incident-sensitivity-classifications", incidentSensitivityClassificationController.GetAllIncidentSensitivityClassifications)
	incidentRoutes.GET("/incident-sensitivity-classifications/:id", incidentSensitivityClassificationController.GetIncidentSensitivityClassificationDetails)
	incidentRoutes.POST("/incident-sensitivity-classifications", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSensitivityClassificationController.CreateIncidentSensitivityClassification)
	incidentRoutes.PUT("/incident-sensitivity-classifications/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSensitivityClassificationController.UpdateIncidentSensitivityClassification)
	incidentRoutes.POST("/incident-sensitivity-classifications/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSensitivityClassificationController.DeleteIncidentSensitivityClassifications)
	incidentRoutes.PUT("/incident-sensitivity-classifications/activate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSensitivityClassificationController.ActivateIncidentSensitivityClassifications)
	incidentRoutes.PUT("/incident-sensitivity-classifications/deactivate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSensitivityClassificationController.DeactivateIncidentSensitivityClassifications)

	// Incident Evidence Type
	incidentRoutes.GET("/incident-evidence-types", incidentEvidenceTypeController.GetAllIncidentEvidenceTypes)
	incidentRoutes.GET("/incident-evidence-types/:id", incidentEvidenceTypeController.GetIncidentEvidenceTypeDetails)
	incidentRoutes.POST("/incident-evidence-types", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentEvidenceTypeController.CreateIncidentEvidenceType)
	incidentRoutes.PUT("/incident-evidence-types/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentEvidenceTypeController.UpdateIncidentEvidenceType)
	incidentRoutes.POST("/incident-evidence-types/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentEvidenceTypeController.DeleteIncidentEvidenceTypes)
	incidentRoutes.PUT("/incident-evidence-types/activate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentEvidenceTypeController.ActivateIncidentEvidenceTypes)
	incidentRoutes.PUT("/incident-evidence-types/deactivate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentEvidenceTypeController.DeactivateIncidentEvidenceTypes)

	// Incident Evidence Source
	incidentRoutes.GET("/incident-evidence-sources", incidentEvidenceSourceController.GetAllIncidentEvidenceSources)
	incidentRoutes.GET("/incident-evidence-sources/:id", incidentEvidenceSourceController.GetIncidentEvidenceSourceDetails)
	incidentRoutes.POST("/incident-evidence-sources", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentEvidenceSourceController.CreateIncidentEvidenceSource)
	incidentRoutes.PUT("/incident-evidence-sources/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentEvidenceSourceController.UpdateIncidentEvidenceSource)
	incidentRoutes.POST("/incident-evidence-sources/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentEvidenceSourceController.DeleteIncidentEvidenceSources)
	incidentRoutes.PUT("/incident-evidence-sources/activate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentEvidenceSourceController.ActivateIncidentEvidenceSources)
	incidentRoutes.PUT("/incident-evidence-sources/deactivate", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentEvidenceSourceController.DeactivateIncidentEvidenceSources)

	// Incident Evidence
	unauthenticatedRoutes.GET("/incidents/:incidentID/evidences/:evidenceId/document-preview", incidentEvidenceController.GetIncidentEvidenceDocumentPreview)
	incidentRoutes.GET("/incidents/:incidentID/evidences/:evidenceId/download", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentEvidenceController.DownloadIncidentEvidence)
	incidentRoutes.GET("/incidents/:incidentID/evidences", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentEvidenceController.GetAllIncidentEvidences)
	incidentRoutes.GET("/incidents/:incidentID/evidences/:evidenceId", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentEvidenceController.GetIncidentEvidenceDetails)
	incidentRoutes.POST("/incidents/:incidentID/evidences", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentEvidenceController.CreateIncidentEvidence)
	incidentRoutes.PUT("/incidents/:incidentID/evidences-update", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentEvidenceController.UpdateIncidentEvidence)
	incidentRoutes.POST("/incidents/:incidentID/evidences/delete", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentEvidenceController.DeleteIncidentEvidences)
	incidentRoutes.GET("/incidents/:incidentID/evidences/:evidenceId/download-histories", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentEvidenceController.GetIncidentEvidenceDownloadHistories)

	// Incident SLA Escalation Matrix
	incidentRoutes.GET("/incident-sla-escalation-matrix/escalation", incidentSLAEscalationMatrixController.GetAllEscalationMatrices)
	incidentRoutes.PUT("/incident-sla-escalation-matrix/update-escalation", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSLAEscalationMatrixController.UpdateEscalationMatrix)
	incidentRoutes.GET("/incident-sla-escalation-matrix/resolution", incidentSLAEscalationMatrixController.GetAllResolutionMatrices)
	incidentRoutes.PUT("/incident-sla-escalation-matrix/update-resolution", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), incidentSLAEscalationMatrixController.UpdateResolutionMatrix)

	// Incident Closure Workflow / Review
	incidentRoutes.GET("/incident-closures/:closureId/workflow-list", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentClosureWorkflowController.GetWorkflowList)
	incidentRoutes.POST("/incident-closures/:closureId/submit", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureReviewController.SubmitIncidentClosure)
	incidentRoutes.POST("/incident-closures/:closureId/approve", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureReviewController.ApproveIncidentClosure)
	incidentRoutes.POST("/incident-closures/:closureId/revert", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureReviewController.RevertIncidentClosure)

	// Incident Closure Review Comments
	unauthenticatedRoutes.GET("/incident-closures/:closureId/review-comments/:commentId/document-download", incidentClosureReviewCommentController.DownloadDocument)
	unauthenticatedRoutes.GET("/incident-closures/:closureId/review-comments/:commentId/document-preview", incidentClosureReviewCommentController.GetDocumentPreview)
	incidentRoutes.GET("/incident-closures/:closureId/review-comments", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentClosureReviewCommentController.GetAllComments)
	incidentRoutes.GET("/incident-closures/:closureId/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "incident_register", "view"), incidentClosureReviewCommentController.GetCommentDetails)
	incidentRoutes.POST("/incident-closures/:closureId/review-comments", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureReviewCommentController.CreateComment)
	incidentRoutes.PUT("/incident-closures/:closureId/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureReviewCommentController.UpdateComment)
	incidentRoutes.POST("/incident-closures/:closureId/review-comments/delete", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureReviewCommentController.DeleteComments)
	incidentRoutes.PUT("/incident-closures/:closureId/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(db, "incident_register", "modify"), incidentClosureReviewCommentController.ResolveComment)

}
