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

	swotPestelQuestionGroupController := app.StrategyModule.SwotPestelQuestionGroupController
	swotPestelQuestionController := app.StrategyModule.SwotPestelQuestionController
	swotPestelCategoryController := app.StrategyModule.SwotPestelCategoryController
	swotPestelQuestionFactorController := app.StrategyModule.SwotPestelQuestionFactorController

	strategyRoutes.GET("/swot-question-groups", swotPestelQuestionGroupController.GetAllSwotPestelQuestionGroups)
	strategyRoutes.GET("/swot-question-groups/:id", swotPestelQuestionGroupController.GetSwotPestelQuestionGroupDetails)
	strategyRoutes.POST("/swot-question-groups", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionGroupController.CreateSwotPestelQuestionGroup)
	strategyRoutes.PUT("/swot-question-groups/:id", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionGroupController.UpdateSwotPestelQuestionGroup)
	strategyRoutes.POST("/swot-question-groups/delete", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionGroupController.DeleteSwotPestelQuestionGroups)
	strategyRoutes.PUT("/swot-question-groups/activate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionGroupController.ActivateSwotPestelQuestionGroups)
	strategyRoutes.PUT("/swot-question-groups/deactivate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionGroupController.DeactivateSwotPestelQuestionGroups)

	strategyRoutes.GET("/swot-questions", swotPestelQuestionController.GetAllSwotPestelQuestions)
	strategyRoutes.GET("/swot-questions/:questionId", swotPestelQuestionController.GetSwotPestelQuestionDetails)
	strategyRoutes.POST("/swot-questions", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionController.CreateSwotPestelQuestion)
	strategyRoutes.PUT("/swot-questions/:questionId", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionController.UpdateSwotPestelQuestion)
	strategyRoutes.POST("/swot-questions/delete", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionController.DeleteSwotPestelQuestions)
	strategyRoutes.PUT("/swot-questions/activate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionController.ActivateSwotPestelQuestions)
	strategyRoutes.PUT("/swot-questions/deactivate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionController.DeactivateSwotPestelQuestions)

	strategyRoutes.GET("/pestel-question-groups", swotPestelQuestionGroupController.GetAllSwotPestelQuestionGroups)
	strategyRoutes.GET("/pestel-question-groups/:id", swotPestelQuestionGroupController.GetSwotPestelQuestionGroupDetails)
	strategyRoutes.POST("/pestel-question-groups", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionGroupController.CreateSwotPestelQuestionGroup)
	strategyRoutes.PUT("/pestel-question-groups/:id", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionGroupController.UpdateSwotPestelQuestionGroup)
	strategyRoutes.POST("/pestel-question-groups/delete", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionGroupController.DeleteSwotPestelQuestionGroups)
	strategyRoutes.PUT("/pestel-question-groups/activate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionGroupController.ActivateSwotPestelQuestionGroups)
	strategyRoutes.PUT("/pestel-question-groups/deactivate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionGroupController.DeactivateSwotPestelQuestionGroups)

	strategyRoutes.GET("/pestel-questions", swotPestelQuestionController.GetAllSwotPestelQuestions)
	strategyRoutes.GET("/pestel-questions/:questionId", swotPestelQuestionController.GetSwotPestelQuestionDetails)
	strategyRoutes.POST("/pestel-questions", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionController.CreateSwotPestelQuestion)
	strategyRoutes.PUT("/pestel-questions/:questionId", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionController.UpdateSwotPestelQuestion)
	strategyRoutes.POST("/pestel-questions/delete", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionController.DeleteSwotPestelQuestions)
	strategyRoutes.PUT("/pestel-questions/activate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionController.ActivateSwotPestelQuestions)
	strategyRoutes.PUT("/pestel-questions/deactivate", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionController.DeactivateSwotPestelQuestions)

	strategyRoutes.PUT("/swot-questions/:questionId/answer", middlewares.RolePermissionMiddleware(db, "swot", "modify"), swotPestelQuestionController.UpdateAnswer)
	strategyRoutes.PUT("/swot-questions/:questionId/answer-invitees", middlewares.RolePermissionMiddleware(db, "swot", "modify"), swotPestelQuestionController.CreateInvitees)

	//Swot Pestel Category
	strategyRoutes.GET("/swot-pestel-categories", swotPestelCategoryController.GetAllSwotPestelCategories)
	strategyRoutes.GET("/swot-pestel-categories/:id", swotPestelCategoryController.GetSwotPestelCategoryDetails)

	// swot pestel question factor
	strategyRoutes.GET("/swot-questions/:questionId/factors", swotPestelQuestionFactorController.GetAllSwotPestelQuestionFactors)
	strategyRoutes.GET("/swot-questions/:questionId/factors/:id", swotPestelQuestionFactorController.GetSwotPestelQuestionFactorDetails)
	strategyRoutes.POST("/swot-questions/:questionId/factors", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionFactorController.CreateSwotPestelQuestionFactor)
	strategyRoutes.PUT("/swot-questions/:questionId/factors/:id", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionFactorController.UpdateSwotPestelQuestionFactor)
	strategyRoutes.POST("/swot-questions/:questionId/factors/delete", middlewares.RolePermissionMiddleware(db, "strategy_setting", "modify"), swotPestelQuestionFactorController.DeleteSwotPestelQuestionFactors)
}
