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
	riskAssociatedItemController := app.RiskModule.RiskAssociatedItemController

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

	riskRoutes.GET("risks/:riskId/inherent-evaluation/:id", riskInherentEvaluationController.GetRiskInherentEvaluation)
	riskRoutes.POST("risks/:riskId/inherent-evaluation", riskInherentEvaluationController.CreateRiskInherentEvaluation)
	riskRoutes.PUT("risks/:riskId/inherent-evaluation/:id", riskInherentEvaluationController.UpdateRiskInherentEvaluation)
	riskRoutes.DELETE("risks/:riskId/inherent-evaluation/:id", riskInherentEvaluationController.DeleteRiskInherentEvaluation)

	riskRoutes.GET("risks/:riskId/controls", riskAssociatedItemController.GetAllRiskControls)

}
