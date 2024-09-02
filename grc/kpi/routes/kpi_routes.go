package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	kpiRoutes := router.Group("/")

	kpiRoutes.Use(middlewares.AuthMiddleware())

	grcDashboardController := app.KPIModule.GRCDashboardController

	//GRCDashboard Summary
	kpiRoutes.GET("/kpis/grc-summary/compliance", grcDashboardController.GetGRCComplianceSummary)
	kpiRoutes.GET("/kpis/grc-summary/risk", grcDashboardController.GetGRCRiskSummary)
	kpiRoutes.GET("/kpis/grc-summary/audit", grcDashboardController.GetGRCAuditSummary)
	kpiRoutes.GET("/kpis/grc-summary/document", grcDashboardController.GetGRCDocumentSummary)
	kpiRoutes.GET("/kpis/grc-summary/organization", grcDashboardController.GetGRCOrganizationSummary)

	kpiRoutes.GET("/kpis/grc-overview/compliance", grcDashboardController.GetGRCComplianceOverview)
	kpiRoutes.GET("/kpis/grc-overview/risk", grcDashboardController.GetGRCRiskOverview)

	kpiRoutes.GET("/kpis/grc-top-five/compliances", grcDashboardController.GetGRCTopCompliances)
	kpiRoutes.GET("/kpis/grc-top-five/risks", grcDashboardController.GetGRCTopRisks)

}
