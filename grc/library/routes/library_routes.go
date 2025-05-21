package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	masterRoutes := router.Group("/")

	masterRoutes.Use(middlewares.AuthMiddleware())

	questionLibraryController := app.LibraryModule.QuestionLibraryController
	KRILibraryController := app.LibraryModule.KriLibraryController
	RiskLibraryController := app.LibraryModule.RiskLibraryController

	masterRoutes.GET("/question-library", questionLibraryController.GetAllQuestionLibrary)
	masterRoutes.POST("/question-library", questionLibraryController.CreateQuestionLibrary)
	masterRoutes.GET("/question-library/:id", questionLibraryController.GetQuestionLibraryDetails)
	masterRoutes.PUT("/question-library/:id", questionLibraryController.UpdateQuestionLibrary)
	masterRoutes.POST("/question-library/delete", questionLibraryController.DeleteQuestionLibraries)
	masterRoutes.PUT("/question-library/activate", questionLibraryController.ActivateQuestionLibraries)
	masterRoutes.PUT("/question-library/deactivate", questionLibraryController.DeactivateQuestionLibraries)

	masterRoutes.GET("/kri-library/export", KRILibraryController.GetKRILibraryExport)
	masterRoutes.GET("/kri-library", KRILibraryController.GetAllKRILibrary)
	masterRoutes.GET("/kri-library/category-count", KRILibraryController.KRICountByCategory)
	masterRoutes.POST("/kri-library", KRILibraryController.CreateKRILibrary)
	masterRoutes.GET("/kri-library/:id", KRILibraryController.GetKRILibraryDetails)
	masterRoutes.PUT("/kri-library/:id", KRILibraryController.UpdateKRILibrary)
	masterRoutes.POST("/kri-library/delete", KRILibraryController.DeleteKRILibraries)
	masterRoutes.PUT("/kri-library/activate", KRILibraryController.ActivateKRILibraries)
	masterRoutes.PUT("/kri-library/deactivate", KRILibraryController.DeactivateKRILibraries)

	masterRoutes.GET("/risk-library/export", RiskLibraryController.GetRiskLibraryExport)
	masterRoutes.GET("/risk-library", RiskLibraryController.GetAllRiskLibrary)
	masterRoutes.GET("/risk-library/category-count", RiskLibraryController.RiskCountByCategory)
	masterRoutes.GET("/risk-library/likelihood-count", RiskLibraryController.RiskCountByLikelihood)
	masterRoutes.GET("/risk-library/impact-count", RiskLibraryController.RiskCountByImpact)
}
