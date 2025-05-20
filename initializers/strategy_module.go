package initializers

import (
	"jethur/grc/strategy/controllers"
	"jethur/grc/strategy/repositories"
	"jethur/grc/strategy/services"

	"gorm.io/gorm"
)

type StrategyModule struct {
	StrategyProfileController                       *controllers.StrategyProfileController
	StrategyProfileFocusAreaController              *controllers.StrategyProfileFocusAreaController
	StrategyProfileGoalController                   *controllers.StrategyProfileGoalController
	StrategyProfileObjectiveController              *controllers.StrategyProfileObjectiveController
	StrategyProfileInitiativeController             *controllers.StrategyProfileInitiativeController
	StrategyProfileProjectController                *controllers.StrategyProfileProjectController
	StrategyProfileMilestoneController              *controllers.StrategyProfileMilestoneController
	StrategyProfileActionPlanController             *controllers.StrategyProfileActionPlanController
	StrategyProfileActionPlanStatusUpdateController *controllers.StrategyProfileActionPlanStatusUpdateController
	StrategyProfileActionPlanCommentController      *controllers.StrategyProfileActionPlanCommentController
	StrategyProfileSettingController                *controllers.StrategyProfileSettingController
	KpiCategoryController                           *controllers.KPICategoryController
	KpiStatusController                             *controllers.KPIStatusController
	KpiRegisterController                           *controllers.KPIRegisterController
}

func NewStrategyModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule,

) *StrategyModule {

	strategyProfileRepo := repositories.NewStrategyProfileRepository(db)
	strategyProfileSvc := services.NewStrategyProfileService(strategyProfileRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, settingsModule.OrganizationUnitSvc)
	strategyProfileController := controllers.NewStrategyProfileController(strategyProfileSvc)

	strategyProfileFocusAreaRepo := repositories.NewStrategyProfileFocusAreaRepository(db)
	strategyProfileFocusAreaSvc := services.NewStrategyProfileFocusAreaService(strategyProfileFocusAreaRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc)
	strategyProfileFocusAreaController := controllers.NewStrategyProfileFocusAreaController(strategyProfileFocusAreaSvc)

	strategyProfileGoalRepo := repositories.NewStrategyProfileGoalRepository(db)
	strategyProfileGoalSvc := services.NewStrategyProfileGoalService(strategyProfileGoalRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc)
	strategyProfileGoalController := controllers.NewStrategyProfileGoalController(strategyProfileGoalSvc)

	strategyProfileObjectiveRepo := repositories.NewStrategyProfileObjectiveRepository(db)
	strategyProfileObjectiveSvc := services.NewStrategyProfileObjectiveService(strategyProfileObjectiveRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc)
	strategyProfileObjectiveController := controllers.NewStrategyProfileObjectiveController(strategyProfileObjectiveSvc)

	strategyProfileInitiativeRepo := repositories.NewStrategyProfileInitiativeRepository(db)
	strategyProfileInitiativeSvc := services.NewStrategyProfileInitiativeService(strategyProfileInitiativeRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc)
	strategyProfileInitiativeController := controllers.NewStrategyProfileInitiativeController(strategyProfileInitiativeSvc)

	strategyProfileProjectRepo := repositories.NewStrategyProfileProjectRepository(db)
	strategyProfileProjectSvc := services.NewStrategyProfileProjectService(strategyProfileProjectRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc)
	strategyProfileProjectController := controllers.NewStrategyProfileProjectController(strategyProfileProjectSvc)

	strategyProfileMilestoneRepo := repositories.NewStrategyProfileMilestoneRepository(db)
	strategyProfileMilestoneSvc := services.NewStrategyProfileMilestoneService(strategyProfileMilestoneRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc)
	strategyProfileMilestoneController := controllers.NewStrategyProfileMilestoneController(strategyProfileMilestoneSvc)

	strategyProfileActionPlanRepo := repositories.NewStrategyProfileActionPlanRepository(db)
	strategyProfileActionPlanSvc := services.NewStrategyProfileActionPlanService(strategyProfileActionPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, strategyProfileSvc, masterModule.ActionPlanStatusSvc, taskModule.TaskService, commonModule.NotificationService, organizationModule.UserSettingSvc)
	strategyProfileActionPlanController := controllers.NewStrategyProfileActionPlanController(strategyProfileActionPlanSvc)

	strategyProfileActionPlanStatusUpdateRepo := repositories.NewStrategyProfileActionPlanStatusUpdateRepository(db)
	strategyProfileActionPlanStatusUpdateSvc := services.NewStrategyProfileActionPlanStatusUpdateService(strategyProfileActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, strategyProfileActionPlanSvc, masterModule.ActionPlanStatusSvc, strategyProfileSvc, organizationModule.UserSvc, commonModule.NotificationService, taskModule.TaskService, taskModule.TaskStatusService)
	strategyProfileActionPlanStatusUpdateController := controllers.NewStrategyProfileActionPlanStatusUpdateController(strategyProfileActionPlanStatusUpdateSvc)

	strategyProfileActionPlanCommentRepository := repositories.NewStrategyProfileActionPlanCommentRepository(db)
	strategyProfileActionPlanCommentSvc := services.NewStrategyProfileActionPlanCommentService(strategyProfileActionPlanCommentRepository, organizationModule.UserSvc, sharedModule.ActivityLogService, strategyProfileActionPlanSvc, commonModule.NotificationService)
	strategyProfileActionPlanCommentController := controllers.NewStrategyProfileActionPlanCommentController(strategyProfileActionPlanCommentSvc)

	strategyProfileSettingRepo := repositories.NewStrategyProfileSettingRepository(db)
	strategyProfileSettingSvc := services.NewStrategyProfileSettingService(strategyProfileSettingRepo)
	strategyProfileSettingController := controllers.NewStrategyProfileSettingController(strategyProfileSettingSvc)

	kpiCategoryRepository := repositories.NewKPICategoryRepository(db)
	kpiCategorySvc := services.NewKPICategoryService(kpiCategoryRepository, sharedModule.ActivityLogService)
	kpiCategoryController := controllers.NewKPICategoryController(kpiCategorySvc)

	kpiStatusRepo := repositories.NewKPIStatusRepository(db)
	kpiStatusSvc := services.NewKPIStatusService(kpiStatusRepo)
	kpiStatusController := controllers.NewKPIStatusController(kpiStatusSvc)

	kpiRegisterRepo := repositories.NewKPIRegisterRepository(db)
	kpiRegisterSvc := services.NewKPIRegisterService(kpiRegisterRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, settingsModule.OrganizationUnitSvc, kpiStatusSvc)
	kpiRegisterController := controllers.NewKPIRegisterController(kpiRegisterSvc)

	return &StrategyModule{
		StrategyProfileController:                       strategyProfileController,
		StrategyProfileFocusAreaController:              strategyProfileFocusAreaController,
		StrategyProfileGoalController:                   strategyProfileGoalController,
		StrategyProfileObjectiveController:              strategyProfileObjectiveController,
		StrategyProfileInitiativeController:             strategyProfileInitiativeController,
		StrategyProfileProjectController:                strategyProfileProjectController,
		StrategyProfileMilestoneController:              strategyProfileMilestoneController,
		StrategyProfileActionPlanController:             strategyProfileActionPlanController,
		StrategyProfileActionPlanStatusUpdateController: strategyProfileActionPlanStatusUpdateController,
		StrategyProfileActionPlanCommentController:      strategyProfileActionPlanCommentController,
		StrategyProfileSettingController:                strategyProfileSettingController,
		KpiCategoryController:                           kpiCategoryController,
		KpiStatusController:                             kpiStatusController,
		KpiRegisterController:                           kpiRegisterController,
	}
}
