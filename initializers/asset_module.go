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
	VulnerabilityActionPlanCommentController      *controllers.VulnerabilityActionPlanCommentController
	ThreatActionPlanController                    *controllers.ThreatActionPlanController
	ThreatActionPlanCommentController             *controllers.ThreatActionPlanCommentController
	ThreatActionPlanStatusUpdateController        *controllers.ThreatActionPlanStatusUpdateController
	AssetMaintenanceController                    *controllers.AssetMaintenanceController
	AssetAssociatedItemController                 *controllers.AssetAssociatedItemController
	AssetCriticalityController                    *controllers.AssetCriticalityController
	AssetPolicyController                         *controllers.AssetPolicyController
	AssetPolicyChecklistController                *controllers.AssetPolicyChecklistController
	AssetPolicyCommentController                  *controllers.AssetPolicyCommentController
	AssetDocumentCategoryController               *controllers.AssetDocumentCategoryController
	AssetOtherDocumentController                  *controllers.AssetOtherDocumentController
	AssetControlAnalysisController                *controllers.AssetControlAnalysisController
	AssetKriController                            *controllers.AssetKriController
	TenableController                             *controllers.TenableController
	SystemRegisterDisruptionLevelController       *controllers.SystemRegisterDisruptionLevelController
	SystemRegisterServiceImportanceController     *controllers.SystemRegisterServiceImportanceController
	AssetDeviceCategoryController                 *controllers.AssetDeviceCategoryController
	SystemRegisterController                      *controllers.SystemRegisterController
	ConfidentialityClassificationController       *controllers.ConfidentialityClassificationController
	AssetEnvironmentController                    *controllers.AssetEnvironmentController
	ApplicationController                         *controllers.ApplicationController
	SLAEscalationMatrixController                 *controllers.SLAEscalationMatrixController
	VulnerabilitySettingsController               *controllers.VulnerabilitySettingsController
	QualysVulnerabilityHistoryController          *controllers.QualysVulnerabilityHistoryController

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
	complianceModule *ComplianceModule,
) *AssetModule {

	assetDeviceCategoryRepo := repositories.NewAssetDeviceCategoryRepository(db)
	assetDeviceCategorySvc := services.NewAssetDeviceCategoryService(assetDeviceCategoryRepo, sharedModule.ActivityLogService)
	assetDeviceCategoryController := controllers.NewAssetDeviceCategoryController(assetDeviceCategorySvc)

	applicationRepo := repositories.NewApplicationRepository(db)
	applicationSvc := services.NewApplicationService(applicationRepo, sharedModule.ActivityLogService)
	applicationController := controllers.NewApplicationController(applicationSvc)

	confidentialityClassificationRepo := repositories.NewConfidentialityClassificationRepository(db)
	confidentialityClassificationSvc := services.NewConfidentialityClassificationService(confidentialityClassificationRepo)
	confidentialityClassificationController := controllers.NewConfidentialityClassificationController(confidentialityClassificationSvc)

	assetEnvironmentRepo := repositories.NewAssetEnvironmentRepository(db)
	assetEnvironmentSvc := services.NewAssetEnvironmentService(assetEnvironmentRepo)
	assetEnvironmentController := controllers.NewAssetEnvironmentController(assetEnvironmentSvc)

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
	assetRegisterSvc := services.NewAssetRegisterService(assetRegisterRepo, sharedModule.ActivityLogService, assetTypeSvc, assetCategorySvc, assetStatusSvc, organizationModule.UserSvc, assetCriticalityMatrixRatingScaleSvc, assetCriticalityMatrixSvc, organizationModule.UserSettingSvc, organizationModule.ControlImplementationStatusSvc, organizationModule.ControlEfficiencyStatusSvc, assetDeviceCategorySvc, assetCriticalityMatrixFactorRatingSvc, confidentialityClassificationSvc, assetEnvironmentSvc)
	assetRegisterController := controllers.NewAssetRegisterController(assetRegisterSvc)

	vulnerabilityTypeRepo := repositories.NewVulnerabilityTypeRepository(db)
	vulnerabilityTypeSvc := services.NewVulnerabilityTypeService(vulnerabilityTypeRepo, sharedModule.ActivityLogService)
	vulnerabilityTypeController := controllers.NewVulnerabilityTypeController(vulnerabilityTypeSvc)

	vulnerabilitySettingsRepo := repositories.NewVulnerabilitySettingsRepository(db)
	vulnerabilitySettingsSvc := services.NewVulnerabilitySettingsService(vulnerabilitySettingsRepo, sharedModule.ActivityLogService)
	vulnerabilitySettingsController := controllers.NewVulnerabilitySettingsController(vulnerabilitySettingsSvc)

	vulnerabilityStatusRepo := repositories.NewVulnerabilityStatusRepository(db)
	vulnerabilityStatusSvc := services.NewVulnerabilityStatusService(vulnerabilityStatusRepo)
	vulnerabilityStatusController := controllers.NewVulnerabilityStatusController(vulnerabilityStatusSvc)

	vulnerabilityRepo := repositories.NewVulnerabilityRepository(db)
	vulnerabilitySvc := services.NewVulnerabilityService(vulnerabilityRepo, sharedModule.ActivityLogService, vulnerabilityTypeSvc, masterModule.SeverityLevelSvc, vulnerabilityStatusSvc, organizationModule.UserSvc, organizationModule.UserSettingSvc, vulnerabilityStatusSvc, vulnerabilitySettingsSvc, organizationModule.ControlSvc, commonModule.NotificationService)
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
		organizationModule.UserSettingSvc,
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
		commonModule.DelegateeUserService,
		taskModule.TaskService,
	)
	vulnerabilityActionPlanController := controllers.NewVulnerabilityActionPlanController(vulnerabilityActionPlanSvc)

	vulnerabilityActionPlanStatusUpdateRepo := repositories.NewVulnerabilityActionPlanStatusUpdateRepository(db)
	vulnerabilityActionPlanStatusUpdateSvc := services.NewVulnerabilityActionPlanStatusUpdateService(
		vulnerabilityActionPlanStatusUpdateRepo,
		sharedModule.ActivityLogService,
		vulnerabilityActionPlanSvc,
		masterModule.ActionPlanStatusSvc,
		vulnerabilitySvc,
		taskModule.TaskService, taskModule.TaskStatusService, organizationModule.UserSvc, commonModule.NotificationService,
	)
	vulnerabilityActionPlanStatusUpdateController := controllers.NewVulnerabilityActionPlanStatusUpdateController(vulnerabilityActionPlanStatusUpdateSvc)

	vulnerabilityeActionPlanCommentRepository := repositories.NewVulnerabilityActionPlanCommentRepository(db)
	vulnerabilityActionPlanCommentSvc := services.NewVulnerabilityActionPlanCommentService(vulnerabilityeActionPlanCommentRepository, organizationModule.UserSvc, sharedModule.ActivityLogService, vulnerabilityActionPlanSvc, commonModule.NotificationService)
	vulnerabilityActionPlanCommentController := controllers.NewVulnerabilityActionPlanCommentController(vulnerabilityActionPlanCommentSvc)

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
		commonModule.DelegateeUserService,
		taskModule.TaskService,
	)
	threatActionPlanController := controllers.NewThreatActionPlanController(threatActionPlanSvc)

	threatActionPlanCommentRepository := repositories.NewThreatActionPlanCommentRepository(db)
	threatActionPlanCommentSvc := services.NewThreatActionPlanCommentService(threatActionPlanCommentRepository, organizationModule.UserSvc, sharedModule.ActivityLogService, threatActionPlanSvc, commonModule.NotificationService)
	threatActionPlanCommentController := controllers.NewThreatActionPlanCommentController(threatActionPlanCommentSvc)

	threatActionPlanStatusUpdateRepo := repositories.NewThreatActionPlanStatusUpdateRepository(db)
	threatActionPlanStatusUpdateSvc := services.NewThreatActionPlanStatusUpdateService(
		threatActionPlanStatusUpdateRepo,
		sharedModule.ActivityLogService,
		threatActionPlanSvc,
		masterModule.ActionPlanStatusSvc,
		threatSvc,
		taskModule.TaskService, taskModule.TaskStatusService, organizationModule.UserSvc, commonModule.NotificationService,
	)
	threatActionPlanStatusUpdateController := controllers.NewThreatActionPlanStatusUpdateController(threatActionPlanStatusUpdateSvc)

	assetMaintenanceRepo := repositories.NewAssetMaintenanceRepository(db)
	assetMaintenanceSvc := services.NewAssetMaintenanceService(assetMaintenanceRepo, sharedModule.ActivityLogService, assetRegisterSvc, organizationModule.UserSvc)
	assetMaintenanceController := controllers.NewAssetMaintenanceController(assetMaintenanceSvc)
	assetAssociatedItemRepo := repositories.NewAssetAssociatedItemRepository(db)
	assetAssociatedItemSvc := services.NewAssetAssociatedItemService(assetAssociatedItemRepo, assetRegisterSvc, organizationModule.ControlImplementationStatusSvc, organizationModule.ControlEfficiencyStatusSvc, vulnerabilitySvc, threatSvc, complianceModule.ComplianceClausevc, commonModule.AssociatedItemValidationService, organizationModule.ControlEfficiencyStatusSvc)
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
		complianceModule.ComplianceStatusSvc,
	)
	assetPolicyController := controllers.NewAssetPolicyController(assetPolicySvc)

	assetPolicyChecklistRepo := repositories.NewAssetPolicyChecklistRepository(db)
	assetPolicyChecklistSvc := services.NewAssetPolicyChecklistService(
		assetPolicyChecklistRepo,
		sharedModule.ActivityLogService,
		assetRegisterSvc,
		documentModule.DocumentRequirementGroupSvc,
		assetPolicySvc,
		complianceModule.ComplianceStatusSvc,
	)
	assetPolicyChecklistController := controllers.NewAssetPolicyChecklistController(assetPolicyChecklistSvc)

	assetPolicyCommentRepo := repositories.NewAssetPolicyCommentRepository(db)
	assetPolicyCommentSvc := services.NewAssetPolicyCommentService(
		assetPolicyCommentRepo,
		sharedModule.ActivityLogService,
		assetRegisterSvc,
		assetPolicySvc,
	)
	assetPolicyCommentController := controllers.NewAssetPolicyCommentController(assetPolicyCommentSvc)

	assetDocumentCategoryRepo := repositories.NewAssetDocumentCategoryRepository(db)
	assetDocumentCategorySvc := services.NewAssetDocumentCategoryService(assetDocumentCategoryRepo, assetTypeSvc, sharedModule.ActivityLogService)
	assetDocumentCategoryController := controllers.NewAssetDocumentCategoryController(assetDocumentCategorySvc)

	assetOtherDocumentRepo := repositories.NewAssetOtherDocumentRepository(db)
	assetOtherDocumentSvc := services.NewAssetOtherDocumentService(
		assetOtherDocumentRepo,
		sharedModule.ActivityLogService,
		assetRegisterSvc,
		assetDocumentCategorySvc,
	)
	assetOtherDocumentController := controllers.NewAssetOtherDocumentController(assetOtherDocumentSvc)
	assetControlAnalysisRepo := repositories.NewAssetControlAnalysisRepository(db)
	assetControlAnalysisSvc := services.NewAssetControlAnalysisService(assetControlAnalysisRepo, sharedModule.ActivityLogService, organizationModule.ControlSvc, organizationModule.ControlImplementationStatusSvc, organizationModule.ControlEfficiencyStatusSvc, assetRegisterSvc)
	assetControlAnalysisController := controllers.NewAssetControlAnalysisController(assetControlAnalysisSvc)

	assetKriRepo := repositories.NewAssetKriRepository(db)
	assetKriSvc := services.NewAssetKriService(
		assetKriRepo,
		sharedModule.ActivityLogService,
		assetRegisterSvc,
	)
	assetKriController := controllers.NewAssetKriController(assetKriSvc)

	tenableRepo := repositories.NewTenableRepository()
	tenableSvc := services.NewTenableService(tenableRepo, assetStatusSvc, settingsModule.TenableSettingsSvc, assetCategorySvc, organizationModule.UserSvc, assetRegisterSvc, assetSubCategorySvc)
	tenableController := controllers.NewTenableController(tenableSvc)

	systemRegisterDisruptionLevelRepo := repositories.NewSystemRegisterDisruptionLevelRepository(db)
	systemRegisterDisruptionLevelSvc := services.NewSystemRegisterDisruptionLevelService(systemRegisterDisruptionLevelRepo, sharedModule.ActivityLogService)
	systemRegisterDisruptionLevelController := controllers.NewSystemRegisterDisruptionLevelController(systemRegisterDisruptionLevelSvc)

	systemRegisterServiceImportanceRepo := repositories.NewSystemRegisterServiceImportanceRepository(db)
	systemRegisterServiceImportanceSvc := services.NewSystemRegisterServiceImportanceService(systemRegisterServiceImportanceRepo, sharedModule.ActivityLogService)
	systemRegisterServiceImportanceController := controllers.NewSystemRegisterServiceImportanceController(systemRegisterServiceImportanceSvc)

	systemRegisterRepo := repositories.NewSystemRegisterRepository(db)
	systemRegisterSvc := services.NewSystemRegisterService(
		systemRegisterRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSettingSvc,
	)
	systemRegisterController := controllers.NewSystemRegisterController(systemRegisterSvc)

	slaEscalationMatrixRepo := repositories.NewSLAEscalationMatrixRepository(db)
	slaEscalationMatrixSvc := services.NewSLAEscalationMatrixService(slaEscalationMatrixRepo, sharedModule.ActivityLogService)
	slaEscalationMatrixController := controllers.NewSLAEscalationMatrixController(slaEscalationMatrixSvc)

	qualysVulnerabilityHistoryRepo := repositories.NewQulaysVulnerabilityHistoryRepository(db)
	qualysVulnerabilityHistorySvc := services.NewQualysVulnerabilityHistoryService(qualysVulnerabilityHistoryRepo)
	qualysVulnerabilityHistoryController := controllers.NewQualysVulnerabilityHistoryController(qualysVulnerabilityHistorySvc)

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
		VulnerabilityActionPlanCommentController:      vulnerabilityActionPlanCommentController,
		ThreatActionPlanController:                    threatActionPlanController,
		ThreatActionPlanCommentController:             threatActionPlanCommentController,
		ThreatActionPlanStatusUpdateController:        threatActionPlanStatusUpdateController,
		AssetMaintenanceController:                    assetMaintenanceController,
		AssetAssociatedItemController:                 assetAssociatedItemController,
		AssetCriticalityController:                    assetCriticalityController,
		AssetPolicyController:                         assetPolicyController,
		AssetPolicyChecklistController:                assetPolicyChecklistController,
		AssetPolicyCommentController:                  assetPolicyCommentController,
		AssetDocumentCategoryController:               assetDocumentCategoryController,
		AssetOtherDocumentController:                  assetOtherDocumentController,
		AssetControlAnalysisController:                assetControlAnalysisController,
		AssetKriController:                            assetKriController,
		TenableController:                             tenableController,
		SystemRegisterDisruptionLevelController:       systemRegisterDisruptionLevelController,
		SystemRegisterServiceImportanceController:     systemRegisterServiceImportanceController,
		SystemRegisterController:                      systemRegisterController,
		AssetDeviceCategoryController:                 assetDeviceCategoryController,
		ConfidentialityClassificationController:       confidentialityClassificationController,
		AssetEnvironmentController:                    assetEnvironmentController,
		ApplicationController:                         applicationController,
		SLAEscalationMatrixController:                 slaEscalationMatrixController,
		VulnerabilitySettingsController:               vulnerabilitySettingsController,
		QualysVulnerabilityHistoryController:          qualysVulnerabilityHistoryController,

		AssetRegisterSvc: assetRegisterSvc,
	}
}
