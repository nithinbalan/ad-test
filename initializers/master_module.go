package initializers

import (
	"jethur/grc/masters/controllers"
	"jethur/grc/masters/repositories"
	"jethur/grc/masters/services"

	"gorm.io/gorm"
)

type MasterModule struct {
	ReviewFrequencyController     *controllers.ReviewFrequencyController
	SeverityLevelController       *controllers.SeverityLevelController
	LifecycleStatusController     *controllers.LifecycleStatusController
	QuestionGroupController       *controllers.QuestionGroupController
	ChecklistController           *controllers.ChecklistController
	ChecklistQuestionController   *controllers.ChecklistQuestionController
	AtionPlanStatusController     *controllers.ActionPlanStatusController
	CurrencyController            *controllers.CurrencyController
	TimezoneController            *controllers.TimezoneController
	UnitController                *controllers.UnitController
	ToolTypeController            *controllers.ToolTypeController
	ExpenseCategoryController     *controllers.ExpenseCategoryController
	SourceController              *controllers.SourceController
	StatusController              *controllers.StatusController
	ConformanceController         *controllers.ConformanceController
	WorkflowStatusController      *controllers.WorkflowStatusController
	IntegrationCategoryController *controllers.IntegrationCategoryController
	ImpactCategoryController      *controllers.ImpactCategoryController
	TimeUnitController            *controllers.TimeUnitController
	ReviewStatusController        *controllers.ReviewStatusController

	ReviewFrequencySvc     services.ReviewFrequencyServiceInterface
	SeverityLevelSvc       services.SeverityLevelServiceInterface
	LifecycleStatusSvc     services.LifecycleStatusServiceInterface
	QuestionGroupService   services.QuestionGroupServiceInterface
	ChecklistSvc           services.ChecklistServiceInterface
	ChecklistQuestionSvc   services.ChecklistQuestionServiceInterface
	ActionPlanStatusSvc    services.ActionPlanStatusServiceInterface
	CurrencySvc            services.CurrencyServiceInterface
	TimezoneSvc            services.TimezoneServiceInterface
	UnitSvc                services.UnitServiceInterface
	ToolTypeSvc            services.ToolTypeServiceInterface
	ExpenseCategorySvc     services.ExpenseCategoryServiceInterface
	SourceSvc              services.SourceServiceInterface
	StatusSvc              services.StatusServiceInterface
	ConformanceSvc         services.ConformanceServiceInterface
	WorkflowStatusSvc      services.WorkflowStatusServiceInterface
	IntegrationCategorySvc services.IntegrationCategoryServiceInterface
	ImpactCategorySvc      services.ImpactCategoryServiceInterface
	TimeUnitSvc            services.TimeUnitServiceInterface
	ReviewStatusSvc        services.ReviewStatusServiceInterface
}

