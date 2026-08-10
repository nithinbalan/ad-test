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
	VulnerabilityAssociatedItemController         *controllers.VulnerabilityAssociatedItemController
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
	QualysController                              *controllers.QualysController
	SystemRegisterDisruptionLevelController       *controllers.SystemRegisterDisruptionLevelController
	SystemRegisterServiceImportanceController     *controllers.SystemRegisterServiceImportanceController
	AssetDeviceCategoryController                 *controllers.AssetDeviceCategoryController
	SystemRegisterController                      *controllers.SystemRegisterController
	ConfidentialityClassificationController       *controllers.ConfidentialityClassificationController
	AssetEnvironmentController                    *controllers.AssetEnvironmentController
	ApplicationController                         *controllers.ApplicationController
	SLAEscalationMatrixController                 *controllers.SLAEscalationMatrixController
	VulnerabilitySettingsController               *controllers.VulnerabilitySettingsController
	VulnerabilityReportController                 *controllers.VulnerabilityReportController
	QualysVulnerabilityHistoryController          *controllers.QualysVulnerabilityHistoryController
	AssetRiskRegisterController                   *controllers.AssetRiskRegisterController
	AssetRiskWorkflowController                   *controllers.AssetRiskWorkflowController
	AssetRiskReviewController                     *controllers.AssetRiskReviewController
	AssetRiskReviewCommentController              *controllers.AssetRiskReviewCommentController
	VulnerabilityClosedWorkflowController         *controllers.VulnerabilityClosureWorkflowController
	VulnerabilityClosedReviewController           *controllers.VulnerabilityClosureReviewController
	VulnerabilityClosedReviewCommentController    *controllers.VulnerabilityClosureReviewCommentController
	AssetRiskMatrixController                     *controllers.AssetRiskMatrixController
	AssetRiskInherentEvaluationController         *controllers.AssetRiskInherentEvaluationController
	AssetRiskInherentControlAnalysisController    *controllers.AssetRiskInherentControlAnalysisController
	AssetRiskResidualEvaluationController         *controllers.AssetRiskResidualEvaluationController
	AssetRiskResponseController                   *controllers.AssetRiskResponseController
	AssetRiskResponseWorkflowController           *controllers.AssetRiskResponseWorkflowController
	AssetRiskResponseReviewController             *controllers.AssetRiskResponseReviewController
	AssetRiskResponseReviewCommentController      *controllers.AssetRiskResponseReviewCommentController
	AssetRiskActionPlanController                 *controllers.AssetRiskActionPlanController
	AssetRiskActionPlanStatusUpdateController     *controllers.AssetRiskActionPlanStatusUpdateController
	AssetRiskActionPlanCommentController          *controllers.AssetRiskActionPlanCommentController
	AssetWorkflowLevelController                  *controllers.AssetWorkflowLevelController

	// controllers which use service of risk module
	AssetRiskRootCauseAnalysisController *controllers.AssetRiskRootCauseAnalysisController
	AssetRiskAssociatedItemController    *controllers.AssetRiskAssociatedItemController
	AssetRiskKriController               *controllers.AssetRiskKriController
	AssetRiskClosureController           *controllers.AssetRiskClosureController

	AssetRegisterSvc services.AssetRegisterServiceInterface
}

