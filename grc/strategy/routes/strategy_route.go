package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	strategyRoutes := router.Group("/")

	strategyRoutes.Use(middlewares.AuthMiddleware(app.DB))

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
	strategicObjectiveCategoryController := app.StrategyModule.StrategicObjectiveCategoryController
	strategicObjectiveController := app.StrategyModule.StrategicObjectiveController
	strategicObjectiveKPIController := app.StrategyModule.StrategicObjectiveKPIController
	strategicObjectivePeriodicReviewController := app.StrategyModule.StrategicObjectivePeriodicReviewController
	focusAreaController := app.StrategyModule.FocusAreaController
	strategicIntiativeController := app.StrategyModule.StrategicIntiativeController
	strategyInitiativeStatusUpdateController := app.StrategyModule.StrategyInitiativeStatusUpdateController
	strategicObjectiveAssociatedItemController := app.StrategyModule.StrategicObjectiveAssociatedItemController
	strategicInitiativeAssociatedItemController := app.StrategyModule.StrategicInitiativeAssociatedItemController
	strategyProfileKpiController := app.StrategyModule.StrategyProfileKpiController
	strategyProfileProjectActionPlanController := app.StrategyModule.StrategyProfileProjectActionPlanController
	strategicInitiativeStatusController := app.StrategyModule.StrategicInitiativeStatusController

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
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/status-updates", strategyProfileActionPlanStatusUpdateController.CreateStrategyProfileActionPlanStatusUpdate)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/status-updates/:id", strategyProfileActionPlanStatusUpdateController.UpdateStrategyProfileActionPlanStatusUpdate)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategyProfileActionPlanStatusUpdateController.GetLatestStrategyProfileActionPlanStatusUpdate)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/approve", strategyProfileActionPlanStatusUpdateController.ApproveStrategyProfileActionPlan)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/action-plans/:actionPlanId/revert", strategyProfileActionPlanStatusUpdateController.RevertStrategyProfileActionPlan)

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
	strategyRoutes.POST("/strategic-objectives", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicObjectiveController.CreateStrategicObjective)
	strategyRoutes.PUT("/strategic-objectives/:id", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicObjectiveController.UpdateStrategicObjective)
	strategyRoutes.POST("/strategic-objectives/delete", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicObjectiveController.DeleteStrategicObjectives)
	strategyRoutes.PUT("/strategic-objectives/activate", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicObjectiveController.ActivateStrategicObjectives)
	strategyRoutes.PUT("/strategic-objectives/deactivate", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicObjectiveController.DeactivateStrategicObjectives)
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
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/focus-areas/:id", strategyProfileFocusAreaController.GetStrategyProfileFocusAreaDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/focus-areas", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileFocusAreaController.CreateStrategyProfileFocusArea)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/focus-areas/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileFocusAreaController.DeleteStrategyProfileFocusArea)

	// objective
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/strategic-objectives", strategyProfileStrategicObjectiveController.GetAllStrategyProfileStrategicObjectives)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/strategic-objectives/:id", strategyProfileStrategicObjectiveController.GetStrategyProfileStrategicObjectiveDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileStrategicObjectiveController.CreateStrategyProfileStrategicObjective)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/strategic-objectives/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileStrategicObjectiveController.DeleteStrategyProfileStrategicObjectives)

	// strategic intiative
	strategyRoutes.GET("/strategic-intiatives", strategicIntiativeController.GetAllStrategicIntiatives)
	strategyRoutes.GET("/strategic-intiatives/export", strategicIntiativeController.GetStrategicInitiativesExport)
	strategyRoutes.GET("/strategic-intiatives/:id", strategicIntiativeController.GetStrategicIntiativeDetails)
	strategyRoutes.POST("/strategic-intiatives", middlewares.RolePermissionMiddleware(db, "strategy_intiative", "modify"), strategicIntiativeController.CreateStrategicIntiative)
	strategyRoutes.PUT("/strategic-intiatives/:id", middlewares.RolePermissionMiddleware(db, "strategy_intiative", "modify"), strategicIntiativeController.UpdateStrategicIntiative)
	strategyRoutes.POST("/strategic-intiatives/delete", middlewares.RolePermissionMiddleware(db, "strategy_intiative", "modify"), strategicIntiativeController.DeleteStrategicIntiatives)
	strategyRoutes.PUT("/strategic-intiatives/:id/status-update", middlewares.RolePermissionMiddleware(db, "strategy_intiative", "modify"), strategicIntiativeController.UpdateStrategicIntiativeStatusAndCompletion)
	strategyRoutes.GET("/strategic-intiatives/:id/status-updates", middlewares.RolePermissionMiddleware(db, "strategy_intiative", "view"), strategyInitiativeStatusUpdateController.GetAllStrategyInitiativeStatusUpdates)
	strategyRoutes.POST("/strategic-intiatives/:id/status-updates", middlewares.RolePermissionMiddleware(db, "strategy_intiative", "modify"), strategyInitiativeStatusUpdateController.CreateStrategyInitiativeStatusUpdate)
	strategyRoutes.PUT("/strategic-intiatives/:id/status-updates/:statusUpdateId", middlewares.RolePermissionMiddleware(db, "strategy_intiative", "modify"), strategyInitiativeStatusUpdateController.UpdateStrategyInitiativeStatusUpdate)
	strategyRoutes.GET("/strategic-intiatives/:id/status-updates/latest", middlewares.RolePermissionMiddleware(db, "strategy_intiative", "view"), strategyInitiativeStatusUpdateController.GetLatestStrategyInitiativeStatusUpdate)

	// strategic objective associated item
	strategyRoutes.GET("/strategic-objectives/:id/clauses", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.GetAllClauses)
	strategyRoutes.PUT("/strategic-objectives/:id/clauses", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.UpdateClauses)
	//kpi
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/kpis", strategyProfileKpiController.GetAllStrategyProfileKpis)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/kpis/:id", strategyProfileKpiController.GetStrategyProfileKpiDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/kpis", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileKpiController.CreateStrategyProfileKpi)
	strategyRoutes.PUT("/strategy-profile-kpis/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileKpiController.UpdateStrategyProfileKpiWeightage)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/kpis/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileKpiController.DeleteStrategyProfileKpi)

	strategyRoutes.GET("/strategic-objectives/:id/risks", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategicObjectiveAssociatedItemController.GetAllStrategicObjectiveRisks)
	strategyRoutes.POST("/strategic-objectives/:id/risks", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.CreateStrategicObjectiveRisk)
	strategyRoutes.DELETE("/strategic-objectives/:id/risks/:riskId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.DeleteStrategicObjectiveRisk)
	strategyRoutes.GET("/strategic-intiatives/:id/risks", middlewares.RolePermissionMiddleware(db, "strategy_intiative", "view"), strategicInitiativeAssociatedItemController.GetAllStrategicIntiativeRisks)
	strategyRoutes.PUT("/strategic-intiatives/:id/risks", middlewares.RolePermissionMiddleware(db, "strategy_intiative", "modify"), strategicInitiativeAssociatedItemController.CreateStrategicIntiativeRisk)
	strategyRoutes.DELETE("/strategic-intiatives/:id/risks/:riskId", middlewares.RolePermissionMiddleware(db, "strategy_intiative", "modify"), strategicInitiativeAssociatedItemController.DeleteStrategicIntiativeRisk)
	// strategic initiative
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/strategic-intiatives", strategyProfileStrategicIntiativeController.GetAllStrategyProfileStrategicIntiatives)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/strategic-intiatives/:id", strategyProfileStrategicIntiativeController.GetStrategyProfileStrategicIntiativeDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/strategic-intiatives", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileStrategicIntiativeController.CreateStrategyProfileStrategicIntiative)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/strategic-intiatives/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileStrategicIntiativeController.DeleteStrategyProfileStrategicIntiative)

	//project
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/projects", strategyProfileProjectController.GetAllStrategyProfileProjects)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/projects/:id", strategyProfileProjectController.GetStrategyProfileProjectDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/projects", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectController.CreateStrategyProfileProject)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/projects/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectController.DeleteStrategyProfileProject)

	strategyRoutes.GET("/strategic-objectives/:id/documents", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategicObjectiveAssociatedItemController.GetAllStrategicObjectiveDocuments)
	strategyRoutes.POST("/strategic-objectives/:id/documents", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.CreateStrategicObjectiveDocuments)
	strategyRoutes.DELETE("/strategic-objectives/:id/documents/:documentId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.DeleteStrategicObjectiveDocument)

	strategyRoutes.GET("/strategic-objectives/:id/controls", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategicObjectiveAssociatedItemController.GetAllStrategicObjectivControls)
	strategyRoutes.POST("/strategic-objectives/:id/controls", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.CreateStrategicObjectiveControl)
	strategyRoutes.DELETE("/strategic-objectives/:id/controls/:controlId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.DeleteStrategicObjectiveControl)

	strategyRoutes.GET("/strategic-objectives/:id/processes", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategicObjectiveAssociatedItemController.GetAllStrategicObjectiveProcesses)
	strategyRoutes.POST("/strategic-objectives/:id/processes", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.CreateStrategicObjectiveProcess)
	strategyRoutes.DELETE("/strategic-objectives/:id/processes/:processId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.DeleteStrategicObjectiveProcess)

	strategyRoutes.GET("/strategic-objectives/:id/assets", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategicObjectiveAssociatedItemController.GetAllStrategicObjectiveAssets)
	strategyRoutes.POST("/strategic-objectives/:id/assets", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.CreateStrategicObjectiveAsset)
	strategyRoutes.DELETE("/strategic-objectives/:id/assets/:assetId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.DeleteStrategicObjectiveAsset)

	strategyRoutes.GET("/strategic-objectives/:id/incidents", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategicObjectiveAssociatedItemController.GetAllStrategicObjectiveIncidents)
	strategyRoutes.POST("/strategic-objectives/:id/incidents", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.CreateStrategicObjectiveIncident)
	strategyRoutes.DELETE("/strategic-objectives/:id/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.DeleteStrategicObjectiveIncident)

	strategyRoutes.GET("/strategic-objectives/:id/associated-kpis", middlewares.RolePermissionMiddleware(db, "strategy_profile", "view"), strategicObjectiveAssociatedItemController.GetAllStrategicObjectiveKPI)
	strategyRoutes.POST("/strategic-objectives/:id/associated-kpis", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.CreateStrategicObjectiveKPI)
	strategyRoutes.DELETE("/strategic-objectives/:id/associated-kpis/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategicObjectiveAssociatedItemController.DeleteStrategicObjectiveKPI)
	// project milestone
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/project-milestones", strategyProfileProjectMilestoneController.GetAllStrategyProfileProjectMilestones)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/project-milestones/:id", strategyProfileProjectMilestoneController.GetStrategyProfileProjectMilestoneDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/project-milestones", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectMilestoneController.CreateStrategyProfileProjectMilestone)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/project-milestones/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectMilestoneController.DeleteStrategyProfileProjectMilestone)

	// project action plan
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/project-action-plans", strategyProfileProjectActionPlanController.GetAllStrategyProfileProjectActionPlans)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/project-action-plans/:id", strategyProfileProjectActionPlanController.GetStrategyProfileProjectActionPlanDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/project-action-plans", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectActionPlanController.CreateStrategyProfileProjectActionPlan)
	strategyRoutes.DELETE("/strategy-profiles/:strategyProfileId/project-action-plans/:id", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectActionPlanController.DeleteStrategyProfileProjectActionPlan)

	strategyRoutes.GET("/strategic-initiative-statuses", strategicInitiativeStatusController.GetAllStrategicInitiativeStatuses)
	strategyRoutes.GET("/strategic-initiative-statuses/:id", strategicInitiativeStatusController.GetStrategicInitiativeStatusDetails)
}
