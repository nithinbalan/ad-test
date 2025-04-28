package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	auditInsightsRoutes := router.Group("/")
	unauthenticatedRoutes := router.Group("/")
	auditInsightsRoutes.Use(middlewares.AuthMiddleware())

	db := app.DB

	auditInsightsAuditPlanController := app.AuditInsightsModule.AuditInsightsAuditPlanController
	auditInsightsAuditController := app.AuditInsightsModule.AuditInsightsAuditController
	auditInsightsDocumentRequestController := app.AuditInsightsModule.AuditInsightsDocumentRequestController
	auditInsightsDocumentRequestConversationController := app.AuditInsightsModule.AuditInsightsDocumentRequestConversationController
	auditInsightsDraftReportController := app.AuditInsightsModule.AuditInsightsDraftReportController
	auditDraftFindingController := app.AuditInsightsModule.AuditInsightsDraftFindingController
	auditInsightsDraftFindingConversationController := app.AuditInsightsModule.AuditInsightsDraftFindingConversationController
	auditInsightsDraftFindingActionPlanController := app.AuditInsightsModule.AuditInsightsDraftFindingActionPlanController
	auditInsightsAuditFindingController := app.AuditInsightsModule.AuditInsightsAuditFindingController
	auditInsightsFindingRootCauseAnalysisController := app.AuditInsightsModule.AuditInsightsFindingRootCauseAnalysisController
	auditInsightsFindingActionPlanController := app.AuditInsightsModule.AuditInsightsFindingActionPlanController
	auditInsightsFindingActionPlanStatusUpdateController := app.AuditInsightsModule.AuditInsightsFindingActionPlanStatusUpdateController
	auditInsightsFinalReportController := app.AuditInsightsModule.AuditInsightsFinalReportController

	// Audit Insights Audit plans
	auditInsightsRoutes.GET("/audit-insights/audit-plans", middlewares.RolePermissionMiddleware(db, "audit_insight_plans", "view"), auditInsightsAuditPlanController.GetAllAuditInsightsAuditPlans)
	auditInsightsRoutes.GET("/audit-insights/audit-plans/export", middlewares.RolePermissionMiddleware(db, "audit_insight_plans", "view"), auditInsightsAuditPlanController.GetAuditInsightsAuditPlansExport)
	auditInsightsRoutes.GET("/audit-insights/audit-plans/:auditPlanId", middlewares.RolePermissionMiddleware(db, "audit_insight_plans", "view"), auditInsightsAuditPlanController.GetAuditInsightsAuditPlanDetails)

	// Audit Insights Audits
	auditInsightsRoutes.GET("/audit-insights/audits", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsAuditController.GetAllAuditInsightsAudits)
	auditInsightsRoutes.GET("/audit-insights/audits/export", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsAuditController.GetAuditInsightsAuditsExport)

	auditInsightsRoutes.GET("/audit-insights/audits/:auditId", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsAuditController.GetAuditInsightsAuditDetails)

	// Audit Insights Audits Document Requests
	auditInsightsRoutes.GET("/audit-insights/document-requests", middlewares.RolePermissionMiddleware(db, "document_requests", "view"), auditInsightsDocumentRequestController.GetAllAuditInsightsDocumentRequests)
	auditInsightsRoutes.GET("/audit-insights/document-requests/:documentRequestId", middlewares.RolePermissionMiddleware(db, "document_requests", "view"), auditInsightsDocumentRequestController.GetAuditInsightsDocumentRequestDetails)

	// Audit Insights Document Requests Add conversations
	auditInsightsRoutes.GET("/audit-insights/document-requests/:documentRequestId/conversations", middlewares.RolePermissionMiddleware(db, "document_requests", "view"), auditInsightsDocumentRequestConversationController.GetAllAuditInsightsDocumentRequestConversations)
	auditInsightsRoutes.POST("/audit-insights/document-requests/:documentRequestId/conversations", middlewares.RolePermissionMiddleware(db, "document_requests", "modify"), auditInsightsDocumentRequestConversationController.CreateAuditInsightsDocumentRequestConversation)
	unauthenticatedRoutes.GET("/audit-insights/document-requests/:documentRequestId/conversations/:conversationId/document-download", auditInsightsDocumentRequestConversationController.DownloadDocumentRequestConversationDocument)
	unauthenticatedRoutes.GET("/audit-insights/document-requests/:documentRequestId/conversations/:conversationId/document-preview", auditInsightsDocumentRequestConversationController.GetDocumentRequestConversationDocumentPreview)

	// Audit Insights Draft Report
	auditInsightsRoutes.GET("/audit-insights/draft-reports", middlewares.RolePermissionMiddleware(db, "draft_reports", "view"), auditInsightsDraftReportController.GetAllAuditInsightsDraftReports)
	auditInsightsRoutes.GET("/audit-insights/draft-reports/export", middlewares.RolePermissionMiddleware(db, "draft_reports", "view"), auditInsightsDraftReportController.GetAuditInsightsDraftReportExport)
	auditInsightsRoutes.GET("/audit-insights/draft-reports/:id", middlewares.RolePermissionMiddleware(db, "draft_reports", "view"), auditInsightsDraftReportController.GetAuditInsightsDraftReportDetails)

	//Audit Insights Draft Findings
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/accept", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditDraftFindingController.AcceptAuditDraftFinding)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/challenge", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditDraftFindingController.ChallengeAuditDraftFinding)

	// Audit Insights  Draft Finding Conversations
	auditInsightsRoutes.GET("/audit-insights/audit-draft-findings/:findingId/conversations", middlewares.RolePermissionMiddleware(db, "draft_reports", "view"), auditInsightsDraftFindingConversationController.GetAllAuditInsightsDraftFindingConversations)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/conversations", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditInsightsDraftFindingConversationController.CreateAuditInsightsDraftFindingConversation)
	unauthenticatedRoutes.GET("/audit-insights/audit-draft-findings/:findingId/conversations/:id/document-download", auditInsightsDraftFindingConversationController.DownloadAuditDraftFindingConversationDocument)
	unauthenticatedRoutes.GET("/audit-insights/audit-draft-findings/:findingId/conversations/:id/document-preview", auditInsightsDraftFindingConversationController.GetAuditDraftFindingConversationDocumentPreview)

	// Audit Insights  Draft Finding Action plans
	auditInsightsRoutes.GET("/audit-insights/audit-draft-findings/:findingId/action-plans", middlewares.RolePermissionMiddleware(db, "draft_reports", "view"), auditInsightsDraftFindingActionPlanController.GetAllAuditInsightsDraftFindingActionPlans)
	auditInsightsRoutes.GET("/audit-insights/audit-draft-findings/:findingId/action-plans/:id", middlewares.RolePermissionMiddleware(db, "draft_reports", "view"), auditInsightsDraftFindingActionPlanController.GetAuditInsightsDraftFindingActionPlanDetails)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/action-plans", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditInsightsDraftFindingActionPlanController.CreateAuditInsightsDraftFindingActionPlan)
	auditInsightsRoutes.PUT("/audit-insights/audit-draft-findings/:findingId/action-plans/:id", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditInsightsDraftFindingActionPlanController.UpdateAuditInsightsDraftFindingActionPlan)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditInsightsDraftFindingActionPlanController.DeleteAuditInsightsDraftFindingActionPlans)

	//Audit Insights  Findings
	auditInsightsRoutes.GET("/audit-insights/audit-findings", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsAuditFindingController.GetAllAuditInsightsAuditFindings)
	auditInsightsRoutes.GET("/audit-insights/audit-findings/export", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsAuditFindingController.GetAuditInsightsAuditFindingsExport)
	auditInsightsRoutes.GET("/audit-insights/audit-findings/:findingId", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsAuditFindingController.GetAuditInsightsAuditFindingDetails)
	auditInsightsRoutes.GET("/audits-insights/audits/:auditId/audit-finding-summary", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsAuditFindingController.GetAuditInsightFindingCountByCategory)

	auditInsightsRoutes.GET("/audit-insights/audit-findings/:findingId/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsFindingRootCauseAnalysisController.GetAllAuditInsightsFindingRootCauseAnalyses)
	auditInsightsRoutes.POST("/audit-insights/audit-findings/:findingId/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingRootCauseAnalysisController.CreateAuditInsightsFindingRootCauseAnalysis)
	auditInsightsRoutes.PUT("/audit-insights/audit-findings/:findingId/root-cause-analysis/:analysisId", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingRootCauseAnalysisController.UpdateAuditInsightsFindingRootCauseAnalysis)
	auditInsightsRoutes.DELETE("/audit-insights/audit-findings/:findingId/root-cause-analysis/:analysisId", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingRootCauseAnalysisController.DeleteAuditInsightsFindingRootCauseAnalyses)

	//Audit Insights action plans
	auditInsightsRoutes.GET("/audit-insights/action-plans", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsFindingActionPlanController.GetAllAuditInsightsFindingActionPlans)
	auditInsightsRoutes.GET("/audit-insights/action-plans/export", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsFindingActionPlanController.GetAuditInsightsFindingActionPlansExport)
	auditInsightsRoutes.GET("/audit-insights/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsFindingActionPlanController.GetAuditInsightsFindingActionPlanDetails)
	auditInsightsRoutes.POST("/audit-insights/action-plans", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingActionPlanController.CreateAuditInsightsFindingActionPlan)
	auditInsightsRoutes.PUT("/audit-insights/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingActionPlanController.UpdateAuditInsightsFindingActionPlan)
	auditInsightsRoutes.POST("/audit-insights/action-plans/delete", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingActionPlanController.DeleteAuditInsightsFindingActionPlans)

	auditInsightsRoutes.GET("/audit-insights/audit-findings/:findingId/action-plans", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsFindingActionPlanController.GetAllAuditInsightsFindingActionPlans)
	auditInsightsRoutes.GET("/audit-insights/audit-findings/:findingId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsFindingActionPlanController.GetAuditInsightsFindingActionPlanDetails)
	auditInsightsRoutes.POST("/audit-insights/audit-findings/:findingId/action-plans", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingActionPlanController.CreateAuditInsightsFindingActionPlan)
	auditInsightsRoutes.PUT("/audit-insights/audit-findings/:findingId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingActionPlanController.UpdateAuditInsightsFindingActionPlan)
	auditInsightsRoutes.POST("/audit-insights/audit-findings/:findingId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingActionPlanController.DeleteAuditInsightsFindingActionPlans)

	auditInsightsRoutes.GET("/audit-insights/audit-findings/:findingId/action-plan-summary", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsFindingActionPlanController.GetAuditInsightsFindingActionPlanCountByStatus)

	auditInsightsRoutes.GET("/audit-insights/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsFindingActionPlanStatusUpdateController.GetAllAuditInsightsFindingActionPlanStatusUpdates)
	auditInsightsRoutes.POST("/audit-insights/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingActionPlanStatusUpdateController.CreateAuditInsightsFindingActionPlanStatusUpdate)
	auditInsightsRoutes.PUT("/audit-insights/action-plans/:actionPlanId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingActionPlanStatusUpdateController.UpdateAuditInsightsFindingActionPlanStatusUpdate)
	auditInsightsRoutes.GET("/audit-insights/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsFindingActionPlanStatusUpdateController.GetLatestAuditInsightsFindingActionPlanStatusUpdate)

	auditInsightsRoutes.GET("/audit-insights/audit-findings/:findingId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsFindingActionPlanStatusUpdateController.GetAllAuditInsightsFindingActionPlanStatusUpdates)
	auditInsightsRoutes.POST("/audit-insights/audit-findings/:findingId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingActionPlanStatusUpdateController.CreateAuditInsightsFindingActionPlanStatusUpdate)
	auditInsightsRoutes.PUT("/audit-insights/audit-findings/:findingId/action-plans/:actionPlanId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "modify"), auditInsightsFindingActionPlanStatusUpdateController.UpdateAuditInsightsFindingActionPlanStatusUpdate)
	auditInsightsRoutes.GET("/audit-insights/audit-findings/:findingId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsFindingActionPlanStatusUpdateController.GetLatestAuditInsightsFindingActionPlanStatusUpdate)
	unauthenticatedRoutes.GET("/audit-insights/audit-findings/:findingId/action-plans/:actionPlanId/status-updates/:id/document-download", auditInsightsFindingActionPlanStatusUpdateController.DownloadAuditInsightsFindingActionPlanStatusUpdateDocument)

	// Audit Insights Draft Report
	auditInsightsRoutes.GET("/audit-insights/final-reports", middlewares.RolePermissionMiddleware(db, "final_report", "view"), auditInsightsFinalReportController.GetAllAuditInsightsFinalReports)
	auditInsightsRoutes.GET("/audit-insights/final-reports/export", middlewares.RolePermissionMiddleware(db, "final_report", "view"), auditInsightsFinalReportController.GetAuditInsightsFinalReportExport)
	auditInsightsRoutes.GET("/audit-insights/final-reports/:id", middlewares.RolePermissionMiddleware(db, "final_report", "view"), auditInsightsFinalReportController.GetAuditInsightsFinalReportDetails)

}
