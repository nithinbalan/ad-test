package routes

import (
	"github.com/gin-gonic/gin"
	"jethur/initializers"
	"jethur/middlewares"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	penetrationTestingRoutes := router.Group("/")
	penetrationTestingRoutes.Use(middlewares.AuthMiddleware(app.DB))
	db := app.DB
	m := app.PenetrationTestingModule
	permission := func() gin.HandlerFunc {
		return middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify")
	}
	penetrationTestingRoutes.GET("/penetration-testing/types", m.PenetrationTestTypeController.GetAllPenetrationTestTypes)
	penetrationTestingRoutes.GET("/penetration-testing/types/:id", m.PenetrationTestTypeController.GetPenetrationTestTypeDetails)
	penetrationTestingRoutes.POST("/penetration-testing/types", permission(), m.PenetrationTestTypeController.CreatePenetrationTestType)
	penetrationTestingRoutes.PUT("/penetration-testing/types/:id", permission(), m.PenetrationTestTypeController.UpdatePenetrationTestType)
	penetrationTestingRoutes.POST("/penetration-testing/types/delete", permission(), m.PenetrationTestTypeController.DeletePenetrationTestTypes)
	penetrationTestingRoutes.PUT("/penetration-testing/types/activate", permission(), m.PenetrationTestTypeController.ActivatePenetrationTestTypes)
	penetrationTestingRoutes.PUT("/penetration-testing/types/deactivate", permission(), m.PenetrationTestTypeController.DeactivatePenetrationTestTypes)
	plan := m.PenetrationTestPlanController
	penetrationTestingRoutes.GET("/penetration-testing/plans", plan.GetAllPenetrationTestPlans)
	penetrationTestingRoutes.GET("/penetration-testing/plans/:id", plan.GetPenetrationTestPlanDetails)
	penetrationTestingRoutes.POST("/penetration-testing/plans", permission(), plan.CreatePenetrationTestPlan)
	penetrationTestingRoutes.PUT("/penetration-testing/plans/:id", permission(), plan.UpdatePenetrationTestPlan)
	penetrationTestingRoutes.POST("/penetration-testing/plans/delete", permission(), plan.DeletePenetrationTestPlans)
	penetrationTestingRoutes.PUT("/penetration-testing/plans/activate", permission(), plan.ActivatePenetrationTestPlans)
	penetrationTestingRoutes.PUT("/penetration-testing/plans/deactivate", permission(), plan.DeactivatePenetrationTestPlans)
	executor := m.PenetrationTestExecutorController
	penetrationTestingRoutes.GET("/penetration-testing/executors", executor.GetAllPenetrationTestExecutors)
	penetrationTestingRoutes.GET("/penetration-testing/executors/:id", executor.GetPenetrationTestExecutorDetails)
	penetrationTestingRoutes.POST("/penetration-testing/executors", permission(), executor.CreatePenetrationTestExecutor)
	penetrationTestingRoutes.PUT("/penetration-testing/executors/:id", permission(), executor.UpdatePenetrationTestExecutor)
	penetrationTestingRoutes.POST("/penetration-testing/executors/delete", permission(), executor.DeletePenetrationTestExecutors)
	penetrationTestingRoutes.PUT("/penetration-testing/executors/activate", permission(), executor.ActivatePenetrationTestExecutors)
	penetrationTestingRoutes.PUT("/penetration-testing/executors/deactivate", permission(), executor.DeactivatePenetrationTestExecutors)
	readiness := m.PenetrationTestEnvironmentReadinessController
	classification := m.PenetrationTestClassificationController
	penetrationTestingRoutes.GET("/penetration-testing/classifications", classification.GetAllPenetrationTestClassifications)
	executionStatus := m.PenetrationTestExecutionStatusController
	penetrationTestingRoutes.GET("/penetration-testing/execution-statuses", executionStatus.GetAllPenetrationTestExecutionStatuses)
	securityResult := m.PenetrationTestSecurityResultController
	penetrationTestingRoutes.GET("/penetration-testing/security-results", securityResult.GetAllPenetrationTestSecurityResults)
	penetrationTestingRoutes.GET("/penetration-testing/environment-readinesses", readiness.GetAllPenetrationTestEnvironmentReadinesses)
	penetrationTestingRoutes.GET("/penetration-testing/environment-readinesses/:id", readiness.GetPenetrationTestEnvironmentReadinessDetails)
	penetrationTestingRoutes.POST("/penetration-testing/environment-readinesses", permission(), readiness.CreatePenetrationTestEnvironmentReadiness)
	penetrationTestingRoutes.PUT("/penetration-testing/environment-readinesses/:id", permission(), readiness.UpdatePenetrationTestEnvironmentReadiness)
	penetrationTestingRoutes.POST("/penetration-testing/environment-readinesses/delete", permission(), readiness.DeletePenetrationTestEnvironmentReadinesses)
	penetrationTestingRoutes.PUT("/penetration-testing/environment-readinesses/activate", permission(), readiness.ActivatePenetrationTestEnvironmentReadinesses)
	penetrationTestingRoutes.PUT("/penetration-testing/environment-readinesses/deactivate", permission(), readiness.DeactivatePenetrationTestEnvironmentReadinesses)
}
