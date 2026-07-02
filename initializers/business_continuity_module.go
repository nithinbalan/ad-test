package initializers

import (
	"jethur/grc/business-continuity/controllers"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/business-continuity/services"

	"gorm.io/gorm"
)

type BusinessContinuityModule struct {
	BusinessFunctionController                       *controllers.BusinessFunctionController
	BusinessImpactAssessmentController               *controllers.BusinessImpactAssessmentController
	BusinessContinuityPlanController                 *controllers.BusinessContinuityPlanController
	BcpBiaMatrixController                           *controllers.BcpBiaMatrixController
	BcpPhaseController                               *controllers.BcpPhaseController
	BcpRecoveryPlanController                        *controllers.BcpRecoveryPlanController
	BcpCommunicationPlanController                   *controllers.BcpCommunicationPlanController
	BcpTestAndExerciseController                     *controllers.BcpTestAndExerciseController
	BcpStatusController                              *controllers.BcpStatusController
	BcpTeamController                                *controllers.BcpTeamController
	TestAndExerciseStatusController                  *controllers.TestAndExerciseStatusController
	TestAndExerciseController                        *controllers.TestAndExerciseController
	CommunicationActionPlanStatusController          *controllers.CommunicationActionPlanStatusController
	TestAndExerciseCommunicationActionPlanController *controllers.TestAndExerciseCommunicationActionPlanController
	TestAndExerciseRecoveryActionPlanController      *controllers.TestAndExerciseRecoveryActionPlanController
	TimelinessController                             *controllers.TimelinessController
	FeedbackOptionController                         *controllers.FeedbackOptionController
	TestResultOptionController                       *controllers.TestResultOptionController
	TestAndExerciseTestResultController              *controllers.TestAndExerciseTestResultController
	TaeRecoveryActionPlanStatusUpdateController      *controllers.TaeRecoveryActionPlanStatusUpdateController
	TaeRecoveryActionPlanCommentController           *controllers.TaeRecoveryActionPlanCommentController
	BcpRiskAnalysisController                        *controllers.BcpRiskAnalysisController
	BcpAssociatedItemController                      *controllers.BcpAssociatedItemController
	BusinessFunctionAssociateItemController          *controllers.BusinessFunctionAssociatedItemController
	BusinessContinuityDurationUnitController         *controllers.BusinessContinuityDurationUnitController
	JobTeamController                                *controllers.JobTeamController
	BCMExecutionRecordTypeController                 *controllers.BCMExecutionRecordTypeController
	BCMExecutionRecordController                     *controllers.BCMExecutionRecordController
	BCMProcessTypeController                         *controllers.BCMProcessTypeController
	BCMManualProcedureFormController                 *controllers.BCMManualProcedureFormController
	BCMImpactAnalysisLevelController                 *controllers.BCMImpactAnalysisLevelController
	BCMBIAOutageBucketController                     *controllers.BCMBIAOutageBucketController
	BCMBIAImpactLevelController                      *controllers.BCMBIAImpactLevelController
	BCMStaffAllocationController                     *controllers.BCMStaffAllocationController

	BusinessFunctionSvc                  services.BusinessFunctionServiceInterface
	BcpTeamSvc                           services.BcpTeamServiceInterface
	businessContinuityPlanSvc            services.BusinessContinuityPlanServiceInterface
	testAndExerciseSvc                   services.TestAndExerciseServiceInterface
	testAndExerciseRecoveryActionPlanSvc services.TestAndExerciseRecoveryActionPlanServiceInterface
}

func NewBusinessContinuityModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	assetModule *AssetModule,
	documentModule *DocumentModule,
	riskModule *RiskModule,
	commonModule *CommonModule,
	taskModule *TaskModule,
) *BusinessContinuityModule {
	businessFunctionRepo := repositories.NewBusinessFunctionRepository(db)

	businessFunctionSvc := services.NewBusinessFunctionService(businessFunctionRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc,
		organizationModule.ProcessAssociatedItemSvc,
		commonModule.AssociatedItemValidationService,
	)
	businessFunctionController := controllers.NewBusinessFunctionController(businessFunctionSvc)

	bcmImpactAnalysisLevelRepo := repositories.NewBCMImpactAnalysisLevelRepository(db)
	bcmImpactAnalysisLevelSvc := services.NewBCMImpactAnalysisLevelService(bcmImpactAnalysisLevelRepo)
	bcmImpactAnalysisLevelController := controllers.NewBCMImpactAnalysisLevelController(bcmImpactAnalysisLevelSvc)

	bcmBiaOutageBucketRepo := repositories.NewBCMBIAOutageBucketRepository(db)
	bcmBiaOutageBucketSvc := services.NewBCMBIAOutageBucketService(bcmBiaOutageBucketRepo)
	bcmBiaOutageBucketController := controllers.NewBCMBIAOutageBucketController(bcmBiaOutageBucketSvc)

	bcmBiaImpactLevelRepo := repositories.NewBCMBIAImpactLevelRepository(db)
	bcmBiaImpactLevelSvc := services.NewBCMBIAImpactLevelService(bcmBiaImpactLevelRepo)
	bcmBiaImpactLevelController := controllers.NewBCMBIAImpactLevelController(bcmBiaImpactLevelSvc)

	businessImpactAssessmentRepo := repositories.NewBusinessImpactAssessmentRepository(db)
	bcmStaffAllocationRepo := repositories.NewBCMStaffAllocationRepository(db)
	bcmStaffAllocationSvc := services.NewBCMStaffAllocationService(bcmStaffAllocationRepo, sharedModule.ActivityLogService)
	bcmStaffAllocationController := controllers.NewBCMStaffAllocationController(bcmStaffAllocationSvc)

	businessImpactAssessmentSvc := services.NewBusinessImpactAssessmentService(businessImpactAssessmentRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, masterModule.ImpactCategorySvc, organizationModule.ProcessSvc, assetModule.AssetRegisterSvc, documentModule.DocumentSvc, organizationModule.ProcessDependencySvc, riskModule.RiskLevelSvc, organizationModule.UserSettingSvc, bcmBiaImpactLevelSvc, bcmBiaOutageBucketSvc, bcmImpactAnalysisLevelSvc, bcmStaffAllocationSvc)
	businessImpactAssessmentController := controllers.NewBusinessImpactAssessmentController(businessImpactAssessmentSvc, organizationModule.ProcessAssociatedItemSvc, organizationModule.ProcessDependencySvc, organizationModule.ProcessSvc)

	bcpStatusRepo := repositories.NewBcpStatusRepository(db)
	bcpStatusSvc := services.NewBcpStatusService(bcpStatusRepo)
	bcpStatusController := controllers.NewBcpStatusController(bcpStatusSvc)

	businessContinuityPlanRepo := repositories.NewBusinessContinuityPlanRepository(db)
	businessContinuityPlanSvc := services.NewBusinessContinuityPlanService(
		businessContinuityPlanRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		businessFunctionSvc,
		organizationModule.UserSettingSvc,
		bcpStatusSvc,
		riskModule.RiskHeatMapSvc,
	)
	businessContinuityPlanController := controllers.NewBusinessContinuityPlanController(businessContinuityPlanSvc)

	bcpRiskAnalysisRepo := repositories.NewBcpRiskAnalysisRepository(db)
	bcpRiskAnalysisSvc := services.NewBcpRiskAnalysisService(bcpRiskAnalysisRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, businessContinuityPlanSvc)
	bcpRiskAnalysisController := controllers.NewBcpRiskAnalysisController(bcpRiskAnalysisSvc)

	bcpBiaMatrixRepo := repositories.NewBcpBiaMatrixRepository(db)
	bcpBiaMatrixSvc := services.NewBcpBiaMatrixService(bcpBiaMatrixRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, businessContinuityPlanSvc, businessFunctionSvc)
	bcpBiaMatrixController := controllers.NewBcpBiaMatrixController(bcpBiaMatrixSvc)

	bcpPhaseRepo := repositories.NewBcpPhaseRepository(db)
	bcpPhaseSvc := services.NewBcpPhaseService(bcpPhaseRepo)
	bcpPhaseController := controllers.NewBcpPhaseController(bcpPhaseSvc)

	bcpRecoveryPlanRepo := repositories.NewBcpRecoveryPlanRepository(db)
	bcpRecoveryPlanSvc := services.NewBcpRecoveryPlanService(bcpRecoveryPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc)
	bcpRecoveryPlanController := controllers.NewBcpRecoveryPlanController(bcpRecoveryPlanSvc)

	bcpCommunicationPlanRepo := repositories.NewBcpCommunicationPlanRepository(db)
	bcpCommunicationPlanSvc := services.NewBcpCommunicationPlanService(bcpCommunicationPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc)
	bcpCommunicationPlanController := controllers.NewBcpCommunicationPlanController(bcpCommunicationPlanSvc)

	bcpTestAndExerciseRepo := repositories.NewBcpTestAndExerciseRepository(db)
	bcpTestAndExerciseSvc := services.NewBcpTestAndExerciseService(bcpTestAndExerciseRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, businessContinuityPlanSvc)
	bcpTestAndExerciseController := controllers.NewBcpTestAndExerciseController(bcpTestAndExerciseSvc)

	bcpTeamRepo := repositories.NewBcpTeamRepository(db)
	bcpTeamSvc := services.NewBcpTeamService(bcpTeamRepo, organizationModule.UserSvc, businessContinuityPlanSvc, sharedModule.ActivityLogService)
	bcpTeamController := controllers.NewBcpTeamController(bcpTeamSvc)

	testAndExerciseStatusRepo := repositories.NewTestAndExerciseStatusRepository(db)
	testAndExerciseStatusSvc := services.NewTestAndExerciseStatusService(testAndExerciseStatusRepo)
	testAndExerciseStatusController := controllers.NewTestAndExerciseStatusController(testAndExerciseStatusSvc)

	communicationActionPlanStatusRepo := repositories.NewCommunicationActionPlanStatusRepository(db)
	communicationActionPlanStatusSvc := services.NewCommunicationActionPlanStatusService(communicationActionPlanStatusRepo)
	communicationActionPlanStatusController := controllers.NewCommunicationActionPlanStatusController(communicationActionPlanStatusSvc)

	testAndExerciseRepo := repositories.NewTestAndExerciseRepository(db)
	testAndExerciseSvc := services.NewTestAndExerciseService(testAndExerciseRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, bcpTestAndExerciseSvc, businessContinuityPlanSvc, testAndExerciseStatusSvc, communicationActionPlanStatusSvc, masterModule.ActionPlanStatusSvc, organizationModule.UserSettingSvc, taskModule.TaskService)
	testAndExerciseController := controllers.NewTestAndExerciseController(testAndExerciseSvc)

	testAndExerciseCommunicationActionPlanRepo := repositories.NewTestAndExerciseCommunicationActionPlanRepository(db)
	testAndExerciseCommunicationActionPlanSvc := services.NewTestAndExerciseCommunicationActionPlanService(testAndExerciseCommunicationActionPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, communicationActionPlanStatusSvc)
	testAndExerciseCommunicationActionPlanController := controllers.NewTestAndExerciseCommunicationActionPlanController(testAndExerciseCommunicationActionPlanSvc)

	testAndExerciseRecoveryActionPlanRepo := repositories.NewTestAndExerciseRecoveryActionPlanRepository(db)
	testAndExerciseRecoveryActionPlanSvc := services.NewTestAndExerciseRecoveryActionPlanService(testAndExerciseRecoveryActionPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, testAndExerciseSvc, masterModule.ActionPlanStatusSvc, taskModule.TaskService, commonModule.NotificationService, organizationModule.UserSettingSvc, commonModule.DelegateeUserService)
	testAndExerciseRecoveryActionPlanController := controllers.NewTestAndExerciseRecoveryActionPlanController(testAndExerciseRecoveryActionPlanSvc)

	timelinessRepo := repositories.NewTimelinessRepository(db)
	timelinessSvc := services.NewTimelinessService(timelinessRepo)
	timelinessController := controllers.NewTimelinessController(timelinessSvc)

	feedbackOptionRepo := repositories.NewFeedbackOptionRepository(db)
	feedbackOptionSvc := services.NewFeedbackOptionService(feedbackOptionRepo)
	feedbackOptionController := controllers.NewFeedbackOptionController(feedbackOptionSvc)

	testResultOptionRepo := repositories.NewTestResultOptionRepository(db)
	testResultOptionSvc := services.NewTestResultOptionService(testResultOptionRepo)
	testResultOptionController := controllers.NewTestResultOptionController(testResultOptionSvc)

	testAndExerciseTestResultRepo := repositories.NewTestAndExerciseTestResultRepository(db)
	testAndExerciseTestResultSvc := services.NewTestAndExerciseTestResultService(testAndExerciseTestResultRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, testAndExerciseSvc)
	testAndExerciseTestResultController := controllers.NewTestAndExerciseTestResultController(testAndExerciseTestResultSvc)

	taeRecoveryActionPlanStatusUpdateRepo := repositories.NewTaeRecoveryActionPlanStatusUpdateRepository(db)
	taeRecoveryActionPlanStatusUpdateSvc := services.NewTaeRecoveryActionPlanStatusUpdateService(taeRecoveryActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, testAndExerciseRecoveryActionPlanSvc, masterModule.ActionPlanStatusSvc, testAndExerciseSvc, organizationModule.UserSvc, commonModule.NotificationService, taskModule.TaskService, taskModule.TaskStatusService)
	taeRecoveryActionPlanStatusUpdateController := controllers.NewTaeRecoveryActionPlanStatusUpdateController(taeRecoveryActionPlanStatusUpdateSvc)

	taeRecoveryActionPlanCommentRepository := repositories.NewTaeRecoveryActionPlanCommentRepository(db)
	taeRecoveryActionPlanCommentSvc := services.NewTaeRecoveryActionPlanCommentService(taeRecoveryActionPlanCommentRepository, organizationModule.UserSvc, sharedModule.ActivityLogService, testAndExerciseRecoveryActionPlanSvc, commonModule.NotificationService)
	taeRecoveryActionPlanCommentController := controllers.NewTaeRecoveryActionPlanCommentController(taeRecoveryActionPlanCommentSvc)

	bcpAssociatedItemRepo := repositories.NewBcpAssociatedItemRepository(db)
	bcpAssociatedItemSvc := services.NewBcpAssociatedItemService(bcpAssociatedItemRepo, businessContinuityPlanSvc, commonModule.AssociatedItemValidationService, bcpRiskAnalysisSvc)
	bcpAssociatedItemController := controllers.NewBcpAssociatedItemController(bcpAssociatedItemSvc)

	businessFunctionAssociatedItemRepo := repositories.NewBusinessFunctionAssociatedItemRepository(db)
	businessFunctionAssociatedItemSvc := services.NewBusinessFunctionAssociatedItemService(businessFunctionAssociatedItemRepo, businessFunctionSvc, commonModule.AssociatedItemValidationService)
	businessFunctionAssociatedItemController := controllers.NewBusinessFunctionAssociatedItemController(businessFunctionAssociatedItemSvc)

	businessContinuityDurationUnitRepo := repositories.NewBusinessContinuityDurationUnitRepository(db)
	businessContinuityDurationUnitSvc := services.NewBusinessContinuityDurationUnitService(businessContinuityDurationUnitRepo)
	businessContinuityDurationUnitController := controllers.NewBusinessContinuityDurationUnitController(businessContinuityDurationUnitSvc)

	jobTeamRepo := repositories.NewJobTeamRepository(db)
	jobTeamSvc := services.NewJobTeamService(jobTeamRepo, sharedModule.ActivityLogService, organizationModule.DesignationSvc)
	jobTeamController := controllers.NewJobTeamController(jobTeamSvc)

	bcmExecutionRecordTypeRepo := repositories.NewBCMExecutionRecordTypeRepository(db)
	bcmExecutionRecordTypeSvc := services.NewBCMExecutionRecordTypeService(bcmExecutionRecordTypeRepo, sharedModule.ActivityLogService)
	bcmExecutionRecordTypeController := controllers.NewBCMExecutionRecordTypeController(bcmExecutionRecordTypeSvc)

	bcmExecutionRecordRepo := repositories.NewBCMExecutionRecordRepository(db)
	bcmExecutionRecordSvc := services.NewBCMExecutionRecordService(bcmExecutionRecordRepo, sharedModule.ActivityLogService)
	bcmExecutionRecordController := controllers.NewBCMExecutionRecordController(bcmExecutionRecordSvc)

	bcmProcessTypeRepo := repositories.NewBCMProcessTypeRepository(db)
	bcmProcessTypeSvc := services.NewBCMProcessTypeService(bcmProcessTypeRepo, sharedModule.ActivityLogService)
	bcmProcessTypeController := controllers.NewBCMProcessTypeController(bcmProcessTypeSvc)

	bcmManualProcedureFormRepo := repositories.NewBCMManualProcedureFormRepository(db)
	bcmManualProcedureFormSvc := services.NewBCMManualProcedureFormService(bcmManualProcedureFormRepo, sharedModule.ActivityLogService)
	bcmManualProcedureFormController := controllers.NewBCMManualProcedureFormController(bcmManualProcedureFormSvc)

	return &BusinessContinuityModule{
		BusinessFunctionController:                       businessFunctionController,
		BusinessImpactAssessmentController:               businessImpactAssessmentController,
		BusinessContinuityPlanController:                 businessContinuityPlanController,
		BcpBiaMatrixController:                           bcpBiaMatrixController,
		BcpPhaseController:                               bcpPhaseController,
		BcpRecoveryPlanController:                        bcpRecoveryPlanController,
		BcpCommunicationPlanController:                   bcpCommunicationPlanController,
		BcpTestAndExerciseController:                     bcpTestAndExerciseController,
		BcpStatusController:                              bcpStatusController,
		BcpTeamController:                                bcpTeamController,
		TestAndExerciseStatusController:                  testAndExerciseStatusController,
		TestAndExerciseController:                        testAndExerciseController,
		CommunicationActionPlanStatusController:          communicationActionPlanStatusController,
		TestAndExerciseCommunicationActionPlanController: testAndExerciseCommunicationActionPlanController,
		TestAndExerciseRecoveryActionPlanController:      testAndExerciseRecoveryActionPlanController,
		TimelinessController:                             timelinessController,
		FeedbackOptionController:                         feedbackOptionController,
		TestResultOptionController:                       testResultOptionController,
		TestAndExerciseTestResultController:              testAndExerciseTestResultController,
		TaeRecoveryActionPlanStatusUpdateController:      taeRecoveryActionPlanStatusUpdateController,
		TaeRecoveryActionPlanCommentController:           taeRecoveryActionPlanCommentController,
		BcpRiskAnalysisController:                        bcpRiskAnalysisController,
		BcpAssociatedItemController:                      bcpAssociatedItemController,
		BusinessFunctionAssociateItemController:          businessFunctionAssociatedItemController,
		BusinessContinuityDurationUnitController:         businessContinuityDurationUnitController,
		JobTeamController:                                jobTeamController,
		BCMExecutionRecordTypeController:                 bcmExecutionRecordTypeController,
		BCMExecutionRecordController:                     bcmExecutionRecordController,
		BCMProcessTypeController:                         bcmProcessTypeController,
		BCMManualProcedureFormController:                 bcmManualProcedureFormController,
		BCMImpactAnalysisLevelController:                 bcmImpactAnalysisLevelController,
		BCMBIAOutageBucketController:                     bcmBiaOutageBucketController,
		BCMBIAImpactLevelController:                      bcmBiaImpactLevelController,
		BCMStaffAllocationController:                     bcmStaffAllocationController,

		BusinessFunctionSvc:                  businessFunctionSvc,
		BcpTeamSvc:                           bcpTeamSvc,
		businessContinuityPlanSvc:            businessContinuityPlanSvc,
		testAndExerciseRecoveryActionPlanSvc: testAndExerciseRecoveryActionPlanSvc,
		testAndExerciseSvc:                   testAndExerciseSvc,
	}
}
