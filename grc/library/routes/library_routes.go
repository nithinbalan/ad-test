package routes

import (
	"jethur/grc/library/controllers"
	"jethur/grc/library/repositories"
	"jethur/grc/library/services"
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	riskRepositories "jethur/grc/risk/repositories"
	riskServices "jethur/grc/risk/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	masterRoutes := router.Group("/")

	masterRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	reviewFrequencyRepo := masterRepositories.NewReviewFrequencyRepository(db)
	reviewFrequencySvc := masterServices.NewReviewFrequencyService(reviewFrequencyRepo)

	KRICategoryRepository := riskRepositories.NewKRICategoryRepository(db)
	KRICategoryService := riskServices.NewKRICategoryService(KRICategoryRepository, activityLogSvc)

	//Question Library
	questionLibraryRepository := repositories.NewQuestionLibraryRepository(db)
	questionLibraryService := services.NewQuestionLibraryService(questionLibraryRepository, activityLogSvc)
	questionLibraryController := controllers.NewQuestionLibraryController(questionLibraryService)

	//KRI Library
	KRILibraryRepository := repositories.NewKRILibraryRepository(db)
	KRILibraryService := services.NewKRILibraryService(KRILibraryRepository, KRICategoryService, reviewFrequencySvc, activityLogSvc)
	KRILibraryController := controllers.NewKRILibraryController(KRILibraryService)

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
