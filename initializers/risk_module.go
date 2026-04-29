package initializers

import (
	"jethur/grc/risk/controllers"
	"jethur/grc/risk/repositories"
	"jethur/grc/risk/services"

	"gorm.io/gorm"
)

type RiskModule struct {
	RiskCategoryController                 *controllers.RiskCategoryController
	PotentialRiskStatusController          *controllers.PotentialRiskStatusController
	RiskLevelController                    *controllers.RiskLevelController
	RiskSubCategoryController              *controllers.RiskSubCategoryController
	LevelOfMeasurementController           *controllers.LevelOfMeasurementController
	DistributionStatusController           *controllers.DistributionStatusController
	PotentialRiskController                *controllers.PotentialRiskController
	PotentialRiskStatusUpdateController    *controllers.PotentialRiskStatusUpdateController
	KriCategoryController                  *controllers.KRICategoryController
	KriRegisterController                  *controllers.KRIRegisterController
	KriValueUpdateController               *controllers.KRIValueUpdateController
	RiskAppetiteStatementController        *controllers.RiskAppetiteStatementController
	RiskImpactCategoryController           *controllers.RiskImpactCategoryController
	RiskMatrixImpactController             *controllers.RiskMatrixImpactController
	RiskMatrixImpactGuidelineController    *controllers.RiskMatrixImpactGuidelineController
	RiskMatrixLikelihoodController         *controllers.RiskMatrixLikelihoodController
	RiskMatrixRiskRatingController         *controllers.RiskMatrixRiskRatingController
	RiskMatrixRiskRatingValueController    *controllers.RiskMatrixRiskRatingValueController
	RiskSourceController                   *controllers.RiskSourceController
	RiskStatusController                   *controllers.RiskStatusController
	RiskRegisterUpdatePlanStatusController *controllers.RiskRegisterUpdatePlanStatusController
	RiskRegisterController                 *controllers.RiskRegisterController
	RiskRootCauseAnalysisController        *controllers.RiskRootCauseAnalysisController
	RiskAppetiteToleranceLevelController   *controllers.RiskAppetiteToleranceLevelController
	RiskInherentEvaluationController       *controllers.RiskInherentEvaluationController
	RiskActionPlanController               *controllers.RiskActionPlanController
	RiskActionPlanStatusUpdateController   *controllers.RiskActionPlanStatusUpdateController
	RiskActionPlanCommentController        *controllers.RiskActionPlanCommentController
	RiskAssociatedItemController           *controllers.RiskAssociatedItemController
	RiskResponseStrategyController         *controllers.RiskResponseStrategyController
	RiskTreatmentGuidelineController       *controllers.RiskTreatmentGuidelineController
	RiskHeatmapController                  *controllers.RiskHeatmapController
	RiskInherentControlAnalysisController  *controllers.RiskInherentControlAnalysisController
	RiskResidualEvaluationController       *controllers.RiskResidualEvaluationController
	RiskResidualControlAnalysisController  *controllers.RiskResidualControlAnalysisController
	RiskKriController                      *controllers.RiskKriController
	RiskTimelineController                 *controllers.RiskTimelineController
	RiskSettingsController                 *controllers.RiskSettingsController
	RiskClosureController                  *controllers.RiskClosureController
	RiskCommentController                  *controllers.RiskCommentController
	RiskReportController                   *controllers.RiskReportController
	KRIAssociatedItemController            *controllers.KRIAssociatedItemController
	KRIRegisterWorkflowController          *controllers.KRIRegisterWorkflowController
	KRIRegisterReviewController            *controllers.KRIRegisterReviewController
	KRIRegisterReviewCommentController     *controllers.KRIRegisterReviewCommentController
	RiskWorkflowController                 *controllers.RiskWorkflowController
	RiskReviewController                   *controllers.RiskReviewController
	RiskReviewCommentController            *controllers.RiskReviewCommentController

	RiskRegisterSvc         services.RiskRegisterServiceInterface
	RiskHeatMapSvc          services.RiskHeatmapServiceInterface
	KriRegisterSvc          services.KRIRegisterServiceInterface
	KRIWorkflowSvc          services.KRIRegisterWorkflowServiceInterface
	KRIReviewSvc            services.KRIRegisterReviewServiceInterface
	KRIReviewCommentSvc     services.KRIRegisterReviewCommentServiceInterface
	RiskActionPlanSvc       services.RiskActionPlanServiceInterface
	RiskLevelSvc            services.RiskLevelServiceInterface
	RiskResponseStrategySvc services.RiskResponseStrategyServiceInterface
	RiskWorkflowSvc         services.RiskWorkflowServiceInterface
	RiskReviewSvc           services.RiskReviewServiceInterface
	RiskReviewCommentSvc    services.RiskReviewCommentServiceInterface
}

func NewRiskModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	libraryModule *LibraryModule,
	organizationModule *OrganizationModule,
	complianceModule *ComplianceModule,
	commonModule *CommonModule,
	taskModule *TaskModule,
	assetModule *AssetModule,
) *RiskModule {
	riskCategoryRepo := repositories.NewRiskCategoryRepository(db)
	riskCategorySvc := services.NewRiskCategoryService(riskCategoryRepo, sharedModule.ActivityLogService)
	riskCategoryController := controllers.NewRiskCategoryController(riskCategorySvc)

	potentialRiskStatusRepo := repositories.NewPotentialRiskStatusRepository(db)
	potentialRiskStatusSvc := services.NewPotentialRiskStatusService(potentialRiskStatusRepo)
	potentialRiskStatusController := controllers.NewPotentialRiskStatusController(potentialRiskStatusSvc)

	riskLevelRepo := repositories.NewRiskLevelRepository(db)
	riskLevelSvc := services.NewRiskLevelService(riskLevelRepo)
	riskLevelController := controllers.NewRiskLevelController(riskLevelSvc)

	riskSubCategoryRepo := repositories.NewRiskSubCategoryRepository(db)
	riskSubCategorySvc := services.NewRiskSubCategoryService(riskSubCategoryRepo, riskCategoryRepo, sharedModule.ActivityLogService)
	riskSubCategoryController := controllers.NewRiskSubCategoryController(riskSubCategorySvc)

	levelOfMeasurementRepo := repositories.NewLevelOfMeasurementRepository(db)
	levelOfMeasurementSvc := services.NewLevelOfMeasurementService(levelOfMeasurementRepo)
	levelOfMeasurementController := controllers.NewLevelOfMeasurementController(levelOfMeasurementSvc)

	distributionStatusRepo := repositories.NewDistributionStatusRepository(db)
	distributionStatusSvc := services.NewDistributionStatusService(distributionStatusRepo)
	distributionStatusController := controllers.NewDistributionStatusController(distributionStatusSvc)

	kriCategoryRepository := repositories.NewKRICategoryRepository(db)
	kriCategorySvc := services.NewKRICategoryService(kriCategoryRepository, sharedModule.ActivityLogService)
	kriCategoryController := controllers.NewKRICategoryController(kriCategorySvc)

	// set kri categories service to kri library service
	libraryModule.KriLibrarySvc.SetKriCategoryService(kriCategorySvc)

	kriRegisterRepo := repositories.NewKRIRegisterRepository(db)
	kriRegisterSvc := services.NewKRIRegisterService(
		kriRegisterRepo,
		sharedModule.ActivityLogService,
		kriCategorySvc,
		organizationModule.UserSvc,
		riskLevelSvc,
		riskSubCategorySvc,
		levelOfMeasurementSvc,
		distributionStatusSvc,
		settingsModule.OrganizationUnitSvc,
		masterModule.ReviewFrequencySvc,
		masterModule.UnitSvc,
		organizationModule.UserSettingSvc,
		commonModule.NotificationService,
	)
	kriRegisterController := controllers.NewKRIRegisterController(kriRegisterSvc)

	kriRegisterWorkflowRepo := repositories.NewKRIRegisterWorkflowRepository(db)
	kriRegisterWorkflowSvc := services.NewKRIRegisterWorkflowService(
		kriRegisterWorkflowRepo,
		sharedModule.ActivityLogService,
		settingsModule.WorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		masterModule.WorkflowContextSvc,
	)
	kriRegisterSvc.SetKRIWorkflowService(kriRegisterWorkflowSvc)
	kriRegisterWorkflowSvc.SetKRIRegisterService(kriRegisterSvc)
	kriRegisterWorkflowController := controllers.NewKRIRegisterWorkflowController(kriRegisterWorkflowSvc)

	kriRegisterReviewRepo := repositories.NewKRIRegisterReviewRepository(db)
	kriRegisterReviewSvc := services.NewKRIRegisterReviewService(
		kriRegisterReviewRepo,
		sharedModule.ActivityLogService,
		kriRegisterSvc,
		kriRegisterWorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		taskModule.TaskService,
	)
	kriRegisterReviewController := controllers.NewKRIRegisterReviewController(kriRegisterReviewSvc)

	kriRegisterReviewCommentRepo := repositories.NewKRIRegisterReviewCommentRepository(db)
	kriRegisterReviewCommentSvc := services.NewKRIRegisterReviewCommentService(
		kriRegisterReviewCommentRepo,
		organizationModule.UserSvc,
		sharedModule.ActivityLogService,
		kriRegisterSvc,
		masterModule.CommentStatusSvc,
	)
	kriRegisterReviewCommentController := controllers.NewKRIRegisterReviewCommentController(kriRegisterReviewCommentSvc)

	kriValueUpdateRepo := repositories.NewKRIValueUpdateRepository(db)
	kriValueUpdateSvc := services.NewKRIValueUpdateService(kriValueUpdateRepo, sharedModule.ActivityLogService, kriRegisterSvc, distributionStatusSvc, masterModule.ReviewFrequencySvc, organizationModule.UserSvc, commonModule.NotificationService)
	kriValueUpdateController := controllers.NewKRIValueUpdateController(kriValueUpdateSvc)

	riskAppetiteStatementRepo := repositories.NewRiskAppetiteStatementRepository(db)
	riskAppetiteStatementSvc := services.NewRiskAppetiteStatementService(
		riskAppetiteStatementRepo,
		sharedModule.ActivityLogService,
		riskCategorySvc,
	)
	riskAppetiteStatementController := controllers.NewRiskAppetiteStatementController(riskAppetiteStatementSvc)

	riskImpactCategoryRepo := repositories.NewRiskImpactCategoryRepository(db)
	riskImpactCategorySvc := services.NewRiskImpactCategoryService(riskImpactCategoryRepo, sharedModule.ActivityLogService)
	riskImpactCategoryController := controllers.NewRiskImpactCategoryController(riskImpactCategorySvc)

	riskMatrixRepo := repositories.NewRiskMatrixRepository(db)
	riskMatrixSvc := services.NewRiskMatrixService(riskMatrixRepo)

	riskMatrixImpactRepo := repositories.NewRiskMatrixImpactRepository(db)
	riskMatrixImpactSvc := services.NewRiskMatrixImpactService(
		riskMatrixImpactRepo,
		sharedModule.ActivityLogService,
		riskMatrixSvc,
	)
	riskMatrixImpactController := controllers.NewRiskMatrixImpactController(riskMatrixImpactSvc)

	riskMatrixImpactGuidelineRepo := repositories.NewRiskMatrixImpactGuidelineRepository(db)
	riskMatrixImpactGuidelineSvc := services.NewRiskMatrixImpactGuidelineService(
		riskMatrixImpactGuidelineRepo,
		sharedModule.ActivityLogService,
		riskImpactCategorySvc,
		riskMatrixImpactSvc,
	)
	riskMatrixImpactGuidelineController := controllers.NewRiskMatrixImpactGuidelineController(riskMatrixImpactGuidelineSvc)

	riskMatrixLikelihoodRepo := repositories.NewRiskMatrixLikelihoodRepository(db)
	riskMatrixLikelihoodSvc := services.NewRiskMatrixLikelihoodService(
		riskMatrixLikelihoodRepo,
		sharedModule.ActivityLogService,
		riskMatrixSvc,
	)
	riskMatrixLikelihoodController := controllers.NewRiskMatrixLikelihoodController(riskMatrixLikelihoodSvc)

	riskMatrixRiskRatingRepo := repositories.NewRiskMatrixRiskRatingRepository(db)
	riskMatrixRiskRatingSvc := services.NewRiskMatrixRiskRatingService(
		riskMatrixRiskRatingRepo,
		sharedModule.ActivityLogService,
		riskMatrixSvc,
	)
	riskMatrixRiskRatingController := controllers.NewRiskMatrixRiskRatingController(riskMatrixRiskRatingSvc)

	riskMatrixRiskRatingValueRepo := repositories.NewRiskMatrixRiskRatingValueRepository(db)
	riskMatrixRiskRatingValueSvc := services.NewRiskMatrixRiskRatingValueService(
		riskMatrixRiskRatingValueRepo,
		sharedModule.ActivityLogService,
		riskMatrixSvc,
		riskMatrixRiskRatingSvc,
	)
	riskMatrixRiskRatingValueController := controllers.NewRiskMatrixRiskRatingValueController(riskMatrixRiskRatingValueSvc)

	riskSourceRepo := repositories.NewRiskSourceRepository(db)
	riskSourceSvc := services.NewRiskSourceService(riskSourceRepo, sharedModule.ActivityLogService)
	riskSourceController := controllers.NewRiskSourceController(riskSourceSvc)

	riskStatusRepo := repositories.NewRiskStatusRepository(db)
	riskStatusSvc := services.NewRiskStatusService(riskStatusRepo)
	riskStatusController := controllers.NewRiskStatusController(riskStatusSvc)

	riskRegisterUpdatePlanStatusRepo := repositories.NewRiskRegisterUpdatePlanStatusRepository(db)
	riskRegisterUpdatePlanStatusSvc := services.NewRiskRegisterUpdatePlanStatusService(riskRegisterUpdatePlanStatusRepo)
	riskRegisterUpdatePlanStatusController := controllers.NewRiskRegisterUpdatePlanStatusController(riskRegisterUpdatePlanStatusSvc)

	riskWorkflowRepo := repositories.NewRiskWorkflowRepository(db)
	riskWorkflowSvc := services.NewRiskWorkflowService(
		riskWorkflowRepo,
		sharedModule.ActivityLogService,
		settingsModule.WorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		masterModule.WorkflowContextSvc,
	)

	riskRegisterRepo := repositories.NewRiskRegisterRepository(db)

	riskTimelineRepo := repositories.NewRiskTimelineRepository(db)
	riskTimelineSvc := services.NewRiskTimelineService(riskTimelineRepo, riskRegisterRepo)
	riskTimelineController := controllers.NewRiskTimelineController(riskTimelineSvc)

	riskSettingsRepo := repositories.NewRiskSettingsRepository(db)
	riskSettingsSvc := services.NewRiskSettingsService(riskSettingsRepo, sharedModule.ActivityLogService)
	riskSettingsController := controllers.NewRiskSettingsController(riskSettingsSvc)

	riskResponseStrategyRepo := repositories.NewRiskResponseStrategyRepository(db)
	riskResponseStrategySvc := services.NewRiskResponseStrategyService(riskResponseStrategyRepo, sharedModule.ActivityLogService)
	riskResponseStrategyController := controllers.NewRiskResponseStrategyController(riskResponseStrategySvc)

	riskReportRepo := repositories.NewRiskReportRepository(db)

	riskRegisterSvc := services.NewRiskRegisterService(
		riskRegisterRepo,
		sharedModule.ActivityLogService,
		riskCategorySvc,
		riskSourceSvc,
		organizationModule.UserSvc,
		riskStatusSvc,
		masterModule.ReviewFrequencySvc,
		organizationModule.ControlSvc,
		complianceModule.ComplianceRequirementSvc,
		organizationModule.ProcessSvc,
		organizationModule.StrategicObjectiveSvc,
		complianceModule.SlaRegisterSvc,
		riskMatrixSvc,
		riskTimelineSvc,
		organizationModule.UserSettingSvc,
		commonModule.NotificationService,
		assetModule.AssetRegisterSvc,
		nil,
		riskSettingsSvc,
		settingsModule.OrganizationUnitSvc,
		organizationModule.OrganizationUnitsSvc,
		riskMatrixLikelihoodSvc,
		riskMatrixImpactSvc,
		riskResponseStrategySvc,
		nil,
		masterModule.SeverityLevelSvc,
		riskReportRepo,
		settingsModule.RiskTemplateSettingSvc,
	)
	riskRegisterController := controllers.NewRiskRegisterController(riskRegisterSvc)

	riskRegisterSvc.SetRiskWorkflowService(riskWorkflowSvc)
	riskWorkflowSvc.SetRiskService(riskRegisterSvc)
	riskWorkflowController := controllers.NewRiskWorkflowController(riskWorkflowSvc)

	riskReviewRepo := repositories.NewRiskReviewRepository(db)
	riskReviewSvc := services.NewRiskReviewService(
		riskReviewRepo,
		sharedModule.ActivityLogService,
		riskRegisterSvc,
		riskWorkflowSvc,
		masterModule.WorkflowStatusSvc,
		riskStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		taskModule.TaskService,
		organizationModule.OrganizationUnitsSvc)
	riskReviewController := controllers.NewRiskReviewController(riskReviewSvc)

	riskReviewCommentRepo := repositories.NewRiskReviewCommentRepository(db)
	riskReviewCommentSvc := services.NewRiskReviewCommentService(riskReviewCommentRepo, organizationModule.UserSvc, sharedModule.ActivityLogService, riskRegisterSvc, masterModule.CommentStatusSvc)
	riskReviewCommentController := controllers.NewRiskReviewCommentController(riskReviewCommentSvc)

	potentialRiskRepo := repositories.NewPotentialRiskRepository(db)
	potentialRiskSvc := services.NewPotentialRiskService(
		potentialRiskRepo,
		sharedModule.ActivityLogService,
		masterModule.SeverityLevelSvc,
		riskCategorySvc,
		potentialRiskStatusSvc,
		organizationModule.UserSettingSvc,
		riskRegisterSvc,
	)
	potentialRiskController := controllers.NewPotentialRiskController(potentialRiskSvc)

	potentialRiskStatusUpdateRepo := repositories.NewPotentialRiskStatusUpdateRepository(db)
	potentialRiskStatusUpdateSvc := services.NewPotentialRiskStatusUpdateService(
		potentialRiskStatusUpdateRepo,
		sharedModule.ActivityLogService,
		potentialRiskSvc,
		potentialRiskStatusSvc)
	potentialRiskStatusUpdateController := controllers.NewPotentialRiskStatusUpdateController(potentialRiskStatusUpdateSvc)

	riskRootCauseAnalysisRepo := repositories.NewRiskRootCauseAnalysisRepository(db)
	riskRootCauseAnalysisSvc := services.NewRiskRootCauseAnalysisService(riskRootCauseAnalysisRepo, sharedModule.ActivityLogService, riskRegisterSvc, riskTimelineSvc)
	riskRootCauseAnalysisController := controllers.NewRiskRootCauseAnalysisController(riskRootCauseAnalysisSvc)

	riskAppetiteToleranceLevelRepo := repositories.NewRiskAppetiteToleranceLevelRepository(db)
	riskAppetiteToleranceLevelSvc := services.NewRiskAppetiteToleranceLevelService(riskAppetiteToleranceLevelRepo)
	riskAppetiteToleranceLevelController := controllers.NewRiskAppetiteToleranceLevelController(riskAppetiteToleranceLevelSvc)

	riskInherentEvaluationRepo := repositories.NewRiskInherentEvaluationRepository(db)
	riskInherentEvaluationSvc := services.NewRiskInherentEvaluationService(riskInherentEvaluationRepo, sharedModule.ActivityLogService, riskRegisterSvc, riskMatrixSvc, riskMatrixRiskRatingValueSvc, riskMatrixImpactSvc, riskMatrixLikelihoodSvc, riskImpactCategorySvc, riskMatrixImpactGuidelineSvc, riskAppetiteStatementSvc, riskTimelineSvc)
	riskInherentEvaluationController := controllers.NewRiskInherentEvaluationController(riskInherentEvaluationSvc)

	riskRegisterSvc.SetRiskInherentEvaluationService(riskInherentEvaluationSvc)

	riskActionPlanRepo := repositories.NewRiskActionPlanRepository(db)
	riskActionPlanSvc := services.NewRiskActionPlanService(
		riskActionPlanRepo,
		sharedModule.ActivityLogService,
		riskRegisterSvc,
		masterModule.ActionPlanStatusSvc,
		masterModule.SeverityLevelSvc,
		organizationModule.UserSettingSvc,
		organizationModule.UserSvc, commonModule.NotificationService, taskModule.TaskService, masterModule.LifecycleStatusSvc, organizationModule.ControlSvc, organizationModule.ControlImplementationStatusSvc, organizationModule.ControlEfficiencyStatusSvc, commonModule.DelegateeUserService,
	)
	riskActionPlanController := controllers.NewRiskActionPlanController(riskActionPlanSvc)

	riskActionPlanStatusUpdateRepo := repositories.NewRiskActionPlanStatusUpdateRepository(db)
	riskActionPlanStatusUpdateSvc := services.NewRiskActionPlanStatusUpdateService(
		riskActionPlanStatusUpdateRepo,
		sharedModule.ActivityLogService,
		riskRegisterSvc,
		riskActionPlanSvc,
		masterModule.ActionPlanStatusSvc,
		settingsModule.ModuleSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		taskModule.TaskService,
		taskModule.TaskStatusService,
	)
	riskActionPlanStatusUpdateController := controllers.NewRiskActionPlanStatusUpdateController(riskActionPlanStatusUpdateSvc)

	riskActionPlanCommentRepo := repositories.NewRiskActionPlanCommentRepository(db)
	riskActionPlanCommentSvc := services.NewRiskActionPlanCommentService(riskActionPlanCommentRepo, organizationModule.UserSvc, sharedModule.ActivityLogService, riskActionPlanSvc, commonModule.NotificationService)
	riskActionPlanCommentController := controllers.NewRiskActionPlanCommentController(riskActionPlanCommentSvc)

	riskAssociatedItemRepo := repositories.NewRiskAssociatedItemRepository(db)
	riskAssociatedItemSvc := services.NewRiskAssociatedItemService(riskAssociatedItemRepo, riskRegisterSvc, complianceModule.ComplianceClausevc, commonModule.AssociatedItemValidationService)
	riskAssociatedItemController := controllers.NewRiskAssociatedItemController(riskAssociatedItemSvc)

	riskTreatmentGuidelineRepo := repositories.NewRiskTreatmentGuidelineRepository(db)
	riskTreatmentGuidelineSvc := services.NewRiskTreatmentGuidelineService(riskTreatmentGuidelineRepo, sharedModule.ActivityLogService, riskMatrixRiskRatingSvc)
	riskTreatmentGuidelineController := controllers.NewRiskTreatmentGuidelineController(riskTreatmentGuidelineSvc)

	riskHeatmapRepo := repositories.NewRiskHeatmapRepository(db)
	riskHeatmapSvc := services.NewRiskHeatmapService(riskHeatmapRepo, riskMatrixSvc, riskMatrixRiskRatingValueSvc, riskRegisterSvc)
	riskHeatmapController := controllers.NewRiskHeatmapController(riskHeatmapSvc)

	riskResidualEvaluationRepo := repositories.NewRiskResidualEvaluationRepository(db)
	riskResidualEvaluationSvc := services.NewRiskResidualEvaluationService(riskResidualEvaluationRepo, sharedModule.ActivityLogService, riskRegisterSvc, riskMatrixSvc, riskMatrixRiskRatingValueSvc, riskMatrixImpactSvc, riskMatrixLikelihoodSvc, riskImpactCategorySvc, riskMatrixImpactGuidelineSvc, riskAppetiteStatementSvc, riskInherentEvaluationSvc, riskTimelineSvc, organizationModule.ControlStatusSvc, organizationModule.ControlSvc)
	riskResidualEvaluationController := controllers.NewRiskResidualEvaluationController(riskResidualEvaluationSvc)

	riskInherentControlAnalysisRepo := repositories.NewRiskInherentControlAnalysisRepository(db)
	riskInherentControlAnalysisSvc := services.NewRiskInherentControlAnalysisService(riskInherentControlAnalysisRepo, sharedModule.ActivityLogService, riskRegisterSvc, organizationModule.ControlSvc, organizationModule.ControlTestSvc, riskTimelineSvc, riskResidualEvaluationSvc)
	riskInherentControlAnalysisController := controllers.NewRiskInherentControlAnalysisController(riskInherentControlAnalysisSvc)

	riskRegisterSvc.SetRiskResidualEvaluationService(riskResidualEvaluationSvc)

	riskResidualControlAnalysisRepo := repositories.NewRiskResidualControlAnalysisRepository(db)
	riskResidualControlAnalysisSvc := services.NewRiskResidualControlAnalysisService(riskResidualControlAnalysisRepo, sharedModule.ActivityLogService, riskRegisterSvc, organizationModule.ControlSvc, organizationModule.ControlTestSvc, riskTimelineSvc)
	riskResidualControlAnalysisController := controllers.NewRiskResidualControlAnalysisController(riskResidualControlAnalysisSvc)

	riskKriRepo := repositories.NewRiskKriRepository(db)
	riskKriSvc := services.NewRiskKriService(
		riskKriRepo,
		sharedModule.ActivityLogService,
		riskRegisterSvc,
		kriRegisterSvc,
		riskTimelineSvc,
	)
	riskKriController := controllers.NewRiskKriController(riskKriSvc)

	riskClosureRepo := repositories.NewRiskClosureRepository(db)
	riskClosureSvc := services.NewRiskClosureService(riskClosureRepo, sharedModule.ActivityLogService, riskRegisterSvc, riskStatusSvc, riskTimelineSvc)
	riskClosureController := controllers.NewRiskClosureController(riskClosureSvc)

	riskCommentRepo := repositories.NewRiskCommentRepository(db)
	riskCommentSvc := services.NewRiskCommentService(riskCommentRepo, riskRegisterSvc, riskStatusSvc, sharedModule.ActivityLogService, organizationModule.UserSvc, riskSettingsSvc, commonModule.NotificationService)
	riskCommentController := controllers.NewRiskCommentController(riskCommentSvc)

	riskReportSvc := services.NewRiskReportService(riskReportRepo, sharedModule.ActivityLogService, settingsModule.RiskTemplateSettingSvc)
	riskReportController := controllers.NewRiskReportController(riskReportSvc)

	kriAssociatedItemRepo := repositories.NewKRIAssociatedItemRepository(db)
	kriAssociatedItemSvc := services.NewKRIAssociatedItemService(kriAssociatedItemRepo, kriRegisterSvc, commonModule.AssociatedItemValidationService)
	kriAssociatedItemController := controllers.NewKRIAssociatedItemController(kriAssociatedItemSvc)

	return &RiskModule{
		RiskCategoryController:                 riskCategoryController,
		PotentialRiskStatusController:          potentialRiskStatusController,
		RiskLevelController:                    riskLevelController,
		RiskSubCategoryController:              riskSubCategoryController,
		LevelOfMeasurementController:           levelOfMeasurementController,
		DistributionStatusController:           distributionStatusController,
		PotentialRiskController:                potentialRiskController,
		PotentialRiskStatusUpdateController:    potentialRiskStatusUpdateController,
		KriCategoryController:                  kriCategoryController,
		KriRegisterController:                  kriRegisterController,
		KriValueUpdateController:               kriValueUpdateController,
		RiskAppetiteStatementController:        riskAppetiteStatementController,
		RiskImpactCategoryController:           riskImpactCategoryController,
		RiskMatrixImpactController:             riskMatrixImpactController,
		RiskMatrixImpactGuidelineController:    riskMatrixImpactGuidelineController,
		RiskMatrixLikelihoodController:         riskMatrixLikelihoodController,
		RiskMatrixRiskRatingController:         riskMatrixRiskRatingController,
		RiskMatrixRiskRatingValueController:    riskMatrixRiskRatingValueController,
		RiskSourceController:                   riskSourceController,
		RiskStatusController:                   riskStatusController,
		RiskRegisterUpdatePlanStatusController: riskRegisterUpdatePlanStatusController,
		RiskRegisterController:                 riskRegisterController,
		RiskRootCauseAnalysisController:        riskRootCauseAnalysisController,
		RiskAppetiteToleranceLevelController:   riskAppetiteToleranceLevelController,
		RiskInherentEvaluationController:       riskInherentEvaluationController,
		RiskActionPlanController:               riskActionPlanController,
		RiskActionPlanStatusUpdateController:   riskActionPlanStatusUpdateController,
		RiskActionPlanCommentController:        riskActionPlanCommentController,
		RiskAssociatedItemController:           riskAssociatedItemController,
		RiskResponseStrategyController:         riskResponseStrategyController,
		RiskTreatmentGuidelineController:       riskTreatmentGuidelineController,
		RiskHeatmapController:                  riskHeatmapController,
		RiskInherentControlAnalysisController:  riskInherentControlAnalysisController,
		RiskResidualEvaluationController:       riskResidualEvaluationController,
		RiskResidualControlAnalysisController:  riskResidualControlAnalysisController,
		RiskKriController:                      riskKriController,
		RiskTimelineController:                 riskTimelineController,
		RiskSettingsController:                 riskSettingsController,
		RiskClosureController:                  riskClosureController,
		RiskCommentController:                  riskCommentController,
		RiskReportController:                   riskReportController,
		KRIAssociatedItemController:            kriAssociatedItemController,
		KRIRegisterWorkflowController:          kriRegisterWorkflowController,
		KRIRegisterReviewController:            kriRegisterReviewController,
		KRIRegisterReviewCommentController:     kriRegisterReviewCommentController,
		RiskWorkflowController:                 riskWorkflowController,
		RiskReviewController:                   riskReviewController,
		RiskReviewCommentController:            riskReviewCommentController,

		RiskRegisterSvc:         riskRegisterSvc,
		RiskHeatMapSvc:          riskHeatmapSvc,
		KriRegisterSvc:          kriRegisterSvc,
		KRIWorkflowSvc:          kriRegisterWorkflowSvc,
		KRIReviewSvc:            kriRegisterReviewSvc,
		KRIReviewCommentSvc:     kriRegisterReviewCommentSvc,
		RiskActionPlanSvc:       riskActionPlanSvc,
		RiskLevelSvc:            riskLevelSvc,
		RiskResponseStrategySvc: riskResponseStrategySvc,
		RiskWorkflowSvc:         riskWorkflowSvc,
		RiskReviewSvc:           riskReviewSvc,
		RiskReviewCommentSvc:    riskReviewCommentSvc,
	}
}
