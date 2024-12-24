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
	businessContinuityRoutes.GET("/business-impact-assessments", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllBusinessImpactAssessments)
	businessContinuityRoutes.GET("/business-impact-assessment-recovery-objectives", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllBusinessImpactAssessmentRecoveryObjectives)
	businessContinuityRoutes.POST("/business-impact-assessments", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.CreateBusinessImpactAssessment)
	businessContinuityRoutes.PUT("/business-impact-assessments", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.UpdateImpactAssessment)
	businessContinuityRoutes.PUT("/business-impact-assessment-recovery-objectives", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.UpdateRecoveryObjective)

	businessContinuityRoutes.GET("business-continuity/processes", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcesses)
	businessContinuityRoutes.GET("business-continuity/processes/:processId/incidents", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessIncidents)
	businessContinuityRoutes.GET("business-continuity/processes/:processId/assets", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessAssets)
	businessContinuityRoutes.GET("business-continuity/processes/:processId/documents", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessDocuments)
}
