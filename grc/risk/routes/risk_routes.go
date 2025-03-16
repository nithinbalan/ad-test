package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	riskRoutes := router.Group("/")

	riskRoutes.Use(middlewares.AuthMiddleware())

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	riskCategoryController := app.RiskModule.RiskCategoryController
	potentialRiskStatusController := app.RiskModule.PotentialRiskStatusController
	riskLevelController := app.RiskModule.RiskLevelController
	potentialRiskController := app.RiskModule.PotentialRiskController
	potentialRiskStatusUpdateController := app.RiskModule.PotentialRiskStatusUpdateController
	kriCategoryController := app.RiskModule.KriCategoryController
	kriRegisterController := app.RiskModule.KriRegisterController
	kriValueUpdateController := app.RiskModule.KriValueUpdateController
	riskAppetiteStatementController := app.RiskModule.RiskAppetiteStatementController
	riskImpactCategoryController := app.RiskModule.RiskImpactCategoryController
	riskMatrixImpactController := app.RiskModule.RiskMatrixImpactController
	riskMatrixImpactGuidelineController := app.RiskModule.RiskMatrixImpactGuidelineController
	riskMatrixLikelihoodController := app.RiskModule.RiskMatrixLikelihoodController
	riskMatrixRiskRatingController := app.RiskModule.RiskMatrixRiskRatingController
	riskMatrixRiskRatingValueController := app.RiskModule.RiskMatrixRiskRatingValueController
	riskSourceController := app.RiskModule.RiskSourceController
	riskStatusController := app.RiskModule.RiskStatusController
	riskRegisterController := app.RiskModule.RiskRegisterController
	riskRootCauseAnalysisController := app.RiskModule.RiskRootCauseAnalysisController
	riskAppetiteToleranceLevelController := app.RiskModule.RiskAppetiteToleranceLevelController
	riskInherentEvaluationController := app.RiskModule.RiskInherentEvaluationController
	riskActionPlanController := app.RiskModule.RiskActionPlanController
	riskActionPlanStatusUpdateController := app.RiskModule.RiskActionPlanStatusUpdateController
	riskAssociatedItemController := app.RiskModule.RiskAssociatedItemController
	riskResponseStrategyController := app.RiskModule.RiskResponseStrategyController
	riskTreatmentGuidelineController := app.RiskModule.RiskTreatmentGuidelineController
	riskHeatmapController := app.RiskModule.RiskHeatmapController
	riskInherentControlAnalysisController := app.RiskModule.RiskInherentControlAnalysisController
	riskResidualEvaluationController := app.RiskModule.RiskResidualEvaluationController
	riskResidualControlAnalysisController := app.RiskModule.RiskResidualControlAnalysisController
	riskKriController := app.RiskModule.RiskKriController
	riskTimelineController := app.RiskModule.RiskTimelineController
	riskSettingsController := app.RiskModule.RiskSettingsController
	riskClosureController := app.RiskModule.RiskClosureController
	riskCommentController := app.RiskModule.RiskCommentController
	riskReportController := app.RiskModule.RiskReportController

	riskRoutes.GET("/risk-categories", riskCategoryController.GetAllRiskCategories)
	riskRoutes.GET("/risk-categories/:id", riskCategoryController.GetRiskCategoryDetails)
	riskRoutes.POST("/risk-categories", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskCategoryController.CreateRiskCategory)
	riskRoutes.PUT("/risk-categories/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskCategoryController.UpdateRiskCategory)
	riskRoutes.POST("/risk-categories/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskCategoryController.DeleteRiskCategories)
	riskRoutes.PUT("/risk-categories/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskCategoryController.ActivateRiskCategories)
	riskRoutes.PUT("/risk-categories/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskCategoryController.DeactivateRiskCategories)

	riskRoutes.GET("/potential-risk-statuses", potentialRiskStatusController.GetAllPotentialRiskStatuses)
	riskRoutes.GET("/potential-risk-statuses/:id", potentialRiskStatusController.GetPotentialRiskStatusDetails)

	riskRoutes.GET("/potential-risks/export", middlewares.RolePermissionMiddleware(db, "potential_risk", "view"), potentialRiskController.GetPotentialRiskExport)
	riskRoutes.GET("/potential-risks", middlewares.RolePermissionMiddleware(db, "potential_risk", "view"), potentialRiskController.GetAllPotentialRisks)
	riskRoutes.GET("/potential-risks/:id", middlewares.RolePermissionMiddleware(db, "potential_risk", "view"), potentialRiskController.GetPotentialRiskDetails)
	riskRoutes.POST("/potential-risks", middlewares.RolePermissionMiddleware(db, "potential_risk", "modify"), potentialRiskController.CreatePotentialRisk)
	riskRoutes.PUT("/potential-risks/:id", middlewares.RolePermissionMiddleware(db, "potential_risk", "modify"), potentialRiskController.UpdatePotentialRisk)
	riskRoutes.POST("/potential-risks/delete", middlewares.RolePermissionMiddleware(db, "potential_risk", "modify"), potentialRiskController.DeletePotentialRisks)

	riskRoutes.GET("/risk-levels", riskLevelController.GetAllRiskLevels)
	riskRoutes.GET("/risk-levels/:id", riskLevelController.GetRiskLevelDetails)

	riskRoutes.GET("/potential-risks/:id/status-updates", middlewares.RolePermissionMiddleware(db, "potential_risk", "view"), potentialRiskStatusUpdateController.GetAllPotentialRiskStatusUpdates)
	riskRoutes.POST("/potential-risks/:id/status-updates", middlewares.RolePermissionMiddleware(db, "potential_risk", "modify"), potentialRiskStatusUpdateController.CreatePotentialRiskStatusUpdate)

	riskRoutes.GET("/risk-appetite-statements", riskAppetiteStatementController.GetAllRiskAppetiteStatements)
	riskRoutes.GET("/risk-appetite-statements-by-category", middlewares.RolePermissionMiddleware(db, "risk_appetite_statement", "view"), riskAppetiteStatementController.GetAllRiskAppetiteStatementByCategory)
	riskRoutes.GET("/risk-appetite-statements/:id", riskAppetiteStatementController.GetRiskAppetiteStatementDetails)
	riskRoutes.POST("/risk-appetite-statements", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskAppetiteStatementController.CreateRiskAppetiteStatement)
	riskRoutes.PUT("/risk-appetite-statements/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskAppetiteStatementController.UpdateRiskAppetiteStatement)
	riskRoutes.POST("/risk-appetite-statements/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskAppetiteStatementController.DeleteRiskAppetiteStatements)
	riskRoutes.PUT("/risk-appetite-statements/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskAppetiteStatementController.ActivateRiskAppetiteStatements)
	riskRoutes.PUT("/risk-appetite-statements/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskAppetiteStatementController.DeactivateRiskAppetiteStatements)

	riskRoutes.GET("/kri-categories", kriCategoryController.GetAllKRICategories)
	riskRoutes.POST("/kri-categories", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), kriCategoryController.CreateKRICategory)
	riskRoutes.GET("/kri-categories/:id", kriCategoryController.GetKRICategoryDetails)
	riskRoutes.PUT("/kri-categories/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), kriCategoryController.UpdateKRICategory)
	riskRoutes.POST("/kri-categories/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), kriCategoryController.DeleteKRICategories)
	riskRoutes.PUT("/kri-categories/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), kriCategoryController.ActivateKRICategories)
	riskRoutes.PUT("/kri-categories/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), kriCategoryController.DeactivateKRICategories)

	riskRoutes.GET("/kri-register/export", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriRegisterController.GetKriRegisterExport)
	riskRoutes.GET("/kri-register", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriRegisterController.GetAllKRIRegisters)
	riskRoutes.GET("/kri-register/:id", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriRegisterController.GetKRIRegisterDetails)
	riskRoutes.POST("/kri-register", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterController.CreateKRIRegister)
	riskRoutes.PUT("/kri-register/:id", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterController.UpdateKRIRegister)
	riskRoutes.POST("/kri-register/delete", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterController.DeleteKRIRegisters)

	unauthenticatedRoutes.GET("/kri-register/:id/value-updates/:updateId/document-preview", kriValueUpdateController.GetKRIValueUpdateDocumentPreview)
	unauthenticatedRoutes.GET("/kri-register/:id/value-updates/:updateId/document-download", kriValueUpdateController.DownloadKRIValueUpdateDocument)
	riskRoutes.GET("/kri-register/:id/value-updates", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriValueUpdateController.GetAllKRIValueUpdates)
	riskRoutes.POST("/kri-register/:id/value-updates", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriValueUpdateController.CreateKRIValueUpdate)
	riskRoutes.GET("/kri-register/:id/value-updates/:updateId", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriValueUpdateController.GetKRIValueUpdateDetails)
	riskRoutes.PUT("/kri-register/:id/value-updates/:updateId", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriValueUpdateController.UpdateKRIValueUpdate)
	riskRoutes.POST("/kri-register/:id/value-updates/delete", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriValueUpdateController.DeleteKRIValueUpdates)

	riskRoutes.GET("/risk-matrix-impacts", riskMatrixImpactController.GetAllRiskMatrixImpacts)
	riskRoutes.GET("/risk-matrix-impacts/:id", riskMatrixImpactController.GetRiskMatrixImpactDetails)
	riskRoutes.POST("/risk-matrix-impacts", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixImpactController.CreateRiskMatrixImpact)
	riskRoutes.PUT("/risk-matrix-impacts/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixImpactController.UpdateRiskMatrixImpact)
	riskRoutes.POST("/risk-matrix-impacts/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixImpactController.DeleteRiskMatrixImpacts)
	riskRoutes.PUT("/risk-matrix-impacts/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixImpactController.ActivateRiskMatrixImpact)
	riskRoutes.PUT("/risk-matrix-impacts/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixImpactController.DeactivateRiskMatrixImpact)

	riskRoutes.GET("/risk-matrix-likelihood", riskMatrixLikelihoodController.GetAllRiskMatrixLikelihoods)
	riskRoutes.GET("/risk-matrix-likelihood/:id", riskMatrixLikelihoodController.GetRiskMatrixLikelihoodDetails)
	riskRoutes.POST("/risk-matrix-likelihood", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixLikelihoodController.CreateRiskMatrixLikelihood)
	riskRoutes.PUT("/risk-matrix-likelihood/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixLikelihoodController.UpdateRiskMatrixLikelihood)
	riskRoutes.POST("/risk-matrix-likelihood/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixLikelihoodController.DeleteRiskMatrixLikelihoods)
	riskRoutes.PUT("/risk-matrix-likelihood/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixLikelihoodController.ActivateRiskMatrixLikelihood)
	riskRoutes.PUT("/risk-matrix-likelihood/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixLikelihoodController.DeactivateRiskMatrixLikelihood)

	riskRoutes.GET("/risk-matrix-risk-ratings", riskMatrixRiskRatingController.GetAllRiskMatrixRiskRatings)
	riskRoutes.GET("/risk-matrix-risk-ratings/:id", riskMatrixRiskRatingController.GetRiskMatrixRiskRatingDetails)
	riskRoutes.POST("/risk-matrix-risk-ratings", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixRiskRatingController.CreateRiskMatrixRiskRating)
	riskRoutes.PUT("/risk-matrix-risk-ratings/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixRiskRatingController.UpdateRiskMatrixRiskRating)
	riskRoutes.POST("/risk-matrix-risk-ratings/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixRiskRatingController.DeleteRiskMatrixRiskRatings)
	riskRoutes.PUT("/risk-matrix-risk-ratings/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixRiskRatingController.ActivateRiskMatrixRiskRating)
	riskRoutes.PUT("/risk-matrix-risk-ratings/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixRiskRatingController.DeactivateRiskMatrixRiskRating)

	riskRoutes.GET("/risk-matrix-risk-rating-value", riskMatrixRiskRatingValueController.GetAllRiskMatrixRiskRatingValue)
	riskRoutes.POST("/risk-matrix-risk-rating-value", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixRiskRatingValueController.CreateRiskMatrixRiskRatingValue)
	riskRoutes.GET("/risk-matrix-configuration", middlewares.RolePermissionMiddleware(db, "risk_matrix", "view"), riskMatrixRiskRatingValueController.GetRiskMatrix)

	riskRoutes.GET("/risk-impact-categories", riskImpactCategoryController.GetAllRiskImpactCategories)
	riskRoutes.GET("/risk-impact-categories/:id", riskImpactCategoryController.GetRiskImpactCategoryDetails)
	riskRoutes.POST("/risk-impact-categories", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskImpactCategoryController.CreateRiskImpactCategory)
	riskRoutes.PUT("/risk-impact-categories/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskImpactCategoryController.UpdateRiskImpactCategory)
	riskRoutes.POST("/risk-impact-categories/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskImpactCategoryController.DeleteRiskImpactCategories)
	riskRoutes.PUT("/risk-impact-categories/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskImpactCategoryController.ActivateRiskImpactCategories)
	riskRoutes.PUT("/risk-impact-categories/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskImpactCategoryController.DeactivateRiskImpactCategories)

	riskRoutes.GET("/risk-matrix-impact-guidelines", riskMatrixImpactGuidelineController.GetAllRiskMatrixImpactGuidelines)
	riskRoutes.GET("/risk-matrix-impact-guidelines/:id", riskMatrixImpactGuidelineController.GetRiskMatrixImpactGuidelineDetails)
	riskRoutes.POST("/risk-matrix-impact-guidelines", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixImpactGuidelineController.CreateRiskMatrixImpactGuideline)
	riskRoutes.PUT("/risk-matrix-impact-guidelines/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixImpactGuidelineController.UpdateRiskMatrixImpactGuideline)
	riskRoutes.POST("/risk-matrix-impact-guidelines/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskMatrixImpactGuidelineController.DeleteRiskMatrixImpactGuidelines)

	// Risk Source
	riskRoutes.GET("/risk-sources", riskSourceController.GetAllRiskSources)
	riskRoutes.GET("/risk-sources/:id", riskSourceController.GetRiskSourceDetails)
	riskRoutes.POST("/risk-sources", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskSourceController.CreateRiskSource)
	riskRoutes.PUT("/risk-sources/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskSourceController.UpdateRiskSource)
	riskRoutes.POST("/risk-sources/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskSourceController.DeleteRiskSources)
	riskRoutes.PUT("/risk-sources/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskSourceController.ActivateRiskSources)
	riskRoutes.PUT("/risk-sources/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskSourceController.DeactivateRiskSources)

	// Risk Status
	riskRoutes.GET("/risk-statuses", riskStatusController.GetAllRiskStatuses)
	riskRoutes.GET("/risk-statuses/:id", riskStatusController.GetRiskStatusDetails)

	// Risk Register
	unauthenticatedRoutes.GET("/risk-register/:id/document-preview", riskRegisterController.GetRiskRegisterDocumentPreview)
	unauthenticatedRoutes.GET("/risk-register/:id/document-download", riskRegisterController.DownloadRiskRegisterDocument)
	riskRoutes.GET("/risk-register/export", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskRegisterController.GetRiskRegisterExport)
	riskRoutes.GET("/risk-register", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskRegisterController.GetAllRiskRegisters)
	riskRoutes.GET("/risk-register/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskRegisterController.GetRiskRegisterDetails)
	riskRoutes.POST("/risk-register", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.CreateRiskRegister)
	riskRoutes.PUT("/risk-register/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.UpdateRiskRegister)
	riskRoutes.POST("/risk-register/delete", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.DeleteRiskRegisters)

	riskRoutes.GET("risks/:riskId/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskRootCauseAnalysisController.GetAllRiskRootCauseAnalyses)
	riskRoutes.POST("risks/:riskId/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRootCauseAnalysisController.CreateRiskRootCauseAnalysis)
	riskRoutes.PUT("risks/:riskId/root-cause-analysis/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRootCauseAnalysisController.UpdateRiskRootCauseAnalysis)
	riskRoutes.DELETE("risks/:riskId/root-cause-analysis/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRootCauseAnalysisController.DeleteRiskRootCauseAnalysis)

	riskRoutes.GET("/risk-appetite-tolerance-levels", riskAppetiteToleranceLevelController.GetAllRiskAppetiteToleranceLevels)
	riskRoutes.GET("/risk-appetite-tolerance-levels/:id", riskAppetiteToleranceLevelController.GetRiskAppetiteToleranceLevelDetails)

	riskRoutes.GET("risks/:riskId/inherent-evaluation", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskInherentEvaluationController.GetRiskInherentEvaluation)
	riskRoutes.POST("risks/:riskId/inherent-evaluation", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskInherentEvaluationController.CreateRiskInherentEvaluation)
	riskRoutes.PUT("risks/:riskId/inherent-evaluation/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskInherentEvaluationController.UpdateRiskInherentEvaluation)
	riskRoutes.DELETE("risks/:riskId/inherent-evaluation/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskInherentEvaluationController.DeleteRiskInherentEvaluation)

	riskRoutes.GET("risks/:riskId/controls", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskControls)
	riskRoutes.GET("risks/:riskId/compliance-requirements", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskComplianceRequirements)
	riskRoutes.GET("risks/:riskId/processes", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskProcesses)
	riskRoutes.GET("risks/:riskId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskStrategicObjectives)
	riskRoutes.GET("risks/:riskId/slas", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskSLAs)
	riskRoutes.GET("risks/:riskId/incidents", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskIncidents)
	riskRoutes.GET("/risks/:riskId/documents", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskDocuments)
	riskRoutes.GET("/risks/:riskId/assets", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskAssets)
	riskRoutes.GET("/risks/:riskId/threats", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskThreats)

	riskRoutes.PUT("/risks/:riskId/response-strategy", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.UpdateResponseStrategy)
	riskRoutes.DELETE("/risks/:riskId/response-strategy", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.DeleteResponseStrategy)

	riskRoutes.GET("/risks/:riskId/control-performance", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskRegisterController.GetControlPerformance)

	// Risk Action Plan
	riskRoutes.GET("/risk-action-plans/export", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "view"), riskActionPlanController.GetRiskActionPlanExport)
	riskRoutes.GET("/risk-action-plans", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "view"), riskActionPlanController.GetAllRiskActionPlans)
	riskRoutes.GET("/risks/:riskId/action-plans", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "view"), riskActionPlanController.GetAllRiskActionPlans)
	riskRoutes.GET("/risks/:riskId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "view"), riskActionPlanController.GetRiskActionPlanDetails)
	riskRoutes.POST("/risks/:riskId/action-plans", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "modify"), riskActionPlanController.CreateRiskActionPlan)
	riskRoutes.PUT("/risks/:riskId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "modify"), riskActionPlanController.UpdateRiskActionPlan)
	riskRoutes.POST("/risks/:riskId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "modify"), riskActionPlanController.DeleteRiskActionPlans)

	// Document Preview and Download for Risk Action Plan Status Update
	unauthenticatedRoutes.GET("/risks/:riskId/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-preview", riskActionPlanStatusUpdateController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/risks/:riskId/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-download", riskActionPlanStatusUpdateController.DownloadRiskActionPlanDocument)

	riskRoutes.GET("/risks/:riskId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "view"), riskActionPlanStatusUpdateController.GetAllRiskActionPlanStatusUpdates)
	riskRoutes.POST("/risks/:riskId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "modify"), riskActionPlanStatusUpdateController.CreateRiskActionPlanStatusUpdate)
	riskRoutes.PUT("/risks/:riskId/action-plans/:actionPlanId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "modify"), riskActionPlanStatusUpdateController.UpdateRiskActionPlanStatusUpdate)
	riskRoutes.GET("/risks/:riskId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "view"), riskActionPlanStatusUpdateController.GetLatestRiskActionPlanStatusUpdate)

	// Risk Treatment Guides
	riskRoutes.GET("/risk-treatment-guidelines", riskTreatmentGuidelineController.GetAllRiskTreatmentGuidelines)
	riskRoutes.GET("/risk-treatment-guidelines/:id", riskTreatmentGuidelineController.GetRiskTreatmentGuidelineDetails)
	riskRoutes.POST("/risk-treatment-guidelines", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskTreatmentGuidelineController.CreateRiskTreatmentGuideline)
	riskRoutes.PUT("/risk-treatment-guidelines/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskTreatmentGuidelineController.UpdateRiskTreatmentGuideline)
	riskRoutes.POST("/risk-treatment-guidelines/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskTreatmentGuidelineController.DeleteRiskTreatmentGuidelines)
	riskRoutes.PUT("/risk-treatment-guidelines/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskTreatmentGuidelineController.ActivateRiskTreatmentGuidelines)
	riskRoutes.PUT("/risk-treatment-guidelines/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskTreatmentGuidelineController.DeactivateRiskTreatmentGuidelines)

	// Risk Response Strategy
	riskRoutes.GET("/risk-response-strategies", riskResponseStrategyController.GetAllRiskResponseStrategies)
	riskRoutes.GET("/risk-response-strategies/:id", riskResponseStrategyController.GetRiskResponseStrategyDetails)
	riskRoutes.POST("/risk-response-strategies", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskResponseStrategyController.CreateRiskResponseStrategy)
	riskRoutes.PUT("/risk-response-strategies/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskResponseStrategyController.UpdateRiskResponseStrategy)
	riskRoutes.POST("/risk-response-strategies/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskResponseStrategyController.DeleteRiskResponseStrategies)
	riskRoutes.PUT("/risk-response-strategies/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskResponseStrategyController.ActivateRiskResponseStrategies)
	riskRoutes.PUT("/risk-response-strategies/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskResponseStrategyController.DeactivateRiskResponseStrategies)

	// Risk Heat Map
	riskRoutes.GET("risks/heatmap-by-responsible-unit", middlewares.RolePermissionMiddleware(db, "risk_heatmap", "view"), riskHeatmapController.GetRiskHeatmapByResponsibleUnit)
	riskRoutes.GET("risks/heatmap-by-risk-category", middlewares.RolePermissionMiddleware(db, "risk_heatmap", "view"), riskHeatmapController.GetRiskHeatmapByRiskCategory)
	riskRoutes.GET("risks/heatmap-by-risk-source", middlewares.RolePermissionMiddleware(db, "risk_heatmap", "view"), riskHeatmapController.GetRiskHeatmapByRiskSource)
	riskRoutes.GET("risks/heatmap-by-risk-matrix", middlewares.RolePermissionMiddleware(db, "risk_heatmap", "view"), riskHeatmapController.GetRiskHeatmapByRiskMatrix)
	riskRoutes.GET("risks/heatmap-by-risk-response-strategy", middlewares.RolePermissionMiddleware(db, "risk_heatmap", "view"), riskHeatmapController.GetRiskHeatmapByResponseStrategy)
	riskRoutes.GET("risks/heatmap-by-risk-status", middlewares.RolePermissionMiddleware(db, "risk_heatmap", "view"), riskHeatmapController.GetRiskHeatmapByRiskStatus)
	riskRoutes.GET("risks/heatmap-by-review-frequency", middlewares.RolePermissionMiddleware(db, "risk_heatmap", "view"), riskHeatmapController.GetRiskHeatmapByReviewFrequency)

	riskRoutes.GET("processes/:processId/risk-heatmap", middlewares.RolePermissionMiddleware(db, "process", "view"), riskHeatmapController.GetProcessRiskHeatmap)
	riskRoutes.GET("controls/:controlId/risk-heatmap", middlewares.RolePermissionMiddleware(db, "controls", "view"), riskHeatmapController.GetControlRiskHeatmap)
	riskRoutes.GET("compliances/:complianceId/risk-heatmap", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), riskHeatmapController.GetComplianceRiskHeatmap)
	riskRoutes.GET("slas/:slaId/risk-heatmap", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), riskHeatmapController.GetSLARiskHeatmap)

	// Risk Control Analysis
	riskRoutes.GET("risks/:riskId/inherent-control-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskInherentControlAnalysisController.GetAllRiskControlAnalyses)
	riskRoutes.GET("risks/:riskId/inherent-control-analysis/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskInherentControlAnalysisController.GetRiskControlAnalysisDetails)
	riskRoutes.POST("risks/:riskId/inherent-control-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskInherentControlAnalysisController.CreateRiskControlAnalysis)
	riskRoutes.PUT("risks/:riskId/inherent-control-analysis/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskInherentControlAnalysisController.UpdateRiskControlAnalysis)
	riskRoutes.POST("risks/:riskId/inherent-control-analysis/delete", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskInherentControlAnalysisController.DeleteRiskControlAnalyses)

	riskRoutes.GET("risks/:riskId/inherent-control-analysis/control-summary", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskInherentControlAnalysisController.GetRiskControlSummary)

	// Risk Residual Evaluation
	riskRoutes.GET("risks/:riskId/residual-evaluation", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskResidualEvaluationController.GetRiskResidualEvaluation)
	riskRoutes.POST("risks/:riskId/residual-evaluation", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResidualEvaluationController.CreateRiskResidualEvaluation)
	riskRoutes.PUT("risks/:riskId/residual-evaluation/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResidualEvaluationController.UpdateRiskResidualEvaluation)
	riskRoutes.DELETE("risks/:riskId/residual-evaluation/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResidualEvaluationController.DeleteRiskResidualEvaluation)

	// Risk Control Analysis
	riskRoutes.GET("risks/:riskId/residual-control-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskResidualControlAnalysisController.GetAllRiskControlAnalyses)
	riskRoutes.GET("risks/:riskId/residual-control-analysis/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskResidualControlAnalysisController.GetRiskResidualControlAnalysisDetails)
	riskRoutes.POST("risks/:riskId/residual-control-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResidualControlAnalysisController.CreateRiskControlAnalysis)
	riskRoutes.PUT("risks/:riskId/residual-control-analysis/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResidualControlAnalysisController.UpdateRiskControlAnalysis)
	riskRoutes.POST("risks/:riskId/residual-control-analysis/delete", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResidualControlAnalysisController.DeleteRiskControlAnalyses)

	riskRoutes.GET("risks/:riskId/residual-control-analysis/control-summary", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskResidualControlAnalysisController.GetRiskControlSummary)

	riskRoutes.GET("/risks/:riskId/kris", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskKriController.GetAllRiskKris)
	riskRoutes.GET("/risks/:riskId/kris/:kriId", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskKriController.GetRiskKriDetails)
	riskRoutes.POST("/risks/:riskId/kris", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskKriController.CreateRiskKri)
	riskRoutes.PUT("/risks/:riskId/kris", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskKriController.UpdateRiskKri)
	riskRoutes.POST("/risks/:riskId/kris/delete", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskKriController.DeleteRiskKris)

	// Risk Timelines
	riskRoutes.GET("/risks/:riskId/timelines", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskTimelineController.GetAllRiskTimelines)

	// Risk Settings
	riskRoutes.GET("/risk-settings", riskSettingsController.GetAllRiskSettings)
	riskRoutes.PUT("/risk-settings", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskSettingsController.UpdateRiskSettings)

	unauthenticatedRoutes.GET("/risks/:riskId/closures/:id/document-preview", riskClosureController.GetRiskClosureDocumentPreview)
	unauthenticatedRoutes.GET("/risks/:riskId/closures/:id/document-download", riskClosureController.DownloadRiskClosureDocument)
	riskRoutes.GET("risks/:riskId/closures", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskClosureController.GetAllRiskClosures)
	riskRoutes.GET("risks/:riskId/closures/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskClosureController.GetRiskClosureDetails)
	riskRoutes.POST("risks/:riskId/closures", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskClosureController.CreateRiskClosure)
	riskRoutes.PUT("risks/:riskId/closures/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskClosureController.UpdateRiskClosure)
	riskRoutes.DELETE("risks/:riskId/closures/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskClosureController.DeleteRiskClosure)

	riskRoutes.GET("/risks/:riskId/comments", middlewares.RolePermissionMiddleware(db, "risk", "view"), riskCommentController.GetAllComments)
	riskRoutes.GET("/risks/:riskId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "risk", "view"), riskCommentController.GetCommentDetails)
	riskRoutes.POST("/risks/:riskId/comments", middlewares.RolePermissionMiddleware(db, "risk", "modify"), riskCommentController.CreateComment)
	riskRoutes.PUT("/risks/:riskId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "risk", "modify"), riskCommentController.UpdateComment)
	riskRoutes.POST("/risks/:riskId/comments/delete", middlewares.RolePermissionMiddleware(db, "risk", "modify"), riskCommentController.DeleteComment)

	riskRoutes.PUT("/risks/:riskId/submit", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.SubmitRisk)
	riskRoutes.PUT("/risks/:riskId/revert", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.RevertRisk)
	riskRoutes.PUT("/risks/:riskId/reject", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.RejectRisk)
	riskRoutes.PUT("/risks/:riskId/approve", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.ApproveRisk)
	riskRoutes.PUT("/risks/:riskId/resubmit", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.ResubmitRisk)

	riskRoutes.GET("/risk-reports", riskReportController.GetAllRiskReports)
	riskRoutes.GET("/risk-reports/:riskReportId", riskReportController.GetRiskReportDetails)
	riskRoutes.POST("/risk-reports", riskReportController.CreateRiskReport)
	riskRoutes.POST("/risk-reports/:riskReportId/:sectionId", riskReportController.UpdateRiskReportSection)
	// riskRoutes.POST("/risk-reports", riskReportController.DeleteRiskReportSection)
	// riskRoutes.POST("/risk-reports", riskReportController.DeleteRiskReport)

}
