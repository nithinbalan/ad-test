package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	masterRoutes := router.Group("/")

	masterRoutes.Use(middlewares.AuthMiddleware(app.DB))

	db := app.DB

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
	statusController := app.MasterModule.StatusController
	conformanceController := app.MasterModule.ConformanceController
	workflowStatusController := app.MasterModule.WorkflowStatusController
	integrationCategoryController := app.MasterModule.IntegrationCategoryController
	impactCategoryController := app.MasterModule.ImpactCategoryController
	timeUnitController := app.MasterModule.TimeUnitController
	reviewStatusController := app.MasterModule.ReviewStatusController
	applicabilityStatusController := app.MasterModule.ApplicabilityStatusController
	requiredDocumentTitleController := app.MasterModule.RequiredDocumentTitleController
	countryController := app.MasterModule.CountryController
	calculationTypeController := app.MasterModule.CalculationTypeController
	projectMilestoneStatusController := app.MasterModule.ProjectMilestoneStatusController
	licenseTypeController := app.MasterModule.LicenseTypeController
	workflowContextController := app.MasterModule.WorkflowContextController
	commentStatusController := app.MasterModule.CommentStatusController
	actionPlanMasterController := app.MasterModule.ActionPlanMasterController
	emailStatusController := app.MasterModule.EmailStatusController
	locationController := app.MasterModule.LocationController

	masterRoutes.GET("/review-frequencies", reviewFrequencyController.GetAllReviewFrequencies)
	masterRoutes.GET("/review-frequencies/:id", reviewFrequencyController.GetReviewFrequencyDetails)

	masterRoutes.GET("/severity-levels", severityLevelController.GetAllSeverityLevels)
	masterRoutes.GET("/severity-levels/:id", severityLevelController.GetSeverityLevelDetails)

	masterRoutes.GET("/lifecycle-statuses", lifecycleStatusController.GetAllLifecycleStatuses)
	masterRoutes.GET("/lifecycle-statuses/:id", lifecycleStatusController.GetLifecycleStatusDetails)

	masterRoutes.GET("/checklists", checklistController.GetAllChecklists)
	masterRoutes.POST("/checklists", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), checklistController.CreateChecklist)
	masterRoutes.GET("/checklists/:id", checklistController.GetChecklistDetails)
	masterRoutes.PUT("/checklists/:id", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), checklistController.UpdateChecklist)
	masterRoutes.POST("/checklists/delete", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), checklistController.DeleteChecklists)
	masterRoutes.PUT("/checklists/activate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), checklistController.ActivateChecklists)
	masterRoutes.PUT("/checklists/deactivate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), checklistController.DeactivateChecklists)

	masterRoutes.GET("/checklist-questions", checklistQuestionController.GetAllChecklistQuestions)
	masterRoutes.POST("/checklist-questions", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), checklistQuestionController.CreateChecklistQuestion)
	masterRoutes.GET("/checklist-questions/:id", checklistQuestionController.GetChecklistQuestionDetails)
	masterRoutes.PUT("/checklist-questions/:id", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), checklistQuestionController.UpdateChecklistQuestion)
	masterRoutes.POST("/checklist-questions/delete", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), checklistQuestionController.DeleteChecklistQuestions)

	masterRoutes.GET("/action-plan-statuses", actionPlanStatusController.GetAllActionPlanStatuses)
	masterRoutes.GET("/action-plan-statuses/:id", actionPlanStatusController.GetActionPlanStatusDetails)

	masterRoutes.GET("/question-groups", questionGroupController.GetAllQuestionGroups)
	masterRoutes.POST("/question-groups", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), questionGroupController.CreateQuestionGroup)
	masterRoutes.GET("/question-groups/:id", questionGroupController.GetQuestionGroupDetails)
	masterRoutes.PUT("/question-groups/:id", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), questionGroupController.UpdateQuestionGroup)
	masterRoutes.POST("/question-groups/delete", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), questionGroupController.DeleteQuestionGroups)
	masterRoutes.PUT("/question-groups/activate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), questionGroupController.ActivateQuestionGroups)
	masterRoutes.PUT("/question-groups/deactivate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), questionGroupController.DeactivateQuestionGroups)

	masterRoutes.GET("/currencies", currencyController.GetAllCurrencies)
	masterRoutes.GET("/currencies/:id", currencyController.GetCurrencyDetails)

	masterRoutes.GET("/timezones", timezoneController.GetAllTimezones)
	masterRoutes.GET("/timezones/:id", timezoneController.GetTimezoneDetails)

	masterRoutes.GET("/units", unitController.GetAllUnits)
	masterRoutes.POST("/units", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), unitController.CreateUnit)
	masterRoutes.GET("/units/:id", unitController.GetUnitDetails)
	masterRoutes.PUT("/units/:id", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), unitController.UpdateUnit)
	masterRoutes.POST("/units/delete", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), unitController.DeleteUnits)
	masterRoutes.PUT("/units/activate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), unitController.ActivateUnits)
	masterRoutes.PUT("/units/deactivate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), unitController.DeactivateUnits)

	masterRoutes.GET("/tool-types", toolTypeController.GetAllToolTypes)
	masterRoutes.POST("/tool-types", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), toolTypeController.CreateToolType)
	masterRoutes.GET("/tool-types/:id", toolTypeController.GetToolTypeDetails)
	masterRoutes.PUT("/tool-types/:id", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), toolTypeController.UpdateToolType)
	masterRoutes.POST("/tool-types/delete", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), toolTypeController.DeleteToolTypes)
	masterRoutes.PUT("/tool-types/activate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), toolTypeController.ActivateToolTypes)
	masterRoutes.PUT("/tool-types/deactivate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), toolTypeController.DeactivateToolTypes)

	masterRoutes.GET("/expense-categories", expenseCategoryController.GetAllExpenseCategories)
	masterRoutes.POST("/expense-categories", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), expenseCategoryController.CreateExpenseCategory)
	masterRoutes.GET("/expense-categories/:id", expenseCategoryController.GetExpenseCategoryDetails)
	masterRoutes.PUT("/expense-categories/:id", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), expenseCategoryController.UpdateExpenseCategory)
	masterRoutes.POST("/expense-categories/delete", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), expenseCategoryController.DeleteExpenseCategories)
	masterRoutes.PUT("/expense-categories/activate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), expenseCategoryController.ActivateExpenseCategories)
	masterRoutes.PUT("/expense-categories/deactivate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), expenseCategoryController.DeactivateExpenseCategories)

	masterRoutes.GET("/sources", sourceController.GetAllSources)
	masterRoutes.GET("/sources/:id", sourceController.GetSourceDetails)

	masterRoutes.GET("/statuses", statusController.GetAllStatuses)
	masterRoutes.GET("/statuses/:id", statusController.GetStatusDetails)

	masterRoutes.GET("/conformances", conformanceController.GetAllConformances)
	masterRoutes.GET("/conformances/:id", conformanceController.GetConformanceDetails)

	masterRoutes.GET("/workflow-statuses", workflowStatusController.GetAllWorkflowStatuses)
	masterRoutes.GET("/workflow-statuses/:id", workflowStatusController.GetWorkflowStatusDetails)

	masterRoutes.GET("/integration-categories", integrationCategoryController.GetAllIntegrationCategories)

	//Impact Category
	masterRoutes.GET("/impact-categories", impactCategoryController.GetAllImpactCategories)
	masterRoutes.GET("/impact-categories/:id", impactCategoryController.GetImpactCategoryDetails)
	masterRoutes.POST("/impact-categories", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), impactCategoryController.CreateImpactCategory)
	masterRoutes.PUT("/impact-categories/:id", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), impactCategoryController.UpdateImpactCategory)
	masterRoutes.POST("/impact-categories/delete", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), impactCategoryController.DeleteImpactCategories)
	masterRoutes.PUT("/impact-categories/activate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), impactCategoryController.ActivateImpactCategories)
	masterRoutes.PUT("/impact-categories/deactivate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), impactCategoryController.DeactivateImpactCategories)

	masterRoutes.GET("/time-units", timeUnitController.GetAllTimeUnits)
	masterRoutes.GET("/time-units/:id", timeUnitController.GetTimeUnitDetails)

	masterRoutes.GET("/review-statuses", reviewStatusController.GetAllReviewStatuses)
	masterRoutes.GET("/review-statuses/:id", reviewStatusController.GetReviewStatusDetails)

	masterRoutes.GET("/applicability-statuses", applicabilityStatusController.GetAllApplicabilityStatuses)
	masterRoutes.GET("/applicability-statuses/:id", applicabilityStatusController.GetApplicabilityStatusDetails)

	masterRoutes.GET("/required-document-titles", requiredDocumentTitleController.GetAllRequiredDocumentTitles)
	masterRoutes.POST("/required-document-titles", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), requiredDocumentTitleController.CreateRequiredDocumentTitle)
	masterRoutes.GET("/required-document-titles/:id", requiredDocumentTitleController.GetRequiredDocumentTitleDetails)
	masterRoutes.PUT("/required-document-titles/:id", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), requiredDocumentTitleController.UpdateRequiredDocumentTitle)
	masterRoutes.POST("/required-document-titles/delete", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), requiredDocumentTitleController.DeleteRequiredDocumentTitles)
	masterRoutes.PUT("/required-document-titles/activate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), requiredDocumentTitleController.ActivateRequiredDocumentTitles)
	masterRoutes.PUT("/required-document-titles/deactivate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), requiredDocumentTitleController.DeactivateRequiredDocumentTitles)

	masterRoutes.GET("/countries", countryController.GetAllCountries)
	masterRoutes.GET("/countries/:id", countryController.GetCountryDetails)

	masterRoutes.GET("/calculation-types", calculationTypeController.GetAllCalculationTypes)
	masterRoutes.GET("/calculation-types/:id", calculationTypeController.GetCalculationTypeDetails)

	masterRoutes.GET("/project-milestone-statuses", projectMilestoneStatusController.GetAllProjectMilestoneStatuses)
	masterRoutes.GET("/project-milestone-statuses/:id", projectMilestoneStatusController.GetProjectMilestoneStatusDetails)

	masterRoutes.GET("/license-types", licenseTypeController.GetAllLicenseTypes)
	masterRoutes.GET("/license-types/:id", licenseTypeController.GetLicenseTypeDetails)

	masterRoutes.GET("/workflow-contexts", workflowContextController.GetAllWorkflowContext)
	masterRoutes.GET("/workflow-contexts/:id", workflowContextController.GetWorkflowContextDetails)

	masterRoutes.GET("/comment-statuses", commentStatusController.GetAllCommentStatuses)
	masterRoutes.GET("/comment-statuses/:id", commentStatusController.GetCommentStatusDetails)

	masterRoutes.GET("/action-plan-masters", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), actionPlanMasterController.GetAllActionPlanMasters)
	masterRoutes.GET("/action-plan-masters/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), actionPlanMasterController.GetActionPlanMasterFullDetails)
	masterRoutes.POST("/action-plan-masters", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), actionPlanMasterController.CreateActionPlanMaster)
	masterRoutes.PUT("/action-plan-masters/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), actionPlanMasterController.UpdateActionPlanMaster)
	masterRoutes.POST("/action-plan-masters/delete", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), actionPlanMasterController.DeleteActionPlanMasters)
	masterRoutes.PUT("/action-plan-masters/activate", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), actionPlanMasterController.ActivateActionPlanMasters)
	masterRoutes.PUT("/action-plan-masters/deactivate", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), actionPlanMasterController.DeactivateActionPlanMasters)

	masterRoutes.GET("/email-statuses", emailStatusController.GetAllEmailStatuses)
	masterRoutes.GET("/email-statuses/:id", emailStatusController.GetEmailStatusDetails)

	// locations
	masterRoutes.GET("/locations", locationController.GetAllLocations)
	masterRoutes.POST("/locations", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), locationController.CreateLocation)
	masterRoutes.GET("/locations/:id", locationController.GetLocationDetails)
	masterRoutes.PUT("/locations/:id", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), locationController.UpdateLocation)
	masterRoutes.POST("/locations/delete", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), locationController.DeleteLocations)
	masterRoutes.PUT("/locations/activate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), locationController.ActivateLocations)
	masterRoutes.PUT("/locations/deactivate", middlewares.RolePermissionMiddleware(db, "common_setting", "modify"), locationController.DeactivateLocations)
}
