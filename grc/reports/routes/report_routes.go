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

	// Compliance Report Summary
	reportsRoutes.GET("/reports/compliance-by-category", complianceReportController.GetComplianceByCategory)
	reportsRoutes.GET("/reports/compliance-registers", complianceReportController.GetAllComplianceRegisters)
	reportsRoutes.GET("/reports/compliance-registers/export", complianceReportController.ComplianceRegisterReportExport)

}
