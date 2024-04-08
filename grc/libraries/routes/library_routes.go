package routes

import (
	libraryController "jethur/grc/libraries/controllers"
	libraryRepo "jethur/grc/libraries/repositories"
	libraryService "jethur/grc/libraries/services"
	activityLogRepository "jethur/grc/logs/repositories"
	activityLogServie "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	middlewares "jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	libraryRoutes := router.Group("/")

	libraryRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := activityLogRepository.NewActivityLogRepository(db)
	activityLogSvc := activityLogServie.NewActivityLogService(activityLogRepo)

	reviewFrequencyRepo := masterRepositories.NewReviewFrequencyRepository(db)
	reviewFrequencySvc := masterServices.NewReviewFrequencyService(reviewFrequencyRepo)

	//Question Groups
	questionGroupRepository := libraryRepo.NewQuestionGroupRepository(db)
	questionGroupService := libraryService.NewQuestionGroupService(questionGroupRepository, activityLogSvc)
	questionGroupController := libraryController.NewQuestionGroupController(questionGroupService)

	libraryRoutes.GET("/question-groups", questionGroupController.GetAllQuestionGroups)
	libraryRoutes.POST("/question-groups", questionGroupController.CreateQuestionGroup)
	libraryRoutes.GET("/question-groups/:id", questionGroupController.GetQuestionGroupDetails)
	libraryRoutes.PUT("/question-groups/:id", questionGroupController.UpdateQuestionGroup)
	libraryRoutes.POST("/question-groups/delete", questionGroupController.DeleteQuestionGroups)
	libraryRoutes.PUT("/question-groups/activate", questionGroupController.ActivateQuestionGroups)
	libraryRoutes.PUT("/question-groups/deactivate", questionGroupController.DeactivateQuestionGroups)

	//Question Library
	questionLibraryRepository := libraryRepo.NewQuestionLibraryRepository(db)
	questionLibraryService := libraryService.NewQuestionLibraryService(questionLibraryRepository, activityLogSvc)
	questionLibraryController := libraryController.NewQuestionLibraryController(questionLibraryService)

	libraryRoutes.GET("/question-library", questionLibraryController.GetAllQuestionLibrary)
	libraryRoutes.POST("/question-library", questionLibraryController.CreateQuestionLibrary)
	libraryRoutes.GET("/question-library/:id", questionLibraryController.GetQuestionLibraryDetails)
	libraryRoutes.PUT("/question-library/:id", questionLibraryController.UpdateQuestionLibrary)
	libraryRoutes.POST("/question-library/delete", questionLibraryController.DeleteQuestionLibraries)
	libraryRoutes.PUT("/question-library/activate", questionLibraryController.ActivateQuestionLibraries)
	libraryRoutes.PUT("/question-library/deactivate", questionLibraryController.DeactivateQuestionLibraries)

	//KRI Categories
	KRICategoryRepository := libraryRepo.NewKRICategoryRepository(db)
	KRICategoryService := libraryService.NewKRICategoryService(KRICategoryRepository, activityLogSvc)
	KRICategoryController := libraryController.NewKRICategoryController(KRICategoryService)

	libraryRoutes.GET("/kri-categories", KRICategoryController.GetAllKRICategories)
	libraryRoutes.POST("/kri-categories", KRICategoryController.CreateKRICategory)
	libraryRoutes.GET("/kri-categories/:id", KRICategoryController.GetKRICategoryDetails)
	libraryRoutes.PUT("/kri-categories/:id", KRICategoryController.UpdateKRICategory)
	libraryRoutes.POST("/kri-categories/delete", KRICategoryController.DeleteKRICategories)
	libraryRoutes.PUT("/kri-categories/activate", KRICategoryController.ActivateKRICategories)
	libraryRoutes.PUT("/kri-categories/deactivate", KRICategoryController.DeactivateKRICategories)

	//KRI Library
	KRILibraryRepository := libraryRepo.NewKRILibraryRepository(db)
	KRILibraryService := libraryService.NewKRILibraryService(KRILibraryRepository, KRICategoryService, reviewFrequencySvc, activityLogSvc)
	KRILibraryController := libraryController.NewKRILibraryController(KRILibraryService)

	libraryRoutes.GET("/kri-library", KRILibraryController.GetAllKRILibrary)
	libraryRoutes.POST("/kri-library", KRILibraryController.CreateKRILibrary)
	libraryRoutes.GET("/kri-library/:id", KRILibraryController.GetKRILibraryDetails)
	libraryRoutes.PUT("/kri-library/:id", KRILibraryController.UpdateKRILibrary)
	libraryRoutes.POST("/kri-library/delete", KRILibraryController.DeleteKRILibraries)
	libraryRoutes.PUT("/kri-library/activate", KRILibraryController.ActivateKRILibraries)
	libraryRoutes.PUT("/kri-library/deactivate", KRILibraryController.DeactivateKRILibraries)

}
