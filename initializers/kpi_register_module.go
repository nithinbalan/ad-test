package initializers

import (
	"jethur/grc/kpi-register/controllers"
	"jethur/grc/kpi-register/repositories"
	"jethur/grc/kpi-register/services"

	"gorm.io/gorm"
)

type KPIRegisterModule struct {
	KpiCategoryController         *controllers.KPICategoryController
	KpiStatusController           *controllers.KPIStatusController
	KpiReviewStatusController     *controllers.KPIReviewStatusController
	EntryTypeController           *controllers.EntryTypeController
	DataSourceController          *controllers.DataSourceController
	KPIRegisterController         *controllers.KPIRegisterController
	KPIDashboardController        *controllers.KPIDashboardController
	KPIValueUpdateController      *controllers.KPIValueUpdateController
	KPIAchievementLogicController *controllers.KPIAchievementLogicController
	KPIAssociatedItemController   *controllers.KPIAssociatedItemController

	KPIRegisterSvc         services.KPIRegisterServiceInterface
	KPIStatusSvc           services.KPIStatusServiceInterface
	KPIReviewStatusSvc     services.KPIReviewStatusServiceInterface
	KPIAchievementLogicSvc services.KPIAchievementLogicServiceInterface
}

func NewKPIRegisterModule(
	db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	commonModule *CommonModule,
	riskModule *RiskModule,
	assetModule *AssetModule,
	tprmModule *TprmModule,
	complianceModule *ComplianceModule,
	documentModule *DocumentModule,
	strategyModule *StrategyModule,
) *KPIRegisterModule {
	kpiCategoryRepository := repositories.NewKPICategoryRepository(db)
	kpiCategorySvc := services.NewKPICategoryService(kpiCategoryRepository, sharedModule.ActivityLogService)
	kpiCategoryController := controllers.NewKPICategoryController(kpiCategorySvc)

	kpiStatusRepo := repositories.NewKPIStatusRepository(db)
	kpiStatusSvc := services.NewKPIStatusService(kpiStatusRepo)
	kpiStatusController := controllers.NewKPIStatusController(kpiStatusSvc)

	kpiReviewStatusRepo := repositories.NewKPIReviewStatusRepository(db)
	kpiReviewStatusSvc := services.NewKPIReviewStatusService(kpiReviewStatusRepo)
	kpiReviewStatusController := controllers.NewKPIReviewStatusController(kpiReviewStatusSvc)

	kpiAchievementLogicRepo := repositories.NewKPIAchievementLogicRepository(db)
	kpiAchievementLogicSvc := services.NewKPIAchievementLogicService(kpiAchievementLogicRepo)
	kpiAchievementLogicController := controllers.NewKPIAchievementLogicController(kpiAchievementLogicSvc)

	entryTypeRepo := repositories.NewEntryTypeRepository(db)
	entryTypeSvc := services.NewEntryTypeService(entryTypeRepo)
	entryTypeController := controllers.NewEntryTypeController(entryTypeSvc)

	dataSourceRepo := repositories.NewDataSourceRepository(db)
	dataSourceSvc := services.NewDataSourceService(dataSourceRepo)
	dataSourceController := controllers.NewDataSourceController(dataSourceSvc)

	kpiRegisterRepo := repositories.NewKPIRegisterRepository(db)
	kpiRegisterSvc := services.NewKPIRegisterService(
		kpiRegisterRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc,
		settingsModule.OrganizationUnitSvc,
		kpiStatusSvc,
		kpiReviewStatusSvc,
		entryTypeSvc,
		dataSourceSvc,
		kpiCategorySvc,
		masterModule.CalculationTypeSvc,
		masterModule.UnitSvc,
		masterModule.ReviewFrequencySvc,
		masterModule.LifecycleStatusSvc,
		kpiAchievementLogicSvc,
		commonModule.NotificationService,
	)
	kpiRegisterController := controllers.NewKPIRegisterController(kpiRegisterSvc)

	kpiDashboardRepo := repositories.NewKPIDashboardRepository(db)
	kpiDashboardSvc := services.NewKPIDashboardService(kpiDashboardRepo, kpiStatusSvc)
	kpiDashboardController := controllers.NewKPIDashboardController(kpiDashboardSvc)

	kpiValueUpdateRepo := repositories.NewKPIValueUpdateRepository(db)
	kpiValueUpdateSvc := services.NewKPIValueUpdateService(
		kpiValueUpdateRepo,
		sharedModule.ActivityLogService,
		kpiRegisterSvc,
		masterModule.ReviewFrequencySvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		kpiStatusSvc,
		kpiReviewStatusSvc,
		organizationModule.UserSettingSvc,
	)
	kpiValueUpdateController := controllers.NewKPIValueUpdateController(kpiValueUpdateSvc)

	// Keep ADF's existing associated-item implementation wired to the upgraded register service.
	kpiAssociatedItemRepo := repositories.NewKPIAssociatedItemRepository(db)
	kpiAssociatedItemSvc := services.NewKPIAssociatedItemService(kpiAssociatedItemRepo, kpiRegisterSvc, organizationModule.ControlSvc, organizationModule.ProcessSvc, assetModule.AssetRegisterSvc, riskModule.RiskRegisterSvc, tprmModule.ThirdPartyRegisterSvc, complianceModule.ComplianceRegisterSvc, documentModule.DocumentSvc, complianceModule.ComplianceClausevc, commonModule.AssociatedItemValidationService)

	_ = strategyModule // retained in the constructor for compatibility with ADF's initializer order

	return &KPIRegisterModule{
		KpiCategoryController:         kpiCategoryController,
		KpiStatusController:           kpiStatusController,
		KpiReviewStatusController:     kpiReviewStatusController,
		EntryTypeController:           entryTypeController,
		DataSourceController:          dataSourceController,
		KPIRegisterController:         kpiRegisterController,
		KPIDashboardController:        kpiDashboardController,
		KPIValueUpdateController:      kpiValueUpdateController,
		KPIAchievementLogicController: kpiAchievementLogicController,
		KPIAssociatedItemController:   controllers.NewKPIAssociatedItemController(kpiAssociatedItemSvc),
		KPIRegisterSvc:                kpiRegisterSvc,
		KPIStatusSvc:                  kpiStatusSvc,
		KPIReviewStatusSvc:            kpiReviewStatusSvc,
		KPIAchievementLogicSvc:        kpiAchievementLogicSvc,
	}
}
