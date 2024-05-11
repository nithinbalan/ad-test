package initializers

import (
	"jethur/grc/masters/controllers"
	"jethur/grc/masters/repositories"
	"jethur/grc/masters/services"

	"gorm.io/gorm"
)

type MasterModule struct {
	ReviewFrequencyController   *controllers.ReviewFrequencyController
	SeverityLevelController     *controllers.SeverityLevelController
	LifecycleStatusController   *controllers.LifecycleStatusController
	QuestionGroupController     *controllers.QuestionGroupController
	ChecklistController         *controllers.ChecklistController
	ChecklistQuestionController *controllers.ChecklistQuestionController
	AtionPlanStatusController   *controllers.ActionPlanStatusController

	ReviewFrequencySvc   services.ReviewFrequencyServiceInterface
	SeverityLevelSvc     services.SeverityLevelServiceInterface
	LifecycleStatusSvc   services.LifecycleStatusServiceInterface
	QuestionGroupService services.QuestionGroupServiceInterface
	ChecklistSvc         services.ChecklistServiceInterface
	ChecklistQuestionSvc services.ChecklistQuestionServiceInterface
	ActionPlanStatusSvc  services.ActionPlanStatusServiceInterface
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

	return &MasterModule{
		ReviewFrequencyController:   reviewFrequencyController,
		SeverityLevelController:     severityLevelController,
		LifecycleStatusController:   lifecycleStatusController,
		QuestionGroupController:     questionGroupController,
		ChecklistController:         checklistController,
		ChecklistQuestionController: checklistQuestionController,
		AtionPlanStatusController:   actionPlanStatusController,

		ReviewFrequencySvc:   reviewFrequencySvc,
		SeverityLevelSvc:     severityLevelSvc,
		LifecycleStatusSvc:   lifecycleStatusSvc,
		QuestionGroupService: questionGroupSvc,
		ChecklistSvc:         checklistSvc,
		ChecklistQuestionSvc: checklistQuestionSvc,
		ActionPlanStatusSvc:  actionPlanStatusSvc,
	}
}
