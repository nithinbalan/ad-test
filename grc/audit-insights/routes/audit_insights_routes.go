package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	auditInsightsRoutes := router.Group("/")

	auditInsightsRoutes.Use(middlewares.AuthMiddleware())

	auditInsightsAuditPlanController := app.AuditInsightsModule.AuditInsightsAuditPlanController
	auditInsightsAuditController := app.AuditInsightsModule.AuditInsightsAuditController
	auditInsightsDocumentRequestController := app.AuditInsightsModule.AuditInsightsDocumentRequestController
	auditInsightsDraftReportController := app.AuditInsightsModule.AuditInsightsDraftReportController

	// Audit Insights Audit plans
	auditInsightsRoutes.GET("/audit-insights/audit-plans", auditInsightsAuditPlanController.GetAllAuditInsightsAuditPlans)
	auditInsightsRoutes.GET("/audit-insights/audit-plans/:id", auditInsightsAuditPlanController.GetAuditInsightsAuditPlanDetails)

	// Audit Insights Audits
	auditInsightsRoutes.GET("/audit-insights/audits", auditInsightsAuditController.GetAllAuditInsightsAudits)
	auditInsightsRoutes.GET("/audit-insights/audits/:id", auditInsightsAuditController.GetAuditInsightsAuditDetails)

	//Audit Insights Audits Document Requests
	auditInsightsRoutes.GET("/audit-insights/document-requests", auditInsightsDocumentRequestController.GetAllAuditInsightsDocumentRequests)
	auditInsightsRoutes.GET("/audit-insights/document-requests/:id", auditInsightsDocumentRequestController.GetAuditInsightsDocumentRequestDetails)

	//Audit Insights Draft Report
	auditInsightsRoutes.GET("/audit-insights/draft-reports", auditInsightsDraftReportController.GetAllAuditInsightsDraftReports)
	auditInsightsRoutes.GET("/audit-insights/audits/:id/draft-reports", auditInsightsDraftReportController.GetAuditInsightsDraftReportDetails)

}