func (am *AssetModule) SetRiskControllers(riskModule *RiskModule) {
	am.AssetRiskRootCauseAnalysisController = controllers.NewAssetRiskRootCauseAnalysisController(riskModule.RiskRootCauseAnalysisSvc)
	am.AssetRiskAssociatedItemController = controllers.NewAssetRiskAssociatedItemController(riskModule.RiskAssociatedItemSvc)
	am.AssetRiskKriController = controllers.NewAssetRiskKriController(riskModule.RiskKriSvc)
	am.AssetRiskClosureController = controllers.NewAssetRiskClosureController(riskModule.RiskClosureSvc)
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

	qualysRepo := repositories.NewQualysRepository()
	qualysSvc := services.NewQualysService(qualysRepo, settingsModule.QualysSettingSvc)
	qualysController := controllers.NewQualysController(qualysSvc)

	slaEscalationMatrixRepo := repositories.NewSLAEscalationMatrixRepository(db)
	slaEscalationMatrixSvc := services.NewSLAEscalationMatrixService(slaEscalationMatrixRepo, sharedModule.ActivityLogService)
	slaEscalationMatrixController := controllers.NewSLAEscalationMatrixController(slaEscalationMatrixSvc)

	vulnerabilityRepo := repositories.NewVulnerabilityRepository(db)
	vulnerabilitySvc := services.NewVulnerabilityService(
		vulnerabilityRepo,
		sharedModule.ActivityLogService,
		vulnerabilityTypeSvc,
		masterModule.SeverityLevelSvc,
		vulnerabilityStatusSvc,
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc,
		vulnerabilityStatusSvc,
		vulnerabilitySettingsSvc,
		qualysSvc,
		assetRegisterSvc,
		organizationModule.ControlSvc,
		commonModule.NotificationService,
		slaEscalationMatrixSvc,
	)
	vulnerabilityController := controllers.NewVulnerabilityController(vulnerabilitySvc)
	vulnerabilityAssociatedItemRepo := repositories.NewVulnerabilityAssociatedItemRepository(db)
	vulnerabilityAssociatedItemSvc := services.NewVulnerabilityAssociatedItemService(
		vulnerabilityAssociatedItemRepo,
		vulnerabilitySvc,
		commonModule.AssociatedItemValidationService,
	)
	vulnerabilityAssociatedItemController := controllers.NewVulnerabilityAssociatedItemController(vulnerabilityAssociatedItemSvc)

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

	vulnerabilityReportRepo := repositories.NewVulnerabilityReportRepository(db)
	vulnerabilityReportSvc := services.NewVulnerabilityReportService(vulnerabilityReportRepo, sharedModule.ActivityLogService, settingsModule.ReportTemplateSettingSvc)
	vulnerabilityReportController := controllers.NewVulnerabilityReportController(vulnerabilityReportSvc)

	qualysVulnerabilityHistoryRepo := repositories.NewQulaysVulnerabilityHistoryRepository(db)
	qualysVulnerabilityHistorySvc := services.NewQualysVulnerabilityHistoryService(qualysVulnerabilityHistoryRepo, vulnerabilitySvc)
	qualysVulnerabilityHistoryController := controllers.NewQualysVulnerabilityHistoryController(qualysVulnerabilityHistorySvc)

	assetRiskRegisterRepo := repositories.NewAssetRiskRegisterRepository(db)
	assetRiskRegisterSvc := services.NewAssetRiskRegisterService(assetRiskRegisterRepo, settingsModule.OrganizationUnitSvc, organizationModule.UserSvc, organizationModule.UserSettingSvc, assetRegisterSvc, vulnerabilitySvc, sharedModule.ActivityLogService, commonModule.NotificationService)
	assetRiskReigsterController := controllers.NewAssetRiskRegisterController(assetRiskRegisterSvc)

	assetRiskWorkflowRepo := repositories.NewAssetRiskWorkflowRepository(db)
	assetRiskWorkflowSvc := services.NewAssetRiskWorkflowService(
		assetRiskWorkflowRepo,
		settingsModule.WorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		masterModule.WorkflowContextSvc,
		assetRiskRegisterSvc,
		organizationModule.OrganizationUnitsSvc,
		settingsModule.OrganizationUnitSvc,
	)
	assetRiskWorkflowController := controllers.NewAssetRiskWorkflowController(assetRiskWorkflowSvc)

	assetRiskReviewRepo := repositories.NewAssetRiskReviewRepository(db)
	assetRiskReviewSvc := services.NewAssetRiskReviewService(
		assetRiskReviewRepo,
		assetRiskRegisterSvc,
		assetRiskWorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		vulnerabilityStatusSvc,
	)
	assetRiskReviewController := controllers.NewAssetRiskReviewController(assetRiskReviewSvc)

	assetRiskReviewCommentRepo := repositories.NewAssetRiskReviewCommentRepository(db)
	assetRiskReviewCommentSvc := services.NewAssetRiskReviewCommentService(
		assetRiskReviewCommentRepo,
		sharedModule.ActivityLogService,
		assetRiskRegisterSvc,
		masterModule.CommentStatusSvc,
	)
	assetRiskReviewCommentController := controllers.NewAssetRiskReviewCommentController(assetRiskReviewCommentSvc)

	vulnerabilityClosureRepo := repositories.NewVulnerabilityClosureRepository(db)
	vulnerabilityClosureSvc := services.NewVulnerabilityClosureService(vulnerabilityClosureRepo)

	vulnerabilityClosedWorkflowRepo := repositories.NewVulnerabilityClosureWorkflowRepository(db)
	vulnerabilityClosedWorkflowSvc := services.NewVulnerabilityClosureWorkflowService(
		vulnerabilityClosedWorkflowRepo,
		sharedModule.ActivityLogService,
		settingsModule.WorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		masterModule.WorkflowContextSvc,
		organizationModule.OrganizationUnitsSvc,
		settingsModule.OrganizationUnitSvc,
	)
	vulnerabilityClosedWorkflowSvc.SetVulnerabilityClosureService(vulnerabilityClosureSvc)
	vulnerabilityClosedWorkflowController := controllers.NewVulnerabilityClosureWorkflowController(vulnerabilityClosedWorkflowSvc)
	vulnerabilityClosedReviewRepo := repositories.NewVulnerabilityClosureReviewRepository(db)
	vulnerabilityClosedReviewSvc := services.NewVulnerabilityClosureReviewService(vulnerabilityClosedReviewRepo, sharedModule.ActivityLogService, vulnerabilitySvc, vulnerabilityClosedWorkflowSvc, vulnerabilityClosureSvc, masterModule.WorkflowStatusSvc, vulnerabilityStatusSvc, organizationModule.UserSvc, commonModule.NotificationService, taskModule.TaskService, organizationModule.OrganizationUnitsSvc)
	vulnerabilityClosedReviewController := controllers.NewVulnerabilityClosureReviewController(vulnerabilityClosedReviewSvc)
	vulnerabilityClosedReviewCommentRepo := repositories.NewVulnerabilityClosureReviewCommentRepository(db)
	vulnerabilityClosedReviewCommentSvc := services.NewVulnerabilityClosureReviewCommentService(vulnerabilityClosedReviewCommentRepo, organizationModule.UserSvc, sharedModule.ActivityLogService, vulnerabilityClosureSvc, masterModule.CommentStatusSvc)
	vulnerabilityClosedReviewCommentController := controllers.NewVulnerabilityClosureReviewCommentController(vulnerabilityClosedReviewCommentSvc)

	assetRiskMatrixRepo := repositories.NewAssetRiskMatrixRepository(db)
	assetRiskMatrixSvc := services.NewAssetRiskMatrixService(assetRiskMatrixRepo)
	assetRiskMatrixController := controllers.NewAssetRiskMatrixController(assetRiskMatrixSvc)

	assetRiskInherentEvaluationRepo := repositories.NewAssetRiskInherentEvaluationRepository(db)
	assetRiskInherentEvaluationSvc := services.NewAssetRiskInherentEvaluationService(assetRiskInherentEvaluationRepo, assetRiskRegisterSvc, sharedModule.ActivityLogService)
	assetRiskInherentEvaluationController := controllers.NewAssetRiskInherentEvaluationController(assetRiskInherentEvaluationSvc)

	assetRiskInherentControlAnalysisRepo := repositories.NewAssetRiskInherentControlAnalysisRepository(db)
	assetRiskInherentControlAnalysisSvc := services.NewAssetRiskInherentControlAnalysisService(assetRiskInherentControlAnalysisRepo, sharedModule.ActivityLogService, assetRiskRegisterSvc, organizationModule.ControlSvc, organizationModule.ControlTestSvc)
	assetRiskInherentControlAnalysisController := controllers.NewAssetRiskInherentControlAnalysisController(assetRiskInherentControlAnalysisSvc)

	assetRiskResidualEvaluationRepo := repositories.NewAssetRiskResidualEvaluationRepository(db)
	assetRiskResidualEvaluationSvc := services.NewAssetRiskResidualEvaluationService(assetRiskResidualEvaluationRepo, assetRiskRegisterSvc, assetRiskInherentEvaluationSvc, sharedModule.ActivityLogService, organizationModule.UserSvc, settingsModule.RoleSvc, commonModule.NotificationService)
	assetRiskResidualEvaluationController := controllers.NewAssetRiskResidualEvaluationController(assetRiskResidualEvaluationSvc)

	assetRiskResponseRepo := repositories.NewAssetRiskResponseRepository(db)
	assetRiskResponseSvc := services.NewAssetRiskResponseService(assetRiskResponseRepo, assetRiskRegisterSvc, sharedModule.ActivityLogService)
	assetRiskResponseController := controllers.NewAssetRiskResponseController(assetRiskResponseSvc)

	assetRiskResponseWorkflowRepo := repositories.NewAssetRiskResponseWorkflowRepository(db)
	assetRiskResponseWorkflowSvc := services.NewAssetRiskResponseWorkflowService(
		assetRiskResponseWorkflowRepo,
		settingsModule.WorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		masterModule.WorkflowContextSvc,
	)
	assetRiskResponseWorkflowController := controllers.NewAssetRiskResponseWorkflowController(assetRiskResponseWorkflowSvc)

	assetRiskActionPlanRepo := repositories.NewAssetRiskActionPlanRepository(db)
	assetRiskActionPlanSvc := services.NewAssetRiskActionPlanService(
		assetRiskActionPlanRepo,
		sharedModule.ActivityLogService,
		assetRiskRegisterSvc,
		masterModule.ActionPlanStatusSvc,
		masterModule.SeverityLevelSvc,
		organizationModule.UserSettingSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		taskModule.TaskService,
		organizationModule.ControlSvc,
		organizationModule.ControlImplementationStatusSvc,
		organizationModule.ControlEfficiencyStatusSvc,
		commonModule.DelegateeUserService,
	)
	assetRiskActionPlanController := controllers.NewAssetRiskActionPlanController(assetRiskActionPlanSvc)

	assetRiskResponseReviewRepo := repositories.NewAssetRiskResponseReviewRepository(db)
	assetRiskResponseReviewSvc := services.NewAssetRiskResponseReviewService(
		assetRiskResponseReviewRepo,
		assetRiskRegisterSvc,
		assetRiskResponseSvc,
		assetRiskResponseWorkflowSvc,
		assetRiskActionPlanSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
	)
	assetRiskResponseReviewController := controllers.NewAssetRiskResponseReviewController(assetRiskResponseReviewSvc)

	assetRiskResponseReviewCommentRepo := repositories.NewAssetRiskResponseReviewCommentRepository(db)
	assetRiskResponseReviewCommentSvc := services.NewAssetRiskResponseReviewCommentService(
		assetRiskResponseReviewCommentRepo,
		sharedModule.ActivityLogService,
		assetRiskResponseSvc,
		masterModule.CommentStatusSvc,
	)
	assetRiskResponseReviewCommentController := controllers.NewAssetRiskResponseReviewCommentController(assetRiskResponseReviewCommentSvc)

	assetRiskActionPlanStatusUpdateRepo := repositories.NewAssetRiskActionPlanStatusUpdateRepository(db)
	assetRiskActionPlanStatusUpdateSvc := services.NewAssetRiskActionPlanStatusUpdateService(
		assetRiskActionPlanStatusUpdateRepo,
		sharedModule.ActivityLogService,
		assetRiskRegisterSvc,
		assetRiskActionPlanSvc,
		masterModule.ActionPlanStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		taskModule.TaskService,
		taskModule.TaskStatusService,
	)
	assetRiskActionPlanStatusUpdateController := controllers.NewAssetRiskActionPlanStatusUpdateController(assetRiskActionPlanStatusUpdateSvc)

	assetRiskActionPlanCommentRepo := repositories.NewAssetRiskActionPlanCommentRepository(db)
	assetRiskActionPlanCommentSvc := services.NewAssetRiskActionPlanCommentService(
		assetRiskActionPlanCommentRepo,
		organizationModule.UserSvc,
		sharedModule.ActivityLogService,
		assetRiskActionPlanSvc,
		commonModule.NotificationService,
	)
	assetRiskActionPlanCommentController := controllers.NewAssetRiskActionPlanCommentController(assetRiskActionPlanCommentSvc)

	assetWorkflowLevelController := controllers.NewAssetWorkflowLevelController(settingsModule.WorkflowSvc)

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
		VulnerabilityAssociatedItemController:         vulnerabilityAssociatedItemController,
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
		QualysController:                              qualysController,
		SystemRegisterDisruptionLevelController:       systemRegisterDisruptionLevelController,
		SystemRegisterServiceImportanceController:     systemRegisterServiceImportanceController,
		SystemRegisterController:                      systemRegisterController,
		AssetDeviceCategoryController:                 assetDeviceCategoryController,
		ConfidentialityClassificationController:       confidentialityClassificationController,
		AssetEnvironmentController:                    assetEnvironmentController,
		ApplicationController:                         applicationController,
		SLAEscalationMatrixController:                 slaEscalationMatrixController,
		VulnerabilitySettingsController:               vulnerabilitySettingsController,
		VulnerabilityReportController:                 vulnerabilityReportController,
		QualysVulnerabilityHistoryController:          qualysVulnerabilityHistoryController,
		AssetRiskRegisterController:                   assetRiskReigsterController,
		AssetRiskWorkflowController:                   assetRiskWorkflowController,
		AssetRiskReviewController:                     assetRiskReviewController,
		AssetRiskReviewCommentController:              assetRiskReviewCommentController,
		VulnerabilityClosedWorkflowController:         vulnerabilityClosedWorkflowController,
		VulnerabilityClosedReviewController:           vulnerabilityClosedReviewController,
		VulnerabilityClosedReviewCommentController:    vulnerabilityClosedReviewCommentController,
		AssetRiskMatrixController:                     assetRiskMatrixController,
		AssetRiskInherentEvaluationController:         assetRiskInherentEvaluationController,
		AssetRiskInherentControlAnalysisController:    assetRiskInherentControlAnalysisController,
		AssetRiskResidualEvaluationController:         assetRiskResidualEvaluationController,
		AssetRiskResponseController:                   assetRiskResponseController,
		AssetRiskResponseWorkflowController:           assetRiskResponseWorkflowController,
		AssetRiskResponseReviewController:             assetRiskResponseReviewController,
		AssetRiskResponseReviewCommentController:      assetRiskResponseReviewCommentController,
		AssetRiskActionPlanController:                 assetRiskActionPlanController,
		AssetRiskActionPlanStatusUpdateController:     assetRiskActionPlanStatusUpdateController,
		AssetRiskActionPlanCommentController:          assetRiskActionPlanCommentController,
		AssetWorkflowLevelController:                  assetWorkflowLevelController,

		AssetRegisterSvc: assetRegisterSvc,
	}
}
