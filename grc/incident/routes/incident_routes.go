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
	incidentPostReviewEffectivenessController := app.IncidentModule.IncidentPostReviewEffectivenessController
	incidentPostReviewController := app.IncidentModule.IncidentPostReviewController
	incidentAssociatedItemController := app.IncidentModule.IncidentAssociatedItemController
	incidentTimelineController := app.IncidentModule.IncidentTimelineController
	incidentActionPlanCommentController := app.IncidentModule.IncidentActionPlanCommentController
	IncidentReportController := app.IncidentModule.IncidentReportController
	violationReportChannelController := app.IncidentModule.ViolationReportChannelController
	violationReportReasonController := app.IncidentModule.ViolationReportReasonController
	violationReportCategoryController := app.IncidentModule.ViolationReportCategoryController
	violationReportSubCategoryController := app.IncidentModule.ViolationReportSubCategoryController
	violationReportStatusController := app.IncidentModule.ViolationReportStatusController
	violationReportController := app.IncidentModule.ViolationReportController

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

	// Violation report channels
	incidentRoutes.GET("/violation-report-channels", violationReportChannelController.GetAllViolationReportChannels)
	incidentRoutes.GET("/violation-report-channels/:id", violationReportChannelController.GetViolationReportChannelDetails)
	incidentRoutes.POST("/violation-report-channels", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportChannelController.CreateViolationReportChannel)
	incidentRoutes.PUT("/violation-report-channels/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportChannelController.UpdateViolationReportChannel)
	incidentRoutes.POST("/violation-report-channels/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportChannelController.DeleteViolationReportChannels)

	// Violation report reasons
	incidentRoutes.GET("/violation-report-reasons", violationReportReasonController.GetAllViolationReportReasons)
	incidentRoutes.GET("/violation-report-reasons/:id", violationReportReasonController.GetViolationReportReasonDetails)
	incidentRoutes.POST("/violation-report-reasons", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportReasonController.CreateViolationReportReason)
	incidentRoutes.PUT("/violation-report-reasons/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportReasonController.UpdateViolationReportReason)
	incidentRoutes.POST("/violation-report-reasons/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportReasonController.DeleteViolationReportReasons)

	// Violation report categories
	incidentRoutes.GET("/violation-report-categories", violationReportCategoryController.GetAllViolationReportCategories)
	incidentRoutes.GET("/violation-report-categories/:id", violationReportCategoryController.GetViolationReportCategoryDetails)
	incidentRoutes.POST("/violation-report-categories", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportCategoryController.CreateViolationReportCategory)
	incidentRoutes.PUT("/violation-report-categories/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportCategoryController.UpdateViolationReportCategory)
	incidentRoutes.POST("/violation-report-categories/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportCategoryController.DeleteViolationReportCategories)

	// Violation report sub categories
	incidentRoutes.GET("/violation-report-sub-categories", violationReportSubCategoryController.GetAllViolationReportSubCategories)
	incidentRoutes.GET("/violation-report-sub-categories/:id", violationReportSubCategoryController.GetViolationReportSubCategoryDetails)
	incidentRoutes.POST("/violation-report-sub-categories", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportSubCategoryController.CreateViolationReportSubCategory)
	incidentRoutes.PUT("/violation-report-sub-categories/:id", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportSubCategoryController.UpdateViolationReportSubCategory)
	incidentRoutes.POST("/violation-report-sub-categories/delete", middlewares.RolePermissionMiddleware(db, "incident_setting", "modify"), violationReportSubCategoryController.DeleteViolationReportSubCategories)

	//Violation report statuses
	incidentRoutes.GET("/violation-report-statuses", violationReportStatusController.GetAllViolationReportStatuses)
	incidentRoutes.GET("/violation-report-statuses/:id", violationReportStatusController.GetViolationReportStatusDetails)

	//Violation Report
	unauthenticatedRoutes.GET("/violation-reports/:id/document-preview", violationReportController.GetViolationReportDocumentPreview)
	unauthenticatedRoutes.GET("/violation-reports/:id/document-download", violationReportController.DownloadViolationReportDocument)
	incidentRoutes.GET("/violation-reports/export", middlewares.RolePermissionMiddleware(db, "violation_report", "view"), violationReportController.GetViolationReportExport)
	incidentRoutes.GET("/violation-reports", middlewares.RolePermissionMiddleware(db, "violation_report", "view"), violationReportController.GetAllViolationReports)
	incidentRoutes.GET("/violation-reports/:id", middlewares.RolePermissionMiddleware(db, "violation_report", "view"), violationReportController.GetViolationReportDetails)
	incidentRoutes.POST("/violation-reports", middlewares.RolePermissionMiddleware(db, "violation_report", "modify"), violationReportController.CreateViolationReport)
	incidentRoutes.PUT("/violation-reports/:id", middlewares.RolePermissionMiddleware(db, "violation_report", "modify"), violationReportController.UpdateViolationReport)
	incidentRoutes.POST("/violation-reports/delete", middlewares.RolePermissionMiddleware(db, "violation_report", "modify"), violationReportController.DeleteViolationReports)
}
