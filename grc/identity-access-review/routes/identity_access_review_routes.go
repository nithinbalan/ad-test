package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	identityAccessReviewRoutes := router.Group("/")
	unauthenticatedRoutes := router.Group("/")

	identityAccessReviewRoutes.Use(middlewares.AuthMiddleware(app.DB))

	reviewTypeController := app.IdentityAccessReviewModule.ReviewTypeController
	reviewCampaignStatusController := app.IdentityAccessReviewModule.ReviewCampaignStatusController
	identityReviewCampaignController := app.IdentityAccessReviewModule.IdentityReviewCampaignController
	identityReviewCampaignReviewController := app.IdentityAccessReviewModule.IdentityReviewCampaignReviewController
	identityReviewCampaignWorkflowController := app.IdentityAccessReviewModule.IdentityReviewCampaignWorkflowController
	identityReviewCampaignReviewCommentController := app.IdentityAccessReviewModule.IdentityReviewCampaignReviewCommentController
	identityReviewUserController := app.IdentityAccessReviewModule.IdentityReviewUserController
	identityReviewUserCommentController := app.IdentityAccessReviewModule.IdentityReviewUserCommentController
	privilegedAccountController := app.IdentityAccessReviewModule.PrivilegedAccountController
	accountTypeController := app.IdentityAccessReviewModule.AccountTypeController
	accountDecisionController := app.IdentityAccessReviewModule.AccountDecisionController
	identityReviewActionPlanController := app.IdentityAccessReviewModule.IdentityReviewActionPlanController
	identityReviewActionPlanStatusUpdateController := app.IdentityAccessReviewModule.IdentityReviewActionPlanStatusUpdateController
	identityReviewActionPlanCommentController := app.IdentityAccessReviewModule.IdentityReviewActionPlanCommentController

	// review types
	identityAccessReviewRoutes.GET("/identity-access-review/review-types", reviewTypeController.GetAllReviewTypes)
	identityAccessReviewRoutes.GET("/identity-access-review/review-types/:id", reviewTypeController.GetReviewTypeDetails)

	// review campaign statuses
	identityAccessReviewRoutes.GET("/identity-access-review/review-campaign-statuses", reviewCampaignStatusController.GetAllReviewCampaignStatuses)
	identityAccessReviewRoutes.GET("/identity-access-review/review-campaign-statuses/:id", reviewCampaignStatusController.GetReviewCampaignStatusDetails)

	// identity review campaigns
	identityAccessReviewRoutes.GET("/identity-access-review/review-campaigns", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewCampaignController.GetAllIdentityReviewCampaigns)
	identityAccessReviewRoutes.GET("/identity-access-review/review-campaigns/:id", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewCampaignController.GetIdentityReviewCampaignDetails)
	identityAccessReviewRoutes.GET("/identity-access-review/review-campaigns/:id/export-pdf", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewCampaignController.ExportIdentityReviewCampaignZip)
	identityAccessReviewRoutes.POST("/identity-access-review/review-campaigns", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignController.CreateIdentityReviewCampaign)
	identityAccessReviewRoutes.PUT("/identity-access-review/review-campaigns/:id", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignController.UpdateIdentityReviewCampaign)
	identityAccessReviewRoutes.POST("/identity-access-review/review-campaigns/delete", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignController.DeleteIdentityReviewCampaigns)

	identityAccessReviewRoutes.POST("/identity-access-review/review-campaigns/:id/submit", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignReviewController.SubmitIdentityReviewCampaign)
	identityAccessReviewRoutes.POST("/identity-access-review/review-campaigns/:id/approve", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignReviewController.ApproveIdentityReviewCampaign)
	identityAccessReviewRoutes.POST("/identity-access-review/review-campaigns/:id/reject", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignReviewController.RejectIdentityReviewCampaign)
	identityAccessReviewRoutes.POST("/identity-access-review/review-campaigns/:id/resubmit", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignReviewController.ResubmitIdentityReviewCampaign)
	identityAccessReviewRoutes.POST("/identity-access-review/review-campaigns/:id/revert", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignReviewController.RevertIdentityReviewCampaign)
	identityAccessReviewRoutes.GET("/identity-access-review/review-campaigns/:id/workflow-list", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewCampaignWorkflowController.GetWorkflowList)

	identityAccessReviewRoutes.GET("/identity-access-review/review-campaigns/:id/review-comments", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewCampaignReviewCommentController.GetAllComments)
	identityAccessReviewRoutes.GET("/identity-access-review/review-campaigns/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewCampaignReviewCommentController.GetCommentDetails)
	identityAccessReviewRoutes.POST("/identity-access-review/review-campaigns/:id/review-comments", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignReviewCommentController.CreateComment)
	identityAccessReviewRoutes.PUT("/identity-access-review/review-campaigns/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignReviewCommentController.UpdateComment)
	identityAccessReviewRoutes.POST("/identity-access-review/review-campaigns/:id/review-comments/delete", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignReviewCommentController.DeleteComment)
	identityAccessReviewRoutes.PUT("/identity-access-review/review-campaigns/:id/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignReviewCommentController.ResolveComment)

	// identity review users
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/template", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserController.GenerateIdentityReviewUserTemplate)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/template/validate", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserController.ValidateIdentityReviewUserTemplate)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/import/validate", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserController.ValidateImportIdentityReviewUsers)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/import", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserController.ImportIdentityReviewUsers)

	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewUserController.GetAllIdentityReviewUsers)
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewUserController.GetIdentityReviewUserDetails)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserController.CreateIdentityReviewUser)
	identityAccessReviewRoutes.PUT("/identity-access-review/identity-review-users/:id", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserController.UpdateIdentityReviewUser)
	identityAccessReviewRoutes.PUT("/identity-access-review/identity-review-users/:id/account-decision", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserController.UpdateIdentityReviewUserAccountDecision)
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id/account-decision-history", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewUserController.GetIdentityReviewUserAccountDecisionHistory)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/delete", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserController.DeleteIdentityReviewUsers)

	// identity review user comments
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id/comments", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewUserCommentController.GetAll)
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id/comments/:commentId", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewUserCommentController.GetDetails)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/:id/comments", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserCommentController.Create)
	identityAccessReviewRoutes.PUT("/identity-access-review/identity-review-users/:id/comments/:commentId", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserCommentController.Update)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/:id/comments/delete", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserCommentController.Delete)

	// identity review action plans
	identityAccessReviewRoutes.GET("/identity-access-review/action-plans/export", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewActionPlanController.ExportIdentityReviewActionPlans)
	identityAccessReviewRoutes.GET("/identity-access-review/action-plans", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewActionPlanController.GetAllIdentityReviewActionPlans)
	identityAccessReviewRoutes.POST("/identity-access-review/action-plans/force-delete", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewActionPlanController.ForceDeleteIdentityReviewActionPlans)
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id/action-plans", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewActionPlanController.GetAllIdentityReviewActionPlans)
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id/action-plan-summary", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewActionPlanController.GetIdentityReviewActionPlanSummary)
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewActionPlanController.GetIdentityReviewActionPlanDetails)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/:id/action-plans", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewActionPlanController.CreateIdentityReviewActionPlan)
	identityAccessReviewRoutes.PUT("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewActionPlanController.UpdateIdentityReviewActionPlan)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/:id/action-plans/delete", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewActionPlanController.DeleteIdentityReviewActionPlans)

	// identity review action plan status updates
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewActionPlanStatusUpdateController.GetAllIdentityReviewActionPlanStatusUpdates)
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewActionPlanStatusUpdateController.GetLatestIdentityReviewActionPlanStatusUpdate)
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewActionPlanStatusUpdateController.GetIdentityReviewActionPlanStatusUpdateDetails)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewActionPlanStatusUpdateController.CreateIdentityReviewActionPlanStatusUpdate)
	identityAccessReviewRoutes.PUT("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewActionPlanStatusUpdateController.UpdateIdentityReviewActionPlanStatusUpdate)
	identityAccessReviewRoutes.PUT("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/approve", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewActionPlanStatusUpdateController.ApproveIdentityReviewActionPlan)
	identityAccessReviewRoutes.PUT("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/revert", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewActionPlanStatusUpdateController.RevertIdentityReviewActionPlan)

	// identity review action plan comments
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewActionPlanCommentController.GetAllComments)
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewActionPlanCommentController.GetCommentDetails)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewActionPlanCommentController.CreateComment)
	identityAccessReviewRoutes.PUT("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewActionPlanCommentController.UpdateComment)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewActionPlanCommentController.DeleteComment)

	// privileged accounts
	identityAccessReviewRoutes.GET("/identity-access-review/privileged-accounts", privilegedAccountController.GetAllPrivilegedAccounts)
	identityAccessReviewRoutes.GET("/identity-access-review/privileged-accounts/:id", privilegedAccountController.GetPrivilegedAccountDetails)
	identityAccessReviewRoutes.POST("/identity-access-review/privileged-accounts", middlewares.RolePermissionMiddleware(app.DB, "identity_access_review_setting", "modify"), privilegedAccountController.CreatePrivilegedAccount)
	identityAccessReviewRoutes.PUT("/identity-access-review/privileged-accounts/:id", middlewares.RolePermissionMiddleware(app.DB, "identity_access_review_setting", "modify"), privilegedAccountController.UpdatePrivilegedAccount)
	identityAccessReviewRoutes.POST("/identity-access-review/privileged-accounts/delete", middlewares.RolePermissionMiddleware(app.DB, "identity_access_review_setting", "modify"), privilegedAccountController.DeletePrivilegedAccounts)
	identityAccessReviewRoutes.PUT("/identity-access-review/privileged-accounts/activate", middlewares.RolePermissionMiddleware(app.DB, "identity_access_review_setting", "modify"), privilegedAccountController.ActivatePrivilegedAccounts)
	identityAccessReviewRoutes.PUT("/identity-access-review/privileged-accounts/deactivate", middlewares.RolePermissionMiddleware(app.DB, "identity_access_review_setting", "modify"), privilegedAccountController.DeactivatePrivilegedAccounts)

	// account types
	identityAccessReviewRoutes.GET("/identity-access-review/account-types", accountTypeController.GetAllAccountTypes)
	identityAccessReviewRoutes.GET("/identity-access-review/account-types/:id", accountTypeController.GetAccountTypeDetails)

	// account decisions
	identityAccessReviewRoutes.GET("/identity-access-review/account-decisions", accountDecisionController.GetAllAccountDecisions)
	identityAccessReviewRoutes.GET("/identity-access-review/account-decisions/:id", accountDecisionController.GetAccountDecisionDetails)

	unauthenticatedRoutes.GET("/identity-access-review/review-campaigns/:id/document-preview", identityReviewCampaignController.GetIdentityReviewCampaignDocumentPreview)
	unauthenticatedRoutes.GET("/identity-access-review/review-campaigns/:id/document-download", identityReviewCampaignController.DownloadIdentityReviewCampaignDocument)
	unauthenticatedRoutes.GET("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-preview", identityReviewActionPlanStatusUpdateController.GetIdentityReviewActionPlanDocumentPreview)
	unauthenticatedRoutes.GET("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-download", identityReviewActionPlanStatusUpdateController.DownloadIdentityReviewActionPlanDocument)
	unauthenticatedRoutes.GET("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/comments/:commentId/document-preview", identityReviewActionPlanCommentController.GetIdentityReviewActionPlanCommentDocumentPreview)
	unauthenticatedRoutes.GET("/identity-access-review/identity-review-users/:id/action-plans/:actionPlanId/comments/:commentId/document-download", identityReviewActionPlanCommentController.DownloadIdentityReviewActionPlanCommentDocument)
}
