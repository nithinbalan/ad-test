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
	commonAIController := app.JethurAIModule.CommonAIController
	riskAIController := app.JethurAIModule.RiskAIController

	jethurAIRoutes.GET("/jethur-ai/documents/:id/summary", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.GetDocumentSummary)
	jethurAIRoutes.GET("/jethur-ai/documents/:id/recommendations", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.GetDocumentRecommendations)
	jethurAIRoutes.POST("/jethur-ai/document-chat", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.ChatWithDocument)

	jethurAIRoutes.POST("/jethur-ai/description-rewriter", commonAIController.RewriteDescription)
	jethurAIRoutes.POST("/jethur-ai/language-translation", commonAIController.TranslateLanguage)

	jethurAIRoutes.POST("/jethur-ai/risk/autofill", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskAutofill)
	jethurAIRoutes.POST("/jethur-ai/risk/rca", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskRCA)
	jethurAIRoutes.POST("/jethur-ai/risk/analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskAnalysis)
	jethurAIRoutes.POST("/jethur-ai/risk/response-strategy-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskResponseStrategyAnalysis)
	jethurAIRoutes.GET("/jethur-ai/risk/:riskId/control-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskControlAnalysis)
	jethurAIRoutes.GET("/jethur-ai/risk/:riskId/risk-kri", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskKriAnalysis)
	jethurAIRoutes.GET("/jethur-ai/risk/:riskId/risk-residual-control-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskResidualControlAnalysis)
	jethurAIRoutes.GET("/jethur-ai/risk/:riskId/residual-risk-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskResidualRiskAnalysis)

}
