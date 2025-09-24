package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	policyRoutes := router.Group("/")

	policyRoutes.Use(middlewares.AuthMiddleware())

	// unauthenticatedRoutes := router.Group("/")

	db := app.DB

	policyFolderController := app.PolicyModule.PolicyFolderController
	policyRegisterController := app.PolicyModule.PolicyRegisterController
	policyClauseController := app.PolicyModule.PolicyClauseController
	policyReviewController := app.PolicyModule.PolicyReviewController
	policyWorkflowController := app.PolicyModule.PolicyWorkflowController
	PolicyClauseCommentController := app.PolicyModule.PolicyClauseCommentController

	//policy folder
	policyRoutes.GET("/policy-folders", policyFolderController.GetAllPolicyFolders)
	policyRoutes.GET("/policy-folders/:id", policyFolderController.GetPolicyFolderDetails)
	policyRoutes.POST("/policy-folders", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyFolderController.CreatePolicyFolder)
	policyRoutes.PUT("/policy-folders/:id", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyFolderController.UpdatePolicyFolder)
	policyRoutes.POST("/policy-folders/delete", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyFolderController.DeletePolicyFolders)

	//policy
	policyRoutes.GET("/policy-register", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), policyRegisterController.GetAllPolicies)
	policyRoutes.GET("/policy-register/:id", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), policyRegisterController.GetPolicyDetails)
	policyRoutes.POST("/policy-register", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyRegisterController.CreatePolicy)
	policyRoutes.PUT("/policy-register/:id", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyRegisterController.UpdatePolicy)
	policyRoutes.POST("/policy-register/delete", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyRegisterController.DeletePolicies)
	policyRoutes.POST("/policy-register/:id/submit", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), middlewares.DocumentAccessCheckingMiddleware(db), policyRegisterController.SubmitPolicy)
	policyRoutes.POST("/policy-register/:id/publish", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), middlewares.DocumentAccessCheckingMiddleware(db), policyRegisterController.PublishPolicy)

	// clause
	policyRoutes.GET("/policy-register/:id/clauses", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), middlewares.DocumentAccessCheckingMiddleware(db), policyClauseController.GetAllPolicyClauses)
	policyRoutes.GET("/policy-register/:id/clauses/:clauseId", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), middlewares.DocumentAccessCheckingMiddleware(db), policyClauseController.GetPolicyClauseDetails)
	policyRoutes.POST("/policy-register/:id/clauses", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), middlewares.DocumentAccessCheckingMiddleware(db), policyClauseController.CreatePolicyClause)
	policyRoutes.PUT("/policy-register/:id/clauses/:clauseId", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), middlewares.DocumentAccessCheckingMiddleware(db), policyClauseController.UpdatePolicyClause)
	policyRoutes.POST("/policy-register/:id/clauses/delete", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), middlewares.DocumentAccessCheckingMiddleware(db), policyClauseController.DeletePolicyClauses)
	policyRoutes.GET("/policy-register/:id/clauses/:clauseId/history", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), middlewares.DocumentAccessCheckingMiddleware(db), policyClauseController.GetClauseHistories)

	// workflow
	policyRoutes.GET("/policy-register/:id/comments", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), middlewares.DocumentAccessCheckingMiddleware(db), policyReviewController.GetComments)
	policyRoutes.POST("/policy-register/:id/comments", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), middlewares.DocumentAccessCheckingMiddleware(db), policyReviewController.AddComment)
	policyRoutes.GET("/policy-register/:id/workflow-list", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), middlewares.DocumentAccessCheckingMiddleware(db), policyWorkflowController.GetWorkflowList)

	// review comment
	policyRoutes.GET("/policy-register/:id/all-comments", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), PolicyClauseCommentController.GetAllDocumentComments)
	policyRoutes.GET("/policy-register/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), middlewares.DocumentAccessCheckingMiddleware(db), PolicyClauseCommentController.GetCommentDetails)
	policyRoutes.POST("/policy-register/:id/review-comments", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), middlewares.DocumentAccessCheckingMiddleware(db), PolicyClauseCommentController.CreateComment)
	policyRoutes.PUT("/policy-register/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), middlewares.DocumentAccessCheckingMiddleware(db), PolicyClauseCommentController.UpdateComment)
	policyRoutes.POST("/policy-register/:id/review-comments/delete", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), middlewares.DocumentAccessCheckingMiddleware(db), PolicyClauseCommentController.DeleteComment)
	policyRoutes.PUT("/policy-register/:id/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), middlewares.DocumentAccessCheckingMiddleware(db), PolicyClauseCommentController.ResolveComment)
}
