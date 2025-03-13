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
	AssetSubCategoryController                    *controllers.AssetSubCategoryController
	AssetStatusController                         *controllers.AssetStatusController
	AssetRegisterController                       *controllers.AssetRegisterController
	AssetCriticalityMatrixController              *controllers.AssetCriticalityMatrixController
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
	AssetMaintenanceController                    *controllers.AssetMaintenanceController
	AssetAssociatedItemController                 *controllers.AssetAssociatedItemController
	AssetCriticalityController                    *controllers.AssetCriticalityController
	AssetPolicyController                         *controllers.AssetPolicyController
	AssetPolicyChecklistController                *controllers.AssetPolicyChecklistController

	AssetRegisterSvc services.AssetRegisterServiceInterface
}

func NewAssetModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
	documentModule *DocumentModule,
) *AssetModule {
	assetTypeRepo := repositories.NewAssetTypeRepository(db)
	assetTypeSvc := services.NewAssetTypeService(assetTypeRepo, sharedModule.ActivityLogService)
	assetTypeController := controllers.NewAssetTypeController(assetTypeSvc)

	assetCategoryRepo := repositories.NewAssetCategoryRepository(db)
	assetCategorySvc := services.NewAssetCategoryService(assetCategoryRepo, assetTypeSvc, sharedModule.ActivityLogService)
	assetCategoryController := controllers.NewAssetCategoryController(assetCategorySvc)

	assetSubCategoryRepo := repositories.NewAssetSubCategoryRepository(db)
	assetSubCategorySvc := services.NewAssetSubCategoryService(assetSubCategoryRepo, assetCategorySvc, sharedModule.ActivityLogService)
	assetSubCategoryController := controllers.NewAssetSubCategoryController(assetSubCategorySvc)

	assetStatusRepo := repositories.NewAssetStatusRepository(db)
	assetStatusSvc := services.NewAssetStatusService(assetStatusRepo)
	assetStatusController := controllers.NewAssetStatusController(assetStatusSvc)

	assetCriticalityMatrixRepo := repositories.NewAssetCriticalityMatrixRepository(db)
	assetCriticalityMatrixSvc := services.NewAssetCriticalityMatrixService(assetCriticalityMatrixRepo, sharedModule.ActivityLogService)
	assetCriticalityMatrixController := controllers.NewAssetCriticalityMatrixController(assetCriticalityMatrixSvc)

	assetCriticalityMatrixFactorRepo := repositories.NewAssetCriticalityMatrixFactorRepository(db)
	assetCriticalityMatrixFactorSvc := services.NewAssetCriticalityMatrixFactorService(assetCriticalityMatrixFactorRepo, sharedModule.ActivityLogService, assetCriticalityMatrixSvc)
	assetCriticalityMatrixFactorController := controllers.NewAssetCriticalityMatrixFactorController(assetCriticalityMatrixFactorSvc)

	assetCriticalityMatrixFactorRatingRepo := repositories.NewAssetCriticalityMatrixFactorRatingRepository(db)
	assetCriticalityMatrixFactorRatingSvc := services.NewAssetCriticalityMatrixFactorRatingService(assetCriticalityMatrixFactorRatingRepo, sharedModule.ActivityLogService, assetCriticalityMatrixSvc)
	assetCriticalityMatrixFactorRatingController := controllers.NewAssetCriticalityMatrixFactorRatingController(assetCriticalityMatrixFactorRatingSvc)

	assetCriticalityMatrixRatingScaleRepo := repositories.NewAssetCriticalityMatrixRatingScaleRepository(db)
	assetCriticalityMatrixRatingScaleSvc := services.NewAssetCriticalityMatrixRatingScaleService(assetCriticalityMatrixRatingScaleRepo, sharedModule.ActivityLogService, assetCriticalityMatrixSvc)
	assetCriticalityMatrixRatingScaleController := controllers.NewAssetCriticalityMatrixRatingScaleController(assetCriticalityMatrixRatingScaleSvc)

	assetRegisterRepo := repositories.NewAssetRegisterRepository(db)
	assetRegisterSvc := services.NewAssetRegisterService(assetRegisterRepo, sharedModule.ActivityLogService, assetTypeSvc, assetCategorySvc, assetStatusSvc, organizationModule.UserSvc, assetCriticalityMatrixRatingScaleSvc, assetCriticalityMatrixSvc)
	assetRegisterController := controllers.NewAssetRegisterController(assetRegisterSvc)

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
		vulnerabilitySvc,
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
		commonModule.NotificationService,
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
		commonModule.NotificationService,
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

	assetMaintenanceRepo := repositories.NewAssetMaintenanceRepository(db)
	assetMaintenanceSvc := services.NewAssetMaintenanceService(assetMaintenanceRepo, sharedModule.ActivityLogService, assetRegisterSvc, organizationModule.UserSvc)
	assetMaintenanceController := controllers.NewAssetMaintenanceController(assetMaintenanceSvc)
	assetAssociatedItemRepo := repositories.NewAssetAssociatedItemRepository(db)
	assetAssociatedItemSvc := services.NewAssetAssociatedItemService(assetAssociatedItemRepo, assetRegisterSvc, organizationModule.ControlImplementationStatusSvc, organizationModule.ControlEfficiencyStatusSvc)
	assetAssociatedItemController := controllers.NewAssetAssociatedItemController(assetAssociatedItemSvc)

	assetCriticalityRepo := repositories.NewAssetCriticalityRepository(db)
	assetCriticalitySvc := services.NewAssetCriticalityService(
		assetCriticalityRepo,
		sharedModule.ActivityLogService,
		assetRegisterSvc,
		assetCriticalityMatrixFactorSvc,
		assetCriticalityMatrixFactorRatingSvc,
	)
	assetCriticalityController := controllers.NewAssetCriticalityController(assetCriticalitySvc)

	assetPolicyRepo := repositories.NewAssetPolicyRepository(db)
	assetPolicySvc := services.NewAssetPolicyService(
		assetPolicyRepo,
		sharedModule.ActivityLogService,
		assetRegisterSvc,
		documentModule.DocumentSvc,
		documentModule.DocumentRequirementGroupSvc,
		documentModule.DocumentRequirementSvc,
	)
	assetPolicyController := controllers.NewAssetPolicyController(assetPolicySvc)

	assetPolicyChecklistRepo := repositories.NewAssetPolicyChecklistRepository(db)
	assetPolicyChecklistSvc := services.NewAssetPolicyChecklistService(
		assetPolicyChecklistRepo,
		sharedModule.ActivityLogService,
		assetRegisterSvc,
		documentModule.DocumentRequirementGroupSvc,
		assetPolicySvc,
	)
	assetPolicyChecklistController := controllers.NewAssetPolicyChecklistController(assetPolicyChecklistSvc)

	return &AssetModule{
		AssetTypeController:                           assetTypeController,
		AssetCategoryController:                       assetCategoryController,
		AssetSubCategoryController:                    assetSubCategoryController,
		AssetStatusController:                         assetStatusController,
		AssetRegisterController:                       assetRegisterController,
		AssetCriticalityMatrixController:              assetCriticalityMatrixController,
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
		AssetMaintenanceController:                    assetMaintenanceController,
		AssetAssociatedItemController:                 assetAssociatedItemController,
		AssetCriticalityController:                    assetCriticalityController,
		AssetPolicyController:                         assetPolicyController,
		AssetPolicyChecklistController:                assetPolicyChecklistController,

		AssetRegisterSvc: assetRegisterSvc,
	}
}
