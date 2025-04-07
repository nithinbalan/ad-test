package initializers

import (
	"jethur/grc/organizational-context/controllers"
	"jethur/grc/organizational-context/repositories"
	"jethur/grc/organizational-context/services"

	"gorm.io/gorm"
)

type OrganizationalContextModule struct {
	SwotPestelQuestionGroupController                  *controllers.SwotPestelQuestionGroupController
	SwotPestelQuestionController                       *controllers.SwotPestelQuestionController
	SwotPestelCategoryController                       *controllers.SwotPestelCategoryController
	SwotPestelQuestionFactorController                 *controllers.SwotPestelQuestionFactorController
	StakeholderController                              *controllers.StakeholderController
	EvaluationStatusController                         *controllers.EvaluationStatusController
	StakeholderNeedsAndExpectationController           *controllers.StakeholderNeedsAndExpectationController
	StakeholderNeedsAndExpectationEvaluationController *controllers.StakeholderNeedsAndExpectationEvaluationController
	StakeholderEngagementController                    *controllers.StakeholderEngagementController
	StakeholderEngagementEvaluationController          *controllers.StakeholderEngagementEvaluationController
	StakeholderActionPlanController                    *controllers.StakeholderActionPlanController
}

func NewOrganizationalContextModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,

) *OrganizationalContextModule {
	swotQuestionGroupRepo := repositories.NewSwotPestelQuestionGroupRepository(db)
	swotQuestionGroupSvc := services.NewSwotPestelQuestionGroupService(swotQuestionGroupRepo, sharedModule.ActivityLogService)
	swotQuestionGroupController := controllers.NewSwotPestelQuestionGroupController(swotQuestionGroupSvc)

	swotQuestionRepo := repositories.NewSwotPestelQuestionRepository(db)
	swotQuestionSvc := services.NewSwotPestelQuestionService(swotQuestionRepo, sharedModule.ActivityLogService)
	swotQuestionController := controllers.NewSwotPestelQuestionController(swotQuestionSvc)

	swotPestelCategoryRepo := repositories.NewSwotPestelCategoryRepository(db)
	swotPestelCategorySvc := services.NewSwotPestelCategoryService(swotPestelCategoryRepo)
	swotPestelCategoryController := controllers.NewSwotPestelCategoryController(swotPestelCategorySvc)

	swotQuestionFactorRepo := repositories.NewSwotPestelQuestionFactorRepository(db)
	swotQuestionFactorSvc := services.NewSwotPestelQuestionFactorService(swotQuestionFactorRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc)
	swotQuestionFactorController := controllers.NewSwotPestelQuestionFactorController(swotQuestionFactorSvc)

	stakeholderRepo := repositories.NewStakeholderRepository(db)
	stakeholderSvc := services.NewStakeholderService(stakeholderRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc)
	stakeholderController := controllers.NewStakeholderController(stakeholderSvc)

	evaluationStatusRepo := repositories.NewEvaluationStatusRepository(db)
	evaluationStatusSvc := services.NewEvaluationStatusService(evaluationStatusRepo)
	evaluationStatusController := controllers.NewEvaluationStatusController(evaluationStatusSvc)

	stakeholderNeedsAndExpectationRepo := repositories.NewStakeholderNeedsAndExpectationRepository(db)
	stakeholderNeedsAndExpectationSvc := services.NewStakeholderNeedsAndExpectationService(
		stakeholderNeedsAndExpectationRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc,
	)
	stakeholderNeedsAndExpectationController := controllers.NewStakeholderNeedsAndExpectationController(stakeholderNeedsAndExpectationSvc)

	stakeholderNeedsAndExpectationEvaluationRepo := repositories.NewStakeholderNeedsAndExpectationEvaluationRepository(db)
	stakeholderNeedsAndExpectationEvaluationSvc := services.NewStakeholderNeedsAndExpectationEvaluationService(stakeholderNeedsAndExpectationEvaluationRepo, sharedModule.ActivityLogService, stakeholderSvc, stakeholderNeedsAndExpectationSvc, evaluationStatusSvc)
	stakeholderNeedsAndExpectationEvaluationController := controllers.NewStakeholderNeedsAndExpectationEvaluationController(stakeholderNeedsAndExpectationEvaluationSvc)

	stakeholderEngagementRepo := repositories.NewStakeholderEngagementRepository(db)
	stakeholderEngagementSvc := services.NewStakeholderEngagementService(
		stakeholderEngagementRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc,
	)
	stakeholderEngagementController := controllers.NewStakeholderEngagementController(stakeholderEngagementSvc)

	stakeholderEngagementEvaluationRepo := repositories.NewStakeholderEngagementEvaluationRepository(db)
	stakeholderEngagementEvaluationSvc := services.NewStakeholderEngagementEvaluationService(stakeholderEngagementEvaluationRepo, sharedModule.ActivityLogService, stakeholderSvc, stakeholderEngagementSvc, evaluationStatusSvc)
	stakeholderEngagementEvaluationController := controllers.NewStakeholderEngagementEvaluationController(stakeholderEngagementEvaluationSvc)

	stakeholderActionPlanRepo := repositories.NewStakeholderActionPlanRepository(db)
	stakeholderActionPlanSvc := services.NewStakeholderActionPlanService(stakeholderActionPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, stakeholderSvc, masterModule.ActionPlanStatusSvc, organizationModule.UserSettingSvc)
	stakeholderActionPlanController := controllers.NewStakeholderActionPlanController(stakeholderActionPlanSvc)

	return &OrganizationalContextModule{
		SwotPestelQuestionGroupController:                  swotQuestionGroupController,
		SwotPestelQuestionController:                       swotQuestionController,
		SwotPestelCategoryController:                       swotPestelCategoryController,
		SwotPestelQuestionFactorController:                 swotQuestionFactorController,
		StakeholderController:                              stakeholderController,
		EvaluationStatusController:                         evaluationStatusController,
		StakeholderNeedsAndExpectationController:           stakeholderNeedsAndExpectationController,
		StakeholderNeedsAndExpectationEvaluationController: stakeholderNeedsAndExpectationEvaluationController,
		StakeholderEngagementController:                    stakeholderEngagementController,
		StakeholderEngagementEvaluationController:          stakeholderEngagementEvaluationController,
		StakeholderActionPlanController:                    stakeholderActionPlanController,
	}
}
