package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	masterRoutes := router.Group("/")

	masterRoutes.Use(middlewares.AuthMiddleware())

	reviewFrequencyController := app.MasterModule.ReviewFrequencyController
	severityLevelController := app.MasterModule.SeverityLevelController
	lifecycleStatusController := app.MasterModule.LifecycleStatusController
	questionGroupController := app.MasterModule.QuestionGroupController
	checklistController := app.MasterModule.ChecklistController
	checklistQuestionController := app.MasterModule.ChecklistQuestionController
	actionPlanStatusController := app.MasterModule.AtionPlanStatusController

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
}
