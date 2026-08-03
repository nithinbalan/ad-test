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
	identityReviewUserController := app.IdentityAccessReviewModule.IdentityReviewUserController
	privilegedAccountController := app.IdentityAccessReviewModule.PrivilegedAccountController
	accountTypeController := app.IdentityAccessReviewModule.AccountTypeController

	// review types
	identityAccessReviewRoutes.GET("/identity-access-review/review-types", reviewTypeController.GetAllReviewTypes)
	identityAccessReviewRoutes.GET("/identity-access-review/review-types/:id", reviewTypeController.GetReviewTypeDetails)

	// review campaign statuses
	identityAccessReviewRoutes.GET("/identity-access-review/review-campaign-statuses", reviewCampaignStatusController.GetAllReviewCampaignStatuses)
	identityAccessReviewRoutes.GET("/identity-access-review/review-campaign-statuses/:id", reviewCampaignStatusController.GetReviewCampaignStatusDetails)

	// identity review campaigns
	identityAccessReviewRoutes.GET("/identity-access-review/review-campaigns", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewCampaignController.GetAllIdentityReviewCampaigns)
	identityAccessReviewRoutes.GET("/identity-access-review/review-campaigns/:id", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewCampaignController.GetIdentityReviewCampaignDetails)
	identityAccessReviewRoutes.POST("/identity-access-review/review-campaigns", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignController.CreateIdentityReviewCampaign)
	identityAccessReviewRoutes.PUT("/identity-access-review/review-campaigns/:id", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignController.UpdateIdentityReviewCampaign)
	identityAccessReviewRoutes.POST("/identity-access-review/review-campaigns/delete", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewCampaignController.DeleteIdentityReviewCampaigns)

	// identity review users
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewUserController.GetAllIdentityReviewUsers)
	identityAccessReviewRoutes.GET("/identity-access-review/identity-review-users/:id", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "view"), identityReviewUserController.GetIdentityReviewUserDetails)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserController.CreateIdentityReviewUser)
	identityAccessReviewRoutes.PUT("/identity-access-review/identity-review-users/:id", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserController.UpdateIdentityReviewUser)
	identityAccessReviewRoutes.POST("/identity-access-review/identity-review-users/delete", middlewares.RolePermissionMiddleware(app.DB, "review_campaign", "modify"), identityReviewUserController.DeleteIdentityReviewUsers)

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

	unauthenticatedRoutes.GET("/identity-access-review/review-campaigns/:id/document-preview", identityReviewCampaignController.GetIdentityReviewCampaignDocumentPreview)
	unauthenticatedRoutes.GET("/identity-access-review/review-campaigns/:id/document-download", identityReviewCampaignController.DownloadIdentityReviewCampaignDocument)
}
