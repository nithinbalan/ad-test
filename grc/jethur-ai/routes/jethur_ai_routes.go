package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	jethurAIRoutes := router.Group("/")
	jethurAIRoutes.Use(middlewares.AuthMiddleware())

	db := app.DB
	documentAIController := app.JethurAIModule.DocumentAIController

	jethurAIRoutes.GET("/jethur-ai/documents/:id/summary", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.GetDocumentSummary)
	jethurAIRoutes.GET("/jethur-ai/documents/:id/recommendations", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.GetDocumentRecommendations)
	jethurAIRoutes.POST("/jethur-ai/document-chat", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.ChatWithDocument)
}
