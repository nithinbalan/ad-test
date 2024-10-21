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

	// Audit Insights Audit plans
	auditInsightsRoutes.GET("/audit-insights/audit-plans", middlewares.RolePermissionMiddleware(db, "audit_insight_plans", "view"), auditInsightsAuditPlanController.GetAllAuditInsightsAuditPlans)
	auditInsightsRoutes.GET("/audit-insights/audit-plans/:auditPlanId", middlewares.RolePermissionMiddleware(db, "audit_insight_plans", "view"), auditInsightsAuditPlanController.GetAuditInsightsAuditPlanDetails)

	// Audit Insights Audits
	auditInsightsRoutes.GET("/audit-insights/audits", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsAuditController.GetAllAuditInsightsAudits)
	auditInsightsRoutes.GET("/audit-insights/audits/:auditId", middlewares.RolePermissionMiddleware(db, "audit_insight_audits", "view"), auditInsightsAuditController.GetAuditInsightsAuditDetails)

	// Audit Insights Audits Document Requests
	auditInsightsRoutes.GET("/audit-insights/document-requests", middlewares.RolePermissionMiddleware(db, "document_requests", "view"), auditInsightsDocumentRequestController.GetAllAuditInsightsDocumentRequests)
	auditInsightsRoutes.GET("/audit-insights/document-requests/:documentRequestId", middlewares.RolePermissionMiddleware(db, "document_requests", "view"), auditInsightsDocumentRequestController.GetAuditInsightsDocumentRequestDetails)

	// Audit Insights Document Requests Add conversations
	auditInsightsRoutes.GET("/audit-insights/document-requests/:documentRequestId/conversations", middlewares.RolePermissionMiddleware(db, "document_requests", "view"), auditInsightsDocumentRequestConversationController.GetAllAuditInsightsDocumentRequestConversations)
	auditInsightsRoutes.POST("/audit-insights/document-requests/:documentRequestId/conversations", middlewares.RolePermissionMiddleware(db, "document_requests", "modify"), auditInsightsDocumentRequestConversationController.CreateAuditInsightsDocumentRequestConversation)
	unauthenticatedRoutes.GET("/audit-insights/document-requests/:documentRequestId/conversations/:conversationId/document-download", auditInsightsDocumentRequestConversationController.DownloadDocumentRequestConversationDocument)

	// Audit Insights Draft Report
	auditInsightsRoutes.GET("/audit-insights/draft-reports", middlewares.RolePermissionMiddleware(db, "draft_reports", "view"), auditInsightsDraftReportController.GetAllAuditInsightsDraftReports)
	auditInsightsRoutes.GET("/audit-insights/draft-reports/:id", middlewares.RolePermissionMiddleware(db, "draft_reports", "view"), auditInsightsDraftReportController.GetAuditInsightsDraftReportDetails)

	//Audit Insights Draft Findings
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/accept", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditDraftFindingController.AcceptAuditDraftFinding)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/challenge", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditDraftFindingController.ChallengeAuditDraftFinding)

	// Audit Insights  Draft Finding Conversations
	auditInsightsRoutes.GET("/audit-insights/audit-draft-findings/:findingId/conversations", middlewares.RolePermissionMiddleware(db, "draft_reports", "view"), auditInsightsDraftFindingConversationController.GetAllAuditInsightsDraftFindingConversations)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/conversations", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditInsightsDraftFindingConversationController.CreateAuditInsightsDraftFindingConversation)

	// Audit Insights  Draft Finding Action plans
	auditInsightsRoutes.GET("/audit-insights/audit-draft-findings/:findingId/action-plans", middlewares.RolePermissionMiddleware(db, "draft_reports", "view"), auditInsightsDraftFindingActionPlanController.GetAllAuditInsightsDraftFindingActionPlans)
	auditInsightsRoutes.GET("/audit-insights/audit-draft-findings/:findingId/action-plans/:id", middlewares.RolePermissionMiddleware(db, "draft_reports", "view"), auditInsightsDraftFindingActionPlanController.GetAuditInsightsDraftFindingActionPlanDetails)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/action-plans", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditInsightsDraftFindingActionPlanController.CreateAuditInsightsDraftFindingActionPlan)
	auditInsightsRoutes.PUT("/audit-insights/audit-draft-findings/:findingId/action-plans/:id", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditInsightsDraftFindingActionPlanController.UpdateAuditInsightsDraftFindingActionPlan)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "draft_reports", "modify"), auditInsightsDraftFindingActionPlanController.DeleteAuditInsightsDraftFindingActionPlans)

}
