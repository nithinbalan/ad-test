package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	organizationalContextRoutes := router.Group("/")

	organizationalContextRoutes.Use(middlewares.AuthMiddleware())

	db := app.DB

	swotPestelQuestionGroupController := app.OrganizationalContextModule.SwotPestelQuestionGroupController
	swotPestelQuestionController := app.OrganizationalContextModule.SwotPestelQuestionController
	swotPestelCategoryController := app.OrganizationalContextModule.SwotPestelCategoryController
	swotPestelQuestionFactorController := app.OrganizationalContextModule.SwotPestelQuestionFactorController

	organizationalContextRoutes.GET("/swot-question-groups", swotPestelQuestionGroupController.GetAllSwotPestelQuestionGroups)
	organizationalContextRoutes.GET("/swot-question-groups/:id", swotPestelQuestionGroupController.GetSwotPestelQuestionGroupDetails)
	organizationalContextRoutes.POST("/swot-question-groups", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.CreateSwotPestelQuestionGroup)
	organizationalContextRoutes.PUT("/swot-question-groups/:id", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.UpdateSwotPestelQuestionGroup)
	organizationalContextRoutes.POST("/swot-question-groups/delete", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.DeleteSwotPestelQuestionGroups)
	organizationalContextRoutes.PUT("/swot-question-groups/activate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.ActivateSwotPestelQuestionGroups)
	organizationalContextRoutes.PUT("/swot-question-groups/deactivate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.DeactivateSwotPestelQuestionGroups)

	organizationalContextRoutes.GET("/swot-questions", swotPestelQuestionController.GetAllSwotPestelQuestions)
	organizationalContextRoutes.GET("/swot-questions/:questionId", swotPestelQuestionController.GetSwotPestelQuestionDetails)
	organizationalContextRoutes.POST("/swot-questions", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.CreateSwotPestelQuestion)
	organizationalContextRoutes.PUT("/swot-questions/:questionId", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.UpdateSwotPestelQuestion)
	organizationalContextRoutes.POST("/swot-questions/delete", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.DeleteSwotPestelQuestions)
	organizationalContextRoutes.PUT("/swot-questions/activate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.ActivateSwotPestelQuestions)
	organizationalContextRoutes.PUT("/swot-questions/deactivate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.DeactivateSwotPestelQuestions)

	organizationalContextRoutes.GET("/pestel-question-groups", swotPestelQuestionGroupController.GetAllSwotPestelQuestionGroups)
	organizationalContextRoutes.GET("/pestel-question-groups/:id", swotPestelQuestionGroupController.GetSwotPestelQuestionGroupDetails)
	organizationalContextRoutes.POST("/pestel-question-groups", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.CreateSwotPestelQuestionGroup)
	organizationalContextRoutes.PUT("/pestel-question-groups/:id", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.UpdateSwotPestelQuestionGroup)
	organizationalContextRoutes.POST("/pestel-question-groups/delete", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.DeleteSwotPestelQuestionGroups)
	organizationalContextRoutes.PUT("/pestel-question-groups/activate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.ActivateSwotPestelQuestionGroups)
	organizationalContextRoutes.PUT("/pestel-question-groups/deactivate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionGroupController.DeactivateSwotPestelQuestionGroups)

	organizationalContextRoutes.GET("/pestel-questions", swotPestelQuestionController.GetAllSwotPestelQuestions)
	organizationalContextRoutes.GET("/pestel-questions/:questionId", swotPestelQuestionController.GetSwotPestelQuestionDetails)
	organizationalContextRoutes.POST("/pestel-questions", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.CreateSwotPestelQuestion)
	organizationalContextRoutes.PUT("/pestel-questions/:questionId", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.UpdateSwotPestelQuestion)
	organizationalContextRoutes.POST("/pestel-questions/delete", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.DeleteSwotPestelQuestions)
	organizationalContextRoutes.PUT("/pestel-questions/activate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.ActivateSwotPestelQuestions)
	organizationalContextRoutes.PUT("/pestel-questions/deactivate", middlewares.RolePermissionMiddleware(db, "organizational_context_setting", "modify"), swotPestelQuestionController.DeactivateSwotPestelQuestions)

	organizationalContextRoutes.PUT("/swot-questions/:questionId/answer", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionController.UpdateAnswer)
	organizationalContextRoutes.PUT("/swot-questions/:questionId/answer-invitees", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionController.CreateInvitees)
	organizationalContextRoutes.PUT("/swot-questions/:questionId/answer-invitees/delete", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionController.DeleteInvitees)

	//Swot Pestel Category
	organizationalContextRoutes.GET("/swot-pestel-categories", swotPestelCategoryController.GetAllSwotPestelCategories)
	organizationalContextRoutes.GET("/swot-pestel-categories/:id", swotPestelCategoryController.GetSwotPestelCategoryDetails)

	// swot pestel question factor
	organizationalContextRoutes.GET("/swot-questions/factors", swotPestelQuestionFactorController.GetAllSwotPestelFactors)
	organizationalContextRoutes.GET("/swot-questions/:questionId/factors", swotPestelQuestionFactorController.GetAllSwotPestelQuestionFactors)
	organizationalContextRoutes.GET("/swot-questions/:questionId/factors/:id", swotPestelQuestionFactorController.GetSwotPestelQuestionFactorDetails)
	organizationalContextRoutes.POST("/swot-questions/:questionId/factors", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionFactorController.CreateSwotPestelQuestionFactor)
	organizationalContextRoutes.PUT("/swot-questions/:questionId/factors/:id", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionFactorController.UpdateSwotPestelQuestionFactor)
	organizationalContextRoutes.POST("/swot-questions/:questionId/factors/delete", middlewares.RolePermissionMiddleware(db, "swot_analysis", "modify"), swotPestelQuestionFactorController.DeleteSwotPestelQuestionFactors)

	organizationalContextRoutes.PUT("/pestel-questions/:questionId/answer", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionController.UpdateAnswer)
	organizationalContextRoutes.PUT("/pestel-questions/:questionId/answer-invitees", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionController.CreateInvitees)
	organizationalContextRoutes.PUT("/pestel-questions/:questionId/answer-invitees/delete", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionController.DeleteInvitees)

	// swot pestel question factor
	organizationalContextRoutes.GET("/pestel-questions/:questionId/factors", swotPestelQuestionFactorController.GetAllSwotPestelQuestionFactors)
	organizationalContextRoutes.GET("/pestel-questions/:questionId/factors/:id", swotPestelQuestionFactorController.GetSwotPestelQuestionFactorDetails)
	organizationalContextRoutes.POST("/pestel-questions/:questionId/factors", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionFactorController.CreateSwotPestelQuestionFactor)
	organizationalContextRoutes.PUT("/pestel-questions/:questionId/factors/:id", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionFactorController.UpdateSwotPestelQuestionFactor)
	organizationalContextRoutes.POST("/pestel-questions/:questionId/factors/delete", middlewares.RolePermissionMiddleware(db, "pestel_analysis", "modify"), swotPestelQuestionFactorController.DeleteSwotPestelQuestionFactors)
}
