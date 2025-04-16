package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	strategyRoutes := router.Group("/")

	strategyRoutes.Use(middlewares.AuthMiddleware())

	db := app.DB

	strategyProfileController := app.StrategyModule.StrategyProfileController
	strategyProfileFocusAreaController := app.StrategyModule.StrategyProfileFocusAreaController
	strategyProfileGoalController := app.StrategyModule.StrategyProfileGoalController
	strategyProfileObjectiveController := app.StrategyModule.StrategyProfileObjectiveController
	strategyProfileInitiativeController := app.StrategyModule.StrategyProfileInitiativeController
	strategyProfileProjectController := app.StrategyModule.StrategyProfileProjectController

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
}
