package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	organizationalContextRoutes := router.Group("/")

	organizationalContextRoutes.Use(middlewares.AuthMiddleware())

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	swotPestelQuestionGroupController := app.OrganizationalContextModule.SwotPestelQuestionGroupController
	swotPestelQuestionController := app.OrganizationalContextModule.SwotPestelQuestionController
	swotPestelCategoryController := app.OrganizationalContextModule.SwotPestelCategoryController
	swotPestelQuestionFactorController := app.OrganizationalContextModule.SwotPestelQuestionFactorController
	stakeholderController := app.OrganizationalContextModule.StakeholderController
	evaluationStatusController := app.OrganizationalContextModule.EvaluationStatusController
	stakeholderNeedsAndExpectationController := app.OrganizationalContextModule.StakeholderNeedsAndExpectationController
	stakeholderNeedsAndExpectationEvaluationController := app.OrganizationalContextModule.StakeholderNeedsAndExpectationEvaluationController
	stakeholderEngagementController := app.OrganizationalContextModule.StakeholderEngagementController
	stakeholderEngagementEvaluationController := app.OrganizationalContextModule.StakeholderEngagementEvaluationController
	stakeholderActionPlanController := app.OrganizationalContextModule.StakeholderActionPlanController
	stakeholderActionPlanStatusUpdateController := app.OrganizationalContextModule.StakeholderActionPlanStatusUpdateController
	stakeholderActionPlanCommentController := app.OrganizationalContextModule.StakeholderActionPlanCommentController

	organizationalContextRoutes.GET("/swot-question-groups", swotPestelQuestionGroupController.GetAllSwotPestelQuestionGroups)
	organizationalContextRoutes.GET("/swot-question-groups/:id", swotPestelQuestionGroupController.GetSwotPestelQuestionGroupDetails)
	organizationalContextRoutes.POST("/swot-question-groups", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.CreateSwotPestelQuestionGroup)
	organizationalContextRoutes.PUT("/swot-question-groups/:id", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.UpdateSwotPestelQuestionGroup)
	organizationalContextRoutes.POST("/swot-question-groups/delete", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.DeleteSwotPestelQuestionGroups)
	organizationalContextRoutes.PUT("/swot-question-groups/activate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.ActivateSwotPestelQuestionGroups)
	organizationalContextRoutes.PUT("/swot-question-groups/deactivate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.DeactivateSwotPestelQuestionGroups)

	organizationalContextRoutes.GET("/swot-questions", swotPestelQuestionController.GetAllSwotPestelQuestions)
	organizationalContextRoutes.GET("/swot-questions/:questionId", swotPestelQuestionController.GetSwotPestelQuestionDetails)
	organizationalContextRoutes.POST("/swot-questions", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.CreateSwotPestelQuestion)
	organizationalContextRoutes.PUT("/swot-questions/:questionId", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.UpdateSwotPestelQuestion)
	organizationalContextRoutes.POST("/swot-questions/delete", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.DeleteSwotPestelQuestions)
	organizationalContextRoutes.PUT("/swot-questions/activate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.ActivateSwotPestelQuestions)
	organizationalContextRoutes.PUT("/swot-questions/deactivate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.DeactivateSwotPestelQuestions)

	organizationalContextRoutes.GET("/pestel-question-groups", swotPestelQuestionGroupController.GetAllSwotPestelQuestionGroups)
	organizationalContextRoutes.GET("/pestel-question-groups/:id", swotPestelQuestionGroupController.GetSwotPestelQuestionGroupDetails)
	organizationalContextRoutes.POST("/pestel-question-groups", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.CreateSwotPestelQuestionGroup)
	organizationalContextRoutes.PUT("/pestel-question-groups/:id", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.UpdateSwotPestelQuestionGroup)
	organizationalContextRoutes.POST("/pestel-question-groups/delete", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.DeleteSwotPestelQuestionGroups)
	organizationalContextRoutes.PUT("/pestel-question-groups/activate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.ActivateSwotPestelQuestionGroups)
	organizationalContextRoutes.PUT("/pestel-question-groups/deactivate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.DeactivateSwotPestelQuestionGroups)

	organizationalContextRoutes.GET("/pestel-questions", swotPestelQuestionController.GetAllSwotPestelQuestions)
	organizationalContextRoutes.GET("/pestel-questions/:questionId", swotPestelQuestionController.GetSwotPestelQuestionDetails)
	organizationalContextRoutes.POST("/pestel-questions", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.CreateSwotPestelQuestion)
	organizationalContextRoutes.PUT("/pestel-questions/:questionId", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.UpdateSwotPestelQuestion)
	organizationalContextRoutes.POST("/pestel-questions/delete", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.DeleteSwotPestelQuestions)
	organizationalContextRoutes.PUT("/pestel-questions/activate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.ActivateSwotPestelQuestions)
	organizationalContextRoutes.PUT("/pestel-questions/deactivate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.DeactivateSwotPestelQuestions)

	organizationalContextRoutes.PUT("/swot-questions/:questionId/answer", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionController.UpdateAnswer)
	organizationalContextRoutes.PUT("/swot-questions/:questionId/answer-invitees", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionController.CreateInvitees)
	organizationalContextRoutes.PUT("/swot-questions/:questionId/answer-invitees/delete", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionController.DeleteInvitees)

	//Swot Pestel Category
	organizationalContextRoutes.GET("/swot-pestel-categories", swotPestelCategoryController.GetAllSwotPestelCategories)
	organizationalContextRoutes.GET("/swot-pestel-categories/:id", swotPestelCategoryController.GetSwotPestelCategoryDetails)

	// swot pestel question factor
	organizationalContextRoutes.GET("/swot-questions/factors", swotPestelQuestionFactorController.GetAllSwotPestelFactors)
	organizationalContextRoutes.GET("/swot-questions/:questionId/factors", swotPestelQuestionFactorController.GetAllSwotPestelQuestionFactors)
	organizationalContextRoutes.GET("/swot-questions/:questionId/factors/:id", swotPestelQuestionFactorController.GetSwotPestelQuestionFactorDetails)
	organizationalContextRoutes.POST("/swot-questions/:questionId/factors", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionFactorController.CreateSwotPestelQuestionFactor)
	organizationalContextRoutes.PUT("/swot-questions/:questionId/factors/:id", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionFactorController.UpdateSwotPestelQuestionFactor)
	organizationalContextRoutes.POST("/swot-questions/:questionId/factors/delete", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionFactorController.DeleteSwotPestelQuestionFactors)

	organizationalContextRoutes.PUT("/pestel-questions/:questionId/answer", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionController.UpdateAnswer)
	organizationalContextRoutes.PUT("/pestel-questions/:questionId/answer-invitees", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionController.CreateInvitees)
	organizationalContextRoutes.PUT("/pestel-questions/:questionId/answer-invitees/delete", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionController.DeleteInvitees)

	// swot pestel question factor
	organizationalContextRoutes.GET("/pestel-questions/factors", swotPestelQuestionFactorController.GetAllSwotPestelFactors)
	organizationalContextRoutes.GET("/pestel-questions/:questionId/factors", swotPestelQuestionFactorController.GetAllSwotPestelQuestionFactors)
	organizationalContextRoutes.GET("/pestel-questions/:questionId/factors/:id", swotPestelQuestionFactorController.GetSwotPestelQuestionFactorDetails)
	organizationalContextRoutes.POST("/pestel-questions/:questionId/factors", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionFactorController.CreateSwotPestelQuestionFactor)
	organizationalContextRoutes.PUT("/pestel-questions/:questionId/factors/:id", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionFactorController.UpdateSwotPestelQuestionFactor)
	organizationalContextRoutes.POST("/pestel-questions/:questionId/factors/delete", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionFactorController.DeleteSwotPestelQuestionFactors)

	// stakeholder
	organizationalContextRoutes.GET("/stakeholders", stakeholderController.GetAllStakeholders)
	organizationalContextRoutes.GET("/stakeholders/:stakeholderId", stakeholderController.GetStakeholderDetails)
	organizationalContextRoutes.POST("/stakeholders", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderController.CreateStakeholder)
	organizationalContextRoutes.PUT("/stakeholders/:stakeholderId", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderController.UpdateStakeholder)
	organizationalContextRoutes.POST("/stakeholders/delete", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderController.DeleteStakeholders)

	// Evaluation Status
	organizationalContextRoutes.GET("/evaluation-statuses", evaluationStatusController.GetAllEvaluationStatuses)
	organizationalContextRoutes.GET("/evaluation-statuses/:id", evaluationStatusController.GetEvaluationStatusDetails)

	// Stakeholder Needs and expectation
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/needs-and-expectations", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderNeedsAndExpectationController.GetAllStakeholderNeedsAndExpectations)
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/needs-and-expectations/:needsAndExpectationId", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderNeedsAndExpectationController.GetStakeholderNeedsAndExpectationDetails)
	organizationalContextRoutes.POST("/stakeholder/:stakeholderId/needs-and-expectations", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderNeedsAndExpectationController.CreateStakeholderNeedsAndExpectation)
	organizationalContextRoutes.PUT("/stakeholder/:stakeholderId/needs-and-expectations/:needsAndExpectationId", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderNeedsAndExpectationController.UpdateStakeholderNeedsAndExpectation)
	organizationalContextRoutes.POST("/stakeholder/:stakeholderId/needs-and-expectations/delete", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderNeedsAndExpectationController.DeleteStakeholderNeedsAndExpectations)

	// Needs and expectation evaluation
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/needs-and-expectations/:needsAndExpectationId/evaluations", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderNeedsAndExpectationEvaluationController.GetAllStakeholderNeedsAndExpectationEvaluations)
	organizationalContextRoutes.POST("/stakeholder/:stakeholderId/needs-and-expectations/:needsAndExpectationId/evaluations", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderNeedsAndExpectationEvaluationController.CreateStakeholderNeedsAndExpectationEvaluation)
	organizationalContextRoutes.PUT("/stakeholder/:stakeholderId/needs-and-expectations/:needsAndExpectationId/evaluations/:id", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderNeedsAndExpectationEvaluationController.UpdateStakeholderNeedsAndExpectationEvaluation)
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/needs-and-expectations/:needsAndExpectationId/evaluations/latest", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderNeedsAndExpectationEvaluationController.GetLatestStakeholderNeedsAndExpectationEvaluation)

	// Stakeholder Engagement
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/engagements", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderEngagementController.GetAllStakeholderEngagements)
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/engagements/:engagementId", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderEngagementController.GetStakeholderEngagementDetails)
	organizationalContextRoutes.POST("/stakeholder/:stakeholderId/engagements", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderEngagementController.CreateStakeholderEngagement)
	organizationalContextRoutes.PUT("/stakeholder/:stakeholderId/engagements/:engagementId", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderEngagementController.UpdateStakeholderEngagement)
	organizationalContextRoutes.POST("/stakeholder/:stakeholderId/engagements/delete", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderEngagementController.DeleteStakeholderEngagements)

	// Engagement evaluation
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/engagements/:engagementId/evaluations", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderEngagementEvaluationController.GetAllStakeholderEngagementEvaluations)
	organizationalContextRoutes.POST("/stakeholder/:stakeholderId/engagements/:engagementId/evaluations", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderEngagementEvaluationController.CreateStakeholderEngagementEvaluation)
	organizationalContextRoutes.PUT("/stakeholder/:stakeholderId/engagements/:engagementId/evaluations/:id", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderEngagementEvaluationController.UpdateStakeholderEngagementEvaluation)
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/engagements/:engagementId/evaluations/latest", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderEngagementEvaluationController.GetLatestStakeholderEngagementEvaluation)

	// stakeholder action plan
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/action-plans", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderActionPlanController.GetAllStakeholderActionPlans)
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderActionPlanController.GetStakeholderActionPlanDetails)
	organizationalContextRoutes.POST("/stakeholder/:stakeholderId/action-plans", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderActionPlanController.CreateStakeholderActionPlan)
	organizationalContextRoutes.PUT("/stakeholder/:stakeholderId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderActionPlanController.UpdateStakeholderActionPlan)
	organizationalContextRoutes.POST("/stakeholder/:stakeholderId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderActionPlanController.DeleteStakeholderActionPlans)

	// Stakeholder Status Update
	unauthenticatedRoutes.GET("/stakeholder/:stakeholderId/action-plans/:actionPlanId/status-updates/:id/document-download", stakeholderActionPlanStatusUpdateController.DownloadStakeholderActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/stakeholder/:stakeholderId/action-plans/:actionPlanId/status-updates/:id/document-preview", stakeholderActionPlanStatusUpdateController.GetStakeholderActionPlanStatusUpdateDocumentPreview)
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "stakeholder_register", "view"), stakeholderActionPlanStatusUpdateController.GetAllStakeholderActionPlanStatusUpdates)
	organizationalContextRoutes.POST("/stakeholder/:stakeholderId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "stakeholder_register", "modify"), stakeholderActionPlanStatusUpdateController.CreateStakeholderActionPlanStatusUpdate)
	organizationalContextRoutes.PUT("/stakeholder/:stakeholderId/action-plans/:actionPlanId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "stakeholder_register", "modify"), stakeholderActionPlanStatusUpdateController.UpdateStakeholderActionPlanStatusUpdate)
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "stakeholder_register", "view"), stakeholderActionPlanStatusUpdateController.GetLatestStakeholderActionPlanStatusUpdate)
	organizationalContextRoutes.PUT("/stakeholder/:stakeholderId/action-plans/:actionPlanId/approve", middlewares.RolePermissionMiddleware(db, "stakeholder_register", "modify"), stakeholderActionPlanStatusUpdateController.ApproveStakeholderActionPlan)
	organizationalContextRoutes.PUT("/stakeholder/:stakeholderId/action-plans/:actionPlanId/revert", middlewares.RolePermissionMiddleware(db, "stakeholder_register", "modify"), stakeholderActionPlanStatusUpdateController.RevertStakeholderActionPlan)

	// Stakeholder Action Plan Comments
	unauthenticatedRoutes.GET("/stakeholder/:stakeholderId/action-plans/:actionPlanId/comments/:commentId/document-download", stakeholderActionPlanCommentController.DownloadStakeholderActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/stakeholder/:stakeholderId/action-plans/:actionPlanId/comments/:commentId/document-preview", stakeholderActionPlanCommentController.GetStakeholderActionPlanCommentDocumentPreview)
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderActionPlanCommentController.GetAllComments)
	organizationalContextRoutes.GET("/stakeholder/:stakeholderId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "view"), stakeholderActionPlanCommentController.GetCommentDetails)
	organizationalContextRoutes.POST("/stakeholder/:stakeholderId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderActionPlanCommentController.CreateComment)
	organizationalContextRoutes.PUT("/stakeholder/:stakeholderId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderActionPlanCommentController.UpdateComment)
	organizationalContextRoutes.POST("/stakeholder/:stakeholderId/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "stakeholder_analysis", "modify"), stakeholderActionPlanCommentController.DeleteComment)

}
