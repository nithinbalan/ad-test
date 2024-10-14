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

	auditInsightsAuditPlanController := app.AuditInsightsModule.AuditInsightsAuditPlanController
	auditInsightsAuditController := app.AuditInsightsModule.AuditInsightsAuditController
	auditInsightsDocumentRequestController := app.AuditInsightsModule.AuditInsightsDocumentRequestController
	auditInsightsDocumentRequestConversationController := app.AuditInsightsModule.AuditInsightsDocumentRequestConversationController
	auditInsightsDraftReportController := app.AuditInsightsModule.AuditInsightsDraftReportController
	auditDraftFindingController := app.AuditInsightsModule.AuditInsightsDraftFindingController
	auditInsightsDraftFindingConversationController := app.AuditInsightsModule.AuditInsightsDraftFindingConversationController
	auditInsightsDraftFindingActionPlanController := app.AuditInsightsModule.AuditInsightsDraftFindingActionPlanController

	// Audit Insights Audit plans
	auditInsightsRoutes.GET("/audit-insights/audit-plans", auditInsightsAuditPlanController.GetAllAuditInsightsAuditPlans)
	auditInsightsRoutes.GET("/audit-insights/audit-plans/:auditPlanId", auditInsightsAuditPlanController.GetAuditInsightsAuditPlanDetails)

	// Audit Insights Audits
	auditInsightsRoutes.GET("/audit-insights/audits", auditInsightsAuditController.GetAllAuditInsightsAudits)
	auditInsightsRoutes.GET("/audit-insights/audits/:auditId", auditInsightsAuditController.GetAuditInsightsAuditDetails)

	// Audit Insights Audits Document Requests
	auditInsightsRoutes.GET("/audit-insights/document-requests", auditInsightsDocumentRequestController.GetAllAuditInsightsDocumentRequests)
	auditInsightsRoutes.GET("/audit-insights/document-requests/:documentRequestId", auditInsightsDocumentRequestController.GetAuditInsightsDocumentRequestDetails)

	// Audit Insights Document Requests Add conversations
	auditInsightsRoutes.GET("/audit-insights/document-requests/:documentRequestId/conversations", auditInsightsDocumentRequestConversationController.GetAllAuditInsightsDocumentRequestConversations)
	auditInsightsRoutes.POST("/audit-insights/document-requests/:documentRequestId/conversations", auditInsightsDocumentRequestConversationController.CreateAuditInsightsDocumentRequestConversation)
	unauthenticatedRoutes.GET("/audit-insights/document-requests/:documentRequestId/conversations/:conversationId/document-download", auditInsightsDocumentRequestConversationController.DownloadDocumentRequestConversationDocument)

	// Audit Insights Draft Report
	auditInsightsRoutes.GET("/audit-insights/draft-reports", auditInsightsDraftReportController.GetAllAuditInsightsDraftReports)
	auditInsightsRoutes.GET("/audit-insights/draft-reports/:id", auditInsightsDraftReportController.GetAuditInsightsDraftReportDetails)

	//Audit Insights Draft Findings
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/accept", auditDraftFindingController.AcceptAuditDraftFinding)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/challenge", auditDraftFindingController.ChallengeAuditDraftFinding)

	// Audit Insights  Draft Finding Conversations
	auditInsightsRoutes.GET("/audit-insights/audit-draft-findings/:findingId/conversations", auditInsightsDraftFindingConversationController.GetAllAuditInsightsDraftFindingConversations)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/conversations", auditInsightsDraftFindingConversationController.CreateAuditInsightsDraftFindingConversation)

	// Audit Insights  Draft Finding Action plans
	auditInsightsRoutes.GET("/audit-insights/audit-draft-findings/:findingId/action-plans", auditInsightsDraftFindingActionPlanController.GetAllAuditInsightsDraftFindingActionPlans)
	auditInsightsRoutes.GET("/audit-insights/audit-draft-findings/:findingId/action-plans/:id", auditInsightsDraftFindingActionPlanController.GetAuditInsightsDraftFindingActionPlanDetails)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/action-plans", auditInsightsDraftFindingActionPlanController.CreateAuditInsightsDraftFindingActionPlan)
	auditInsightsRoutes.PUT("/audit-insights/audit-draft-findings/:findingId/action-plans/:id", auditInsightsDraftFindingActionPlanController.UpdateAuditInsightsDraftFindingActionPlan)
	auditInsightsRoutes.POST("/audit-insights/audit-draft-findings/:findingId/action-plans/delete", auditInsightsDraftFindingActionPlanController.DeleteAuditInsightsDraftFindingActionPlans)

}
