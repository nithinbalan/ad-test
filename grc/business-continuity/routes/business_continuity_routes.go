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
	bcpRecoveryPlanController := app.BusinessContinuityModule.BcpRecoveryPlanController
	bcpCommunicationPlanController := app.BusinessContinuityModule.BcpCommunicationPlanController
	bcpRolesController := app.BusinessContinuityModule.BcpRolesController

	// bcp phase
	businessContinuityRoutes.GET("/bcp-phases", bcpPhaseController.GetAllBcpPhases)
	businessContinuityRoutes.GET("/bcp-phases/:id", bcpPhaseController.GetBcpPhaseDetails)

	// Business Function
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

	// Business Continuity Plan
	businessContinuityRoutes.GET("/business-continuity-plans", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetAllBusinessContinuityPlans)
	businessContinuityRoutes.GET("/business-continuity-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetBusinessContinuityPlanDetails)
	businessContinuityRoutes.POST("/business-continuity-plans", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.CreateBusinessContinuityPlan)
	businessContinuityRoutes.PUT("/business-continuity-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.UpdateBusinessContinuityPlan)
	businessContinuityRoutes.POST("/business-continuity-plans/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.DeleteBusinessContinuityPlans)

	// bcp bia matrix
	businessContinuityRoutes.GET("/bcp/:bcpId/key-processes", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpBiaMatrixController.GetAllBcpKeyProcesses)
	businessContinuityRoutes.GET("/bcp/:bcpId/bia-matrics", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpBiaMatrixController.GetBcpBiaMatrixDetails)
	businessContinuityRoutes.PUT("/bcp/:bcpId/bia-matrics", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpBiaMatrixController.UpdateBcpBiaMatrix)

	// Bcp Recovery Plan
	businessContinuityRoutes.GET("/bcp/:bcpId/bcp-recovery-plans", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpRecoveryPlanController.GetAllBcpRecoveryPlans)
	businessContinuityRoutes.GET("/bcp/:bcpId/bcp-recovery-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpRecoveryPlanController.GetBcpRecoveryPlanDetails)
	businessContinuityRoutes.POST("/bcp/:bcpId/bcp-recovery-plans", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpRecoveryPlanController.CreateBcpRecoveryPlan)
	businessContinuityRoutes.PUT("/bcp/:bcpId/bcp-recovery-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpRecoveryPlanController.UpdateBcpRecoveryPlan)
	businessContinuityRoutes.POST("/bcp/:bcpId/bcp-recovery-plans/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpRecoveryPlanController.DeleteBcpRecoveryPlans)

	// Bcp Recovery Plan
	businessContinuityRoutes.GET("/bcp/:bcpId/bcp-communication-plans", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpCommunicationPlanController.GetAllBcpCommunicationPlans)
	businessContinuityRoutes.GET("/bcp/:bcpId/bcp-communication-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpCommunicationPlanController.GetBcpCommunicationPlanDetails)
	businessContinuityRoutes.POST("/bcp/:bcpId/bcp-communication-plans", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpCommunicationPlanController.CreateBcpCommunicationPlan)
	businessContinuityRoutes.PUT("/bcp/:bcpId/bcp-communication-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpCommunicationPlanController.UpdateBcpCommunicationPlan)
	businessContinuityRoutes.POST("/bcp/:bcpId/bcp-communication-plans/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpCommunicationPlanController.DeleteBcpCommunicationPlans)
	businessContinuityRoutes.GET("/bia-matrix", middlewares.RolePermissionMiddleware(db, "bia_matrics", "view"), businessImpactAssessmentController.GetBIAMatrix)

	// teams and roles
	businessContinuityRoutes.GET("/bcp/:bcpId/teams", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpRolesController.GetAllBcpRoles)
	businessContinuityRoutes.POST("/bcp/:bcpId/teams", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpRolesController.CreateBcpRole)
	businessContinuityRoutes.PUT("/bcp/:bcpId/teams/:id", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpRolesController.UpdateBcpRole)
	businessContinuityRoutes.POST("/bcp/teams/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpRolesController.DeleteBcpRoles)
}
