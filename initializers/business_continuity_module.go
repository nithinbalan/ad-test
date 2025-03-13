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
) *BusinessContinuityModule {
	businessFunctionRepo := repositories.NewBusinessFunctionRepository(db)

	businessFunctionSvc := services.NewBusinessFunctionService(businessFunctionRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc,
	)
	businessFunctionController := controllers.NewBusinessFunctionController(businessFunctionSvc)

	businessImpactAssessmentRepo := repositories.NewBusinessImpactAssessmentRepository(db)
	businessImpactAssessmentSvc := services.NewBusinessImpactAssessmentService(businessImpactAssessmentRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, masterModule.ImpactCategorySvc, organizationModule.ProcessSvc, assetModule.AssetRegisterSvc, documentModule.DocumentSvc, riskModule.RiskLevelSvc, organizationModule.UserSettingSvc)
	businessImpactAssessmentController := controllers.NewBusinessImpactAssessmentController(businessImpactAssessmentSvc, organizationModule.ProcessAssociatedItemSvc, organizationModule.ProcessSvc)

	bcpStatusRepo := repositories.NewBcpStatusRepository(db)
	bcpStatusSvc := services.NewBcpStatusService(bcpStatusRepo)
	bcpStatusController := controllers.NewBcpStatusController(bcpStatusSvc)

	businessContinuityPlanRepo := repositories.NewBusinessContinuityPlanRepository(db)
	businessContinuityPlanSvc := services.NewBusinessContinuityPlanService(businessContinuityPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, businessFunctionSvc, organizationModule.UserSettingSvc, bcpStatusSvc)
	businessContinuityPlanController := controllers.NewBusinessContinuityPlanController(businessContinuityPlanSvc)

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
	testAndExerciseSvc := services.NewTestAndExerciseService(testAndExerciseRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, bcpTestAndExerciseSvc, businessContinuityPlanSvc, testAndExerciseStatusSvc, communicationActionPlanStatusSvc, masterModule.ActionPlanStatusSvc, organizationModule.UserSettingSvc)
	testAndExerciseController := controllers.NewTestAndExerciseController(testAndExerciseSvc)

	testAndExerciseCommunicationActionPlanRepo := repositories.NewTestAndExerciseCommunicationActionPlanRepository(db)
	testAndExerciseCommunicationActionPlanSvc := services.NewTestAndExerciseCommunicationActionPlanService(testAndExerciseCommunicationActionPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, communicationActionPlanStatusSvc)
	testAndExerciseCommunicationActionPlanController := controllers.NewTestAndExerciseCommunicationActionPlanController(testAndExerciseCommunicationActionPlanSvc)

	testAndExerciseRecoveryActionPlanRepo := repositories.NewTestAndExerciseRecoveryActionPlanRepository(db)
	testAndExerciseRecoveryActionPlanSvc := services.NewTestAndExerciseRecoveryActionPlanService(testAndExerciseRecoveryActionPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, testAndExerciseSvc, masterModule.ActionPlanStatusSvc)
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
	taeRecoveryActionPlanStatusUpdateSvc := services.NewTaeRecoveryActionPlanStatusUpdateService(taeRecoveryActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, testAndExerciseRecoveryActionPlanSvc, masterModule.ActionPlanStatusSvc)
	taeRecoveryActionPlanStatusUpdateController := controllers.NewTaeRecoveryActionPlanStatusUpdateController(taeRecoveryActionPlanStatusUpdateSvc)

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

		BusinessFunctionSvc:                  businessFunctionSvc,
		BcpTeamSvc:                           bcpTeamSvc,
		businessContinuityPlanSvc:            businessContinuityPlanSvc,
		testAndExerciseRecoveryActionPlanSvc: testAndExerciseRecoveryActionPlanSvc,
		testAndExerciseSvc:                   testAndExerciseSvc,
	}
}
