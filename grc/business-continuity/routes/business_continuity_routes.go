package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	businessContinuityRoutes := router.Group("/")

	businessContinuityRoutes.Use(middlewares.AuthMiddleware(app.DB))

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
	bcmImpactAnalysisLevelController := app.BusinessContinuityModule.BCMImpactAnalysisLevelController
	bcmBiaOutageBucketController := app.BusinessContinuityModule.BCMBIAOutageBucketController
	bcmBiaImpactLevelController := app.BusinessContinuityModule.BCMBIAImpactLevelController
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
	bcpAssociatedItemController := app.BusinessContinuityModule.BcpAssociatedItemController
	businessContinuityDurationUnitController := app.BusinessContinuityModule.BusinessContinuityDurationUnitController
	jobTeamController := app.BusinessContinuityModule.JobTeamController
	bcmExecutionRecordTypeController := app.BusinessContinuityModule.BCMExecutionRecordTypeController
	bcmExecutionRecordController := app.BusinessContinuityModule.BCMExecutionRecordController
	bcmProcessTypeController := app.BusinessContinuityModule.BCMProcessTypeController
	bcmManualProcedureFormController := app.BusinessContinuityModule.BCMManualProcedureFormController
	bcmStaffAllocationController := app.BusinessContinuityModule.BCMStaffAllocationController
	biaProcessController := app.BusinessContinuityModule.BiaProcessController
	biaProcessReviewController := app.BusinessContinuityModule.BiaProcessReviewController
	biaProcessWorkflowController := app.BusinessContinuityModule.BiaProcessWorkflowController
	biaProcessReviewCommentController := app.BusinessContinuityModule.BiaProcessReviewCommentController
	businessContinuityWorkflowLevelController := app.BusinessContinuityModule.BusinessContinuityWorkflowLevelController

	// businessfunctionAssociatedItemController := app.BusinessContinuityModule.BusinessFunctionAssociateItemController

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

	businessContinuityRoutes.GET("/business-functions/:id/overall-test-results", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetBusinessFunctionTestResults)

	businessContinuityRoutes.GET("/business-functions/:id/processes/:processId/associated-risks", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetProcessAssociatedRisks)
	businessContinuityRoutes.GET("/business-functions/:id/processes/:processId/associated-incidents", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetProcessAssociatedIncidents)
	businessContinuityRoutes.GET("/business-functions/:id/processes/:processId/associated-controls", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.GetProcessAssociatedControls)
	businessContinuityRoutes.POST("/business-functions/:id/processes", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.CreateBusinessFunctionProcess)
	businessContinuityRoutes.DELETE("/business-functions/:id/processes/:processId", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessFunctionController.DeleteBusinessFunctionProcess)
	// Business Impact Assessment
	businessContinuityRoutes.GET("/bia/processes", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcesses)
	businessContinuityRoutes.GET("/bia/processes/export", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessesExport)
	businessContinuityRoutes.GET("/bia/export", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllBIAExport)
	businessContinuityRoutes.GET("/bia/processes/:processId", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetProcessDetails)
	businessContinuityRoutes.GET("/bia/processes/:processId/incidents", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessIncidents)
	businessContinuityRoutes.GET("/bia/processes/:processId/assets", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessAssets)
	businessContinuityRoutes.GET("/bia/processes/:processId/documents", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessDocuments)
	businessContinuityRoutes.GET("/bia/processes/:processId/risks", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllProcessRisks)
	businessContinuityRoutes.GET("/bia/processes/:processId/dependencies", businessImpactAssessmentController.GetAllProcessDependencies)

	businessContinuityRoutes.POST("/business-impact-assessments", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.CreateBusinessImpactAssessment)
	businessContinuityRoutes.GET("/bia-impact-analysis", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllBIAImpactAnalysis)
	businessContinuityRoutes.GET("/bia-recovery-objectives", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetAllBIARecoveryObjectives)
	businessContinuityRoutes.PUT("/bia-impact-analysis", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.UpdateImpactAnalysis)
	businessContinuityRoutes.PUT("/bia-impact-analysis-ratings/:id", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.UpdateBIARating)
	businessContinuityRoutes.PUT("/bia-recovery-objectives", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.UpdateRecoveryObjective)
	businessContinuityRoutes.PUT("/bia-process-criticality", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.UpdateProcessCriticality)
	businessContinuityRoutes.POST("/business-impact-assessments/delete", middlewares.RolePermissionMiddleware(db, "bia", "modify"), businessImpactAssessmentController.DeleteBusinessImpactAssessments)
	businessContinuityRoutes.GET("/bia-matrix", middlewares.RolePermissionMiddleware(db, "bia_matrics", "view"), businessImpactAssessmentController.GetBIAMatrix)
	businessContinuityRoutes.GET("/bia-matrix/export", middlewares.RolePermissionMiddleware(db, "bia_matrics", "view"), businessImpactAssessmentController.GetBiaMatrixExport)
	businessContinuityRoutes.GET("/bia/critical-systems", middlewares.RolePermissionMiddleware(db, "bia", "view"), businessImpactAssessmentController.GetCriticalSystems)

	// Business Continuity Plan
	unauthenticatedRoutes.GET("/business-continuity-plans/:id/document-preview", businessContinuityPlanController.GetBusinessContinuityPlanDocumentPreview)
	unauthenticatedRoutes.GET("/business-continuity-plans/:id/document-download", businessContinuityPlanController.DownloadBusinessContinuityPlanDocument)
	businessContinuityRoutes.GET("/business-continuity-plans", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetAllBusinessContinuityPlans)
	businessContinuityRoutes.GET("/business-continuity-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetBusinessContinuityPlanDetails)
	businessContinuityRoutes.POST("/business-continuity-plans", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.CreateBusinessContinuityPlan)
	businessContinuityRoutes.PUT("/business-continuity-plans/:id", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.UpdateBusinessContinuityPlan)
	businessContinuityRoutes.PUT("/business-continuity-plans/:id/publish", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.MarkAsCompleted)
	businessContinuityRoutes.POST("/business-continuity-plans/delete", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.DeleteBusinessContinuityPlans)
	businessContinuityRoutes.GET("/business-continuity-plans/export", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetAllBusinessContinuityPlanExport)
	businessContinuityRoutes.GET("/business-continuity-plans/:id/test-summary", businessContinuityPlanController.GetTestSummary)
	businessContinuityRoutes.POST("/business-continuity-plans/create-version", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), businessContinuityPlanController.CreateVersion)
	businessContinuityRoutes.GET("/business-continuity-plans/:id/version-history", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.GetVersionHistory)
	businessContinuityRoutes.GET("/business-continuity-plans/:id/export-pdf", middlewares.RolePermissionMiddleware(db, "bcp", "view"), businessContinuityPlanController.ExportBusinessContinuityPlanToPDF)

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
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/mark-as-completed", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseController.MarkAsCompleted)
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
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/export", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), testAndExerciseRecoveryActionPlanController.GetAllTestAndExerciseRecoveryActionPlansExport)
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
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/test-result-summary", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), testAndExerciseController.GetTestResultSummary)
	businessContinuityRoutes.PUT("/test-result-system-registers/:id", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseTestResultController.UpdateTestResultSystemRegister)

	// Test Result Performance Measures
	businessContinuityRoutes.POST("/test-results/:testResultId/performance-measures", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseTestResultController.CreatePerformanceMeasure)
	businessContinuityRoutes.GET("/test-results/:testResultId/performance-measures", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), testAndExerciseTestResultController.GetPerformanceMeasures)
	businessContinuityRoutes.PUT("/performance-measures/:id", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseTestResultController.UpdatePerformanceMeasure)
	businessContinuityRoutes.DELETE("/performance-measures/:id", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), testAndExerciseTestResultController.DeletePerformanceMeasure)

	// Recovery Action Plan
	unauthenticatedRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates/:id/document-download", taeRecoveryActionPlanStatusUpdateController.DownloadTaeRecoveryActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates/:id/document-preview", taeRecoveryActionPlanStatusUpdateController.GetTaeRecoveryActionPlanStatusUpdateDocumentPreview)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), taeRecoveryActionPlanStatusUpdateController.GetAllTaeRecoveryActionPlanStatusUpdates)
	businessContinuityRoutes.POST("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates", taeRecoveryActionPlanStatusUpdateController.CreateTaeRecoveryActionPlanStatusUpdate)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates/:id", taeRecoveryActionPlanStatusUpdateController.UpdateTaeRecoveryActionPlanStatusUpdate)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), taeRecoveryActionPlanStatusUpdateController.GetLatestTaeRecoveryActionPlanStatusUpdate)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/approve", taeRecoveryActionPlanStatusUpdateController.ApproveTaeRecoveryActionPlan)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/revert", taeRecoveryActionPlanStatusUpdateController.RevertTaeRecoveryActionPlan)

	// TaeRecovery Action Plan Comments
	unauthenticatedRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments/:commentId/document-download", taeRecoveryActionPlanCommentController.DownloadTaeRecoveryActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments/:commentId/document-preview", taeRecoveryActionPlanCommentController.GetTaeRecoveryActionPlanCommentDocumentPreview)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), taeRecoveryActionPlanCommentController.GetAllComments)
	businessContinuityRoutes.GET("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "view"), taeRecoveryActionPlanCommentController.GetCommentDetails)
	businessContinuityRoutes.POST("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), taeRecoveryActionPlanCommentController.CreateComment)
	businessContinuityRoutes.PUT("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), taeRecoveryActionPlanCommentController.UpdateComment)
	businessContinuityRoutes.POST("/test-and-exercises/:testAndExerciseId/recovery-action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "test_and_exercise", "modify"), taeRecoveryActionPlanCommentController.DeleteComment)

	businessContinuityRoutes.GET("/business-continuity-plans/:id/documents", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpAssociatedItemController.GetAllBcpDocuments)
	businessContinuityRoutes.POST("/business-continuity-plans/:id/documents", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.CreateBcpDocument)
	businessContinuityRoutes.DELETE("/business-continuity-plans/:id/documents/:documentId", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.DeleteBcpDocument)

	businessContinuityRoutes.GET("/business-continuity-plans/:id/controls", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpAssociatedItemController.GetAllBcpControl)
	businessContinuityRoutes.POST("/business-continuity-plans/:id/controls", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.CreateBcpControl)
	businessContinuityRoutes.DELETE("/business-continuity-plans/:id/controls/:controlId", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.DeleteBcpControl)

	businessContinuityRoutes.GET("/business-continuity-plans/:id/kpis", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpAssociatedItemController.GetAllBcpKPI)
	businessContinuityRoutes.POST("/business-continuity-plans/:id/kpis", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.CreateBcpKPI)
	businessContinuityRoutes.DELETE("/business-continuity-plans/:id/kpis/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.DeleteBcpKPI)

	businessContinuityRoutes.GET("/business-continuity-plans/:id/risks", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpAssociatedItemController.GetAllBcpRisk)
	businessContinuityRoutes.POST("/business-continuity-plans/:id/risks", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.CreateBcpRisk)
	businessContinuityRoutes.DELETE("/business-continuity-plans/:id/risks/:riskId", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.DeleteBcpRisk)

	// businessContinuityRoutes.GET("/business-functions/:id/processes", middlewares.RolePermissionMiddleware(db, "business_functions", "view"), businessfunctionAssociatedItemController.GetAllBusinessFunctionProcesses)
	// businessContinuityRoutes.POST("/business-functions/:id/processes", middlewares.RolePermissionMiddleware(db, "business_functions", "modify"), businessfunctionAssociatedItemController.CreateBusinessFunctionProcess)
	// businessContinuityRoutes.DELETE("/business-functions/:id/processes/:processId", middlewares.RolePermissionMiddleware(db, "business_functions", "modify"), businessfunctionAssociatedItemController.DeleteBusinessFunctionProcess)

	businessContinuityRoutes.GET("/business-continuity-plans/:id/assets", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpAssociatedItemController.GetAllBcpAssets)
	businessContinuityRoutes.POST("/business-continuity-plans/:id/assets", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.CreateBcpAsset)
	businessContinuityRoutes.DELETE("/business-continuity-plans/:id/assets/:assetId", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.DeleteBcpAsset)

	businessContinuityRoutes.GET("/business-continuity-plans/:id/incidents", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpAssociatedItemController.GetAllBcpIncidents)
	businessContinuityRoutes.POST("/business-continuity-plans/:id/incidents", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.CreateBcpIncident)
	businessContinuityRoutes.DELETE("/business-continuity-plans/:id/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.DeleteBcpIncident)

	businessContinuityRoutes.GET("/business-continuity-plans/:id/strategic-objectives", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpAssociatedItemController.GetAllBcpStrategicObjectives)
	businessContinuityRoutes.POST("/business-continuity-plans/:id/strategic-objectives", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.CreateBcpStrategicObjectives)
	businessContinuityRoutes.DELETE("/business-continuity-plans/:id/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.DeleteBcpStrategicObjectives)

	businessContinuityRoutes.GET("/business-continuity-plans/:id/clauses", middlewares.RolePermissionMiddleware(db, "bcp", "view"), bcpAssociatedItemController.GetAllClauses)
	businessContinuityRoutes.PUT("/business-continuity-plans/:id/clauses", middlewares.RolePermissionMiddleware(db, "bcp", "modify"), bcpAssociatedItemController.UpdateClauses)

	// business continuity duration units
	businessContinuityRoutes.GET("/business-continuity-duration-units", businessContinuityDurationUnitController.GetAllBusinessContinuityDurationUnits)
	businessContinuityRoutes.GET("/business-continuity-duration-units/:id", businessContinuityDurationUnitController.GetBusinessContinuityDurationUnitDetails)

	// job teams
	businessContinuityRoutes.GET("/job-teams", jobTeamController.GetAllJobTeams)
	businessContinuityRoutes.GET("/job-teams/:id", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "view"), jobTeamController.GetJobTeamDetails)
	businessContinuityRoutes.POST("/job-teams", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), jobTeamController.CreateJobTeam)
	businessContinuityRoutes.PUT("/job-teams/:id", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), jobTeamController.UpdateJobTeam)
	businessContinuityRoutes.POST("/job-teams/delete", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), jobTeamController.DeleteJobTeams)

	// bcm execution record types
	businessContinuityRoutes.GET("/bcm-execution-record-types", bcmExecutionRecordTypeController.GetAllBCMExecutionRecordTypes)
	businessContinuityRoutes.GET("/bcm-execution-record-types/:id", bcmExecutionRecordTypeController.GetBCMExecutionRecordTypeDetails)
	businessContinuityRoutes.POST("/bcm-execution-record-types", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmExecutionRecordTypeController.CreateBCMExecutionRecordType)
	businessContinuityRoutes.PUT("/bcm-execution-record-types/:id", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmExecutionRecordTypeController.UpdateBCMExecutionRecordType)
	businessContinuityRoutes.POST("/bcm-execution-record-types/delete", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmExecutionRecordTypeController.DeleteBCMExecutionRecordTypes)
	businessContinuityRoutes.PUT("/bcm-execution-record-types/activate", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmExecutionRecordTypeController.ActivateBCMExecutionRecordTypes)
	businessContinuityRoutes.PUT("/bcm-execution-record-types/deactivate", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmExecutionRecordTypeController.DeactivateBCMExecutionRecordTypes)

	// bcm execution records
	businessContinuityRoutes.GET("/bcm-execution-records", bcmExecutionRecordController.GetAllBCMExecutionRecords)
	businessContinuityRoutes.GET("/bcm-execution-records/:id", bcmExecutionRecordController.GetBCMExecutionRecordFullDetails)
	businessContinuityRoutes.POST("/bcm-execution-records", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmExecutionRecordController.CreateBCMExecutionRecord)
	businessContinuityRoutes.PUT("/bcm-execution-records/:id", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmExecutionRecordController.UpdateBCMExecutionRecord)
	businessContinuityRoutes.POST("/bcm-execution-records/delete", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmExecutionRecordController.DeleteBCMExecutionRecords)

	// bcm process types
	businessContinuityRoutes.GET("/bcm-process-types", bcmProcessTypeController.GetAllBCMProcessTypes)
	businessContinuityRoutes.GET("/bcm-process-types/:id", bcmProcessTypeController.GetBCMProcessTypeDetails)
	businessContinuityRoutes.POST("/bcm-process-types", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmProcessTypeController.CreateBCMProcessType)
	businessContinuityRoutes.PUT("/bcm-process-types/:id", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmProcessTypeController.UpdateBCMProcessType)
	businessContinuityRoutes.POST("/bcm-process-types/delete", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmProcessTypeController.DeleteBCMProcessTypes)
	businessContinuityRoutes.PUT("/bcm-process-types/activate", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmProcessTypeController.ActivateBCMProcessTypes)
	businessContinuityRoutes.PUT("/bcm-process-types/deactivate", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmProcessTypeController.DeactivateBCMProcessTypes)

	// bcm manual procedure forms
	businessContinuityRoutes.GET("/bcm-manual-procedure-forms", bcmManualProcedureFormController.GetAllBCMManualProcedureForms)
	businessContinuityRoutes.GET("/bcm-manual-procedure-forms/:id", bcmManualProcedureFormController.GetBCMManualProcedureFormFullDetails)
	businessContinuityRoutes.POST("/bcm-manual-procedure-forms", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmManualProcedureFormController.CreateBCMManualProcedureForm)
	businessContinuityRoutes.PUT("/bcm-manual-procedure-forms/:id", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmManualProcedureFormController.UpdateBCMManualProcedureForm)
	businessContinuityRoutes.POST("/bcm-manual-procedure-forms/delete", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmManualProcedureFormController.DeleteBCMManualProcedureForms)

	// bcm impact analysis level
	businessContinuityRoutes.GET("/bcm-impact-analysis-levels", bcmImpactAnalysisLevelController.GetAllBCMImpactAnalysisLevels)
	businessContinuityRoutes.GET("/bcm-impact-analysis-levels/:id", bcmImpactAnalysisLevelController.GetBCMImpactAnalysisLevelDetails)

	// bcm bia outage bucket
	businessContinuityRoutes.GET("/bcm-bia-outage-buckets", bcmBiaOutageBucketController.GetAllBCMBIAOutageBuckets)
	businessContinuityRoutes.GET("/bcm-bia-outage-buckets/:id", bcmBiaOutageBucketController.GetBCMBIAOutageBucketDetails)

	// bcm bia impact level
	businessContinuityRoutes.GET("/bcm-bia-impact-levels", bcmBiaImpactLevelController.GetAllBCMBIAImpactLevels)
	businessContinuityRoutes.GET("/bcm-bia-impact-levels/:id", bcmBiaImpactLevelController.GetBCMBIAImpactLevelFullDetails)
	businessContinuityRoutes.GET("/bcm-bia-impact-levels/analysis-level/:analysisLevelId", bcmBiaImpactLevelController.GetBCMBIAImpactLevelByImpactAnalysisLevelID)
	businessContinuityRoutes.PUT("/bcm-bia-impact-levels", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), bcmBiaImpactLevelController.UpdateBCMBIAImpactLevels)

	// Staff Allocations (auto-generated, user updates is_remote per slot)
	businessContinuityRoutes.GET("/bia/processes/:processId/staff-allocations", middlewares.RolePermissionMiddleware(db, "bia", "view"), bcmStaffAllocationController.GetFullDetailsByProcessID)
	businessContinuityRoutes.PUT("/bia/processes/:processId/staff-allocations", middlewares.RolePermissionMiddleware(db, "bia", "modify"), bcmStaffAllocationController.UpdateAllocation)

	// BIA Process
	businessContinuityRoutes.GET("/bia-process-details/:processId", middlewares.RolePermissionMiddleware(db, "bia", "view"), biaProcessController.GetBiaProcessDetails)

	businessContinuityRoutes.POST("/bia-process/:id/submit", middlewares.RolePermissionMiddleware(db, "bia", "modify"), biaProcessReviewController.SubmitBiaProcess)
	businessContinuityRoutes.POST("/bia-process/:id/approve", middlewares.RolePermissionMiddleware(db, "bia", "modify"), biaProcessReviewController.ApproveBiaProcess)
	businessContinuityRoutes.POST("/bia-process/:id/revert", middlewares.RolePermissionMiddleware(db, "bia", "modify"), biaProcessReviewController.RejectBiaProcess)
	businessContinuityRoutes.POST("/bia-process/:id/reopen", middlewares.RolePermissionMiddleware(db, "bia", "modify"), biaProcessReviewController.ReopenBiaProcess)
	businessContinuityRoutes.GET("/bia-process/:id/workflow-list", middlewares.RolePermissionMiddleware(db, "bia", "view"), biaProcessWorkflowController.GetWorkflowList)

	businessContinuityRoutes.GET("/bia-process/:id/review-comments", middlewares.RolePermissionMiddleware(db, "bia", "view"), biaProcessReviewCommentController.GetAllComments)
	businessContinuityRoutes.GET("/bia-process/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "bia", "view"), biaProcessReviewCommentController.GetCommentDetails)
	businessContinuityRoutes.POST("/bia-process/:id/review-comments", middlewares.RolePermissionMiddleware(db, "bia", "modify"), biaProcessReviewCommentController.CreateComment)
	businessContinuityRoutes.PUT("/bia-process/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "bia", "modify"), biaProcessReviewCommentController.UpdateComment)
	businessContinuityRoutes.POST("/bia-process/:id/review-comments/delete", middlewares.RolePermissionMiddleware(db, "bia", "modify"), biaProcessReviewCommentController.DeleteComments)
	businessContinuityRoutes.PUT("/bia-process/:id/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(db, "bia", "modify"), biaProcessReviewCommentController.ResolveComment)

	unauthenticatedRoutes.GET("/bia-process/:id/review-comments/:commentId/document-preview", biaProcessReviewCommentController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/bia-process/:id/review-comments/:commentId/document-download", biaProcessReviewCommentController.DownloadDocument)

	// Business Continuity Workflow Level
	businessContinuityRoutes.POST("/business-continuity-workflow/:workflowId/levels", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), businessContinuityWorkflowLevelController.CreateWorkflowLevel)
	businessContinuityRoutes.DELETE("/business-continuity-workflow/:workflowId/levels/:levelId", middlewares.RolePermissionMiddleware(db, "business_continuity_setting", "modify"), businessContinuityWorkflowLevelController.DeleteWorkflowLevel)

}
