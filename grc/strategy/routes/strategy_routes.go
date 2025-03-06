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

	swotQuestionGroupController := app.StrategyModule.SwotQuestionGroupController
	swotQuestionController := app.StrategyModule.SwotQuestionController

	strategyRoutes.GET("/swot-question-groups", swotQuestionGroupController.GetAllSwotQuestionGroups)
	strategyRoutes.GET("/swot-question-groups/:id", swotQuestionGroupController.GetSwotQuestionGroupDetails)
	strategyRoutes.POST("/swot-question-groups", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotQuestionGroupController.CreateSwotQuestionGroup)
	strategyRoutes.PUT("/swot-question-groups/:id", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotQuestionGroupController.UpdateSwotQuestionGroup)
	strategyRoutes.POST("/swot-question-groups/delete", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotQuestionGroupController.DeleteSwotQuestionGroups)
	strategyRoutes.PUT("/swot-question-groups/activate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotQuestionGroupController.ActivateSwotQuestionGroups)
	strategyRoutes.PUT("/swot-question-groups/deactivate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotQuestionGroupController.DeactivateSwotQuestionGroups)

	strategyRoutes.GET("/swot-questions", swotQuestionController.GetAllSwotQuestions)
	strategyRoutes.GET("/swot-questions/:id", swotQuestionController.GetSwotQuestionDetails)
	strategyRoutes.POST("/swot-questions", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotQuestionController.CreateSwotQuestion)
	strategyRoutes.PUT("/swot-questions/:id", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotQuestionController.UpdateSwotQuestion)
	strategyRoutes.POST("/swot-questions/delete", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotQuestionController.DeleteSwotQuestions)
	strategyRoutes.PUT("/swot-questions/activate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotQuestionController.ActivateSwotQuestions)
	strategyRoutes.PUT("/swot-questions/deactivate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotQuestionController.DeactivateSwotQuestions)
}
