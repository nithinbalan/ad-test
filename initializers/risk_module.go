package initializers

import (
	"jethur/grc/risk/controllers"
	"jethur/grc/risk/repositories"
	"jethur/grc/risk/services"

	"gorm.io/gorm"
)

type RiskModule struct {
	RiskCategoryController               *controllers.RiskCategoryController
	PotentialRiskStatusController        *controllers.PotentialRiskStatusController
	RiskLevelController                  *controllers.RiskLevelController
	PotentialRiskController              *controllers.PotentialRiskController
	PotentialRiskStatusUpdateController  *controllers.PotentialRiskStatusUpdateController
	KriCategoryController                *controllers.KRICategoryController
	KriRegisterController                *controllers.KRIRegisterController
	KriValueUpdateController             *controllers.KRIValueUpdateController
	RiskAppetiteStatementController      *controllers.RiskAppetiteStatementController
	RiskImpactCategoryController         *controllers.RiskImpactCategoryController
	RiskMatrixImpactController           *controllers.RiskMatrixImpactController
	RiskMatrixImpactGuidelineController  *controllers.RiskMatrixImpactGuidelineController
	RiskMatrixLikelihoodController       *controllers.RiskMatrixLikelihoodController
	RiskMatrixRiskRatingController       *controllers.RiskMatrixRiskRatingController
	RiskMatrixRiskRatingValueController  *controllers.RiskMatrixRiskRatingValueController
	RiskSourceController                 *controllers.RiskSourceController
	RiskStatusController                 *controllers.RiskStatusController
	RiskRegisterController               *controllers.RiskRegisterController
	RiskRootCauseAnalysisController      *controllers.RiskRootCauseAnalysisController
	RiskAppetiteToleranceLevelController *controllers.RiskAppetiteToleranceLevelController
	RiskInherentEvaluationController     *controllers.RiskInherentEvaluationController
}

func NewRiskModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	libraryModule *LibraryModule,
	organizationModule *OrganizationModule) *RiskModule {

	riskCategoryRepo := repositories.NewRiskCategoryRepository(db)
	riskCategorySvc := services.NewRiskCategoryService(riskCategoryRepo, sharedModule.ActivityLogService)
	riskCategoryController := controllers.NewRiskCategoryController(riskCategorySvc)

	potentialRiskStatusRepo := repositories.NewPotentialRiskStatusRepository(db)
	potentialRiskStatusSvc := services.NewPotentialRiskStatusService(potentialRiskStatusRepo)
	potentialRiskStatusController := controllers.NewPotentialRiskStatusController(potentialRiskStatusSvc)

	riskLevelRepo := repositories.NewRiskLevelRepository(db)
	riskLevelSvc := services.NewRiskLevelService(riskLevelRepo)
	riskLevelController := controllers.NewRiskLevelController(riskLevelSvc)

	potentialRiskRepo := repositories.NewPotentialRiskRepository(db)
	potentialRiskSvc := services.NewPotentialRiskService(
		potentialRiskRepo,
		sharedModule.ActivityLogService,
		masterModule.SeverityLevelSvc,
		riskCategorySvc,
		potentialRiskStatusSvc,
	)
	potentialRiskController := controllers.NewPotentialRiskController(potentialRiskSvc)

	potentialRiskStatusUpdateRepo := repositories.NewPotentialRiskStatusUpdateRepository(db)
	potentialRiskStatusUpdateSvc := services.NewPotentialRiskStatusUpdateService(
		potentialRiskStatusUpdateRepo,
		sharedModule.ActivityLogService,
		potentialRiskSvc,
		potentialRiskStatusSvc)
	potentialRiskStatusUpdateController := controllers.NewPotentialRiskStatusUpdateController(potentialRiskStatusUpdateSvc)

	kriCategoryRepository := repositories.NewKRICategoryRepository(db)
	kriCategorySvc := services.NewKRICategoryService(kriCategoryRepository, sharedModule.ActivityLogService)
	kriCategoryController := controllers.NewKRICategoryController(kriCategorySvc)

	//set kri categories service to kri library service
	libraryModule.KriLibrarySvc.SetKriCategoryService(kriCategorySvc)

	kriRegisterRepo := repositories.NewKRIRegisterRepository(db)
	kriRegisterSvc := services.NewKRIRegisterService(
		kriRegisterRepo,
		sharedModule.ActivityLogService,
		kriCategorySvc,
		organizationModule.UserSvc,
		riskLevelSvc,
		masterModule.ReviewFrequencySvc,
	)
	kriRegisterController := controllers.NewKRIRegisterController(kriRegisterSvc)

	kriValueUpdateRepo := repositories.NewKRIValueUpdateRepository(db)
	kriValueUpdateSvc := services.NewKRIValueUpdateService(kriValueUpdateRepo, sharedModule.ActivityLogService, kriRegisterSvc, masterModule.ReviewFrequencySvc)
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

	riskRegisterRepo := repositories.NewRiskRegisterRepository(db)
	riskRegisterSvc := services.NewRiskRegisterService(
		riskRegisterRepo,
		sharedModule.ActivityLogService,
		riskCategorySvc,
		riskSourceSvc,
		organizationModule.UserSvc,
		riskStatusSvc,
		masterModule.ReviewFrequencySvc)
	riskRegisterController := controllers.NewRiskRegisterController(riskRegisterSvc)

	riskRootCauseAnalysisRepo := repositories.NewRiskRootCauseAnalysisRepository(db)
	riskRootCauseAnalysisSvc := services.NewRiskRootCauseAnalysisService(riskRootCauseAnalysisRepo, sharedModule.ActivityLogService, riskRegisterSvc)
	riskRootCauseAnalysisController := controllers.NewRiskRootCauseAnalysisController(riskRootCauseAnalysisSvc)

	riskAppetiteToleranceLevelRepo := repositories.NewRiskAppetiteToleranceLevelRepository(db)
	riskAppetiteToleranceLevelSvc := services.NewRiskAppetiteToleranceLevelService(riskAppetiteToleranceLevelRepo)
	riskAppetiteToleranceLevelController := controllers.NewRiskAppetiteToleranceLevelController(riskAppetiteToleranceLevelSvc)

	// Update variable names
	riskInherentEvaluationRepo := repositories.NewRiskInherentEvaluationRepository(db)
	riskInherentEvaluationSvc := services.NewRiskInherentEvaluationService(riskInherentEvaluationRepo, sharedModule.ActivityLogService, riskRegisterSvc, riskMatrixImpactSvc, riskMatrixLikelihoodSvc, riskImpactCategorySvc, riskMatrixImpactGuidelineSvc, riskAppetiteStatementSvc)
	riskInherentEvaluationController := controllers.NewRiskInherentEvaluationController(riskInherentEvaluationSvc)

	return &RiskModule{
		RiskCategoryController:               riskCategoryController,
		PotentialRiskStatusController:        potentialRiskStatusController,
		RiskLevelController:                  riskLevelController,
		PotentialRiskController:              potentialRiskController,
		PotentialRiskStatusUpdateController:  potentialRiskStatusUpdateController,
		KriCategoryController:                kriCategoryController,
		KriRegisterController:                kriRegisterController,
		KriValueUpdateController:             kriValueUpdateController,
		RiskAppetiteStatementController:      riskAppetiteStatementController,
		RiskImpactCategoryController:         riskImpactCategoryController,
		RiskMatrixImpactController:           riskMatrixImpactController,
		RiskMatrixImpactGuidelineController:  riskMatrixImpactGuidelineController,
		RiskMatrixLikelihoodController:       riskMatrixLikelihoodController,
		RiskMatrixRiskRatingController:       riskMatrixRiskRatingController,
		RiskMatrixRiskRatingValueController:  riskMatrixRiskRatingValueController,
		RiskSourceController:                 riskSourceController,
		RiskStatusController:                 riskStatusController,
		RiskRegisterController:               riskRegisterController,
		RiskRootCauseAnalysisController:      riskRootCauseAnalysisController,
		RiskAppetiteToleranceLevelController: riskAppetiteToleranceLevelController,
		RiskInherentEvaluationController:     riskInherentEvaluationController,
	}

}
