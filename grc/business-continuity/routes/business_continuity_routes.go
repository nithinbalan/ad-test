package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	businessContinuityRoutes := router.Group("/")

	businessContinuityRoutes.Use(middlewares.AuthMiddleware())

	db := app.DB

	businessFunctionController := app.BusinessContinuityModule.BusinessFunctionController
	businessImpactAssessmentController := app.BusinessContinuityModule.BusinessImpactAssessmentController

	//Business Function
	businessContinuityRoutes.GET("/business-functions", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetAllBusinessFunctions)
	businessContinuityRoutes.GET("/business-functions/:id", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetBusinessFunctionDetails)
	businessContinuityRoutes.POST("/business-functions", middlewares.RolePermissionMiddleware(db, "business_functions", "modify"), businessFunctionController.CreateBusinessFunction)
	businessContinuityRoutes.PUT("/business-functions/:id", middlewares.RolePermissionMiddleware(db, "business_functions", "modify"), businessFunctionController.UpdateBusinessFunction)
	businessContinuityRoutes.POST("/business-functions/delete", middlewares.RolePermissionMiddleware(db, "business_functions", "modify"), businessFunctionController.DeleteBusinessFunctions)

	// Business Impact Assessment
	businessContinuityRoutes.GET("bia/processes", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcesses)
	businessContinuityRoutes.GET("bia/processes/:processId/incidents", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessIncidents)
	businessContinuityRoutes.GET("bia/processes/:processId/assets", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessAssets)
	businessContinuityRoutes.GET("bia/processes/:processId/documents", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessDocuments)

	businessContinuityRoutes.POST("/business-impact-assessments", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.CreateBusinessImpactAssessment)
	businessContinuityRoutes.GET("/bia-impact-analysis", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllBIAImpactAnalysis)
	businessContinuityRoutes.GET("/bia-recovery-objectives", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllBIARecoveryObjectives)
	businessContinuityRoutes.PUT("/bia-impact-analysis", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.UpdateImpactAnalysis)
	businessContinuityRoutes.PUT("/bia-recovery-objectives", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.UpdateRecoveryObjective)

}
