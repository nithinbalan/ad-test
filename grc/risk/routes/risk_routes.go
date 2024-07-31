package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	riskRoutes := router.Group("/")

	riskRoutes.Use(middlewares.AuthMiddleware())

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

	riskRoutes.GET("/risk-categories", riskCategoryController.GetAllRiskCategories)
	riskRoutes.GET("/risk-categories/:id", riskCategoryController.GetRiskCategoryDetails)
	riskRoutes.POST("/risk-categories", riskCategoryController.CreateRiskCategory)
	riskRoutes.PUT("/risk-categories/:id", riskCategoryController.UpdateRiskCategory)
	riskRoutes.POST("/risk-categories/delete", riskCategoryController.DeleteRiskCategories)
	riskRoutes.PUT("/risk-categories/activate", riskCategoryController.ActivateRiskCategories)
	riskRoutes.PUT("/risk-categories/deactivate", riskCategoryController.DeactivateRiskCategories)

	riskRoutes.GET("/potential-risk-statuses", potentialRiskStatusController.GetAllPotentialRiskStatuses)
	riskRoutes.GET("/potential-risk-statuses/:id", potentialRiskStatusController.GetPotentialRiskStatusDetails)

	riskRoutes.GET("/potential-risks", potentialRiskController.GetAllPotentialRisks)
	riskRoutes.GET("/potential-risks/:id", potentialRiskController.GetPotentialRiskDetails)
	riskRoutes.POST("/potential-risks", potentialRiskController.CreatePotentialRisk)
	riskRoutes.PUT("/potential-risks/:id", potentialRiskController.UpdatePotentialRisk)
	riskRoutes.POST("/potential-risks/delete", potentialRiskController.DeletePotentialRisks)

	riskRoutes.GET("/risk-levels", riskLevelController.GetAllRiskLevels)
	riskRoutes.GET("/risk-levels/:id", riskLevelController.GetRiskLevelDetails)

	riskRoutes.GET("/potential-risks/:id/status-updates", potentialRiskStatusUpdateController.GetAllPotentialRiskStatusUpdates)
	riskRoutes.POST("/potential-risks/:id/status-updates", potentialRiskStatusUpdateController.CreatePotentialRiskStatusUpdate)

	riskRoutes.GET("/risk-appetite-statements", riskAppetiteStatementController.GetAllRiskAppetiteStatements)
	riskRoutes.GET("/risk-appetite-statements/:id", riskAppetiteStatementController.GetRiskAppetiteStatementDetails)
	riskRoutes.POST("/risk-appetite-statements", riskAppetiteStatementController.CreateRiskAppetiteStatement)
	riskRoutes.PUT("/risk-appetite-statements/:id", riskAppetiteStatementController.UpdateRiskAppetiteStatement)
	riskRoutes.POST("/risk-appetite-statements/delete", riskAppetiteStatementController.DeleteRiskAppetiteStatements)
	riskRoutes.PUT("/risk-appetite-statements/activate", riskAppetiteStatementController.ActivateRiskAppetiteStatements)
	riskRoutes.PUT("/risk-appetite-statements/deactivate", riskAppetiteStatementController.DeactivateRiskAppetiteStatements)

	riskRoutes.GET("/kri-categories", kriCategoryController.GetAllKRICategories)
	riskRoutes.POST("/kri-categories", kriCategoryController.CreateKRICategory)
	riskRoutes.GET("/kri-categories/:id", kriCategoryController.GetKRICategoryDetails)
	riskRoutes.PUT("/kri-categories/:id", kriCategoryController.UpdateKRICategory)
	riskRoutes.POST("/kri-categories/delete", kriCategoryController.DeleteKRICategories)
	riskRoutes.PUT("/kri-categories/activate", kriCategoryController.ActivateKRICategories)
	riskRoutes.PUT("/kri-categories/deactivate", kriCategoryController.DeactivateKRICategories)

	riskRoutes.GET("/kri-register", kriRegisterController.GetAllKRIRegisters)
	riskRoutes.GET("/kri-register/:id", kriRegisterController.GetKRIRegisterDetails)
	riskRoutes.POST("/kri-register", kriRegisterController.CreateKRIRegister)
	riskRoutes.PUT("/kri-register/:id", kriRegisterController.UpdateKRIRegister)
	riskRoutes.POST("/kri-register/delete", kriRegisterController.DeleteKRIRegisters)

	riskRoutes.GET("/kri-register/:id/value-updates", kriValueUpdateController.GetAllKRIValueUpdates)
	riskRoutes.POST("/kri-register/:id/value-updates", kriValueUpdateController.CreateKRIValueUpdate)
	riskRoutes.GET("/kri-register/:id/value-updates/:updateId", kriValueUpdateController.GetKRIValueUpdateDetails)
	riskRoutes.PUT("/kri-register/:id/value-updates/:updateId", kriValueUpdateController.UpdateKRIValueUpdate)
	riskRoutes.POST("/kri-register/:id/value-updates/delete", kriValueUpdateController.DeleteKRIValueUpdates)

	riskRoutes.GET("/risk-matrix-impacts", riskMatrixImpactController.GetAllRiskMatrixImpacts)
	riskRoutes.GET("/risk-matrix-impacts/:id", riskMatrixImpactController.GetRiskMatrixImpactDetails)
	riskRoutes.POST("/risk-matrix-impacts", riskMatrixImpactController.CreateRiskMatrixImpact)
	riskRoutes.PUT("/risk-matrix-impacts/:id", riskMatrixImpactController.UpdateRiskMatrixImpact)
	riskRoutes.POST("/risk-matrix-impacts/delete", riskMatrixImpactController.DeleteRiskMatrixImpacts)
	riskRoutes.PUT("/risk-matrix-impacts/activate", riskMatrixImpactController.ActivateRiskMatrixImpact)
	riskRoutes.PUT("/risk-matrix-impacts/deactivate", riskMatrixImpactController.DeactivateRiskMatrixImpact)

	riskRoutes.GET("/risk-matrix-likelihood", riskMatrixLikelihoodController.GetAllRiskMatrixLikelihoods)
	riskRoutes.GET("/risk-matrix-likelihood/:id", riskMatrixLikelihoodController.GetRiskMatrixLikelihoodDetails)
	riskRoutes.POST("/risk-matrix-likelihood", riskMatrixLikelihoodController.CreateRiskMatrixLikelihood)
	riskRoutes.PUT("/risk-matrix-likelihood/:id", riskMatrixLikelihoodController.UpdateRiskMatrixLikelihood)
	riskRoutes.POST("/risk-matrix-likelihood/delete", riskMatrixLikelihoodController.DeleteRiskMatrixLikelihoods)
	riskRoutes.PUT("/risk-matrix-likelihood/activate", riskMatrixLikelihoodController.ActivateRiskMatrixLikelihood)
	riskRoutes.PUT("/risk-matrix-likelihood/deactivate", riskMatrixLikelihoodController.DeactivateRiskMatrixLikelihood)

	riskRoutes.GET("/risk-matrix-risk-ratings", riskMatrixRiskRatingController.GetAllRiskMatrixRiskRatings)
	riskRoutes.GET("/risk-matrix-risk-ratings/:id", riskMatrixRiskRatingController.GetRiskMatrixRiskRatingDetails)
	riskRoutes.POST("/risk-matrix-risk-ratings", riskMatrixRiskRatingController.CreateRiskMatrixRiskRating)
	riskRoutes.PUT("/risk-matrix-risk-ratings/:id", riskMatrixRiskRatingController.UpdateRiskMatrixRiskRating)
	riskRoutes.POST("/risk-matrix-risk-ratings/delete", riskMatrixRiskRatingController.DeleteRiskMatrixRiskRatings)
	riskRoutes.PUT("/risk-matrix-risk-ratings/activate", riskMatrixRiskRatingController.ActivateRiskMatrixRiskRating)
	riskRoutes.PUT("/risk-matrix-risk-ratings/deactivate", riskMatrixRiskRatingController.DeactivateRiskMatrixRiskRating)

	riskRoutes.GET("/risk-matrix-risk-rating-value", riskMatrixRiskRatingValueController.GetAllRiskMatrixRiskRatingValue)
	riskRoutes.POST("/risk-matrix-risk-rating-value", riskMatrixRiskRatingValueController.CreateRiskMatrixRiskRatingValue)
	riskRoutes.GET("/risk-matrix-configuration", riskMatrixRiskRatingValueController.GetRiskMatrix)

	riskRoutes.GET("/risk-impact-categories", riskImpactCategoryController.GetAllRiskImpactCategories)
	riskRoutes.GET("/risk-impact-categories/:id", riskImpactCategoryController.GetRiskImpactCategoryDetails)
	riskRoutes.POST("/risk-impact-categories", riskImpactCategoryController.CreateRiskImpactCategory)
	riskRoutes.PUT("/risk-impact-categories/:id", riskImpactCategoryController.UpdateRiskImpactCategory)
	riskRoutes.POST("/risk-impact-categories/delete", riskImpactCategoryController.DeleteRiskImpactCategories)
	riskRoutes.PUT("/risk-impact-categories/activate", riskImpactCategoryController.ActivateRiskImpactCategories)
	riskRoutes.PUT("/risk-impact-categories/deactivate", riskImpactCategoryController.DeactivateRiskImpactCategories)

	riskRoutes.GET("/risk-matrix-impact-guidelines", riskMatrixImpactGuidelineController.GetAllRiskMatrixImpactGuidelines)
	riskRoutes.GET("/risk-matrix-impact-guidelines/:id", riskMatrixImpactGuidelineController.GetRiskMatrixImpactGuidelineDetails)
	riskRoutes.POST("/risk-matrix-impact-guidelines", riskMatrixImpactGuidelineController.CreateRiskMatrixImpactGuideline)
	riskRoutes.PUT("/risk-matrix-impact-guidelines/:id", riskMatrixImpactGuidelineController.UpdateRiskMatrixImpactGuideline)
	riskRoutes.POST("/risk-matrix-impact-guidelines/delete", riskMatrixImpactGuidelineController.DeleteRiskMatrixImpactGuidelines)

	// Risk Source
	riskRoutes.GET("/risk-sources", riskSourceController.GetAllRiskSources)
	riskRoutes.GET("/risk-sources/:id", riskSourceController.GetRiskSourceDetails)
	riskRoutes.POST("/risk-sources", riskSourceController.CreateRiskSource)
	riskRoutes.PUT("/risk-sources/:id", riskSourceController.UpdateRiskSource)
	riskRoutes.POST("/risk-sources/delete", riskSourceController.DeleteRiskSources)
	riskRoutes.PUT("/risk-sources/activate", riskSourceController.ActivateRiskSources)
	riskRoutes.PUT("/risk-sources/deactivate", riskSourceController.DeactivateRiskSources)

	// Risk Status
	riskRoutes.GET("/risk-statuses", riskStatusController.GetAllRiskStatuses)
	riskRoutes.GET("/risk-statuses/:id", riskStatusController.GetRiskStatusDetails)

	// Risk Register
	riskRoutes.GET("/risk-register", riskRegisterController.GetAllRiskRegisters)
	riskRoutes.GET("/risk-register/:id", riskRegisterController.GetRiskRegisterDetails)
	riskRoutes.POST("/risk-register", riskRegisterController.CreateRiskRegister)
	riskRoutes.PUT("/risk-register/:id", riskRegisterController.UpdateRiskRegister)
	riskRoutes.POST("/risk-register/delete", riskRegisterController.DeleteRiskRegisters)

	riskRoutes.GET("risks/:riskId/root-cause-analysis", riskRootCauseAnalysisController.GetAllRiskRootCauseAnalyses)
	riskRoutes.POST("risks/:riskId/root-cause-analysis", riskRootCauseAnalysisController.CreateRiskRootCauseAnalysis)
	riskRoutes.PUT("risks/:riskId/root-cause-analysis/:id", riskRootCauseAnalysisController.UpdateRiskRootCauseAnalysis)
	riskRoutes.DELETE("risks/:riskId/root-cause-analysis/:id", riskRootCauseAnalysisController.DeleteRiskRootCauseAnalysis)

	riskRoutes.GET("/risk-appetite-tolerance-levels", riskAppetiteToleranceLevelController.GetAllRiskAppetiteToleranceLevels)
	riskRoutes.GET("/risk-appetite-tolerance-levels/:id", riskAppetiteToleranceLevelController.GetRiskAppetiteToleranceLevelDetails)

	riskRoutes.GET("risks/:riskId/inherent-evaluation", riskInherentEvaluationController.GetRiskInherentEvaluation)
	riskRoutes.POST("risks/:riskId/inherent-evaluation", riskInherentEvaluationController.CreateRiskInherentEvaluation)
	riskRoutes.PUT("risks/:riskId/inherent-evaluation/:id", riskInherentEvaluationController.UpdateRiskInherentEvaluation)
	riskRoutes.DELETE("risks/:riskId/inherent-evaluation/:id", riskInherentEvaluationController.DeleteRiskInherentEvaluation)

	riskRoutes.GET("risks/:riskId/controls", riskAssociatedItemController.GetAllRiskControls)
	riskRoutes.GET("risks/:riskId/compliance-requirements", riskAssociatedItemController.GetAllRiskComplianceRequirements)
	riskRoutes.GET("risks/:riskId/processes", riskAssociatedItemController.GetAllRiskProcesses)
	riskRoutes.GET("risks/:riskId/strategic-objectives", riskAssociatedItemController.GetAllRiskStrategicObjectives)
	riskRoutes.GET("risks/:riskId/slas", riskAssociatedItemController.GetAllRiskSLAs)
	riskRoutes.GET("risks/:riskId/incidents", riskAssociatedItemController.GetAllRiskIncidents)

	riskRoutes.PUT("/risks/:riskId/response-strategy", riskRegisterController.UpdateResponseStrategy)
	// Risk Action Plan
	riskRoutes.GET("/risk-action-plans", riskActionPlanController.GetAllRiskActionPlans)
	riskRoutes.GET("/risks/:riskId/action-plans", riskActionPlanController.GetAllRiskActionPlans)
	riskRoutes.GET("/risks/:riskId/action-plans/:actionPlanId", riskActionPlanController.GetRiskActionPlanDetails)
	riskRoutes.POST("/risks/:riskId/action-plans", riskActionPlanController.CreateRiskActionPlan)
	riskRoutes.PUT("/risks/:riskId/action-plans/:actionPlanId", riskActionPlanController.UpdateRiskActionPlan)
	riskRoutes.POST("/risks/:riskId/action-plans/delete", riskActionPlanController.DeleteRiskActionPlans)

	riskRoutes.GET("/risks/:riskId/action-plans/:actionPlanId/status-updates", riskActionPlanStatusUpdateController.GetAllRiskActionPlanStatusUpdates)
	riskRoutes.POST("/risks/:riskId/action-plans/:actionPlanId/status-updates", riskActionPlanStatusUpdateController.CreateRiskActionPlanStatusUpdate)
	riskRoutes.PUT("/risks/:riskId/action-plans/:actionPlanId/status-updates/:id", riskActionPlanStatusUpdateController.UpdateRiskActionPlanStatusUpdate)
	riskRoutes.GET("/risks/:riskId/action-plans/:actionPlanId/status-updates/latest", riskActionPlanStatusUpdateController.GetLatestRiskActionPlanStatusUpdate)


	// Risk Treatment Guides
	riskRoutes.GET("/risk-treatment-guidelines", riskTreatmentGuidelineController.GetAllRiskTreatmentGuidelines)
	riskRoutes.GET("/risk-treatment-guidelines/:id", riskTreatmentGuidelineController.GetRiskTreatmentGuidelineDetails)
	riskRoutes.POST("/risk-treatment-guidelines", riskTreatmentGuidelineController.CreateRiskTreatmentGuideline)
	riskRoutes.PUT("/risk-treatment-guidelines/:id", riskTreatmentGuidelineController.UpdateRiskTreatmentGuideline)
	riskRoutes.POST("/risk-treatment-guidelines/delete", riskTreatmentGuidelineController.DeleteRiskTreatmentGuidelines)
	riskRoutes.PUT("/risk-treatment-guidelines/activate", riskTreatmentGuidelineController.ActivateRiskTreatmentGuidelines)
	riskRoutes.PUT("/risk-treatment-guidelines/deactivate", riskTreatmentGuidelineController.DeactivateRiskTreatmentGuidelines)

	// Risk Response Strategy
	riskRoutes.GET("/risk-response-strategies", riskResponseStrategyController.GetAllRiskResponseStrategies)
	riskRoutes.GET("/risk-response-strategies/:id", riskResponseStrategyController.GetRiskResponseStrategyDetails)
	riskRoutes.POST("/risk-response-strategies", riskResponseStrategyController.CreateRiskResponseStrategy)
	riskRoutes.PUT("/risk-response-strategies/:id", riskResponseStrategyController.UpdateRiskResponseStrategy)
	riskRoutes.POST("/risk-response-strategies/delete", riskResponseStrategyController.DeleteRiskResponseStrategies)
	riskRoutes.PUT("/risk-response-strategies/activate", riskResponseStrategyController.ActivateRiskResponseStrategies)
	riskRoutes.PUT("/risk-response-strategies/deactivate", riskResponseStrategyController.DeactivateRiskResponseStrategies)

	// Risk Heat Map
	riskRoutes.GET("risks/heatmap-by-owner-unit", riskHeatmapController.GetRiskHeatmapByResponsibleUnit)
	riskRoutes.GET("processes/:processId/risk-heatmap", riskHeatmapController.GetProcessRiskHeatmap)
	riskRoutes.GET("controls/:controlId/risk-heatmap", riskHeatmapController.GetControlRiskHeatmap)
	riskRoutes.GET("compliances/:complianceId/risk-heatmap", riskHeatmapController.GetComplianceRiskHeatmap)

	//Risk Control Analysis
	riskRoutes.GET("risks/:riskId/inherent-control-analysis", riskInherentControlAnalysisController.GetAllRiskControlAnalyses)
	riskRoutes.GET("risks/:riskId/inherent-control-analysis/:id", riskInherentControlAnalysisController.GetRiskControlAnalysisDetails)
	riskRoutes.POST("risks/:riskId/inherent-control-analysis", riskInherentControlAnalysisController.CreateRiskControlAnalysis)
	riskRoutes.PUT("risks/:riskId/inherent-control-analysis/:id", riskInherentControlAnalysisController.UpdateRiskControlAnalysis)
	riskRoutes.POST("risks/:riskId/inherent-control-analysis/delete", riskInherentControlAnalysisController.DeleteRiskControlAnalyses)

	riskRoutes.GET("risks/:riskId/inherent-control-analysis/control-summary", riskInherentControlAnalysisController.GetRiskControlSummary)

	//Risk Residual Evaluation
	riskRoutes.GET("risks/:riskId/residual-evaluation", riskResidualEvaluationController.GetRiskResidualEvaluation)
	riskRoutes.POST("risks/:riskId/residual-evaluation", riskResidualEvaluationController.CreateRiskResidualEvaluation)
	riskRoutes.PUT("risks/:riskId/residual-evaluation/:id", riskResidualEvaluationController.UpdateRiskResidualEvaluation)
	riskRoutes.DELETE("risks/:riskId/residual-evaluation/:id", riskResidualEvaluationController.DeleteRiskResidualEvaluation)

	//Risk Control Analysis
	riskRoutes.GET("risks/:riskId/residual-control-analysis", riskResidualControlAnalysisController.GetAllRiskControlAnalyses)
	riskRoutes.GET("risks/:riskId/residual-control-analysis/:id", riskResidualControlAnalysisController.GetRiskResidualControlAnalysisDetails)
	riskRoutes.POST("risks/:riskId/residual-control-analysis", riskResidualControlAnalysisController.CreateRiskControlAnalysis)
	riskRoutes.PUT("risks/:riskId/residual-control-analysis/:id", riskResidualControlAnalysisController.UpdateRiskControlAnalysis)
	riskRoutes.POST("risks/:riskId/residual-control-analysis/delete", riskResidualControlAnalysisController.DeleteRiskControlAnalyses)

	riskRoutes.GET("risks/:riskId/residual-control-analysis/control-summary", riskResidualControlAnalysisController.GetRiskControlSummary)

	riskRoutes.GET("/risks/:riskId/kris", riskKriController.GetAllRiskKris)
	riskRoutes.GET("/risks/:riskId/kris/:kriId", riskKriController.GetRiskKriDetails)
	riskRoutes.POST("/risks/:riskId/kris", riskKriController.CreateRiskKri)
	riskRoutes.PUT("/risks/:riskId/kris", riskKriController.UpdateRiskKri)
	riskRoutes.POST("/risks/:riskId/kris/delete", riskKriController.DeleteRiskKris)

}
