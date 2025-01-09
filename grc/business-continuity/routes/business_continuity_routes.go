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
	businessContinuityPlanController := app.BusinessContinuityModule.BusinessContinuityPlanController
	bcpBiaMatrixController := app.BusinessContinuityModule.BcpBiaMatrixController
	bcpPhaseController := app.BusinessContinuityModule.BcpPhaseController

	// bcp phase
	businessContinuityRoutes.GET("/bcp-phases", bcpPhaseController.GetAllBcpPhases)
	businessContinuityRoutes.GET("/bcp-phases/:id", bcpPhaseController.GetBcpPhaseDetails)

	//Business Function
	businessContinuityRoutes.GET("/business-functions", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetAllBusinessFunctions)
	businessContinuityRoutes.GET("/business-functions/:id", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetBusinessFunctionDetails)
	businessContinuityRoutes.POST("/business-functions", middlewares.RolePermissionMiddleware(db, "business_functions", "modify"), businessFunctionController.CreateBusinessFunction)
	businessContinuityRoutes.PUT("/business-functions/:id", middlewares.RolePermissionMiddleware(db, "business_functions", "modify"), businessFunctionController.UpdateBusinessFunction)
	businessContinuityRoutes.POST("/business-functions/delete", middlewares.RolePermissionMiddleware(db, "business_functions", "modify"), businessFunctionController.DeleteBusinessFunctions)

	// Business Impact Assessment
	businessContinuityRoutes.GET("bia/processes", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcesses)
	businessContinuityRoutes.GET("bia/processes/:processId", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetProcessDetails)
	businessContinuityRoutes.GET("bia/processes/:processId/incidents", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessIncidents)
	businessContinuityRoutes.GET("bia/processes/:processId/assets", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessAssets)
	businessContinuityRoutes.GET("bia/processes/:processId/documents", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessDocuments)
	businessContinuityRoutes.GET("bia/processes/:processId/risks", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessRisks)

	businessContinuityRoutes.POST("/business-impact-assessments", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.CreateBusinessImpactAssessment)
	businessContinuityRoutes.GET("/bia-impact-analysis", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllBIAImpactAnalysis)
	businessContinuityRoutes.GET("/bia-recovery-objectives", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllBIARecoveryObjectives)
	businessContinuityRoutes.PUT("/bia-impact-analysis", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.UpdateImpactAnalysis)
	businessContinuityRoutes.PUT("/bia-recovery-objectives", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.UpdateRecoveryObjective)

	//Business Continuity Plan
	businessContinuityRoutes.GET("/business-continuity-plans", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetAllBusinessContinuityPlans)
	businessContinuityRoutes.GET("/business-continuity-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetBusinessContinuityPlanDetails)
	businessContinuityRoutes.POST("/business-continuity-plans", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.CreateBusinessContinuityPlan)
	businessContinuityRoutes.PUT("/business-continuity-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.UpdateBusinessContinuityPlan)
	businessContinuityRoutes.POST("/business-continuity-plans/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.DeleteBusinessContinuityPlans)

	// bcp bia matrix
	businessContinuityRoutes.GET("/bcp/:id/bia-matrics", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpBiaMatrixController.GetBcpBiaMatrixDetails)
	businessContinuityRoutes.PUT("/bcp/:id/bia-matrics", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpBiaMatrixController.UpdateBcpBiaMatrix)
}
