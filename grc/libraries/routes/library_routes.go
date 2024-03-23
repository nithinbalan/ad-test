package routes

import (
	libraryController "jethur/grc/libraries/controllers"
	libraryRepo "jethur/grc/libraries/repositories"
	libraryService "jethur/grc/libraries/services"
	activityLogRepository "jethur/grc/logs/repositories"
	activityLogServie "jethur/grc/logs/services"
	middlewares "jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	libraryRoutes := router.Group("/")

	libraryRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := activityLogRepository.NewActivityLogRepository(db)
	activityLogSvc := activityLogServie.NewActivityLogService(activityLogRepo)

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

}
