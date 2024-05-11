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

	masterRoutes.GET("/question-library", questionLibraryController.GetAllQuestionLibrary)
	masterRoutes.POST("/question-library", questionLibraryController.CreateQuestionLibrary)
	masterRoutes.GET("/question-library/:id", questionLibraryController.GetQuestionLibraryDetails)
	masterRoutes.PUT("/question-library/:id", questionLibraryController.UpdateQuestionLibrary)
	masterRoutes.POST("/question-library/delete", questionLibraryController.DeleteQuestionLibraries)
	masterRoutes.PUT("/question-library/activate", questionLibraryController.ActivateQuestionLibraries)
	masterRoutes.PUT("/question-library/deactivate", questionLibraryController.DeactivateQuestionLibraries)

	masterRoutes.GET("/kri-library", KRILibraryController.GetAllKRILibrary)
	masterRoutes.POST("/kri-library", KRILibraryController.CreateKRILibrary)
	masterRoutes.GET("/kri-library/:id", KRILibraryController.GetKRILibraryDetails)
	masterRoutes.PUT("/kri-library/:id", KRILibraryController.UpdateKRILibrary)
	masterRoutes.POST("/kri-library/delete", KRILibraryController.DeleteKRILibraries)
	masterRoutes.PUT("/kri-library/activate", KRILibraryController.ActivateKRILibraries)
	masterRoutes.PUT("/kri-library/deactivate", KRILibraryController.DeactivateKRILibraries)
}
