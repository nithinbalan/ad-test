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
	riskAIController := app.JethurAIModule.RiskAIController

	jethurAIRoutes.GET("/jethur-ai/documents/:id/summary", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.GetDocumentSummary)
	jethurAIRoutes.GET("/jethur-ai/documents/:id/recommendations", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.GetDocumentRecommendations)
	jethurAIRoutes.POST("/jethur-ai/document-chat", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.ChatWithDocument)

	jethurAIRoutes.POST("/jethur-ai/risk/autofill", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskAutofill)
	jethurAIRoutes.POST("/jethur-ai/risk/rca", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskRCA)
}
