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
	reportsRoutes.GET("/reports/compliance-by-category", complianceReportController.GetComplianceByCategory)

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

	reportsRoutes.GET("/reports/controls-by-control-type", organizationReportController.GetControlsByControlType)
	reportsRoutes.GET("/reports/controls-by-control-category", organizationReportController.GetControlsByCategory)
	reportsRoutes.GET("/reports/controls-by-compliance-register", organizationReportController.GetControlsByComplianceRegister)
	reportsRoutes.GET("/reports/controls-by-risk-rating", organizationReportController.GetControlsByRiskRating)
	reportsRoutes.GET("/reports/controls-by-implementation-status", organizationReportController.GetControlsByImplementationStatus)
	reportsRoutes.GET("/reports/controls-by-efficiency-status", organizationReportController.GetControlsByEfficiencyStatus)
	reportsRoutes.GET("/reports/controls-by-applicability", organizationReportController.GetControlsByApplicability)
}
