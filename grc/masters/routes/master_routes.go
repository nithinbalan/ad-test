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
	currencyController := app.MasterModule.CurrencyController
	timezoneController := app.MasterModule.TimezoneController
	unitController := app.MasterModule.UnitController
	toolTypeController := app.MasterModule.ToolTypeController
	expenseCategoryController := app.MasterModule.ExpenseCategoryController
	sourceController := app.MasterModule.SourceController

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

	masterRoutes.GET("/currencies", currencyController.GetAllCurrencies)
	masterRoutes.GET("/currencies/:id", currencyController.GetCurrencyDetails)

	masterRoutes.GET("/timezones", timezoneController.GetAllTimezones)
	masterRoutes.GET("/timezones/:id", timezoneController.GetTimezoneDetails)

	masterRoutes.GET("/units", unitController.GetAllUnits)
	masterRoutes.POST("/units", unitController.CreateUnit)
	masterRoutes.GET("/units/:id", unitController.GetUnitDetails)
	masterRoutes.PUT("/units/:id", unitController.UpdateUnit)
	masterRoutes.POST("/units/delete", unitController.DeleteUnits)
	masterRoutes.PUT("/units/activate", unitController.ActivateUnits)
	masterRoutes.PUT("/units/deactivate", unitController.DeactivateUnits)

	masterRoutes.GET("/tool-types", toolTypeController.GetAllToolTypes)
	masterRoutes.POST("/tool-types", toolTypeController.CreateToolType)
	masterRoutes.GET("/tool-types/:id", toolTypeController.GetToolTypeDetails)
	masterRoutes.PUT("/tool-types/:id", toolTypeController.UpdateToolType)
	masterRoutes.POST("/tool-types/delete", toolTypeController.DeleteToolTypes)
	masterRoutes.PUT("/tool-types/activate", toolTypeController.ActivateToolTypes)
	masterRoutes.PUT("/tool-types/deactivate", toolTypeController.DeactivateToolTypes)

	masterRoutes.GET("/expense-categories", expenseCategoryController.GetAllExpenseCategories)
	masterRoutes.POST("/expense-categories", expenseCategoryController.CreateExpenseCategory)
	masterRoutes.GET("/expense-categories/:id", expenseCategoryController.GetExpenseCategoryDetails)
	masterRoutes.PUT("/expense-categories/:id", expenseCategoryController.UpdateExpenseCategory)
	masterRoutes.POST("/expense-categories/delete", expenseCategoryController.DeleteExpenseCategories)
	masterRoutes.PUT("/expense-categories/activate", expenseCategoryController.ActivateExpenseCategories)
	masterRoutes.PUT("/expense-categories/deactivate", expenseCategoryController.DeactivateExpenseCategories)

	masterRoutes.GET("/sources", sourceController.GetAllSources)
	masterRoutes.GET("/sources/:id", sourceController.GetSourceDetails)

}
