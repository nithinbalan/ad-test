package routes

import (
	questionGroupController "jethur/grc/libraries/controllers"
	questionGroupRepo "jethur/grc/libraries/repositories"
	questionGroupSvc "jethur/grc/libraries/services"
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
	questionGroupRepository := questionGroupRepo.NewQuestionGroupRepository(db)
	questionGroupService := questionGroupSvc.NewQuestionGroupService(questionGroupRepository, activityLogSvc)
	questionGroupController := questionGroupController.NewQuestionGroupController(questionGroupService)

	libraryRoutes.GET("/question-groups", questionGroupController.GetAllQuestionGroups)
	libraryRoutes.POST("/question-groups", questionGroupController.CreateQuestionGroup)
	libraryRoutes.GET("/question-groups/:id", questionGroupController.GetQuestionGroupDetails)
	libraryRoutes.PUT("/question-groups/:id", questionGroupController.UpdateQuestionGroup)
	libraryRoutes.POST("/question-groups/delete", questionGroupController.DeleteQuestionGroups)
	libraryRoutes.POST("/question-groups/activate", questionGroupController.ActivateQuestionGroups)
	libraryRoutes.POST("/question-groups/deactivate", questionGroupController.DeactivateQuestionGroups)

}
