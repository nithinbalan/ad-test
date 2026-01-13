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
	complianceAIController := app.JethurAIModule.ComplianceAIController
	policyAIController :=app.JethurAIModule.PolicyAIController

	jethurAIRoutes.GET("/jethur-ai/documents/:id/summary", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.GetDocumentSummary)
	jethurAIRoutes.GET("/jethur-ai/documents/:id/recommendations", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.GetDocumentRecommendations)
	jethurAIRoutes.POST("/jethur-ai/document-chat", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.ChatWithDocument)
	jethurAIRoutes.GET("/jethur-ai/documents/:id/requirements", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAIController.GetDocumentRequirements)

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
	jethurAIRoutes.POST("/jethur-ai/risk-chat", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.ChatWithRisk)
	jethurAIRoutes.GET("/jethur-ai/risk-chat/:contextId", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskChatsByContext)
	jethurAIRoutes.GET("/jethur-ai/risks/duplicate", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetDuplicateRisks)
	jethurAIRoutes.GET("/jethur-ai/risks/:riskId/closure/autofill", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAIController.GetRiskClosureAutofill)

	jethurAIRoutes.POST("/jethur-ai/compliances/autofill", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceAutofill)
	jethurAIRoutes.POST("/jethur-ai/compliance/compliance-register-import", middlewares.RolePermissionMiddleware(db, "compliance", "create"), complianceAIController.ImportComplianceRegister)

	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/compliance-summary", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceSummary)
	jethurAIRoutes.GET("/jethur-ai/compliances/clauses/:clauseId/clause-guidance", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceClauseGuidance)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceStrategicObjectivesJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/suggest-kpi", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceSuggestKPIJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/suggest-process", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceSuggestProcessJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/suggest-risks", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceSuggestRisksJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/suggest-controls", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceSuggestControlsJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/suggest-associates", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceSuggestAssociatesJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/job-status/:jobId", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceJobStatus)

	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/clauses/:clauseId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceClauseStrategicObjectivesJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/clauses/:clauseId/suggest-kpi", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceClauseSuggestKPIJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/clauses/:clauseId/suggest-process", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceClauseSuggestProcessJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/clauses/:clauseId/suggest-risks", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceClauseSuggestRisksJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/clauses/:clauseId/suggest-controls", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceClauseSuggestControlsJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/clauses/:clauseId/suggest-incidents", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceClauseSuggestIncidentsJob)
	jethurAIRoutes.GET("/jethur-ai/compliances/:complianceId/clauses/:clauseId/suggest-associates", middlewares.RolePermissionMiddleware(db, "compliance", "view"), complianceAIController.GetComplianceClauseSuggestAssociatesJob)

	jethurAIRoutes.POST("/jethur-ai/policy/autofill", middlewares.RolePermissionMiddleware(db, "policy", "view"), policyAIController.GetPolicyAutofill)
}
