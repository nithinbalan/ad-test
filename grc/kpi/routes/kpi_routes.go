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

}