func NewMasterModule(db *gorm.DB, sharedModule *SharedModule) *MasterModule {
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
	questionGroupSvc := services.NewQuestionGroupService(questionGroupRepository, sharedModule.ActivityLogService)
	questionGroupController := controllers.NewQuestionGroupController(questionGroupSvc)

	//Checklists
	checklistRepo := repositories.NewChecklistRepository(db)
	checklistSvc := services.NewChecklistService(checklistRepo, sharedModule.ActivityLogService)
	checklistController := controllers.NewChecklistController(checklistSvc)

	//Checklist Questions
	checklistQuestionRepo := repositories.NewChecklistQuestionRepository(db)
	checklistQuestionSvc := services.NewChecklistQuestionService(checklistQuestionRepo, sharedModule.ActivityLogService, questionGroupSvc, checklistSvc)
	checklistQuestionController := controllers.NewChecklistQuestionController(checklistQuestionSvc)

	// Action Plan Status
	actionPlanStatusRepo := repositories.NewActionPlanStatusRepository(db)
	actionPlanStatusSvc := services.NewActionPlanStatusService(actionPlanStatusRepo)
	actionPlanStatusController := controllers.NewActionPlanStatusController(actionPlanStatusSvc)

	currencyRepo := repositories.NewCurrencyRepository(db)
	currencySvc := services.NewCurrencyService(currencyRepo)
	currencyController := controllers.NewCurrencyController(currencySvc)

	timezoneRepo := repositories.NewTimezoneRepository(db)
	timezoneSvc := services.NewTimezoneService(timezoneRepo)
	timezoneController := controllers.NewTimezoneController(timezoneSvc)

	//Units
	unitRepo := repositories.NewUnitRepository(db)
	unitSvc := services.NewUnitService(unitRepo, sharedModule.ActivityLogService)
	unitController := controllers.NewUnitController(unitSvc)

	//ToolTypes
	toolTypeRepo := repositories.NewToolTypeRepository(db)
	toolTypeSvc := services.NewToolTypeService(toolTypeRepo, sharedModule.ActivityLogService)
	toolTypeController := controllers.NewToolTypeController(toolTypeSvc)

	//ToolTypes
	expenseCategoryRepo := repositories.NewExpenseCategoryRepository(db)
	expenseCategorySvc := services.NewExpenseCategoryService(expenseCategoryRepo, sharedModule.ActivityLogService)
	expenseCategoryController := controllers.NewExpenseCategoryController(expenseCategorySvc)

	//Sources
	sourceRepo := repositories.NewSourceRepository(db)
	sourceSvc := services.NewSourceService(sourceRepo)
	sourceController := controllers.NewSourceController(sourceSvc)

	//Status
	statusRepo := repositories.NewStatusRepository(db)
	statusSvc := services.NewStatusService(statusRepo)
	statusController := controllers.NewStatusController(statusSvc)

	//Conformance
	conformanceRepo := repositories.NewConformanceRepository(db)
	conformanceSvc := services.NewConformanceService(conformanceRepo)
	conformanceController := controllers.NewConformanceController(conformanceSvc)
	//Workflow Status
	workflowStatusRepo := repositories.NewWorkflowStatusRepository(db)
	workflowStatusSvc := services.NewWorkflowStatusService(workflowStatusRepo)
	workflowStatusController := controllers.NewWorkflowStatusController(workflowStatusSvc)

	// Integration category
	integrationCategoryRepo := repositories.NewIntegrationCategoryRepository(db)
	integrationCategorySvc := services.NewIntegrationCategoryService(integrationCategoryRepo)
	integrationCategoryController := controllers.NewIntegrationCategoryController(integrationCategorySvc)

	// Impact category
	impactCategoryRepo := repositories.NewImpactCategoryRepository(db)
	impactCategorySvc := services.NewImpactCategoryService(impactCategoryRepo, sharedModule.ActivityLogService)
	impactCategoryController := controllers.NewImpactCategoryController(impactCategorySvc)

	// Impact category
	timeUnitRepo := repositories.NewTimeUnitRepository(db)
	timeUnitSvc := services.NewTimeUnitService(timeUnitRepo)
	timeUnitController := controllers.NewTimeUnitController(timeUnitSvc)

	//review status
	reviewStatusRepo := repositories.NewReviewStatusRepository(db)
	reviewStatusSvc := services.NewReviewStatusService(reviewStatusRepo)
	reviewStatusController := controllers.NewReviewStatusController(reviewStatusSvc)

	return &MasterModule{
		ReviewFrequencyController:     reviewFrequencyController,
		SeverityLevelController:       severityLevelController,
		LifecycleStatusController:     lifecycleStatusController,
		QuestionGroupController:       questionGroupController,
		ChecklistController:           checklistController,
		ChecklistQuestionController:   checklistQuestionController,
		AtionPlanStatusController:     actionPlanStatusController,
		CurrencyController:            currencyController,
		TimezoneController:            timezoneController,
		UnitController:                unitController,
		ToolTypeController:            toolTypeController,
		ExpenseCategoryController:     expenseCategoryController,
		SourceController:              sourceController,
		StatusController:              statusController,
		ConformanceController:         conformanceController,
		WorkflowStatusController:      workflowStatusController,
		IntegrationCategoryController: integrationCategoryController,
		ImpactCategoryController:      impactCategoryController,
		TimeUnitController:            timeUnitController,
		ReviewStatusController:        reviewStatusController,

		ReviewFrequencySvc:     reviewFrequencySvc,
		SeverityLevelSvc:       severityLevelSvc,
		LifecycleStatusSvc:     lifecycleStatusSvc,
		QuestionGroupService:   questionGroupSvc,
		ChecklistSvc:           checklistSvc,
		ChecklistQuestionSvc:   checklistQuestionSvc,
		ActionPlanStatusSvc:    actionPlanStatusSvc,
		CurrencySvc:            currencySvc,
		TimezoneSvc:            timezoneSvc,
		UnitSvc:                unitSvc,
		ToolTypeSvc:            toolTypeSvc,
		ExpenseCategorySvc:     expenseCategorySvc,
		SourceSvc:              sourceSvc,
		StatusSvc:              statusSvc,
		ConformanceSvc:         conformanceSvc,
		WorkflowStatusSvc:      workflowStatusSvc,
		IntegrationCategorySvc: integrationCategorySvc,
		ImpactCategorySvc:      impactCategorySvc,
		TimeUnitSvc:            timeUnitSvc,
		ReviewStatusSvc:        reviewStatusSvc,
	}
}
