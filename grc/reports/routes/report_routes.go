package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	reportsRoutes := router.Group("/")

	reportsRoutes.Use(middlewares.AuthMiddleware())

	complianceReportController := app.ReportsModule.ComplianceReportController
	organizationReportController := app.ReportsModule.OrganizationReportController

	reportsRoutes.GET("/reports/compliance-registers", complianceReportController.GetAllComplianceRegisters)
	reportsRoutes.GET("/reports/compliance-registers/export", complianceReportController.ComplianceRegisterReportExport)
	reportsRoutes.GET("/reports/compliance-by-compliance-category", complianceReportController.GetComplianceByComplianceCategory)
	reportsRoutes.GET("/reports/compliance-by-severity-level", complianceReportController.GetComplianceBySeverityLevel)
	reportsRoutes.GET("/reports/compliance-by-compliance-status", complianceReportController.GetComplianceByComplianceStatus)
	reportsRoutes.GET("/reports/compliance-by-responsible-unit", complianceReportController.GetComplianceByResponsibleUnit)
	reportsRoutes.GET("/reports/compliance-by-responsible-user", complianceReportController.GetComplianceByResponsibleUser)

	reportsRoutes.GET("/reports/compliance-requirements", complianceReportController.GetAllComplianceRequirements)
	reportsRoutes.GET("/reports/compliance-requirement-by-compliance-status", complianceReportController.GetComplianceRequirementByComplianceStatus)
	reportsRoutes.GET("/reports/compliance-requirement-by-severity-level", complianceReportController.GetComplianceRequirementBySeverityLevel)
	reportsRoutes.GET("/reports/compliance-requirement-by-category", complianceReportController.GetComplianceRequirementByCategory)
	reportsRoutes.GET("/reports/compliance-requirement-by-responsible-user", complianceReportController.GetComplianceRequirementByResponsibleUser)

	reportsRoutes.GET("/reports/compliance-action-plans", complianceReportController.GetAllComplianceActionPlans)
	reportsRoutes.GET("/reports/compliance-action-plan-by-status", complianceReportController.GetComplianceActionPlanByStatus)
	reportsRoutes.GET("/reports/compliance-action-plan-by-responsible-user", complianceReportController.GetComplianceActionPlanByResponsibleUser)

	reportsRoutes.GET("/reports/slas", complianceReportController.GetAllSLARegisters)
	reportsRoutes.GET("/reports/slas/export", complianceReportController.GetSLARegisterExport)
	reportsRoutes.GET("/reports/sla-by-category", complianceReportController.GetSLAByCategory)
	reportsRoutes.GET("/reports/sla-by-severity-level", complianceReportController.GetSLABySeverityLevel)
	reportsRoutes.GET("/reports/sla-by-compliance-status", complianceReportController.GetSLAByComplianceStatus)

	reportsRoutes.GET("/reports/users", organizationReportController.GetAllUsers)
	reportsRoutes.GET("/reports/users/export", organizationReportController.GetUsersExport)
	reportsRoutes.GET("/reports/users-by-status", organizationReportController.GetUsersByStatus)
	reportsRoutes.GET("/reports/users-by-role", organizationReportController.GetUsersByRole)
	reportsRoutes.GET("/reports/users-by-designation", organizationReportController.GetUsersByDesignation)

	reportsRoutes.GET("/reports/processes", organizationReportController.GetAllProcesses)
	reportsRoutes.GET("/reports/processes/export", organizationReportController.GetProcessesExport)
	reportsRoutes.GET("/reports/processes-by-responsible-unit", organizationReportController.GetProcessesByResponsibleUnit)
	reportsRoutes.GET("/reports/processes-by-risk-rating", organizationReportController.GetProcessesByRiskRating)
	reportsRoutes.GET("/reports/processes-by-responsible-user", organizationReportController.GetProcessesByResponsibleUser)
	reportsRoutes.GET("/reports/processes-by-process-group", organizationReportController.GetProcessesByProcessGroup)

	reportsRoutes.GET("/reports/controls", organizationReportController.GetAllControls)
	reportsRoutes.GET("/reports/controls/export", organizationReportController.GetControlsExport)
	reportsRoutes.GET("/reports/controls-by-control-type", organizationReportController.GetControlsByControlType)
	reportsRoutes.GET("/reports/controls-by-control-category", organizationReportController.GetControlsByCategory)
	reportsRoutes.GET("/reports/controls-by-compliance-register", organizationReportController.GetControlsByComplianceRegister)
	reportsRoutes.GET("/reports/controls-by-risk-rating", organizationReportController.GetControlsByRiskRating)
	reportsRoutes.GET("/reports/controls-by-implementation-status", organizationReportController.GetControlsByImplementationStatus)
	reportsRoutes.GET("/reports/controls-by-efficiency-status", organizationReportController.GetControlsByEfficiencyStatus)
	reportsRoutes.GET("/reports/controls-by-applicability", organizationReportController.GetControlsByApplicability)
}
