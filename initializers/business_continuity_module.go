package initializers

import (
	"jethur/grc/business-continuity/controllers"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/business-continuity/services"

	"gorm.io/gorm"
)

type BusinessContinuityModule struct {
	BusinessFunctionController         *controllers.BusinessFunctionController
	BusinessImpactAssessmentController *controllers.BusinessImpactAssessmentController
	BusinessContinuityPlanController   *controllers.BusinessContinuityPlanController
	BcpBiaMatrixController             *controllers.BcpBiaMatrixController
	BcpPhaseController                 *controllers.BcpPhaseController
	BcpRecoveryPlanController          *controllers.BcpRecoveryPlanController
	BcpCommunicationPlanController     *controllers.BcpCommunicationPlanController
	BcpTestAndExerciseController       *controllers.BcpTestAndExerciseController

	BusinessFunctionSvc services.BusinessFunctionServiceInterface
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
	)
	businessFunctionController := controllers.NewBusinessFunctionController(businessFunctionSvc)

	businessImpactAssessmentRepo := repositories.NewBusinessImpactAssessmentRepository(db)
	businessImpactAssessmentSvc := services.NewBusinessImpactAssessmentService(businessImpactAssessmentRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, masterModule.ImpactCategorySvc, organizationModule.ProcessSvc, assetModule.AssetRegisterSvc, documentModule.DocumentSvc, riskModule.RiskLevelSvc)
	businessImpactAssessmentController := controllers.NewBusinessImpactAssessmentController(businessImpactAssessmentSvc, organizationModule.ProcessAssociatedItemSvc, organizationModule.ProcessSvc)

	businessContinuityPlanRepo := repositories.NewBusinessContinuityPlanRepository(db)
	businessContinuityPlanSvc := services.NewBusinessContinuityPlanService(businessContinuityPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, businessFunctionSvc, organizationModule.UserSettingSvc)
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
	bcpTestAndExerciseSvc := services.NewBcpTestAndExerciseService(bcpTestAndExerciseRepo, sharedModule.ActivityLogService, organizationModule.UserSvc)
	bcpTestAndExerciseController := controllers.NewBcpTestAndExerciseController(bcpTestAndExerciseSvc)

	return &BusinessContinuityModule{
		BusinessFunctionController:         businessFunctionController,
		BusinessImpactAssessmentController: businessImpactAssessmentController,
		BusinessContinuityPlanController:   businessContinuityPlanController,
		BcpBiaMatrixController:             bcpBiaMatrixController,
		BcpPhaseController:                 bcpPhaseController,
		BcpRecoveryPlanController:          bcpRecoveryPlanController,
		BcpCommunicationPlanController:     bcpCommunicationPlanController,
		BcpTestAndExerciseController:       bcpTestAndExerciseController,

		BusinessFunctionSvc: businessFunctionSvc,
	}
}
