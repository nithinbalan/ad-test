package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	riskRoutes := router.Group("/")

	riskRoutes.Use(middlewares.AuthMiddleware(app.DB))

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	riskCategoryController := app.RiskModule.RiskCategoryController
	potentialRiskStatusController := app.RiskModule.PotentialRiskStatusController
	riskLevelController := app.RiskModule.RiskLevelController
	riskSubCategoryController := app.RiskModule.RiskSubCategoryController
	levelOfMeasurementController := app.RiskModule.LevelOfMeasurementController
	distributionStatusController := app.RiskModule.DistributionStatusController
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
	riskRegisterUpdatePlanStatusController := app.RiskModule.RiskRegisterUpdatePlanStatusController
	riskRegisterController := app.RiskModule.RiskRegisterController
	riskRootCauseAnalysisController := app.RiskModule.RiskRootCauseAnalysisController
	riskAppetiteToleranceLevelController := app.RiskModule.RiskAppetiteToleranceLevelController
	riskInherentEvaluationController := app.RiskModule.RiskInherentEvaluationController
	riskActionPlanController := app.RiskModule.RiskActionPlanController
	riskActionPlanStatusUpdateController := app.RiskModule.RiskActionPlanStatusUpdateController
	riskActionPlanCommentController := app.RiskModule.RiskActionPlanCommentController
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
	kriAssociatedItemController := app.RiskModule.KRIAssociatedItemController
	riskWorkflowController := app.RiskModule.RiskWorkflowController
	riskReviewController := app.RiskModule.RiskReviewController
	riskReviewCommentController := app.RiskModule.RiskReviewCommentController
	kriRegisterWorkflowController := app.RiskModule.KRIRegisterWorkflowController
	kriRegisterReviewController := app.RiskModule.KRIRegisterReviewController
	kriRegisterReviewCommentController := app.RiskModule.KRIRegisterReviewCommentController
	riskResponseController := app.RiskModule.RiskResponseController
	riskResponseWorkflowController := app.RiskModule.RiskResponseWorkflowController
	riskResponseReviewController := app.RiskModule.RiskResponseReviewController
	riskResponseReviewCommentController := app.RiskModule.RiskResponseReviewCommentController
	riskRegisterUpdatePlanController := app.RiskModule.RiskRegisterUpdatePlanController
	riskUpdatePlanTimelineTypeController := app.RiskModule.RiskUpdatePlanTimelineTypeController
	riskUpdatePlanIntroductoryMeetingController := app.RiskModule.RiskUpdatePlanIntroductoryMeetingController
	riskUpdatePlanSelfAwarenessController := app.RiskModule.RiskUpdatePlanSelfAwarenessController
	riskUpdatePlanBusinessContextController := app.RiskModule.RiskUpdatePlanBusinessContextController
	riskUpdatePlanERMDeptController := app.RiskModule.RiskUpdatePlanERMDeptController
	riskUpdatePlanLeadershipAssessmentController := app.RiskModule.RiskUpdatePlanLeadershipAssessmentController
	riskUpdatePlanRiskTeamAssessmentController := app.RiskModule.RiskUpdatePlanRiskTeamAssessmentController
	riskUpdatePlanRiskIdentificationController := app.RiskModule.RiskUpdatePlanRiskIdentificationController
	riskUpdatePlanControlRegisterUpdateController := app.RiskModule.RiskUpdatePlanControlRegisterUpdateController
	riskUpdatePlanRiskAssessmentWorkshopController := app.RiskModule.RiskUpdatePlanRiskAssessmentWorkshopController
	riskUpdatePlanRiskWorkshopAssessmentController := app.RiskModule.RiskUpdatePlanRiskWorkshopAssessmentController
	riskUpdatePlanRiskAnalysisController := app.RiskModule.RiskUpdatePlanRiskAnalysisController
	riskUpdatePlanRiskControlAssociationController := app.RiskModule.RiskUpdatePlanRiskControlAssociationController
	riskUpdatePlanRiskMitigationController := app.RiskModule.RiskUpdatePlanRiskMitigationController

	riskRoutes.GET("/risk-categories", riskCategoryController.GetAllRiskCategories)
	riskRoutes.GET("/risk-categories/:id", riskCategoryController.GetRiskCategoryDetails)
	riskRoutes.POST("/risk-categories", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskCategoryController.CreateRiskCategory)
	riskRoutes.PUT("/risk-categories/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskCategoryController.UpdateRiskCategory)
	riskRoutes.POST("/risk-categories/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskCategoryController.DeleteRiskCategories)
	riskRoutes.PUT("/risk-categories/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskCategoryController.ActivateRiskCategories)
	riskRoutes.PUT("/risk-categories/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskCategoryController.DeactivateRiskCategories)

	riskRoutes.GET("/risk-sub-categories", riskSubCategoryController.GetAllRiskSubCategories)
	riskRoutes.GET("/risk-sub-categories/:id", riskSubCategoryController.GetRiskSubCategoryDetails)
	riskRoutes.POST("/risk-sub-categories", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskSubCategoryController.CreateRiskSubCategory)
	riskRoutes.PUT("/risk-sub-categories/:id", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskSubCategoryController.UpdateRiskSubCategory)
	riskRoutes.POST("/risk-sub-categories/delete", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskSubCategoryController.DeleteRiskSubCategories)
	riskRoutes.PUT("/risk-sub-categories/activate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskSubCategoryController.ActivateRiskSubCategories)
	riskRoutes.PUT("/risk-sub-categories/deactivate", middlewares.RolePermissionMiddleware(db, "risk_setting", "modify"), riskSubCategoryController.DeactivateRiskSubCategories)

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

	riskRoutes.GET("/level-of-measurements", levelOfMeasurementController.GetAllLevelOfMeasurements)
	riskRoutes.GET("/level-of-measurements/:id", levelOfMeasurementController.GetLevelOfMeasurementDetails)

	riskRoutes.GET("/distribution-statuses", distributionStatusController.GetAllDistributionStatuses)
	riskRoutes.GET("/distribution-statuses/:id", distributionStatusController.GetDistributionStatusDetails)

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

	riskRoutes.GET("/risk-appetite-statements/all", middlewares.RolePermissionMiddleware(db, "risk_appetite_statement", "view"), riskAppetiteStatementController.GetAllRiskAppetiteStatments)
	riskRoutes.GET("/risk-appetite-statements/summary", middlewares.RolePermissionMiddleware(db, "risk_appetite_statement", "view"), riskAppetiteStatementController.GetRiskAppetiteStatementSummary)
	riskRoutes.GET("/risk-appetite-statements/:id/risks", middlewares.RolePermissionMiddleware(db, "risk_appetite_statement", "view"), riskAppetiteStatementController.GetAllRisks)
	riskRoutes.GET("/risk-appetite-statements/:id/summary", riskAppetiteStatementController.GetRiskAppetiteStatementDetailsSummary)

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
	riskRoutes.GET("/kri-register/:id/trends", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriRegisterController.GetKRITrendAnalysis)
	riskRoutes.POST("/kri-register/delete", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterController.DeleteKRIRegisters)
	riskRoutes.POST("/kri-register/force-delete", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterController.DeleteForceKRIRegisters)
	riskRoutes.POST("/kri-register/:id/submit", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterReviewController.SubmitKRIRegister)
	riskRoutes.POST("/kri-register/:id/approve", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterReviewController.ApproveKRIRegister)
	riskRoutes.POST("/kri-register/:id/reject", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterReviewController.RejectKRIRegister)
	riskRoutes.GET("/kri-register/:id/workflow-list", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriRegisterWorkflowController.GetWorkflowList)

	riskRoutes.GET("/kri-register/:id/review-comments", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriRegisterReviewCommentController.GetAllComments)
	riskRoutes.GET("/kri-register/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriRegisterReviewCommentController.GetCommentDetails)
	riskRoutes.POST("/kri-register/:id/review-comments", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterReviewCommentController.CreateComment)
	riskRoutes.PUT("/kri-register/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterReviewCommentController.UpdateComment)
	riskRoutes.POST("/kri-register/:id/review-comments/delete", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterReviewCommentController.DeleteComment)
	riskRoutes.PUT("/kri-register/:id/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriRegisterReviewCommentController.ResolveComment)

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

	riskRoutes.GET("/risk-register-update-plan-statuses", riskRegisterUpdatePlanStatusController.GetAllRiskRegisterUpdatePlanStatuses)
	riskRoutes.GET("/risk-register-update-plan-statuses/:id", riskRegisterUpdatePlanStatusController.GetRiskRegisterUpdatePlanStatusDetails)

	// Risk Register
	unauthenticatedRoutes.GET("/risk-register/:id/document-preview", riskRegisterController.GetRiskRegisterDocumentPreview)
	unauthenticatedRoutes.GET("/risk-register/:id/document-download", riskRegisterController.DownloadRiskRegisterDocument)
	riskRoutes.GET("/risk-register/export", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskRegisterController.GetRiskRegisterExport)
	riskRoutes.GET("/risk-register", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskRegisterController.GetAllRiskRegisters)
	riskRoutes.GET("/risk-register/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskRegisterController.GetRiskRegisterDetails)
	riskRoutes.POST("/risk-register", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.CreateRiskRegister)
	riskRoutes.PUT("/risk-register/:id", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.UpdateRiskRegister)
	riskRoutes.POST("/risk-register/delete", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.DeleteRiskRegisters)
	riskRoutes.POST("/risk-register/force-delete", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.DeleteForceRiskRegisters)

	unauthenticatedRoutes.GET("/identified-risks/:id/document-preview", riskRegisterController.GetRiskRegisterDocumentPreview)
	unauthenticatedRoutes.GET("/identified-risks/:id/document-download", riskRegisterController.DownloadRiskRegisterDocument)
	riskRoutes.GET("/identified-risks/export", middlewares.RolePermissionMiddleware(db, "identified_risk", "view"), riskRegisterController.GetIdentifiedRiskRegisterExport)
	riskRoutes.GET("/identified-risks", middlewares.RolePermissionMiddleware(db, "identified_risk", "view"), riskRegisterController.GetAllIdentifiedRisks)
	riskRoutes.GET("/identified-risks/:id", middlewares.RolePermissionMiddleware(db, "identified_risk", "view"), riskRegisterController.GetRiskRegisterDetails)
	riskRoutes.POST("/identified-risks", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskRegisterController.CreateRiskRegister)
	riskRoutes.PUT("/identified-risks/:id", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskRegisterController.UpdateRiskRegister)
	riskRoutes.POST("/identified-risks/delete", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskRegisterController.DeleteRiskRegisters)

	riskRoutes.GET("/identified-risks/template", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskRegisterController.GenerateRiskTemplate)
	riskRoutes.POST("/identified-risks/template/validate", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskRegisterController.ValidateRiskTemplate)
	riskRoutes.POST("/identified-risks/import/validate", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskRegisterController.ValidateImportRisks)
	riskRoutes.POST("/identified-risks/import", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskRegisterController.ImportRisks)

	riskRoutes.GET("risks/:riskId/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskRootCauseAnalysisController.GetAllRiskRootCauseAnalyses)
	riskRoutes.POST("risks/:riskId/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRootCauseAnalysisController.CreateRiskRootCauseAnalysis)
	riskRoutes.POST("risks/:riskId/root-cause-analyses", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRootCauseAnalysisController.CreateRiskRootCauseAnalyses)
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
	riskRoutes.GET("/risks/:riskId/clauses", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllClauses)
	riskRoutes.PUT("/risks/:riskId/clauses", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.UpdateClauses)

	riskRoutes.POST("risks/:riskId/processes", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.CreatRiskProcess)
	riskRoutes.DELETE("risks/:riskId/processes/:processId", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.DeletetRiskProcess)

	riskRoutes.POST("risks/:riskId/documents", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.CreatRiskDocument)
	riskRoutes.DELETE("risks/:riskId/documents/:documentId", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.DeletetRiskDocument)

	riskRoutes.POST("risks/:riskId/incidents", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.CreatRiskIncident)
	riskRoutes.DELETE("risks/:riskId/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.DeletetRiskIncident)

	riskRoutes.POST("risks/:riskId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.CreatRiskStrategicObjective)
	riskRoutes.DELETE("risks/:riskId/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.DeletetRiskStrategicObjective)

	riskRoutes.GET("risks/:riskId/kpis", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskKPI)
	riskRoutes.POST("risks/:riskId/kpis", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.CreatRiskKPI)
	riskRoutes.DELETE("risks/:riskId/kpis/:kpiId", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.DeletetRiskKPI)

	riskRoutes.GET("risks/:riskId/bcp", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskBCP)
	riskRoutes.POST("risks/:riskId/bcp", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.CreatRiskBCP)
	riskRoutes.DELETE("risks/:riskId/bcp/:bcpId", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.DeletetRiskBCP)

	riskRoutes.GET("risks/:riskId/projects", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskAssociatedItemController.GetAllRiskProject)
	riskRoutes.POST("risks/:riskId/projects", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.CreatRiskProject)
	riskRoutes.DELETE("risks/:riskId/projects/:projectId", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.DeletetRiskProject)

	riskRoutes.POST("risks/:riskId/controls", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.CreatRiskControl)
	riskRoutes.DELETE("risks/:riskId/controls/:controlId", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.DeletetRiskControl)

	riskRoutes.POST("risks/:riskId/assets", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.CreatRiskAsset)
	riskRoutes.DELETE("risks/:riskId/assets/:assetId", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskAssociatedItemController.DeletetRiskAsset)

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
	riskRoutes.POST("/risks/:riskId/action-plans/:actionPlanId/status-updates", riskActionPlanStatusUpdateController.CreateRiskActionPlanStatusUpdate)
	riskRoutes.PUT("/risks/:riskId/action-plans/:actionPlanId/status-updates/:id", riskActionPlanStatusUpdateController.UpdateRiskActionPlanStatusUpdate)
	riskRoutes.GET("/risks/:riskId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "view"), riskActionPlanStatusUpdateController.GetLatestRiskActionPlanStatusUpdate)
	riskRoutes.PUT("/risks/:riskId/action-plans/:actionPlanId/approve", riskActionPlanStatusUpdateController.ApproveRiskActionPlan)
	riskRoutes.PUT("/risks/:riskId/action-plans/:actionPlanId/revert", riskActionPlanStatusUpdateController.RevertRiskActionPlan)

	//Risk Action Plan Comment
	unauthenticatedRoutes.GET("/risks/:riskId/action-plans/:actionPlanId/comments/:commentId/document-download", riskActionPlanCommentController.DownloadRiskActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/risks/:riskId/action-plans/:actionPlanId/comments/:commentId/document-preview", riskActionPlanCommentController.GetRiskActionPlanCommentDocumentPreview)
	riskRoutes.GET("/risks/:riskId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "view"), riskActionPlanCommentController.GetAllComments)
	riskRoutes.GET("/risks/:riskId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "view"), riskActionPlanCommentController.GetCommentDetails)
	riskRoutes.POST("/risks/:riskId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "modify"), riskActionPlanCommentController.CreateComment)
	riskRoutes.PUT("/risks/:riskId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "modify"), riskActionPlanCommentController.UpdateComment)
	riskRoutes.POST("/risks/:riskId/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "risk_mitigation", "modify"), riskActionPlanCommentController.DeleteComment)

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
	riskRoutes.GET("assets/:id/risk-heatmap", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), riskHeatmapController.GetAssetRiskHeatmap)
	riskRoutes.GET("strategic-objectives/:id/risk-heatmap", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), riskHeatmapController.GetStrategicObjectiveRiskHeatmap)
	riskRoutes.GET("kpis/:id/risk-heatmap", middlewares.RolePermissionMiddleware(db, "kpi", "view"), riskHeatmapController.GetKPIRiskHeatmap)

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
	riskRoutes.GET("risks/:riskId/residual-evaluation-history", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskResidualEvaluationController.GetRiskResidualEvaluationHistory)

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

	riskRoutes.GET("/risks/:riskId/comments", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskCommentController.GetAllComments)
	riskRoutes.GET("/risks/:riskId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskCommentController.GetCommentDetails)
	riskRoutes.POST("/risks/:riskId/comments", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.CreateComment)
	riskRoutes.PUT("/risks/:riskId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.UpdateComment)
	riskRoutes.POST("/risks/:riskId/comments/delete", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.DeleteComment)

	riskRoutes.PUT("/risks/:riskId/submit", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.SubmitRisk)
	riskRoutes.PUT("/risks/:riskId/revert", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.RevertRisk)
	riskRoutes.PUT("/risks/:riskId/reject", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.RejectRisk)
	riskRoutes.PUT("/risks/:riskId/approve", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.ApproveRisk)
	riskRoutes.PUT("/risks/:riskId/resubmit", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskCommentController.ResubmitRisk)

	unauthenticatedRoutes.GET("/risk-reports/:riskReportId/document-preview", riskReportController.GetRiskReportDocumentPreview)
	unauthenticatedRoutes.GET("/risk-reports/:riskReportId/document-download", riskReportController.DownloadRiskReportDocument)
	riskRoutes.GET("/risk-reports", riskReportController.GetAllRiskReports)
	riskRoutes.GET("/risk-reports/:riskReportId", riskReportController.GetRiskReportDetails)
	riskRoutes.GET("/risk-reports/:riskReportId/export-pdf", riskReportController.GetRiskReportExportZip)
	riskRoutes.POST("/risk-reports", riskReportController.CreateRiskReport)
	riskRoutes.GET("/risk-reports/:riskReportId/sections/:sectionId", riskReportController.GetRiskReportSection)
	riskRoutes.PUT("/risk-reports/:riskReportId/sections/:sectionId", riskReportController.UpdateRiskReportSection)
	riskRoutes.DELETE("/risk-reports/:riskReportId/sections/:sectionId", riskReportController.DeleteRiskReportSection)
	riskRoutes.DELETE("/risk-reports/:riskReportId", riskReportController.DeleteRiskReport)
	riskRoutes.POST("/risk-reports/:riskReportId/documents", riskReportController.CreateRiskReportDocuments)
	riskRoutes.DELETE("/risk-reports/:riskReportId/documents/:documentId", riskReportController.DeleteRiskReportDocument)

	riskRoutes.GET("/risks/template", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.GenerateRiskTemplate)
	riskRoutes.POST("/risks/template/validate", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.ValidateRiskTemplate)
	riskRoutes.POST("/risks/import/validate", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.ValidateImportRisks)
	riskRoutes.POST("/risks/import", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskRegisterController.ImportRisks)

	riskRoutes.GET("/kri-register/:id/documents", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriAssociatedItemController.GetAllKRIDocuments)
	riskRoutes.POST("/kri-register/:id/documents", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.CreateKRIDocument)
	riskRoutes.DELETE("/kri-register/:id/documents/:documentId", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.DeleteKRIDocument)

	riskRoutes.GET("/kri-register/:id/controls", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriAssociatedItemController.GetAllKRIControls)
	riskRoutes.POST("/kri-register/:id/controls", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.CreateKRIControl)
	riskRoutes.DELETE("/kri-register/:id/controls/:controlId", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.DeleteKRIControl)

	riskRoutes.GET("/kri-register/:id/assets", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriAssociatedItemController.GetAllKRIAssets)
	riskRoutes.POST("/kri-register/:id/assets", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.CreateKRIAsset)
	riskRoutes.DELETE("/kri-register/:id/assets/:assetId", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.DeleteKRIAsset)

	riskRoutes.GET("/kri-register/:id/incidents", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriAssociatedItemController.GetAllKRIIncidents)
	riskRoutes.POST("/kri-register/:id/incidents", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.CreateKRIIncident)
	riskRoutes.DELETE("/kri-register/:id/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.DeleteKRIIncident)

	riskRoutes.GET("/kri-register/:id/risks", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriAssociatedItemController.GetAllKRIRisks)
	riskRoutes.POST("/kri-register/:id/risks", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.CreateKRIRisk)
	riskRoutes.DELETE("/kri-register/:id/risks/:riskId", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.DeleteKRIRisk)

	riskRoutes.GET("/kri-register/:id/kpis", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriAssociatedItemController.GetAllKRIKPI)
	riskRoutes.POST("/kri-register/:id/kpis", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.CreateKRIKPI)
	riskRoutes.DELETE("/kri-register/:id/kpis/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.DeleteKRIKPI)

	riskRoutes.GET("/kri-register/:id/processes", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriAssociatedItemController.GetAllKRIProcesses)
	riskRoutes.POST("/kri-register/:id/processes", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.CreateKRIProcess)
	riskRoutes.DELETE("/kri-register/:id/processes/:processId", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.DeleteKRIProcess)

	riskRoutes.GET("/kri-register/:id/clauses", middlewares.RolePermissionMiddleware(db, "kri", "view"), kriAssociatedItemController.GetAllClauses)
	riskRoutes.PUT("/kri-register/:id/clauses", middlewares.RolePermissionMiddleware(db, "kri", "modify"), kriAssociatedItemController.UpdateClauses)

	riskRoutes.POST("/identified-risks/:id/submit", middlewares.RolePermissionMiddleware(db, "identified_risk", "view"), riskReviewController.SubmitRisk)
	riskRoutes.POST("/identified-risks/:id/approve", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskReviewController.ApproveRisk)
	riskRoutes.POST("/identified-risks/:id/reject", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskReviewController.RejectRisk)
	riskRoutes.GET("/identified-risks/:id/workflow-list", middlewares.RolePermissionMiddleware(db, "identified_risk", "view"), riskWorkflowController.GetWorkflowList)

	// risk review comment
	riskRoutes.GET("/identified-risks/:id/review-comments", middlewares.RolePermissionMiddleware(db, "identified_risk", "view"), riskReviewCommentController.GetAllComments)
	riskRoutes.GET("/identified-risks/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "identified_risk", "view"), riskReviewCommentController.GetCommentDetails)
	riskRoutes.POST("/identified-risks/:id/review-comments", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskReviewCommentController.CreateComment)
	riskRoutes.PUT("/identified-risks/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskReviewCommentController.UpdateComment)
	riskRoutes.POST("/identified-risks/:id/review-comments/delete", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskReviewCommentController.DeleteComment)
	riskRoutes.PUT("/identified-risks/:id/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(db, "identified_risk", "modify"), riskReviewCommentController.ResolveComment)

	riskRoutes.GET("/risks/:riskId/export-pdf", riskRegisterController.ExportRiskRegisterToPDF)

	//Risk Response
	riskRoutes.GET("/risk-response-details/:riskId", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskResponseController.GetRiskResponseDetails)

	riskRoutes.POST("/risk-response/:id/submit", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResponseReviewController.SubmitRiskResponse)
	riskRoutes.POST("/risk-response/:id/approve", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResponseReviewController.ApproveRiskResponse)
	riskRoutes.POST("/risk-response/:id/revert", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResponseReviewController.RejectRiskResponse)
	riskRoutes.POST("/risk-response/:id/reopen", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResponseReviewController.ReopenRiskResponse)
	riskRoutes.GET("/risk-response/:id/workflow-list", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskResponseWorkflowController.GetWorkflowList)

	riskRoutes.GET("/risk-response/:id/review-comments", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskResponseReviewCommentController.GetAllComments)
	riskRoutes.GET("/risk-response/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "risk_register", "view"), riskResponseReviewCommentController.GetCommentDetails)
	riskRoutes.POST("/risk-response/:id/review-comments", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResponseReviewCommentController.CreateComment)
	riskRoutes.PUT("/risk-response/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResponseReviewCommentController.UpdateComment)
	riskRoutes.POST("/risk-response/:id/review-comments/delete", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResponseReviewCommentController.DeleteComments)
	riskRoutes.PUT("/risk-response/:id/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(db, "risk_register", "modify"), riskResponseReviewCommentController.ResolveComment)

	unauthenticatedRoutes.GET("/risk-response/:id/review-comments/:commentId/document-preview", riskResponseReviewCommentController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/risk-response/:id/review-comments/:commentId/document-download", riskResponseReviewCommentController.DownloadDocument)

	// Risk register update plan
	riskRoutes.GET("/risk-register-update-plans", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskRegisterUpdatePlanController.GetAllRiskRegisterUpdatePlan)
	riskRoutes.GET("/risk-register-update-plans/:id", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskRegisterUpdatePlanController.GetRiskRegisterUpdatePlanDetailsByID)
	riskRoutes.POST("/risk-register-update-plans", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskRegisterUpdatePlanController.CreateRiskRegisterUpdatePlan)
	riskRoutes.PUT("/risk-register-update-plans/:id", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskRegisterUpdatePlanController.UpdateRiskRegisterUpdatePlan)
	riskRoutes.POST("/risk-register-update-plans/delete", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskRegisterUpdatePlanController.DeleteRiskRegisterUpdatePlan)
	// ERM Department Evaluation
	riskRoutes.POST("/risk-register-update-plans/:id/create-erm-evaluation", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanERMDeptController.CreateERMDepartmentEvaluation)
	// ERM Department Evaluation - Risks
	riskRoutes.POST("/risk-register-update-plans/:id/update-erm-evaluation-risks", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanERMDeptController.UpdateERMDepartmentEvaluationRisks)
	riskRoutes.GET("/risk-register-update-plans/:id/erm-evaluations-risks", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanERMDeptController.GetERMDepartmentEvaluationRisks)
	// ERM Department Evaluation - Processes
	riskRoutes.POST("/risk-register-update-plans/:id/update-erm-evaluation-processes", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanERMDeptController.UpdateERMDepartmentEvaluationProcesses)
	riskRoutes.GET("/risk-register-update-plans/:id/erm-evaluations-processes", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanERMDeptController.GetERMDepartmentEvaluationProcesses)
	// ERM Department Evaluation - Controls
	riskRoutes.POST("/risk-register-update-plans/:id/update-erm-evaluation-controls", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanERMDeptController.UpdateERMDepartmentEvaluationControls)
	riskRoutes.GET("/risk-register-update-plans/:id/erm-evaluations-controls", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanERMDeptController.GetERMDepartmentEvaluationControls)
	// ERM Department Evaluation - Operations
	riskRoutes.POST("/risk-register-update-plans/:id/mark-erm-evaluation-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanERMDeptController.MarkERMDeptEvaluationAsNotApplicable)
	// ERM Department Evaluation - Documents
	riskRoutes.GET("/risk-register-update-plans/:id/erm-dept-evaluation-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanERMDeptController.GetERMDeptEvaluationDocuments)
	riskRoutes.POST("/risk-register-update-plans/:id/erm-dept-evaluation-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanERMDeptController.CreateERMDeptEvaluationDocuments)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/erm-dept-evaluation-documents/preview", riskUpdatePlanERMDeptController.GetERMDeptEvaluationDocumentPreview)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/erm-dept-evaluation-documents/download", riskUpdatePlanERMDeptController.DownloadERMDeptEvaluationDocument)

	// Introductory Meeting Documents
	riskRoutes.GET("/risk-register-update-plans/:id/introductory-meeting-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanIntroductoryMeetingController.GetIntroductoryMeetingDocuments)
	riskRoutes.POST("/risk-register-update-plans/:id/introductory-meeting-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanIntroductoryMeetingController.CreateIntroductoryMeetingDocuments)
	// Introductory Meeting - Operations
	riskRoutes.POST("/risk-register-update-plans/:id/mark-introductory-meeting-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanIntroductoryMeetingController.MarkIntroductoryMeetingAsNotApplicable)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/introductory-meeting-documents/preview", riskUpdatePlanIntroductoryMeetingController.GetIntroductoryMeetingDocumentPreview)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/introductory-meeting-documents/download", riskUpdatePlanIntroductoryMeetingController.DownloadIntroductoryMeetingDocument)

	// Self-Awareness Training Documents
	riskRoutes.GET("/risk-register-update-plans/:id/self-awareness-trainings", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanSelfAwarenessController.GetSelfAwarenessTrainingDocuments)
	riskRoutes.POST("/risk-register-update-plans/:id/self-awareness-trainings", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanSelfAwarenessController.CreateSelfAwarenessTrainingDocuments)
	// Self-Awareness Training - Operations
	riskRoutes.POST("/risk-register-update-plans/:id/mark-self-awareness-training-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanSelfAwarenessController.MarkSelfAwarenessTrainingAsNotApplicable)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/self-awareness-trainings/preview", riskUpdatePlanSelfAwarenessController.GetSelfAwarenessTrainingDocumentPreview)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/self-awareness-trainings/download", riskUpdatePlanSelfAwarenessController.DownloadSelfAwarenessTrainingDocument)

	// Business Context Operations
	riskRoutes.GET("/risk-register-update-plans/:id/business-context", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanBusinessContextController.GetBusinessContext)
	riskRoutes.POST("/risk-register-update-plans/:id/business-context", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanBusinessContextController.CreateBusinessContext)
	riskRoutes.PUT("/risk-update-plan-business-context/:businessContextId", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanBusinessContextController.UpdateBusinessContext)
	riskRoutes.PUT("/risk-register-update-plans/:id/business-context-status", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanBusinessContextController.UpdateBusinessContextStatus)
	riskRoutes.POST("/risk-register-update-plans/:id/mark-business-context-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanBusinessContextController.MarkBusinessContextAsNotApplicable)

	// Leadership Assessment Operations
	riskRoutes.GET("/risk-register-update-plans/:id/leadership-assessment-users", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanLeadershipAssessmentController.GetLeadershipAssessmentUsers)
	riskRoutes.POST("/risk-register-update-plans/:id/leadership-assessments", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanLeadershipAssessmentController.CreateLeadershipAssessment)
	riskRoutes.POST("/risk-register-update-plans/:id/mark-leadership-assessment-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanLeadershipAssessmentController.MarkLeadershipAssessmentAsNotApplicable)
	riskRoutes.GET("/risk-register-update-plans/:id/leadership-assessment-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanLeadershipAssessmentController.GetLeadershipAssessmentDocuments)
	riskRoutes.POST("/risk-register-update-plans/:id/leadership-assessment-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanLeadershipAssessmentController.CreateLeadershipAssessmentDocuments)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/leadership-assessment-documents/preview", riskUpdatePlanLeadershipAssessmentController.GetLeadershipAssessmentDocumentPreview)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/leadership-assessment-documents/download", riskUpdatePlanLeadershipAssessmentController.DownloadLeadershipAssessmentDocument)

	// Risk Team Assessment Operations
	riskRoutes.GET("/risk-register-update-plans/:id/risk-team-assessment-user", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanRiskTeamAssessmentController.GetRiskTeamAssessmentUser)
	riskRoutes.POST("/risk-register-update-plans/:id/risk-team-assessments", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskTeamAssessmentController.CreateRiskTeamAssessment)
	riskRoutes.POST("/risk-register-update-plans/:id/mark-risk-team-assessment-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskTeamAssessmentController.MarkRiskTeamAssessmentAsNotApplicable)
	riskRoutes.GET("/risk-register-update-plans/:id/risk-team-assessment-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanRiskTeamAssessmentController.GetRiskTeamAssessmentDocuments)
	riskRoutes.POST("/risk-register-update-plans/:id/risk-team-assessment-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskTeamAssessmentController.CreateRiskTeamAssessmentDocuments)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/risk-team-assessment-documents/preview", riskUpdatePlanRiskTeamAssessmentController.GetRiskTeamAssessmentDocumentPreview)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/risk-team-assessment-documents/download", riskUpdatePlanRiskTeamAssessmentController.DownloadRiskTeamAssessmentDocument)

	// Risk Identification Operations
	riskRoutes.GET("/risk-register-update-plans/:id/risk-identifications", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanRiskIdentificationController.GetAllRiskIdentifications)
	riskRoutes.GET("/risk-register-update-plans/:id/risk-identifications/:riskId", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanRiskIdentificationController.GetRiskIdentificationDetails)
	riskRoutes.POST("/risk-register-update-plans/:id/risk-identifications", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskIdentificationController.CreateRiskIdentification)
	riskRoutes.PUT("/risk-register-update-plans/:id/risk-identifications/:riskId", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskIdentificationController.UpdateRiskIdentification)
	riskRoutes.POST("/risk-register-update-plans/:id/risk-identifications/delete", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskIdentificationController.DeleteRiskIdentification)
	riskRoutes.POST("/risk-register-update-plans/:id/mark-risk-identification-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskIdentificationController.MarkRiskIdentificationAsNotApplicable)

	// Control Register Update Operations
	riskRoutes.GET("/risk-register-update-plans/:id/control-counts", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanControlRegisterUpdateController.GetControlCountsByImplementationStatus)
	riskRoutes.POST("/risk-register-update-plans/:id/control-tests", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanControlRegisterUpdateController.CreateControlTest)
	riskRoutes.POST("/risk-register-update-plans/:id/mark-control-register-update-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanControlRegisterUpdateController.MarkControlRegisterUpdateAsNotApplicable)

	// Risk Assessment Workshop Documents
	riskRoutes.GET("/risk-register-update-plans/:id/risk-assessment-workshop-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanRiskAssessmentWorkshopController.GetRiskAssessmentWorkshopDocuments)
	riskRoutes.POST("/risk-register-update-plans/:id/risk-assessment-workshop-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskAssessmentWorkshopController.CreateRiskAssessmentWorkshopDocuments)
	// Risk Assessment Workshop - Operations
	riskRoutes.POST("/risk-register-update-plans/:id/mark-risk-assessment-workshop-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskAssessmentWorkshopController.MarkRiskAssessmentWorkshopAsNotApplicable)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/risk-assessment-workshop-documents/preview", riskUpdatePlanRiskAssessmentWorkshopController.GetRiskAssessmentWorkshopDocumentPreview)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/risk-assessment-workshop-documents/download", riskUpdatePlanRiskAssessmentWorkshopController.DownloadRiskAssessmentWorkshopDocument)

	// Risk Workshop Assessment Documents
	riskRoutes.GET("/risk-register-update-plans/:id/risk-workshop-assessment-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanRiskWorkshopAssessmentController.GetRiskWorkshopAssessmentDocuments)
	riskRoutes.POST("/risk-register-update-plans/:id/risk-workshop-assessment-documents", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskWorkshopAssessmentController.CreateRiskWorkshopAssessmentDocuments)
	// Risk Workshop Assessment - Operations
	riskRoutes.POST("/risk-register-update-plans/:id/mark-risk-workshop-assessment-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskWorkshopAssessmentController.MarkRiskWorkshopAssessmentAsNotApplicable)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/risk-workshop-assessment-documents/preview", riskUpdatePlanRiskWorkshopAssessmentController.GetRiskWorkshopAssessmentDocumentPreview)
	unauthenticatedRoutes.GET("/risk-register-update-plans/:id/risk-workshop-assessment-documents/download", riskUpdatePlanRiskWorkshopAssessmentController.DownloadRiskWorkshopAssessmentDocument)

	// Risk Analysis Operations
	riskRoutes.GET("/risk-register-update-plans/:id/risk-analysis", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanRiskAnalysisController.GetRiskAnalysisList)
	riskRoutes.POST("/risk-register-update-plans/:id/risk-analysis", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskAnalysisController.CreateRiskAnalysis)
	riskRoutes.POST("/risk-register-update-plans/:id/mark-risk-analysis-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskAnalysisController.MarkRiskAnalysisAsNotApplicable)

	// Risk Control Association
	riskRoutes.GET("/risk-register-update-plans/:id/risk-control-associations", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanRiskControlAssociationController.GetRiskControlAssociationList)
	riskRoutes.POST("/risk-register-update-plans/:id/risk-control-association", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskControlAssociationController.CreateRiskControlAssociation)
	riskRoutes.POST("/risk-register-update-plans/:id/mark-risk-control-association-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskControlAssociationController.MarkRiskControlAssociationAsNotApplicable)

	// Risk Mitigation Operations
	riskRoutes.POST("/risk-register-update-plans/:id/risk-mitigation", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskMitigationController.UpdateRiskMitigation)
	riskRoutes.DELETE("/risk-register-update-plans/:id/risk-mitigation/:riskId", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskMitigationController.DeleteRiskMitigation)
	riskRoutes.POST("/risk-register-update-plans/:id/mark-risk-mitigation-not-applicable", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskMitigationController.MarkRiskMitigationAsNotApplicable)
	// Risk Mitigation - Action Plans
	riskRoutes.GET("/risk-register-update-plans/:id/risks/:riskId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "view"), riskUpdatePlanRiskMitigationController.GetActionPlanDetails)
	riskRoutes.POST("/risk-register-update-plans/:id/risks/:riskId/action-plans", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskMitigationController.CreateActionPlan)
	riskRoutes.PUT("/risk-register-update-plans/:id/risks/:riskId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskMitigationController.UpdateActionPlan)
	riskRoutes.POST("/risk-register-update-plans/:id/risks/:riskId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "risk_register_update_plan", "modify"), riskUpdatePlanRiskMitigationController.DeleteActionPlan)

	// Risk Update Plan Timeline Types
	riskRoutes.GET("/risk-update-plan-timeline-types", riskUpdatePlanTimelineTypeController.GetAllRiskUpdatePlanTimelineTypes)
	riskRoutes.GET("/risk-update-plan-timeline-types/:id", riskUpdatePlanTimelineTypeController.GetRiskUpdatePlanTimelineTypeDetails)
}
