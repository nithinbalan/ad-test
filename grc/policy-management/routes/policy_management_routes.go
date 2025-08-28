package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	policyManagementRoutes := router.Group("/")

	policyManagementRoutes.Use(middlewares.AuthMiddleware())

	// unauthenticatedRoutes := router.Group("/")

	db := app.DB

	policyRegisterController := app.PolicyManagementModule.PolicyRegisterController

	//policy management
	policyManagementRoutes.GET("/policy-register", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), policyRegisterController.GetAllDocuments)
	policyManagementRoutes.GET("/policy-register/:id", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), policyRegisterController.GetDocumentDetails)
	policyManagementRoutes.POST("/policy-register", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyRegisterController.CreateDocument)
	policyManagementRoutes.PUT("/policy-register/:id", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyRegisterController.UpdateDocument)
	policyManagementRoutes.POST("/policy-register/delete", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyRegisterController.DeleteDocuments)
}
