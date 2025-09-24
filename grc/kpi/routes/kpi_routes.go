package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	kpiRoutes := router.Group("/")

	kpiRoutes.Use(middlewares.AuthMiddleware())

	db := app.DB

	grcDashboardController := app.KPIModule.GRCDashboardController
	ermDashboardController := app.KPIModule.ERMDashboardController

	//GRCDashboard Summary
	kpiRoutes.GET("/kpis/grc-summary/compliance", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCComplianceSummary)
	kpiRoutes.GET("/kpis/grc-summary/risk", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCRiskSummary)
	kpiRoutes.GET("/kpis/grc-summary/audit", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCAuditSummary)
	kpiRoutes.GET("/kpis/grc-summary/document", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCDocumentSummary)
	kpiRoutes.GET("/kpis/grc-summary/organization", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCOrganizationSummary)
	kpiRoutes.GET("/kpis/grc-summary/incident", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCIncidentSummary)

	kpiRoutes.GET("/kpis/grc-overview/compliance", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCComplianceOverview)
	kpiRoutes.GET("/kpis/grc-overview/risk", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCRiskOverview)
	kpiRoutes.GET("/kpis/grc-overview/document", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCDocumentOverview)

	kpiRoutes.GET("/kpis/grc-top-five/compliances", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCTopCompliances)
	kpiRoutes.GET("/kpis/grc-top-five/risks", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCTopRisks)
	kpiRoutes.GET("/kpis/grc-top-five/documents", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCTopPublishedDocuments)

	kpiRoutes.GET("/kpis/grc-chart/documents-by-type", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCDocumentTypeChart)

	kpiRoutes.GET("/kpis/grc-risk-heat-map", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCRiskHeatMap)

	// ERM Dashboard
	kpiRoutes.GET("/kpis/erm/risk-overview", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskOverview)
	kpiRoutes.GET("/kpis/erm/risk-distribution", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskDistribution)
	kpiRoutes.GET("/kpis/erm/overdue-risks", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMOverdueRisks)
	kpiRoutes.GET("/kpis/erm/risks-above-appetite", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRisksAboveAppetite)
	kpiRoutes.GET("/kpis/erm/leadership-attention", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMLeadershipAttentionMetrics)
	kpiRoutes.GET("/kpis/erm/control-efficiency-summary", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMControlEfficiencySummary)
	kpiRoutes.GET("/kpis/erm/control-efficiency-breakdown", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMControlEfficiencyBreakdown)
	kpiRoutes.GET("/kpis/erm/controls-overview", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMControlsOverview)
	kpiRoutes.GET("/kpis/erm/risk-appetite-overview", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskAppetiteOverview)
	kpiRoutes.GET("/kpis/erm/risk-appetite-summary", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskAppetiteSummary)
}
