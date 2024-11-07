package initializers

import (
	"jethur/grc/asset/controllers"
	"jethur/grc/asset/repositories"
	"jethur/grc/asset/services"

	"gorm.io/gorm"
)

type AssetModule struct {
	AssetTypeController                           *controllers.AssetTypeController
	AssetCategoryController                       *controllers.AssetCategoryController
	AssetStatusController                         *controllers.AssetStatusController
	AssetRegisterController                       *controllers.AssetRegisterController
	AssetCriticalityMatrixFactorController        *controllers.AssetCriticalityMatrixFactorController
	AssetCriticalityMatrixFactorRatingController  *controllers.AssetCriticalityMatrixFactorRatingController
	AssetCriticalityMatrixRatingScaleController   *controllers.AssetCriticalityMatrixRatingScaleController
	VulnerabilityTypeController                   *controllers.VulnerabilityTypeController
	VulnerabilityStatusController                 *controllers.VulnerabilityStatusController
	VulnerabilityController                       *controllers.VulnerabilityController
	ThreatTypeController                          *controllers.ThreatTypeController
	ThreatStatusController                        *controllers.ThreatStatusController
	ThreatController                              *controllers.ThreatController
	VulnerabilityActionPlanController             *controllers.VulnerabilityActionPlanController
	VulnerabilityActionPlanStatusUpdateController *controllers.VulnerabilityActionPlanStatusUpdateController
	ThreatActionPlanController                    *controllers.ThreatActionPlanController
	ThreatActionPlanStatusUpdateController        *controllers.ThreatActionPlanStatusUpdateController

	AssetRegisterSvc services.AssetRegisterServiceInterface
}

func NewAssetModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule) *AssetModule {

	assetTypeRepo := repositories.NewAssetTypeRepository(db)
	assetTypeSvc := services.NewAssetTypeService(assetTypeRepo, sharedModule.ActivityLogService)
	assetTypeController := controllers.NewAssetTypeController(assetTypeSvc)

	assetCategoryRepo := repositories.NewAssetCategoryRepository(db)
	assetCategorySvc := services.NewAssetCategoryService(assetCategoryRepo, assetTypeSvc, sharedModule.ActivityLogService)
	assetCategoryController := controllers.NewAssetCategoryController(assetCategorySvc)

	assetStatusRepo := repositories.NewAssetStatusRepository(db)
	assetStatusSvc := services.NewAssetStatusService(assetStatusRepo)
	assetStatusController := controllers.NewAssetStatusController(assetStatusSvc)

	assetRegisterRepo := repositories.NewAssetRegisterRepository(db)
	assetRegisterSvc := services.NewAssetRegisterService(assetRegisterRepo, sharedModule.ActivityLogService, assetTypeSvc, assetCategorySvc, assetStatusSvc, organizationModule.UserSvc)
	assetRegisterController := controllers.NewAssetRegisterController(assetRegisterSvc)

	assetCriticalityMatrixRepo := repositories.NewAssetCriticalityMatrixRepository(db)
	assetCriticalityMatrixSvc := services.NewAssetCriticalityMatrixService(assetCriticalityMatrixRepo)

	assetCriticalityMatrixFactorRepo := repositories.NewAssetCriticalityMatrixFactorRepository(db)
	assetCriticalityMatrixFactorSvc := services.NewAssetCriticalityMatrixFactorService(assetCriticalityMatrixFactorRepo, sharedModule.ActivityLogService, assetCriticalityMatrixSvc)
	assetCriticalityMatrixFactorController := controllers.NewAssetCriticalityMatrixFactorController(assetCriticalityMatrixFactorSvc)

	assetCriticalityMatrixFactorRatingRepo := repositories.NewAssetCriticalityMatrixFactorRatingRepository(db)
	assetCriticalityMatrixFactorRatingSvc := services.NewAssetCriticalityMatrixFactorRatingService(assetCriticalityMatrixFactorRatingRepo, sharedModule.ActivityLogService, assetCriticalityMatrixSvc)
	assetCriticalityMatrixFactorRatingController := controllers.NewAssetCriticalityMatrixFactorRatingController(assetCriticalityMatrixFactorRatingSvc)

	assetCriticalityMatrixRatingScaleRepo := repositories.NewAssetCriticalityMatrixRatingScaleRepository(db)
	assetCriticalityMatrixRatingScaleSvc := services.NewAssetCriticalityMatrixRatingScaleService(assetCriticalityMatrixRatingScaleRepo, sharedModule.ActivityLogService, assetCriticalityMatrixSvc)
	assetCriticalityMatrixRatingScaleController := controllers.NewAssetCriticalityMatrixRatingScaleController(assetCriticalityMatrixRatingScaleSvc)

	vulnerabilityTypeRepo := repositories.NewVulnerabilityTypeRepository(db)
	vulnerabilityTypeSvc := services.NewVulnerabilityTypeService(vulnerabilityTypeRepo, sharedModule.ActivityLogService)
	vulnerabilityTypeController := controllers.NewVulnerabilityTypeController(vulnerabilityTypeSvc)

	vulnerabilityStatusRepo := repositories.NewVulnerabilityStatusRepository(db)
	vulnerabilityStatusSvc := services.NewVulnerabilityStatusService(vulnerabilityStatusRepo)
	vulnerabilityStatusController := controllers.NewVulnerabilityStatusController(vulnerabilityStatusSvc)

	vulnerabilityRepo := repositories.NewVulnerabilityRepository(db)
	vulnerabilitySvc := services.NewVulnerabilityService(vulnerabilityRepo, sharedModule.ActivityLogService, vulnerabilityTypeSvc, masterModule.SeverityLevelSvc, vulnerabilityStatusSvc, organizationModule.UserSvc)
	vulnerabilityController := controllers.NewVulnerabilityController(vulnerabilitySvc)

	threatTypeRepo := repositories.NewThreatTypeRepository(db)
	threatTypeSvc := services.NewThreatTypeService(threatTypeRepo, sharedModule.ActivityLogService)
	threatTypeController := controllers.NewThreatTypeController(threatTypeSvc)

	threatStatusRepo := repositories.NewThreatStatusRepository(db)
	threatStatusSvc := services.NewThreatStatusService(threatStatusRepo)
	threatStatusController := controllers.NewThreatStatusController(threatStatusSvc)

	threatRepo := repositories.NewThreatRepository(db)
	threatSvc := services.NewThreatService(
		threatRepo,
		sharedModule.ActivityLogService,
		threatTypeSvc,
		masterModule.SeverityLevelSvc,
		threatStatusSvc,
		organizationModule.UserSvc,
	)
	threatController := controllers.NewThreatController(threatSvc)

	vulnerabilityActionPlanRepo := repositories.NewVulnerabilityActionPlanRepository(db)
	vulnerabilityActionPlanSvc := services.NewVulnerabilityActionPlanService(
		vulnerabilityActionPlanRepo,
		sharedModule.ActivityLogService,
		vulnerabilitySvc,
		masterModule.ActionPlanStatusSvc,
		masterModule.SeverityLevelSvc,
		organizationModule.UserSettingSvc,
		organizationModule.UserSvc,
	)
	vulnerabilityActionPlanController := controllers.NewVulnerabilityActionPlanController(vulnerabilityActionPlanSvc)

	vulnerabilityActionPlanStatusUpdateRepo := repositories.NewVulnerabilityActionPlanStatusUpdateRepository(db)
	vulnerabilityActionPlanStatusUpdateSvc := services.NewVulnerabilityActionPlanStatusUpdateService(
		vulnerabilityActionPlanStatusUpdateRepo,
		sharedModule.ActivityLogService,
		vulnerabilityActionPlanSvc,
		masterModule.ActionPlanStatusSvc,
	)
	vulnerabilityActionPlanStatusUpdateController := controllers.NewVulnerabilityActionPlanStatusUpdateController(vulnerabilityActionPlanStatusUpdateSvc)

	threatActionPlanRepo := repositories.NewThreatActionPlanRepository(db)
	threatActionPlanSvc := services.NewThreatActionPlanService(
		threatActionPlanRepo,
		sharedModule.ActivityLogService,
		threatSvc,
		masterModule.ActionPlanStatusSvc,
		masterModule.SeverityLevelSvc,
		organizationModule.UserSettingSvc,
		organizationModule.UserSvc,
	)
	threatActionPlanController := controllers.NewThreatActionPlanController(threatActionPlanSvc)

	threatActionPlanStatusUpdateRepo := repositories.NewThreatActionPlanStatusUpdateRepository(db)
	threatActionPlanStatusUpdateSvc := services.NewThreatActionPlanStatusUpdateService(
		threatActionPlanStatusUpdateRepo,
		sharedModule.ActivityLogService,
		threatActionPlanSvc,
		masterModule.ActionPlanStatusSvc,
	)
	threatActionPlanStatusUpdateController := controllers.NewThreatActionPlanStatusUpdateController(threatActionPlanStatusUpdateSvc)

	return &AssetModule{
		AssetTypeController:                           assetTypeController,
		AssetCategoryController:                       assetCategoryController,
		AssetStatusController:                         assetStatusController,
		AssetRegisterController:                       assetRegisterController,
		AssetCriticalityMatrixFactorController:        assetCriticalityMatrixFactorController,
		AssetCriticalityMatrixFactorRatingController:  assetCriticalityMatrixFactorRatingController,
		AssetCriticalityMatrixRatingScaleController:   assetCriticalityMatrixRatingScaleController,
		VulnerabilityTypeController:                   vulnerabilityTypeController,
		VulnerabilityStatusController:                 vulnerabilityStatusController,
		VulnerabilityController:                       vulnerabilityController,
		ThreatTypeController:                          threatTypeController,
		ThreatStatusController:                        threatStatusController,
		ThreatController:                              threatController,
		VulnerabilityActionPlanController:             vulnerabilityActionPlanController,
		VulnerabilityActionPlanStatusUpdateController: vulnerabilityActionPlanStatusUpdateController,
		ThreatActionPlanController:                    threatActionPlanController,
		ThreatActionPlanStatusUpdateController:        threatActionPlanStatusUpdateController,

		AssetRegisterSvc: assetRegisterSvc,
	}
}
