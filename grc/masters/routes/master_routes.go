package routes

import (
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/masters/controllers"
	"jethur/grc/masters/repositories"
	"jethur/grc/masters/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	masterRoutes := router.Group("/")

	masterRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	//Review Frequency
	reviewFrequencyRepo := repositories.NewReviewFrequencyRepository(db)
	reviewFrequencySvc := services.NewReviewFrequencyService(reviewFrequencyRepo)
	reviewFrequencyController := controllers.NewReviewFrequencyController(reviewFrequencySvc)

	//Severity Level
	severityLevelRepo := repositories.NewSeverityLevelRepository(db)
	severityLevelSvc := services.NewSeverityLevelService(severityLevelRepo)
	severityLevelController := controllers.NewSeverityLevelController(severityLevelSvc)

	//Lifecycle Status
	lifecycleStatusRepo := repositories.NewLifecycleStatusRepository(db)
	lifecycleStatusSvc := services.NewLifecycleStatusService(lifecycleStatusRepo)
	lifecycleStatusController := controllers.NewLifecycleStatusController(lifecycleStatusSvc)

	//Question Groups
	questionGroupRepository := repositories.NewQuestionGroupRepository(db)
	questionGroupService := services.NewQuestionGroupService(questionGroupRepository, activityLogSvc)
	questionGroupController := controllers.NewQuestionGroupController(questionGroupService)

	//Question Library
	questionLibraryRepository := repositories.NewQuestionLibraryRepository(db)
	questionLibraryService := services.NewQuestionLibraryService(questionLibraryRepository, activityLogSvc)
	questionLibraryController := controllers.NewQuestionLibraryController(questionLibraryService)

	//Checklists
	checklistRepo := repositories.NewChecklistRepository(db)
	checklistSvc := services.NewChecklistService(checklistRepo, activityLogSvc)
	checklistController := controllers.NewChecklistController(checklistSvc)

	//Checklist Questions
	checklistQuestionRepo := repositories.NewChecklistQuestionRepository(db)
	checklistQuestionSvc := services.NewChecklistQuestionService(checklistQuestionRepo, activityLogSvc, questionGroupService, checklistSvc, questionLibraryService)
	checklistQuestionController := controllers.NewChecklistQuestionController(checklistQuestionSvc)

	// Action Plan Status
	actionPlanStatusRepo := repositories.NewActionPlanStatusRepository(db)
	actionPlanStatusSvc := services.NewActionPlanStatusService(actionPlanStatusRepo)
	actionPlanStatusController := controllers.NewActionPlanStatusController(actionPlanStatusSvc)

	//KRI Categories
	KRICategoryRepository := repositories.NewKRICategoryRepository(db)
	KRICategoryService := services.NewKRICategoryService(KRICategoryRepository, activityLogSvc)
	KRICategoryController := controllers.NewKRICategoryController(KRICategoryService)

	//KRI Library
	KRILibraryRepository := repositories.NewKRILibraryRepository(db)
	KRILibraryService := services.NewKRILibraryService(KRILibraryRepository, KRICategoryService, reviewFrequencySvc, activityLogSvc)
	KRILibraryController := controllers.NewKRILibraryController(KRILibraryService)

	masterRoutes.GET("/review-frequencies", reviewFrequencyController.GetAllReviewFrequencies)
	masterRoutes.GET("/review-frequencies/:id", reviewFrequencyController.GetReviewFrequencyDetails)

	masterRoutes.GET("/severity-levels", severityLevelController.GetAllSeverityLevels)
	masterRoutes.GET("/severity-levels/:id", severityLevelController.GetSeverityLevelDetails)

	masterRoutes.GET("/lifecycle-statuses", lifecycleStatusController.GetAllLifecycleStatuses)
	masterRoutes.GET("/lifecycle-statuses/:id", lifecycleStatusController.GetLifecycleStatusDetails)

	masterRoutes.GET("/checklists", checklistController.GetAllChecklists)
	masterRoutes.POST("/checklists", checklistController.CreateChecklist)
	masterRoutes.GET("/checklists/:id", checklistController.GetChecklistDetails)
	masterRoutes.PUT("/checklists/:id", checklistController.UpdateChecklist)
	masterRoutes.POST("/checklists/delete", checklistController.DeleteChecklists)
	masterRoutes.PUT("/checklists/activate", checklistController.ActivateChecklists)
	masterRoutes.PUT("/checklists/deactivate", checklistController.DeactivateChecklists)

	masterRoutes.GET("/checklist-questions", checklistQuestionController.GetAllChecklistQuestions)
	masterRoutes.POST("/checklist-questions", checklistQuestionController.CreateChecklistQuestion)
	masterRoutes.GET("/checklist-questions/:id", checklistQuestionController.GetChecklistQuestionDetails)
	masterRoutes.PUT("/checklist-questions/:id", checklistQuestionController.UpdateChecklistQuestion)
	masterRoutes.POST("/checklist-questions/delete", checklistQuestionController.DeleteChecklistQuestions)

	masterRoutes.GET("/action-plan-statuses", actionPlanStatusController.GetAllActionPlanStatuses)
	masterRoutes.GET("/action-plan-statuses/:id", actionPlanStatusController.GetActionPlanStatusDetails)

	masterRoutes.GET("/question-groups", questionGroupController.GetAllQuestionGroups)
	masterRoutes.POST("/question-groups", questionGroupController.CreateQuestionGroup)
	masterRoutes.GET("/question-groups/:id", questionGroupController.GetQuestionGroupDetails)
	masterRoutes.PUT("/question-groups/:id", questionGroupController.UpdateQuestionGroup)
	masterRoutes.POST("/question-groups/delete", questionGroupController.DeleteQuestionGroups)
	masterRoutes.PUT("/question-groups/activate", questionGroupController.ActivateQuestionGroups)
	masterRoutes.PUT("/question-groups/deactivate", questionGroupController.DeactivateQuestionGroups)

	masterRoutes.GET("/question-library", questionLibraryController.GetAllQuestionLibrary)
	masterRoutes.POST("/question-library", questionLibraryController.CreateQuestionLibrary)
	masterRoutes.GET("/question-library/:id", questionLibraryController.GetQuestionLibraryDetails)
	masterRoutes.PUT("/question-library/:id", questionLibraryController.UpdateQuestionLibrary)
	masterRoutes.POST("/question-library/delete", questionLibraryController.DeleteQuestionLibraries)
	masterRoutes.PUT("/question-library/activate", questionLibraryController.ActivateQuestionLibraries)
	masterRoutes.PUT("/question-library/deactivate", questionLibraryController.DeactivateQuestionLibraries)

	masterRoutes.GET("/kri-categories", KRICategoryController.GetAllKRICategories)
	masterRoutes.POST("/kri-categories", KRICategoryController.CreateKRICategory)
	masterRoutes.GET("/kri-categories/:id", KRICategoryController.GetKRICategoryDetails)
	masterRoutes.PUT("/kri-categories/:id", KRICategoryController.UpdateKRICategory)
	masterRoutes.POST("/kri-categories/delete", KRICategoryController.DeleteKRICategories)
	masterRoutes.PUT("/kri-categories/activate", KRICategoryController.ActivateKRICategories)
	masterRoutes.PUT("/kri-categories/deactivate", KRICategoryController.DeactivateKRICategories)

	masterRoutes.GET("/kri-library", KRILibraryController.GetAllKRILibrary)
	masterRoutes.POST("/kri-library", KRILibraryController.CreateKRILibrary)
	masterRoutes.GET("/kri-library/:id", KRILibraryController.GetKRILibraryDetails)
	masterRoutes.PUT("/kri-library/:id", KRILibraryController.UpdateKRILibrary)
	masterRoutes.POST("/kri-library/delete", KRILibraryController.DeleteKRILibraries)
	masterRoutes.PUT("/kri-library/activate", KRILibraryController.ActivateKRILibraries)
	masterRoutes.PUT("/kri-library/deactivate", KRILibraryController.DeactivateKRILibraries)
}
