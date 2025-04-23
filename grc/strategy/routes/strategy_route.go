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
	strategyProfileObjectiveController := app.StrategyModule.StrategyProfileObjectiveController
	strategyProfileInitiativeController := app.StrategyModule.StrategyProfileInitiativeController
	strategyProfileProjectController := app.StrategyModule.StrategyProfileProjectController
	strategyProfileMilestoneController := app.StrategyModule.StrategyProfileMilestoneController
	strategyProfileActionPlanController := app.StrategyModule.StrategyProfileActionPlanController
	strategyProfileActionPlanStatusUpdateController := app.StrategyModule.StrategyProfileActionPlanStatusUpdateController
	strategyProfileActionPlanCommentController := app.StrategyModule.StrategyProfileActionPlanCommentController
	strategyProfileSettingController := app.StrategyModule.StrategyProfileSettingController
	strategyProfileMappingController := app.StrategyModule.StrategyProfileMappingController

	// strategy profile
	strategyRoutes.GET("/strategy-profiles", strategyProfileController.GetAllStrategyProfiles)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId", strategyProfileController.GetStrategyProfileDetails)
	strategyRoutes.POST("/strategy-profiles", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileController.CreateStrategyProfile)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileController.UpdateStrategyProfile)
	strategyRoutes.POST("/strategy-profiles/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileController.DeleteStrategyProfiles)

	// focus area
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/focus-areas", strategyProfileFocusAreaController.GetAllStrategyProfileFocusAreas)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/focus-areas/:focusAreaId", strategyProfileFocusAreaController.GetStrategyProfileFocusAreaDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/focus-areas", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileFocusAreaController.CreateStrategyProfileFocusArea)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/focus-areas/:focusAreaId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileFocusAreaController.UpdateStrategyProfileFocusArea)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/focus-areas/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileFocusAreaController.DeleteStrategyProfileFocusAreas)

	// goal
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/goals", strategyProfileGoalController.GetAllStrategyProfileGoals)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/goals/:goalId", strategyProfileGoalController.GetStrategyProfileGoalDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/goals", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileGoalController.CreateStrategyProfileGoal)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/goals/:goalId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileGoalController.UpdateStrategyProfileGoal)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/goals/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileGoalController.DeleteStrategyProfileGoals)

	// objective
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/objectives", strategyProfileObjectiveController.GetAllStrategyProfileObjectives)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/objectives/:objectiveId", strategyProfileObjectiveController.GetStrategyProfileObjectiveDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/objectives", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileObjectiveController.CreateStrategyProfileObjective)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/objectives/:objectiveId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileObjectiveController.UpdateStrategyProfileObjective)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/objectives/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileObjectiveController.DeleteStrategyProfileObjectives)

	// initiative
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/initiatives", strategyProfileInitiativeController.GetAllStrategyProfileInitiatives)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/initiatives/:initiativeId", strategyProfileInitiativeController.GetStrategyProfileInitiativeDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/initiatives", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileInitiativeController.CreateStrategyProfileInitiative)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/initiatives/:initiativeId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileInitiativeController.UpdateStrategyProfileInitiative)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/initiatives/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileInitiativeController.DeleteStrategyProfileInitiatives)

	// project
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/projects", strategyProfileProjectController.GetAllStrategyProfileProjects)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/projects/:projectId", strategyProfileProjectController.GetStrategyProfileProjectDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/projects", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectController.CreateStrategyProfileProject)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/projects/:projectId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectController.UpdateStrategyProfileProject)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/projects/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileProjectController.DeleteStrategyProfileProjects)

	// milestone
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/milestones", strategyProfileMilestoneController.GetAllStrategyProfileMilestones)
	strategyRoutes.GET("/strategy-profiles/:strategyProfileId/milestones/:milestoneId", strategyProfileMilestoneController.GetStrategyProfileMilestoneDetails)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/milestones", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileMilestoneController.CreateStrategyProfileMilestone)
	strategyRoutes.PUT("/strategy-profiles/:strategyProfileId/milestones/:milestoneId", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileMilestoneController.UpdateStrategyProfileMilestone)
	strategyRoutes.POST("/strategy-profiles/:strategyProfileId/milestones/delete", middlewares.RolePermissionMiddleware(db, "strategy_profile", "modify"), strategyProfileMilestoneController.DeleteStrategyProfileMilestones)

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
}
