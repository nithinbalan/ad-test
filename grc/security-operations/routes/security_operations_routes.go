package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	securityOperationsRoutes := router.Group("/")

	securityOperationsRoutes.Use(middlewares.AuthMiddleware(app.DB))

	db := app.DB

	// TODO: Pull controllers from app.SecurityOperationsModule here, e.g.:
	// someController := app.SecurityOperationsModule.SomeController

	// TODO: Add routes here, e.g.:
	// securityOperationsRoutes.GET("/security-operations/items", someController.GetAllItems)
	// securityOperationsRoutes.GET("/security-operations/items/:id", someController.GetItemDetails)
	// securityOperationsRoutes.POST("/security-operations/items", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), someController.CreateItem)
	// securityOperationsRoutes.PUT("/security-operations/items/:id", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), someController.UpdateItem)
	// securityOperationsRoutes.POST("/security-operations/items/delete", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), someController.DeleteItems)

	_ = db
}
