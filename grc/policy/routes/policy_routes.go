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

	policyRegisterController := app.PolicyModule.PolicyRegisterController

	//policy
	policyRoutes.GET("/policy-register", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), policyRegisterController.GetAllDocuments)
	policyRoutes.GET("/policy-register/:id", middlewares.RolePermissionMiddleware(db, "policy_register", "view"), policyRegisterController.GetDocumentDetails)
	policyRoutes.POST("/policy-register", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyRegisterController.CreateDocument)
	policyRoutes.PUT("/policy-register/:id", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyRegisterController.UpdateDocument)
	policyRoutes.POST("/policy-register/delete", middlewares.RolePermissionMiddleware(db, "policy_register", "modify"), policyRegisterController.DeleteDocuments)
}
