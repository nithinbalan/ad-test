package routes

import (
	"github.com/gin-gonic/gin"
	"jethur/initializers"
	"jethur/middlewares"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	penetrationTestingRoutes := router.Group("/")
	penetrationTestingRoutes.Use(middlewares.AuthMiddleware(app.DB))
	unauthenticatedRoutes := router.Group("/")
	db := app.DB
	m := app.PenetrationTestingModule
	penetrationTest := m.PenetrationTestController
	penetrationTestingRoutes.GET("/penetration-testing/pentest-logs", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), penetrationTest.GetAllPenetrationTests)
	penetrationTestingRoutes.GET("/penetration-testing/pentest-logs/:id", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), penetrationTest.GetPenetrationTestDetails)
	penetrationTestingRoutes.GET("/penetration-testing/pentest-logs/:id/severity-breakdown", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), penetrationTest.GetPenetrationTestSeverityBreakdown)
	penetrationTestingRoutes.POST("/penetration-testing/pentest-logs", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), penetrationTest.CreatePenetrationTest)
	penetrationTestingRoutes.PUT("/penetration-testing/pentest-logs/:id", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), penetrationTest.UpdatePenetrationTest)
	environmentPreparation := m.PenetrationTestEnvironmentPreparationController
	penetrationTestingRoutes.GET("/penetration-testing/pentest-logs/:id/environment-preparation", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), environmentPreparation.GetPenetrationTestEnvironmentPreparationDetails)
	penetrationTestingRoutes.POST("/penetration-testing/pentest-logs/:id/environment-preparation", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), environmentPreparation.CreatePenetrationTestEnvironmentPreparation)
	penetrationTestingRoutes.PUT("/penetration-testing/pentest-logs/:id/environment-preparation", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), environmentPreparation.UpdatePenetrationTestEnvironmentPreparation)
	penetrationTestingRoutes.DELETE("/penetration-testing/pentest-logs/:id/environment-preparation", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), environmentPreparation.DeletePenetrationTestEnvironmentPreparations)
	execution := m.PenetrationTestExecutionController
	penetrationTestingRoutes.GET("/penetration-testing/pentest-logs/:id/execution", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), execution.GetPenetrationTestExecutionDetails)
	penetrationTestingRoutes.POST("/penetration-testing/pentest-logs/:id/execution", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), execution.CreatePenetrationTestExecution)
	penetrationTestingRoutes.PUT("/penetration-testing/pentest-logs/:id/execution", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), execution.UpdatePenetrationTestExecution)
	penetrationTestingRoutes.DELETE("/penetration-testing/pentest-logs/:id/execution", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), execution.DeletePenetrationTestExecutions)
	managementSharing := m.PenetrationTestManagementSharingController
	penetrationTestingRoutes.GET("/penetration-testing/pentest-logs/:id/management-sharing", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), managementSharing.GetPenetrationTestManagementSharingDetails)
	penetrationTestingRoutes.POST("/penetration-testing/pentest-logs/:id/management-sharing", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), managementSharing.CreatePenetrationTestManagementSharing)
	penetrationTestingRoutes.PUT("/penetration-testing/pentest-logs/:id/management-sharing", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), managementSharing.UpdatePenetrationTestManagementSharing)
	penetrationTestingRoutes.DELETE("/penetration-testing/pentest-logs/:id/management-sharing", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), managementSharing.DeletePenetrationTestManagementSharings)
	penetrationTestingRoutes.POST("/penetration-testing/pentest-logs/delete", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), penetrationTest.DeletePenetrationTests)
	finding := m.PenetrationTestFindingController
	penetrationTestingRoutes.GET("/penetration-test-findings", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), finding.GetAllPenetrationTestFindings)
	penetrationTestingRoutes.GET("/penetration-test-findings/:id", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), finding.GetPenetrationTestFindingDetails)
	penetrationTestingRoutes.POST("/penetration-test-findings", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), finding.CreatePenetrationTestFinding)
	penetrationTestingRoutes.PUT("/penetration-test-findings/:id", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), finding.UpdatePenetrationTestFinding)
	penetrationTestingRoutes.POST("/penetration-test-findings/delete", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), finding.DeletePenetrationTestFindings)
	findingStatusUpdate := m.PenetrationTestFindingStatusUpdateController
	penetrationTestingRoutes.GET("/penetration-test-findings/:id/status-updates", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), findingStatusUpdate.GetAllPenetrationTestFindingStatusUpdates)
	penetrationTestingRoutes.GET("/penetration-test-findings/:id/status-updates/latest", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), findingStatusUpdate.GetLatestPenetrationTestFindingStatusUpdate)
	penetrationTestingRoutes.POST("/penetration-test-findings/:id/status-updates", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), findingStatusUpdate.CreatePenetrationTestFindingStatusUpdate)
	penetrationTestingRoutes.PUT("/penetration-test-findings/:id/status-updates/:statusUpdateId", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), findingStatusUpdate.UpdatePenetrationTestFindingStatusUpdate)
	unauthenticatedRoutes.GET("/penetration-test-findings/:id/status-updates/:statusUpdateId/document-preview", findingStatusUpdate.GetPenetrationTestFindingStatusUpdateDocumentPreview)
	unauthenticatedRoutes.GET("/penetration-test-findings/:id/status-updates/:statusUpdateId/document-download", findingStatusUpdate.DownloadPenetrationTestFindingStatusUpdateDocument)
	retest := m.PenetrationTestRetestController
	penetrationTestingRoutes.GET("/penetration-test-retests", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), retest.GetAllPenetrationTestRetests)
	penetrationTestingRoutes.GET("/penetration-test-retests/:id", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), retest.GetPenetrationTestRetestDetails)
	penetrationTestingRoutes.POST("/penetration-test-retests", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), retest.CreatePenetrationTestRetest)
	penetrationTestingRoutes.PUT("/penetration-test-retests/:id", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), retest.UpdatePenetrationTestRetest)
	penetrationTestingRoutes.GET("/penetration-test-retests/:id/eligible-findings", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), retest.GetEligiblePenetrationTestRetestFindings)
	penetrationTestingRoutes.PUT("/penetration-test-retests/:id/findings", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), retest.UpdatePenetrationTestRetestFindings)
	penetrationTestingRoutes.DELETE("/penetration-test-retests/:id/findings/:retestFindingId", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), retest.DeletePenetrationTestRetestFinding)
	penetrationTestingRoutes.PUT("/penetration-test-retests/:id/findings/:retestFindingId/save", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), retest.SavePenetrationTestRetestFinding)
	penetrationTestingRoutes.POST("/penetration-test-retests/:id/findings/:retestFindingId/result", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), retest.DecidePenetrationTestRetestFinding)
	unauthenticatedRoutes.GET("/penetration-test-retests/:id/findings/:retestFindingId/document-preview", retest.GetPenetrationTestRetestFindingDocumentPreview)
	unauthenticatedRoutes.GET("/penetration-test-retests/:id/findings/:retestFindingId/document-download", retest.DownloadPenetrationTestRetestFindingDocument)
	penetrationTestingRoutes.POST("/penetration-test-retests/:id/complete", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), retest.CompletePenetrationTestRetest)
	penetrationTestingRoutes.POST("/penetration-test-retests/delete", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), retest.DeletePenetrationTestRetests)
	supportingDocument := m.PenetrationTestSupportingDocumentController
	unauthenticatedRoutes.GET("/penetration-testing/pentest-logs/:id/supporting-documents/:supportingDocumentId/documents/:documentId/download", supportingDocument.DownloadDocument)
	unauthenticatedRoutes.GET("/penetration-testing/pentest-logs/:id/supporting-documents/:supportingDocumentId/documents/:documentId/preview", supportingDocument.GetDocumentPreview)
	penetrationTestingRoutes.GET("/penetration-testing/pentest-logs/:id/supporting-documents", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), supportingDocument.GetAllSupportingDocuments)
	penetrationTestingRoutes.GET("/penetration-testing/pentest-logs/:id/supporting-documents/:supportingDocumentId", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), supportingDocument.GetSupportingDocumentDetails)
	penetrationTestingRoutes.POST("/penetration-testing/pentest-logs/:id/supporting-documents", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), supportingDocument.CreateSupportingDocument)
	penetrationTestingRoutes.PUT("/penetration-testing/pentest-logs/:id/supporting-documents/:supportingDocumentId", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), supportingDocument.UpdateSupportingDocument)
	penetrationTestingRoutes.POST("/penetration-testing/pentest-logs/:id/supporting-documents/delete", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), supportingDocument.DeleteSupportingDocuments)
	closure := m.PenetrationTestClosureController
	penetrationTestingRoutes.GET("/penetration-testing/pentest-logs/:id/closure", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "view"), closure.GetPenetrationTestClosureDetails)
	penetrationTestingRoutes.POST("/penetration-testing/pentest-logs/:id/closure", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), closure.CreatePenetrationTestClosure)
	penetrationTestingRoutes.PUT("/penetration-testing/pentest-logs/:id/closure", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), closure.UpdatePenetrationTestClosure)
	penetrationTestingRoutes.DELETE("/penetration-testing/pentest-logs/:id/closure", middlewares.RolePermissionMiddleware(db, "penetration_testing_pentest_log", "modify"), closure.DeletePenetrationTestClosures)
	penetrationTestingRoutes.GET("/penetration-testing/types", m.PenetrationTestTypeController.GetAllPenetrationTestTypes)
	penetrationTestingRoutes.GET("/penetration-testing/types/:id", m.PenetrationTestTypeController.GetPenetrationTestTypeDetails)
	penetrationTestingRoutes.POST("/penetration-testing/types", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), m.PenetrationTestTypeController.CreatePenetrationTestType)
	penetrationTestingRoutes.PUT("/penetration-testing/types/:id", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), m.PenetrationTestTypeController.UpdatePenetrationTestType)
	penetrationTestingRoutes.POST("/penetration-testing/types/delete", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), m.PenetrationTestTypeController.DeletePenetrationTestTypes)
	penetrationTestingRoutes.PUT("/penetration-testing/types/activate", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), m.PenetrationTestTypeController.ActivatePenetrationTestTypes)
	penetrationTestingRoutes.PUT("/penetration-testing/types/deactivate", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), m.PenetrationTestTypeController.DeactivatePenetrationTestTypes)
	plan := m.PenetrationTestPlanController
	penetrationTestingRoutes.GET("/penetration-testing/plans", plan.GetAllPenetrationTestPlans)
	penetrationTestingRoutes.GET("/penetration-testing/plans/:id", plan.GetPenetrationTestPlanDetails)
	penetrationTestingRoutes.POST("/penetration-testing/plans", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), plan.CreatePenetrationTestPlan)
	penetrationTestingRoutes.PUT("/penetration-testing/plans/:id", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), plan.UpdatePenetrationTestPlan)
	penetrationTestingRoutes.POST("/penetration-testing/plans/delete", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), plan.DeletePenetrationTestPlans)
	penetrationTestingRoutes.PUT("/penetration-testing/plans/activate", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), plan.ActivatePenetrationTestPlans)
	penetrationTestingRoutes.PUT("/penetration-testing/plans/deactivate", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), plan.DeactivatePenetrationTestPlans)
	executor := m.PenetrationTestExecutorController
	penetrationTestingRoutes.GET("/penetration-testing/executors", executor.GetAllPenetrationTestExecutors)
	penetrationTestingRoutes.GET("/penetration-testing/executors/:id", executor.GetPenetrationTestExecutorDetails)
	penetrationTestingRoutes.POST("/penetration-testing/executors", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), executor.CreatePenetrationTestExecutor)
	penetrationTestingRoutes.PUT("/penetration-testing/executors/:id", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), executor.UpdatePenetrationTestExecutor)
	penetrationTestingRoutes.POST("/penetration-testing/executors/delete", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), executor.DeletePenetrationTestExecutors)
	penetrationTestingRoutes.PUT("/penetration-testing/executors/activate", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), executor.ActivatePenetrationTestExecutors)
	penetrationTestingRoutes.PUT("/penetration-testing/executors/deactivate", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), executor.DeactivatePenetrationTestExecutors)
	readiness := m.PenetrationTestEnvironmentReadinessController
	classification := m.PenetrationTestClassificationController
	penetrationTestingRoutes.GET("/penetration-testing/classifications", classification.GetAllPenetrationTestClassifications)
	executionStatus := m.PenetrationTestExecutionStatusController
	penetrationTestingRoutes.GET("/penetration-testing/execution-statuses", executionStatus.GetAllPenetrationTestExecutionStatuses)
	securityResult := m.PenetrationTestSecurityResultController
	penetrationTestingRoutes.GET("/penetration-testing/security-results", securityResult.GetAllPenetrationTestSecurityResults)
	penetrationTestingRoutes.GET("/penetration-testing/environment-readinesses", readiness.GetAllPenetrationTestEnvironmentReadinesses)
	penetrationTestingRoutes.GET("/penetration-testing/environment-readinesses/:id", readiness.GetPenetrationTestEnvironmentReadinessDetails)
	penetrationTestingRoutes.POST("/penetration-testing/environment-readinesses", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), readiness.CreatePenetrationTestEnvironmentReadiness)
	penetrationTestingRoutes.PUT("/penetration-testing/environment-readinesses/:id", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), readiness.UpdatePenetrationTestEnvironmentReadiness)
	penetrationTestingRoutes.POST("/penetration-testing/environment-readinesses/delete", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), readiness.DeletePenetrationTestEnvironmentReadinesses)
	penetrationTestingRoutes.PUT("/penetration-testing/environment-readinesses/activate", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), readiness.ActivatePenetrationTestEnvironmentReadinesses)
	penetrationTestingRoutes.PUT("/penetration-testing/environment-readinesses/deactivate", middlewares.RolePermissionMiddleware(db, "penetration_testing_setting", "modify"), readiness.DeactivatePenetrationTestEnvironmentReadinesses)
}
