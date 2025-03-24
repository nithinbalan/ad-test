package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	businessContinuityRoutes := router.Group("/")

	businessContinuityRoutes.Use(middlewares.AuthMiddleware())

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	businessFunctionController := app.BusinessContinuityModule.BusinessFunctionController
	businessImpactAssessmentController := app.BusinessContinuityModule.BusinessImpactAssessmentController
	businessContinuityPlanController := app.BusinessContinuityModule.BusinessContinuityPlanController
	bcpBiaMatrixController := app.BusinessContinuityModule.BcpBiaMatrixController
	bcpPhaseController := app.BusinessContinuityModule.BcpPhaseController
	bcpRecoveryPlanController := app.BusinessContinuityModule.BcpRecoveryPlanController
	bcpCommunicationPlanController := app.BusinessContinuityModule.BcpCommunicationPlanController
	bcpTestAndExerciseController := app.BusinessContinuityModule.BcpTestAndExerciseController
	bcpStatusController := app.BusinessContinuityModule.BcpStatusController
	BcpTeamController := app.BusinessContinuityModule.BcpTeamController
	testAndExerciseStatusController := app.BusinessContinuityModule.TestAndExerciseStatusController
	testAndExerciseController := app.BusinessContinuityModule.TestAndExerciseController
	communicationActionPlanStatusController := app.BusinessContinuityModule.CommunicationActionPlanStatusController
	testAndExerciseCommunicationActionPlanController := app.BusinessContinuityModule.TestAndExerciseCommunicationActionPlanController
	testAndExerciseRecoveryActionPlanController := app.BusinessContinuityModule.TestAndExerciseRecoveryActionPlanController
	taeRecoveryActionPlanStatusUpdateController := app.BusinessContinuityModule.TaeRecoveryActionPlanStatusUpdateController
	timelinessController := app.BusinessContinuityModule.TimelinessController
	feedbackOptionController := app.BusinessContinuityModule.FeedbackOptionController
	testResultOptionController := app.BusinessContinuityModule.TestResultOptionController
	testAndExerciseTestResultController := app.BusinessContinuityModule.TestAndExerciseTestResultController
	taeRecoveryActionPlanCommentController := app.BusinessContinuityModule.TaeRecoveryActionPlanCommentController
	bcpRiskAnalysisController := app.BusinessContinuityModule.BcpRiskAnalysisController

	// bcp phase
	businessContinuityRoutes.GET("/bcp-phases", bcpPhaseController.GetAllBcpPhases)
	businessContinuityRoutes.GET("/bcp-phases/:id", bcpPhaseController.GetBcpPhaseDetails)

	// bcp status
	businessContinuityRoutes.GET("/bcp-statuses", bcpStatusController.GetAllBcpStatuses)
	businessContinuityRoutes.GET("/bcp-statuses/:id", bcpStatusController.GetBcpStatusDetails)

	// test and exercise status
	businessContinuityRoutes.GET("/test-and-exercise-statuses", testAndExerciseStatusController.GetAllTestAndExerciseStatuses)
	businessContinuityRoutes.GET("/test-and-exercise-statuses/:id", testAndExerciseStatusController.GetTestAndExerciseStatusDetails)

	// communication action plan status
	businessContinuityRoutes.GET("/communication-action-plan-statuses", communicationActionPlanStatusController.GetAllCommunicationActionPlanStatuses)
	businessContinuityRoutes.GET("/communication-action-plan-statuses/:id", communicationActionPlanStatusController.GetCommunicationActionPlanStatusDetails)

	// timeliness
	businessContinuityRoutes.GET("/timeliness", timelinessController.GetAllTimeliness)
	businessContinuityRoutes.GET("/timeliness/:id", timelinessController.GetTimelinessDetails)

	// feedback option
	businessContinuityRoutes.GET("/feedback-options", feedbackOptionController.GetAllFeedbackOptions)
	businessContinuityRoutes.GET("/feedback-options/:id", feedbackOptionController.GetFeedbackOptionDetails)

	// test result option
	businessContinuityRoutes.GET("/test-result-options", testResultOptionController.GetAllTestResultOptions)
	businessContinuityRoutes.GET("/test-result-options/:id", testResultOptionController.GetTestResultOptionDetails)

	// Business Function
	businessContinuityRoutes.GET("/business-functions", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetAllBusinessFunctions)
	businessContinuityRoutes.GET("/business-functions/:id", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetBusinessFunctionDetails)
	businessContinuityRoutes.POST("/business-functions", middlewares.RolePermissionMiddleware(db, "business_functions", "modify"), businessFunctionController.CreateBusinessFunction)
	businessContinuityRoutes.PUT("/business-functions/:id", middlewares.RolePermissionMiddleware(db, "business_functions", "modify"), businessFunctionController.UpdateBusinessFunction)
	businessContinuityRoutes.POST("/business-functions/delete", middlewares.RolePermissionMiddleware(db, "business_functions", "modify"), businessFunctionController.DeleteBusinessFunctions)
	businessContinuityRoutes.GET("/business-functions/export", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetAllBusinessFunctionExport)

	businessContinuityRoutes.GET("/business-functions/:id/processes/:processId/associated-risks", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetProcessAssociatedRisks)
	businessContinuityRoutes.GET("/business-functions/:id/processes/:processId/associated-incidents", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetProcessAssociatedIncidents)
	businessContinuityRoutes.GET("/business-functions/:id/processes/:processId/associated-controls", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetProcessAssociatedControls)
	// Business Impact Assessment
	businessContinuityRoutes.GET("bia/processes", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcesses)
	businessContinuityRoutes.GET("bia/processes/export", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessesExport)
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
	businessContinuityRoutes.POST("/business-impact-assessments/delete", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.DeleteBusinessImpactAssessments)
	businessContinuityRoutes.GET("/bia-matrix", middlewares.RolePermissionMiddleware(db, "bia_matrics", "view"), businessImpactAssessmentController.GetBIAMatrix)
	businessContinuityRoutes.GET("/bia-matrix/export", middlewares.RolePermissionMiddleware(db, "bia_matrics", "view"), businessImpactAssessmentController.GetBiaMatrixExport)

	// Business Continuity Plan
	unauthenticatedRoutes.GET("/business-continuity-plans/:id/document-preview", businessContinuityPlanController.GetBusinessContinuityPlanDocumentPreview)
	unauthenticatedRoutes.GET("/business-continuity-plans/:id/document-download", businessContinuityPlanController.DownloadBusinessContinuityPlanDocument)
	businessContinuityRoutes.GET("/business-continuity-plans", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetAllBusinessContinuityPlans)
	businessContinuityRoutes.GET("/business-continuity-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetBusinessContinuityPlanDetails)
	businessContinuityRoutes.POST("/business-continuity-plans", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.CreateBusinessContinuityPlan)
	businessContinuityRoutes.PUT("/business-continuity-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.UpdateBusinessContinuityPlan)
	businessContinuityRoutes.POST("/business-continuity-plans/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.DeleteBusinessContinuityPlans)
	businessContinuityRoutes.GET("/business-continuity-plans/export", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetAllBusinessContinuityPlanExport)

	//bcp heatmap
	businessContinuityRoutes.GET("/business-continuity-plans/:id/heatmap", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetBcpHeatmap)

	// bcp risk analysis
	businessContinuityRoutes.POST("/bcp/:bcpId/risk-analysis/risks", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpRiskAnalysisController.AddRiskToBcp)
	businessContinuityRoutes.GET("/bcp/:bcpId/risk-analysis/risks", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpRiskAnalysisController.GetAllRisksByBcpId)
	businessContinuityRoutes.POST("/bcp/:bcpId/risk-analysis/risks/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpRiskAnalysisController.DeleteRiskFromBcp)
	businessContinuityRoutes.GET("/bcp/:bcpId/risk-analysis/processes-risks", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpRiskAnalysisController.GetAllProcessRisks)

	// bcp bia matrix
	businessContinuityRoutes.GET("/bcp/:bcpId/key-processes", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpBiaMatrixController.GetAllBcpKeyProcesses)
	businessContinuityRoutes.GET("/bcp/:bcpId/bia-matrics", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpBiaMatrixController.GetBcpBiaMatrixDetails)
	businessContinuityRoutes.PUT("/bcp/:bcpId/bia-matrics", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpBiaMatrixController.UpdateBcpBiaMatrix)

	// Bcp Recovery Plan
	businessContinuityRoutes.GET("/bcp/:bcpId/recovery-plans", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpRecoveryPlanController.GetAllBcpRecoveryPlans)
	businessContinuityRoutes.GET("/bcp/:bcpId/recovery-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpRecoveryPlanController.GetBcpRecoveryPlanDetails)
	businessContinuityRoutes.POST("/bcp/:bcpId/recovery-plans", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpRecoveryPlanController.CreateBcpRecoveryPlan)
	businessContinuityRoutes.PUT("/bcp/:bcpId/recovery-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpRecoveryPlanController.UpdateBcpRecoveryPlan)
	businessContinuityRoutes.POST("/bcp/:bcpId/recovery-plans/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpRecoveryPlanController.DeleteBcpRecoveryPlans)

	// Bcp Recovery Plan
	businessContinuityRoutes.GET("/bcp/:bcpId/communication-plans", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpCommunicationPlanController.GetAllBcpCommunicationPlans)
	businessContinuityRoutes.GET("/bcp/:bcpId/communication-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpCommunicationPlanController.GetBcpCommunicationPlanDetails)
	businessContinuityRoutes.POST("/bcp/:bcpId/communication-plans", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpCommunicationPlanController.CreateBcpCommunicationPlan)
	businessContinuityRoutes.PUT("/bcp/:bcpId/communication-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpCommunicationPlanController.UpdateBcpCommunicationPlan)
	businessContinuityRoutes.POST("/bcp/:bcpId/communication-plans/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpCommunicationPlanController.DeleteBcpCommunicationPlans)

	// Bcp Test and Exercise
	businessContinuityRoutes.GET("/bcp/:bcpId/test-and-exercises", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpTestAndExerciseController.GetAllBcpTestAndExercises)
	businessContinuityRoutes.GET("/bcp/:bcpId/test-and-exercises/:id", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpTestAndExerciseController.GetBcpTestAndExerciseDetails)
	businessContinuityRoutes.POST("/bcp/:bcpId/test-and-exercises", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpTestAndExerciseController.CreateBcpTestAndExercise)
	businessContinuityRoutes.PUT("/bcp/:bcpId/test-and-exercises/:id", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpTestAndExerciseController.UpdateBcpTestAndExercise)
	businessContinuityRoutes.POST("/bcp/:bcpId/test-and-exercises/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpTestAndExerciseController.DeleteBcpTestAndExercises)

	// teams and roles
	businessContinuityRoutes.GET("/bcp/:bcpId/teams", middlewares.RolePermissionMiddleware(db, "bcp", "view"), BcpTeamController.GetAllBcpTeams)
	businessContinuityRoutes.GET("/bcp/:bcpId/teams/:id", middlewares.RolePermissionMiddleware(db, "bcp", "view"), BcpTeamController.GetBcpTeamDetails)
	businessContinuityRoutes.POST("/bcp/:bcpId/teams", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), BcpTeamController.CreateBcpTeam)
	businessContinuityRoutes.PUT("/bcp/:bcpId/teams/:id", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), BcpTeamController.UpdateBcpTeam)
	businessContinuityRoutes.POST("/bcp/:bcpId/teams/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), BcpTeamController.DeleteBcpTeams)

	// Test and Exercise
	businessContinuityRoutes.GET("/test-and-exercises", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), testAndExerciseController.GetAllTestAndExercises)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), testAndExerciseController.GetTestAndExerciseDetails)
	businessContinuityRoutes.POST("/test-and-exercises", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseController.CreateTestAndExercise)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseController.UpdateTestAndExercise)
	businessContinuityRoutes.POST("/test-and-exercises/delete", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseController.DeleteTestAndExercises)
	businessContinuityRoutes.GET("/test-and-exercises/export", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), testAndExerciseController.GetAllTestAndExerciseExport)

	// Test and exercise communication action plan
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/communication-action-plans", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), testAndExerciseCommunicationActionPlanController.GetAllTestAndExerciseCommunicationActionPlans)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/communication-action-plans/:id", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), testAndExerciseCommunicationActionPlanController.GetTestAndExerciseCommunicationActionPlanDetails)
	businessContinuityRoutes.POST("/test-and-exercises/:testAndExerciseId/communication-action-plans", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseCommunicationActionPlanController.CreateTestAndExerciseCommunicationActionPlan)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/communication-action-plans/:id", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseCommunicationActionPlanController.UpdateTestAndExerciseCommunicationActionPlan)
	businessContinuityRoutes.POST("/test-and-exercises/:testAndExerciseId/communication-action-plans/delete", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseCommunicationActionPlanController.DeleteTestAndExerciseCommunicationActionPlans)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/communication-action-plans/:id/status-update", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseCommunicationActionPlanController.StatusUpdate)

	// test and exercise recovery action plan
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), testAndExerciseRecoveryActionPlanController.GetAllTestAndExerciseRecoveryActionPlans)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), testAndExerciseRecoveryActionPlanController.GetTestAndExerciseRecoveryActionPlanDetails)
	businessContinuityRoutes.POST("/test-and-exercises/:testAndExerciseId/recovery-action-plans", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseRecoveryActionPlanController.CreateTestAndExerciseRecoveryActionPlan)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseRecoveryActionPlanController.UpdateTestAndExerciseRecoveryActionPlan)
	businessContinuityRoutes.POST("/test-and-exercises/:testAndExerciseId/recovery-action-plans/delete", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseRecoveryActionPlanController.DeleteTestAndExerciseRecoveryActionPlans)

	// Test and exercise test result
	unauthenticatedRoutes.GET("/test-and-exercises/:testAndExerciseId/test-result-analysis/document-preview", testAndExerciseTestResultController.GetTestAndExerciseTestResultDocumentPreview)
	unauthenticatedRoutes.GET("/test-and-exercises/:testAndExerciseId/test-result-analysis/document-download", testAndExerciseTestResultController.DownloadTestAndExerciseTestResultDocument)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/test-result-analysis", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), testAndExerciseTestResultController.GetTestAndExerciseTestResultDetails)
	businessContinuityRoutes.POST("/test-and-exercises/:testAndExerciseId/test-result-analysis", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseTestResultController.CreateTestAndExerciseTestResult)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/test-result-analysis", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseTestResultController.UpdateTestAndExerciseTestResult)
	businessContinuityRoutes.DELETE("/test-and-exercises/:testAndExerciseId/test-result-analysis", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseTestResultController.DeleteTestAndExerciseTestResults)

	// Recovery Action Plan
	unauthenticatedRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates/:id/document-download", taeRecoveryActionPlanStatusUpdateController.DownloadTaeRecoveryActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates/:id/document-preview", taeRecoveryActionPlanStatusUpdateController.GetTaeRecoveryActionPlanStatusUpdateDocumentPreview)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), taeRecoveryActionPlanStatusUpdateController.GetAllTaeRecoveryActionPlanStatusUpdates)
	businessContinuityRoutes.POST("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), taeRecoveryActionPlanStatusUpdateController.CreateTaeRecoveryActionPlanStatusUpdate)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), taeRecoveryActionPlanStatusUpdateController.UpdateTaeRecoveryActionPlanStatusUpdate)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), taeRecoveryActionPlanStatusUpdateController.GetLatestTaeRecoveryActionPlanStatusUpdate)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/approve", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), taeRecoveryActionPlanStatusUpdateController.ApproveTaeRecoveryActionPlan)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/revert", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), taeRecoveryActionPlanStatusUpdateController.RevertTaeRecoveryActionPlan)

	// TaeRecovery Action Plan Comments
	unauthenticatedRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments/:commentId/document-download", taeRecoveryActionPlanCommentController.DownloadTaeRecoveryActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments/:commentId/document-preview", taeRecoveryActionPlanCommentController.GetTaeRecoveryActionPlanCommentDocumentPreview)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), taeRecoveryActionPlanCommentController.GetAllComments)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), taeRecoveryActionPlanCommentController.GetCommentDetails)
	businessContinuityRoutes.POST("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), taeRecoveryActionPlanCommentController.CreateComment)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), taeRecoveryActionPlanCommentController.UpdateComment)
	businessContinuityRoutes.POST("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), taeRecoveryActionPlanCommentController.DeleteComment)

}
