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
}
