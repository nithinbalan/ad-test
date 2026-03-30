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
	StrategyProfileStrategicObjectiveController     *controllers.StrategyProfileStrategicObjectiveController
	StrategyProfileStrategicIntiativeController     *controllers.StrategyProfileStrategicIntiativeController
	StrategyProfileProjectController                *controllers.StrategyProfileProjectController
	StrategyProfileProjectMilestoneController       *controllers.StrategyProfileProjectMilestoneController
	StrategyProfileActionPlanController             *controllers.StrategyProfileActionPlanController
	StrategyProfileActionPlanStatusUpdateController *controllers.StrategyProfileActionPlanStatusUpdateController
	StrategyProfileActionPlanCommentController      *controllers.StrategyProfileActionPlanCommentController
	StrategyProfileSettingController                *controllers.StrategyProfileSettingController
	StrategyProfileMappingController                *controllers.StrategyProfileMappingController
	KpiCategoryController                           *controllers.KPICategoryController
	KpiStatusController                             *controllers.KPIStatusController
	EntryTypeController                             *controllers.EntryTypeController
	DataSourceController                            *controllers.DataSourceController
	KpiRegisterController                           *controllers.KPIRegisterController
	StrategicObjectiveCategoryController            *controllers.StrategicObjectiveCategoryController
	StrategicObjectiveController                    *controllers.StrategicObjectiveController
	StrategicObjectiveKPIController                 *controllers.StrategicObjectiveKPIController
	StrategicObjectivePeriodicReviewController      *controllers.StrategicObjectivePeriodicReviewController
	KpiValueUpdateController                        *controllers.KPIValueUpdateController
	KpiAssociatedItemController                     *controllers.KPIAssociatedItemController
	FocusAreaController                             *controllers.FocusAreaController
	StrategicIntiativeController                    *controllers.StrategicIntiativeController
	StrategicObjectiveAssociatedItemController      *controllers.StrategicObjectiveAssociatedItemController
	StrategyProfileKpiController                    *controllers.StrategyProfileKpiController
	StrategyProfileProjectActionPlanController      *controllers.StrategyProfileProjectActionPlanController
	StrategicInitiativeStatusController             *controllers.StrategicInitiativeStatusController
}

func NewStrategyModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
	riskModule *RiskModule,
	assetModule *AssetModule,
	tprmModule *TprmModule,
	complianceModule *ComplianceModule,
	documentModule *DocumentModule,
	projectModule *ProjectModule,

) *StrategyModule {

	strategyProfileRepo := repositories.NewStrategyProfileRepository(db)
	strategyProfileSvc := services.NewStrategyProfileService(strategyProfileRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, settingsModule.OrganizationUnitSvc)
	strategyProfileController := controllers.NewStrategyProfileController(strategyProfileSvc)

	strategyProfileGoalRepo := repositories.NewStrategyProfileGoalRepository(db)
	strategyProfileGoalSvc := services.NewStrategyProfileGoalService(strategyProfileGoalRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc)
	strategyProfileGoalController := controllers.NewStrategyProfileGoalController(strategyProfileGoalSvc)

	strategyProfileActionPlanRepo := repositories.NewStrategyProfileActionPlanRepository(db)
	strategyProfileActionPlanSvc := services.NewStrategyProfileActionPlanService(strategyProfileActionPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, strategyProfileSvc, masterModule.ActionPlanStatusSvc, taskModule.TaskService, commonModule.NotificationService, organizationModule.UserSettingSvc, commonModule.DelegateeUserService)
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

	strategyProfileMappingRepo := repositories.NewStrategyProfileMappingRepository(db)
	strategyProfileMappingSvc := services.NewStrategyProfileMappingService(strategyProfileMappingRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, settingsModule.OrganizationUnitSvc, strategyProfileSettingSvc)
	strategyProfileMappingController := controllers.NewStrategyProfileMappingController(strategyProfileMappingSvc)

	kpiCategoryRepository := repositories.NewKPICategoryRepository(db)
	kpiCategorySvc := services.NewKPICategoryService(kpiCategoryRepository, sharedModule.ActivityLogService)
	kpiCategoryController := controllers.NewKPICategoryController(kpiCategorySvc)

	kpiStatusRepo := repositories.NewKPIStatusRepository(db)
	kpiStatusSvc := services.NewKPIStatusService(kpiStatusRepo)
	kpiStatusController := controllers.NewKPIStatusController(kpiStatusSvc)

	entryTypeRepo := repositories.NewEntryTypeRepository(db)
	entryTypeSvc := services.NewEntryTypeService(entryTypeRepo)
	entryTypeController := controllers.NewEntryTypeController(entryTypeSvc)

	dataSourceRepo := repositories.NewDataSourceRepository(db)
	dataSourceSvc := services.NewDataSourceService(dataSourceRepo)
	dataSourceController := controllers.NewDataSourceController(dataSourceSvc)

	kpiRegisterRepo := repositories.NewKPIRegisterRepository(db)
	kpiRegisterSvc := services.NewKPIRegisterService(kpiRegisterRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, settingsModule.OrganizationUnitSvc, kpiStatusSvc, entryTypeSvc, dataSourceSvc)
	kpiRegisterController := controllers.NewKPIRegisterController(kpiRegisterSvc)

	strategicObjectiveCategoryRepository := repositories.NewStrategicObjectiveCategoryRepository(db)
	strategicObjectiveCategorySvc := services.NewStrategicObjectiveCategoryService(strategicObjectiveCategoryRepository, sharedModule.ActivityLogService)
	strategicObjectiveCategoryController := controllers.NewStrategicObjectiveCategoryController(strategicObjectiveCategorySvc)

	strategicObjectiveRepository := repositories.NewStrategicObjectiveRepository(db)
	strategicObjectiveSvc := services.NewStrategicObjectiveService(strategicObjectiveRepository,
		sharedModule.ActivityLogService,
		strategicObjectiveCategorySvc,
		settingsModule.OrganizationUnitSvc,
		organizationModule.UserSvc,
		masterModule.ReviewFrequencySvc,
		masterModule.ReviewStatusSvc,
	)
	strategicObjectiveController := controllers.NewStrategicObjectiveController(strategicObjectiveSvc)

	strategicObjectiveKPIRepository := repositories.NewStrategicObjectiveKPIRepository(db)
	strategicObjectiveKPISvc := services.NewStrategicObjectiveKPIService(strategicObjectiveKPIRepository,
		sharedModule.ActivityLogService,
		strategicObjectiveSvc,
		kpiRegisterSvc,
	)
	strategicObjectiveKPIController := controllers.NewStrategicObjectiveKPIController(strategicObjectiveKPISvc)

	strategicObjectivePeriodicReviewRepository := repositories.NewStrategicObjectivePeriodicReviewRepository(db)
	strategicObjectivePeriodicReviewSvc := services.NewStrategicObjectivePeriodicReviewService(strategicObjectivePeriodicReviewRepository,
		sharedModule.ActivityLogService,
		strategicObjectiveSvc,
		masterModule.ReviewStatusSvc,
		masterModule.ReviewFrequencySvc,
	)

	kpiAssociatedItemRepo := repositories.NewKPIAssociatedItemRepository(db)
	kpiAssociatedItemSvc := services.NewKPIAssociatedItemService(kpiAssociatedItemRepo, kpiRegisterSvc, organizationModule.ControlSvc, organizationModule.ProcessSvc, assetModule.AssetRegisterSvc, riskModule.RiskRegisterSvc, tprmModule.ThirdPartyRegisterSvc, complianceModule.ComplianceRegisterSvc, documentModule.DocumentSvc, complianceModule.ComplianceClausevc, commonModule.AssociatedItemValidationService)
	kpiAssociatedItemController := controllers.NewKPIAssociatedItemController(kpiAssociatedItemSvc)

	focusAreaRepository := repositories.NewFocusAreaRepository(db)
	focusAreaSvc := services.NewFocusAreaService(focusAreaRepository,
		sharedModule.ActivityLogService)
	focusAreaController := controllers.NewFocusAreaController(focusAreaSvc)

	strategyProfileFocusAreaRepo := repositories.NewStrategyProfileFocusAreaRepository(db)
	strategyProfileFocusAreaSvc := services.NewStrategyProfileFocusAreaService(strategyProfileFocusAreaRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc, focusAreaSvc)
	strategyProfileFocusAreaController := controllers.NewStrategyProfileFocusAreaController(strategyProfileFocusAreaSvc)

	strategyProfileStrategicObjectiveRepo := repositories.NewStrategyProfileStrategicObjectiveRepository(db)
	strategyProfileStrategicObjectiveSvc := services.NewStrategyProfileStrategicObjectiveService(strategyProfileStrategicObjectiveRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc, strategyProfileFocusAreaSvc, strategyProfileGoalSvc, strategicObjectiveSvc)
	strategyProfileStrategicObjectiveController := controllers.NewStrategyProfileStrategicObjectiveController(strategyProfileStrategicObjectiveSvc)

	strategicIntiativeRepo := repositories.NewStrategicIntiativeRepository(db)
	strategicInitiativeStatusRepo := repositories.NewStrategicInitiativeStatusRepository(db)
	strategicInitiativeStatusSvc := services.NewStrategicInitiativeStatusService(strategicInitiativeStatusRepo)
	strategicInitiativeStatusController := controllers.NewStrategicInitiativeStatusController(strategicInitiativeStatusSvc)
	strategicIntiativeSvc := services.NewStrategicIntiativeService(strategicIntiativeRepo, strategicInitiativeStatusRepo, sharedModule.ActivityLogService, settingsModule.OrganizationUnitSvc, organizationModule.UserSvc)
	strategicIntiativeController := controllers.NewStrategicIntiativeController(strategicIntiativeSvc)
	strategicObjectiveAssociatedItemRepo := repositories.NewStrategicObectiveAssociatedItemRepository(db)
	strategicObjectiveAssociatedItemSvc := services.NewstrategicObjectiveAssociatedItemService(strategicObjectiveAssociatedItemRepo, strategicObjectiveSvc, complianceModule.ComplianceClausevc, commonModule.AssociatedItemValidationService)
	strategicObjectiveAssociatedItemController := controllers.NewStrategicObjectiveAssociatedItemController(strategicObjectiveAssociatedItemSvc)

	strategyProfileKpiRepo := repositories.NewStrategyProfileKpiRepository(db)
	strategyProfileKpiSvc := services.NewStrategyProfileKpiService(strategyProfileKpiRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc, strategyProfileStrategicObjectiveSvc, kpiRegisterSvc)
	strategyProfileKpiController := controllers.NewStrategyProfileKpiController(strategyProfileKpiSvc)

	strategyProfileStrategicIntiativeRepo := repositories.NewStrategyProfileStrategicIntiativeRepository(db)
	strategyProfileStrategicIntiativeSvc := services.NewStrategyProfileStrategicIntiativeService(strategyProfileStrategicIntiativeRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc, strategyProfileStrategicObjectiveSvc, strategicIntiativeSvc)
	strategyProfileStrategicIntiativeController := controllers.NewStrategyProfileStrategicIntiativeController(strategyProfileStrategicIntiativeSvc)

	strategyProfileProjectRepo := repositories.NewStrategyProfileProjectRepository(db)
	strategyProfileProjectSvc := services.NewStrategyProfileProjectService(strategyProfileProjectRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc, projectModule.projectSvc, strategyProfileStrategicObjectiveSvc, strategyProfileStrategicIntiativeSvc)
	strategyProfileProjectController := controllers.NewStrategyProfileProjectController(strategyProfileProjectSvc)

	strategyProfileProjectMilestoneRepo := repositories.NewStrategyProfileProjectMilestoneRepository(db)
	strategyProfileProjectMilestoneSvc := services.NewStrategyProfileProjectMilestoneService(strategyProfileProjectMilestoneRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc, strategyProfileProjectSvc, projectModule.projectMilestoneSvc)
	strategyProfileProjectMilestoneController := controllers.NewStrategyProfileProjectMilestoneController(strategyProfileProjectMilestoneSvc)

	strategyProfileProjectActionPlanRepo := repositories.NewStrategyProfileProjectActionPlanRepository(db)
	strategyProfileProjectActionPlanSvc := services.NewStrategyProfileProjectActionPlanService(strategyProfileProjectActionPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc, strategyProfileProjectSvc, strategyProfileProjectMilestoneSvc, projectModule.projectActionPlanSvc)
	strategyProfileProjectActionPlanController := controllers.NewStrategyProfileProjectActionPlanController(strategyProfileProjectActionPlanSvc)

	strategicObjectivePeriodicReviewController := controllers.NewStrategicObjectivePeriodicReviewController(strategicObjectivePeriodicReviewSvc)
	kpiValueUpdateRepo := repositories.NewKPIValueUpdateRepository(db)
	kpiValueUpdateSvc := services.NewKPIValueUpdateService(kpiValueUpdateRepo, sharedModule.ActivityLogService, kpiRegisterSvc, masterModule.ReviewFrequencySvc, organizationModule.UserSvc, commonModule.NotificationService, strategyProfileKpiSvc, strategyProfileStrategicObjectiveSvc)
	kpiValueUpdateController := controllers.NewKPIValueUpdateController(kpiValueUpdateSvc)

	return &StrategyModule{
		StrategyProfileController:                       strategyProfileController,
		StrategyProfileFocusAreaController:              strategyProfileFocusAreaController,
		StrategyProfileGoalController:                   strategyProfileGoalController,
		StrategyProfileStrategicObjectiveController:     strategyProfileStrategicObjectiveController,
		StrategyProfileStrategicIntiativeController:     strategyProfileStrategicIntiativeController,
		StrategyProfileProjectController:                strategyProfileProjectController,
		StrategyProfileProjectMilestoneController:       strategyProfileProjectMilestoneController,
		StrategyProfileActionPlanController:             strategyProfileActionPlanController,
		StrategyProfileActionPlanStatusUpdateController: strategyProfileActionPlanStatusUpdateController,
		StrategyProfileActionPlanCommentController:      strategyProfileActionPlanCommentController,
		StrategyProfileSettingController:                strategyProfileSettingController,
		StrategyProfileMappingController:                strategyProfileMappingController,
		KpiCategoryController:                           kpiCategoryController,
		KpiStatusController:                             kpiStatusController,
		EntryTypeController:                             entryTypeController,
		DataSourceController:                            dataSourceController,
		KpiRegisterController:                           kpiRegisterController,
		KpiValueUpdateController:                        kpiValueUpdateController,
		KpiAssociatedItemController:                     kpiAssociatedItemController,
		StrategicObjectiveCategoryController:            strategicObjectiveCategoryController,
		StrategicObjectiveController:                    strategicObjectiveController,
		StrategicObjectiveKPIController:                 strategicObjectiveKPIController,
		StrategicObjectivePeriodicReviewController:      strategicObjectivePeriodicReviewController,
		FocusAreaController:                             focusAreaController,
		StrategicIntiativeController:                    strategicIntiativeController,
		StrategicObjectiveAssociatedItemController:      strategicObjectiveAssociatedItemController,
		StrategyProfileKpiController:                    strategyProfileKpiController,
		StrategyProfileProjectActionPlanController:      strategyProfileProjectActionPlanController,
		StrategicInitiativeStatusController:             strategicInitiativeStatusController,
	}
}
