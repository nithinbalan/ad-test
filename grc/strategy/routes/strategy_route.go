package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	strategyRoutes := router.Group("/")

	strategyRoutes.Use(middlewares.AuthMiddleware())

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	strategyProfileController := app.StrategyModule.StrategyProfileController
	strategyProfileFocusAreaController := app.StrategyModule.StrategyProfileFocusAreaController
	strategyProfileGoalController := app.StrategyModule.StrategyProfileGoalController
	strategyProfileStrategicObjectiveController := app.StrategyModule.StrategyProfileStrategicObjectiveController
	strategyProfileStrategicIntiativeController := app.StrategyModule.StrategyProfileStrategicIntiativeController
	strategyProfileProjectController := app.StrategyModule.StrategyProfileProjectController
	strategyProfileProjectMilestoneController := app.StrategyModule.StrategyProfileProjectMilestoneController
	strategyProfileActionPlanController := app.StrategyModule.StrategyProfileActionPlanController
	strategyProfileActionPlanStatusUpdateController := app.StrategyModule.StrategyProfileActionPlanStatusUpdateController
	strategyProfileActionPlanCommentController := app.StrategyModule.StrategyProfileActionPlanCommentController
	strategyProfileSettingController := app.StrategyModule.StrategyProfileSettingController
	strategyProfileMappingController := app.StrategyModule.StrategyProfileMappingController
	kpiCategoryController := app.StrategyModule.KpiCategoryController
	kpiStatusController := app.StrategyModule.KpiStatusController
	entryTypeController := app.StrategyModule.EntryTypeController
	dataSourceController := app.StrategyModule.DataSourceController
	kpiRegisterController := app.StrategyModule.KpiRegisterController
	strategicObjectiveCategoryController := app.StrategyModule.StrategicObjectiveCategoryController
	strategicObjectiveController := app.StrategyModule.StrategicObjectiveController
	strategicObjectiveKPIController := app.StrategyModule.StrategicObjectiveKPIController
	strategicObjectivePeriodicReviewController := app.StrategyModule.StrategicObjectivePeriodicReviewController
	kpiValueUpdateController := app.StrategyModule.KpiValueUpdateController
	kpiAssociatedItemController := app.StrategyModule.KpiAssociatedItemController
	focusAreaController := app.StrategyModule.FocusAreaController
	strategicIntiativeController := app.StrategyModule.StrategicIntiativeController
	strategicObjectiveAssociatedItemController := app.StrategyModule.StrategicObjectiveAssociatedItemController
	strategyProfileKpiController := app.StrategyModule.StrategyProfileKpiController
	strategyProfileProjectActionPlanController := app.StrategyModule.StrategyProfileProjectActionPlanController

	// strategy profile
	strategyRoutes.GET("/strategy-profiles", strategyProfileController.GetAllStrategyProfiles)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId", strategyProfileController.GetStrategyProfileDetails)
	strategyRoutes.POST("/strategy-profiles", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileController.CreateStrategyProfile)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileController.UpdateStrategyProfile)
	strategyRoutes.POST("/strategy-profiles/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileController.DeleteStrategyProfiles)

	// goal
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/goals", strategyProfileGoalController.GetAllStrategyProfileGoals)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/goals/:goalId", strategyProfileGoalController.GetStrategyProfileGoalDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/goals", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileGoalController.CreateStrategyProfileGoal)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/goals/:goalId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileGoalController.UpdateStrategyProfileGoal)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/goals/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileGoalController.DeleteStrategyProfileGoals)

	// action plan
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/action-plans", strategyProfileActionPlanController.GetAllStrategyProfileActionPlans)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId", strategyProfileActionPlanController.GetStrategyProfileActionPlanDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/action-plans", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileActionPlanController.CreateStrategyProfileActionPlan)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileActionPlanController.UpdateStrategyProfileActionPlan)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileActionPlanController.DeleteStrategyProfileActionPlans)

	// status update
	unauthenticatedRoutes.GET("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/status-updates/:id/document-download", strategyProfileActionPlanStatusUpdateController.DownloadStrategyProfileActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/status-updates/:id/document-preview", strategyProfileActionPlanStatusUpdateController.GetStrategyProfileActionPlanStatusUpdateDocumentPreview)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategyProfileActionPlanStatusUpdateController.GetAllStrategyProfileActionPlanStatusUpdates)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileActionPlanStatusUpdateController.CreateStrategyProfileActionPlanStatusUpdate)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileActionPlanStatusUpdateController.UpdateStrategyProfileActionPlanStatusUpdate)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategyProfileActionPlanStatusUpdateController.GetLatestStrategyProfileActionPlanStatusUpdate)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/approve", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileActionPlanStatusUpdateController.ApproveStrategyProfileActionPlan)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/revert", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileActionPlanStatusUpdateController.RevertStrategyProfileActionPlan)

	//  Comments
	unauthenticatedRoutes.GET("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/comments/:commentId/document-download", strategyProfileActionPlanCommentController.DownloadStrategyProfileActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/comments/:commentId/document-preview", strategyProfileActionPlanCommentController.GetStrategyProfileActionPlanCommentDocumentPreview)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategyProfileActionPlanCommentController.GetAllComments)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategyProfileActionPlanCommentController.GetCommentDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileActionPlanCommentController.CreateComment)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileActionPlanCommentController.UpdateComment)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileActionPlanCommentController.DeleteComment)

	// profile setting
	strategyRoutes.GET("/strategy-profile-setting", strategyProfileSettingController.GetAllStrategyProfileSettings)
	strategyRoutes.GET("/strategy-profile-setting/:settingId", strategyProfileSettingController.GetStrategyProfileSettingDetail)
	strategyRoutes.PUT("/strategy-profile-setting/:settingId/enable", strategyProfileSettingController.EnableStrategyProfileSetting)
	strategyRoutes.PUT("/strategy-profile-setting/:settingId/disable", strategyProfileSettingController.DisableStrategyProfileSetting)

	// mapping
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/mapping", strategyProfileMappingController.GetStrategyProfile)
	// kpi category
	strategyRoutes.GET("/kpi-categories", kpiCategoryController.GetAllKPICategories)
	strategyRoutes.POST("/kpi-categories", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), kpiCategoryController.CreateKPICategory)
	strategyRoutes.GET("/kpi-categories/:id", kpiCategoryController.GetKPICategoryDetails)
	strategyRoutes.PUT("/kpi-categories/:id", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), kpiCategoryController.UpdateKPICategory)
	strategyRoutes.POST("/kpi-categories/delete", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), kpiCategoryController.DeleteKPICategories)
	strategyRoutes.PUT("/kpi-categories/activate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), kpiCategoryController.ActivateKPICategories)
	strategyRoutes.PUT("/kpi-categories/deactivate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), kpiCategoryController.DeactivateKPICategories)

	// kpi status
	strategyRoutes.GET("/kpi-statuses", kpiStatusController.GetAllKPIStatuses)
	strategyRoutes.GET("/kpi-statuses/:id", kpiStatusController.GetKPIStatusDetails)

	// entry type
	strategyRoutes.GET("/entry-types", entryTypeController.GetAllEntryTypees)
	strategyRoutes.GET("/entry-types/:id", entryTypeController.GetEntryTypeDetails)

	// data source
	strategyRoutes.GET("/data-sources", dataSourceController.GetAllDataSourcees)
	strategyRoutes.GET("/data-sources/:id", dataSourceController.GetDataSourceDetails)

	// kpi register
	strategyRoutes.GET("/kpi-register", kpiRegisterController.GetAllKPIRegister)
	strategyRoutes.GET("/kpi-register/:kpiRegisterId", kpiRegisterController.GetKPIRegisterDetails)
	strategyRoutes.POST("/kpi-register", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.CreateKPIRegister)
	strategyRoutes.PUT("/kpi-register/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.UpdateKPIRegister)
	strategyRoutes.POST("/kpi-register/delete", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.DeleteKPIRegister)
	strategyRoutes.GET("/kpi-register/:kpiRegisterId/trends", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiRegisterController.GetKPITrendAnalysis)
	strategyRoutes.GET("/kpi-register/:kpiRegisterId/latest-value-details", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiRegisterController.GetKPILatestValueUpdate)

	unauthenticatedRoutes.GET("/kpi-register/:kpiRegisterId/value-updates/:updateId/document-preview", kpiValueUpdateController.GetKPIValueUpdateDocumentPreview)
	unauthenticatedRoutes.GET("/kpi-register/:kpiRegisterId/value-updates/:updateId/document-download", kpiValueUpdateController.DownloadKPIValueUpdateDocument)
	strategyRoutes.GET("/kpi-register/:kpiRegisterId/value-updates", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiValueUpdateController.GetAllKPIValueUpdates)
	strategyRoutes.POST("/kpi-register/:kpiRegisterId/value-updates", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiValueUpdateController.CreateKPIValueUpdate)
	strategyRoutes.GET("/kpi-register/:kpiRegisterId/value-updates/:updateId", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiValueUpdateController.GetKPIValueUpdateDetails)
	strategyRoutes.PUT("/kpi-register/:kpiRegisterId/value-updates/:updateId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiValueUpdateController.UpdateKPIValueUpdate)
	strategyRoutes.POST("/kpi-register/:kpiRegisterId/value-updates/delete", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiValueUpdateController.DeleteKPIValueUpdates)

	// KPI Associated Items
	strategyRoutes.GET("/kpi-register/:kpiRegisterId/controls", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIControls)
	strategyRoutes.POST("/kpi-register/:kpiRegisterId/controls", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIControls)
	strategyRoutes.DELETE("/kpi-register/:kpiRegisterId/controls/:controlId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIControl)

	strategyRoutes.GET("/kpi-register/:kpiRegisterId/processes", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIProcesses)
	strategyRoutes.POST("/kpi-register/:kpiRegisterId/processes", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIProcesses)
	strategyRoutes.DELETE("/kpi-register/:kpiRegisterId/processes/:processId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIProcess)

	strategyRoutes.GET("/kpi-register/:kpiRegisterId/compliances", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPICompliances)
	strategyRoutes.POST("/kpi-register/:kpiRegisterId/compliances", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPICompliances)
	strategyRoutes.DELETE("/kpi-register/:kpiRegisterId/compliances/:complianceId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPICompliance)

	strategyRoutes.GET("/kpi-register/:kpiRegisterId/assets", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIAssets)
	strategyRoutes.POST("/kpi-register/:kpiRegisterId/assets", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIAssets)
	strategyRoutes.DELETE("/kpi-register/:kpiRegisterId/assets/:assetId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIAsset)

	strategyRoutes.GET("/kpi-register/:kpiRegisterId/risks", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIRisks)
	strategyRoutes.POST("/kpi-register/:kpiRegisterId/risks", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIRisks)
	strategyRoutes.DELETE("/kpi-register/:kpiRegisterId/risks/:riskId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIRisk)

	strategyRoutes.GET("/kpi-register/:kpiRegisterId/third-party-register", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIThirdParties)
	strategyRoutes.POST("/kpi-register/:kpiRegisterId/third-party-register", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIThirdParties)
	strategyRoutes.DELETE("/kpi-register/:kpiRegisterId/third-party-register/:thirdPartyRegisterId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIThirdParty)

	strategyRoutes.GET("/kpi-register/:kpiRegisterId/documents", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIPolicies)
	strategyRoutes.POST("/kpi-register/:kpiRegisterId/documents", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIPolicies)
	strategyRoutes.DELETE("/kpi-register/:kpiRegisterId/documents/:documentId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIPolicy)

	strategyRoutes.GET("/kpi-register/:kpiRegisterId/clauses", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllClauses)
	strategyRoutes.PUT("/kpi-register/:kpiRegisterId/clauses", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.UpdateClauses)
	// strategic objective category
	strategyRoutes.GET("/strategic-objective-categories", strategicObjectiveCategoryController.GetAllStrategicObjectiveCategories)
	strategyRoutes.POST("/strategic-objective-categories", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveCategoryController.CreateStrategicObjectiveCategory)
	strategyRoutes.GET("/strategic-objective-categories/:id", strategicObjectiveCategoryController.GetStrategicObjectiveCategoryDetails)
	strategyRoutes.PUT("/strategic-objective-categories/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveCategoryController.UpdateStrategicObjectiveCategory)
	strategyRoutes.POST("/strategic-objective-categories/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveCategoryController.DeleteStrategicObjectiveCategories)
	strategyRoutes.PUT("/strategic-objective-categories/activate", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveCategoryController.ActivateStrategicObjectiveCategories)
	strategyRoutes.PUT("/strategic-objective-categories/deactivate", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveCategoryController.DeactivateStrategicObjectiveCategories)

	strategyRoutes.GET("/strategic-objectives", strategicObjectiveController.GetAllStrategicObjectives)
	strategyRoutes.GET("/strategic-objectives/:id", strategicObjectiveController.GetStrategicObjectiveDetails)
	strategyRoutes.POST("/strategic-objectives", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveController.CreateStrategicObjective)
	strategyRoutes.PUT("/strategic-objectives/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveController.UpdateStrategicObjective)
	strategyRoutes.POST("/strategic-objectives/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveController.DeleteStrategicObjectives)
	strategyRoutes.PUT("/strategic-objectives/activate", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveController.ActivateStrategicObjectives)
	strategyRoutes.PUT("/strategic-objectives/deactivate", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveController.DeactivateStrategicObjectives)
	strategyRoutes.GET("/strategic-objectives/:id/focus-areas", strategicObjectiveController.GetFocusAreasByStrategicObjectiveID)
	strategyRoutes.GET("/strategic-objectives/:id/risk-summary", strategicObjectiveController.GetRiskSummary)

	strategyRoutes.GET("/strategic-objectives/:id/kpis", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveKPIController.GetAllStrategicObjectiveKPI)
	strategyRoutes.POST("/strategic-objectives/:id/kpis", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveKPIController.CreateStrategicObjectiveKPI)
	strategyRoutes.DELETE("/strategic-objectives/:id/kpis/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveKPIController.DeleteStrategicObjectiveKPI)

	//review
	strategyRoutes.GET("/strategic-objectives/:id/reviews", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategicObjectivePeriodicReviewController.GetAllStrategicObjectivePeriodicReviews)
	strategyRoutes.GET("/strategic-objectives/:id/reviews/:reviewId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategicObjectivePeriodicReviewController.GetStrategicObjectivePeriodicReviewDetails)
	strategyRoutes.POST("/strategic-objectives/:id/reviews", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectivePeriodicReviewController.CreateStrategicObjectivePeriodicReview)
	strategyRoutes.PUT("/strategic-objectives/:id/reviews/:reviewId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectivePeriodicReviewController.UpdateStrategicObjectivePeriodicReview)
	strategyRoutes.POST("/strategic-objectives/:id/reviews/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectivePeriodicReviewController.DeleteStrategicObjectivePeriodicReview)

	unauthenticatedRoutes.GET("/strategic-objectives/:id/reviews/:reviewId/document-download", strategicObjectivePeriodicReviewController.DownloadStrategicObjectivePeriodicReviewDocument)

	strategyRoutes.GET("/focus-areas", focusAreaController.GetAllFocusAreas)
	strategyRoutes.GET("/focus-areas/:id", focusAreaController.GetFocusAreaDetails)
	strategyRoutes.POST("/focus-areas", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), focusAreaController.CreateFocusArea)
	strategyRoutes.PUT("/focus-areas/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), focusAreaController.UpdateFocusArea)
	strategyRoutes.POST("/focus-areas/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), focusAreaController.DeleteFocusAreas)
	strategyRoutes.PUT("/focus-areas/activate", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), focusAreaController.ActivateFocusAreas)
	strategyRoutes.PUT("/focus-areas/deactivate", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), focusAreaController.DeactivateFocusAreas)

	// focus area
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/focus-areas", strategyProfileFocusAreaController.GetAllStrategyProfileFocusAreas)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/focus-areas", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileFocusAreaController.CreateStrategyProfileFocusArea)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/focus-areas/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileFocusAreaController.DeleteStrategyProfileFocusArea)

	// objective
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/strategic-objectives", strategyProfileStrategicObjectiveController.GetAllStrategyProfileStrategicObjectives)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileStrategicObjectiveController.CreateStrategyProfileStrategicObjective)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/strategic-objectives/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileStrategicObjectiveController.DeleteStrategyProfileStrategicObjectives)

	// strategic intiative
	strategyRoutes.GET("/strategic-intiatives", strategicIntiativeController.GetAllStrategicIntiatives)
	strategyRoutes.GET("/strategic-intiatives/:id", strategicIntiativeController.GetStrategicIntiativeDetails)
	strategyRoutes.POST("/strategic-intiatives", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicIntiativeController.CreateStrategicIntiative)
	strategyRoutes.PUT("/strategic-intiatives/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicIntiativeController.UpdateStrategicIntiative)
	strategyRoutes.POST("/strategic-intiatives/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicIntiativeController.DeleteStrategicIntiatives)

	// strategic objective associated item
	strategyRoutes.GET("/strategic-objectives/:id/clauses", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.GetAllClauses)
	strategyRoutes.PUT("/strategic-objectives/:id/clauses", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.UpdateClauses)
	//kpi
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/kpis", strategyProfileKpiController.GetAllStrategyProfileKpis)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/kpis", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileKpiController.CreateStrategyProfileKpi)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/kpis/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileKpiController.DeleteStrategyProfileKpi)

	strategyRoutes.GET("/strategic-objectives/:id/risks", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategicObjectiveAssociatedItemController.GetAllStrategicObjectiveRisks)
	strategyRoutes.POST("/strategic-objectives/:id/risks", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.CreateStrategicObjectiveRisk)
	strategyRoutes.DELETE("/strategic-objectives/:id/risks/:riskId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.DeleteStrategicObjectiveRisk)
	// strategic initiative
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/strategic-intiatives", strategyProfileStrategicIntiativeController.GetAllStrategyProfileStrategicIntiatives)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/strategic-intiatives", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileStrategicIntiativeController.CreateStrategyProfileStrategicIntiative)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/strategic-intiatives/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileStrategicIntiativeController.DeleteStrategyProfileStrategicIntiative)

	//project
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/projects", strategyProfileProjectController.GetAllStrategyProfileProjects)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/projects", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectController.CreateStrategyProfileProject)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/projects/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectController.DeleteStrategyProfileProject)

	// project milestone
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/project-milestones", strategyProfileProjectMilestoneController.GetAllStrategyProfileProjectMilestones)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/project-milestones", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectMilestoneController.CreateStrategyProfileProjectMilestone)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/project-milestones/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectMilestoneController.DeleteStrategyProfileProjectMilestone)

	// project action plan
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/project-action-plans", strategyProfileProjectActionPlanController.GetAllStrategyProfileProjectActionPlans)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/project-action-plans", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectActionPlanController.CreateStrategyProfileProjectActionPlan)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/project-action-plans/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectActionPlanController.DeleteStrategyProfileProjectActionPlan)

}
