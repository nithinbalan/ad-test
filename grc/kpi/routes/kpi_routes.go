package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	kpiRoutes := router.Group("/")

	kpiRoutes.Use(middlewares.AuthMiddleware(app.DB))

	db := app.DB

	grcDashboardController := app.KPIModule.GRCDashboardController
	ermDashboardController := app.KPIModule.ERMDashboardController
	complianceDashboardController := app.KPIModule.ComplianceDashboardController
	kriDashboardController := app.KPIModule.KRIDashboardController
	awarenessDashboardController := app.KPIModule.AwarenessDashboardController

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

	kpiRoutes.GET("/kpis/grc-policy-documents", middlewares.RolePermissionMiddleware(db, "grc_dashboard", "view"), grcDashboardController.GetGRCPolicyDocument)

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
	kpiRoutes.GET("/kpis/erm/risk-aging-status", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskAgingStatus)
	kpiRoutes.GET("/kpis/erm/risk-treatment-plan-progress", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskTreatmentPlanProgress)
	kpiRoutes.GET("/kpis/erm/strategic-objectives-risk-overview", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMStrategicObjectivesRiskOverview)
	kpiRoutes.GET("/kpis/erm/strategic-objectives-risk-summary", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMStrategicObjectivesRiskSummary)
	kpiRoutes.GET("/kpis/erm/risk-summary-by-organization-unit", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskSummaryByOrganizationUnit)
	kpiRoutes.GET("/kpis/erm/risk-heat-map", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskHeatMap)

	// ERM Dashbaord Risk Details
	kpiRoutes.GET("/kpis/erm-details/category", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskCategorySummary)
	kpiRoutes.GET("/kpis/erm-details/source", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskSourceSummary)
	kpiRoutes.GET("/kpis/erm-details/status", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskStatusSummary)
	kpiRoutes.GET("/kpis/erm-details/organization-unit", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskOrganizationUnitSummary)
	kpiRoutes.GET("/kpis/erm-details/identified", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRisksIdentifiedSummary)
	kpiRoutes.GET("/kpis/erm-details/response-strategy", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskResponseStrategySummary)
	kpiRoutes.GET("/kpis/erm-details/control-efficiency", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskControlEfficiencySummary)
	kpiRoutes.GET("/kpis/erm-details/mitigation-plan-status", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMRiskMitigationPlanStatusSummary)
	kpiRoutes.GET("/kpis/erm-details/open-risks", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetERMOpenRisksSummary)

	kpiRoutes.GET("/kpis/erm/top-twenty-risks", middlewares.RolePermissionMiddleware(db, "erm_dashboard", "view"), ermDashboardController.GetTopRiskRegisters)

	// Compliance Dashboard
	kpiRoutes.GET("/kpis/compliance/overview", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetComplianceOverview)
	kpiRoutes.GET("/kpis/compliance/leadership-attention", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetLeadershipAttentionSummary)
	kpiRoutes.GET("/kpis/compliance/breakdown", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetComplianceBreakdown)
	kpiRoutes.GET("/kpis/compliance/high-performance", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetHighPerformanceCompliances)
	kpiRoutes.GET("/kpis/compliance/least-performance", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetLeastPerformanceCompliances)
	kpiRoutes.GET("/kpis/compliance/high-performance/clauses", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetHighPerformanceComplianceClauses)
	kpiRoutes.GET("/kpis/compliance/least-performance/clauses", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetLeastPerformanceComplianceClauses)
	kpiRoutes.GET("/kpis/compliance/overview/policies", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetPolicyOverview)
	kpiRoutes.GET("/kpis/compliance/summary/policy-statuses", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetPolicyStatusBreakdown)
	kpiRoutes.GET("/kpis/compliance/summary/policy-organization-unit", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetPolicyByOrganizationUnit)
	kpiRoutes.GET("/kpis/compliance/summary/controls", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetControlSummary)
	kpiRoutes.GET("/kpis/compliance/summary/control-efficiency", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetComplianceEfficiencySummary)
	kpiRoutes.GET("/kpis/compliance/overview/controls", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetControlOverview)
	kpiRoutes.GET("/kpis/compliance/summary/control-organization-unit", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetOrganizationUnitWiseAssessmentCoverage)
	kpiRoutes.GET("/kpis/compliance/overview/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetActionPlanOverview)
	kpiRoutes.GET("/kpis/compliance/summary/action-plan-statuses", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetActionPlanStatusSummary)
	kpiRoutes.GET("/kpis/compliance/summary/action-plan-organization-units", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetActionPlansByOrganizationUnit)

	kpiRoutes.GET("/kpis/compliance/top-twenty-compliances", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetTopCompliances)
	kpiRoutes.GET("/kpis/compliance/summary/rating-distribution", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetComplianceCountByRating)
	kpiRoutes.GET("/kpis/compliance/summary/register-classification", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDashboardController.GetComplianceRegisterClassificationSummary)

	// KRI Dashboard
	kpiRoutes.GET("/kpis/kri/overview", middlewares.RolePermissionMiddleware(db, "kri_dashboard", "view"), kriDashboardController.GetKRIOverview)
	kpiRoutes.GET("/kpis/kri/review-status", middlewares.RolePermissionMiddleware(db, "kri_dashboard", "view"), kriDashboardController.GetKRIReviewStatus)
	kpiRoutes.GET("/kpis/kri/status-category", middlewares.RolePermissionMiddleware(db, "kri_dashboard", "view"), kriDashboardController.GetKRIStatusCategory)
	kpiRoutes.GET("/kpis/kri/severity-breakdown", middlewares.RolePermissionMiddleware(db, "kri_dashboard", "view"), kriDashboardController.GetKRISeverityBreakdown)
	kpiRoutes.GET("/kpis/kri/top-risk-exposure", middlewares.RolePermissionMiddleware(db, "kri_dashboard", "view"), kriDashboardController.GetTopKRIRiskExposure)
	kpiRoutes.GET("/kpis/kri/severity-organization-unit", middlewares.RolePermissionMiddleware(db, "kri_dashboard", "view"), kriDashboardController.GetKRISeverityByOrganizationUnit)
	kpiRoutes.GET("/kpis/kri/severity-category", middlewares.RolePermissionMiddleware(db, "kri_dashboard", "view"), kriDashboardController.GetKRISeverityByCategory)
	kpiRoutes.GET("/kpis/kri/breaches", middlewares.RolePermissionMiddleware(db, "kri_dashboard", "view"), kriDashboardController.GetKRIMostFrequentBreaches)

	// Awarness Dashboard
	kpiRoutes.GET("/kpis/awareness-reporting/summary", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetAwarenessPlanSummary)
	kpiRoutes.GET("/kpis/awareness-reporting/campaigns-by-type", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetCampaignsByType)
	kpiRoutes.GET("/kpis/awareness-reporting/priority-distribution", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetPriorityDistribution)
	kpiRoutes.GET("/kpis/awareness-reporting/campaigns-per-month", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetCampaignsPerMonth)
	kpiRoutes.GET("/kpis/awareness-reporting/channel-mix-per-month", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetChannelMixPerMonth)
	kpiRoutes.GET("/kpis/awareness-reporting/coverage-by-audience", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetCoverageByAudience)

	kpiRoutes.GET("/kpis/awareness-reporting/progress/summary", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetAwarenessProgressSummary)
	kpiRoutes.GET("/kpis/awareness-reporting/progress/completion-by-channel", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetCompletionByChannel)
	kpiRoutes.GET("/kpis/awareness-reporting/progress/completion-per-campaign", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetCompletionPerCampaign)
	kpiRoutes.GET("/kpis/awareness-reporting/progress/survey-completion", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetSurveyCompletionStatus)
	kpiRoutes.GET("/kpis/awareness-reporting/progress/completion-by-department", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetCompletionByDepartment)
	kpiRoutes.GET("/kpis/awareness-reporting/progress/phishing-results", middlewares.RolePermissionMiddleware(db, "awareness_reporting_dashboard", "view"), awarenessDashboardController.GetPhishingResults)
}
