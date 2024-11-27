package initializers

import (
	"jethur/grc/compliance/controllers"
	"jethur/grc/compliance/repositories"
	"jethur/grc/compliance/services"

	"gorm.io/gorm"
)

type ComplianceModule struct {
	ComplianceCategoryController                  *controllers.ComplianceCategoryController
	ComplianceStatusController                    *controllers.ComplianceStatusController
	ComplianceReviewStatusController              *controllers.ComplianceReviewStatusController
	ComplianceRegisterController                  *controllers.ComplianceRegisterController
	ComplianceAssociatedItemController            *controllers.ComplianceAssociatedItemController
	ComplianceRequirementAssociatedItemController *controllers.ComplianceRequirementAssociatedItemController
	ComplianceRequirementCategoryController       *controllers.ComplianceRequirementCategoryController
	ComplianceRequirementController               *controllers.ComplianceRequirementController
	ComplianceActionPlanController                *controllers.ComplianceActionPlanController
	ComplianceStatusUpdateController              *controllers.ComplianceStatusUpdateController
	ComplianceRequirementStatusUpdateController   *controllers.ComplianceRequirementStatusUpdateController
	ComplianceActionPlanStatusUpdateController    *controllers.ComplianceActionPlanStatusUpdateController
	ComplianceReviewController                    *controllers.ComplianceReviewController
	SLACategoryController                         *controllers.SLACategoryController
	SLARegisterController                         *controllers.SLARegisterController
	SLACheckQuestionController                    *controllers.SLACheckQuestionController
	SLACheckAnswerController                      *controllers.SLACheckAnswerController
	SLARenewController                            *controllers.SLARenewController
	SLAActionPlanController                       *controllers.SLAActionPlanController
	SLAActionPlanStatusUpdateController           *controllers.SLAActionPlanStatusUpdateController
	SLAAssociatedItemController                   *controllers.SLAAssociatedItemController

	ComplianceRegisterSvc    services.ComplianceRegisterServiceInterface
	ComplianceRequirementSvc services.ComplianceRequirementServiceInterface
	ComplianceActionPlanSvc  services.ComplianceActionPlanServiceInterface
	SlaRegisterSvc           services.SLARegisterServiceInterface
}

func NewComplianceModule(db *gorm.DB, sharedModule *SharedModule, masterModule *MasterModule, organizationModule *OrganizationModule, taskModule *TaskModule, settingsModule *SettingsModule, commonModule *CommonModule) *ComplianceModule {
	complianceCategoryRepo := repositories.NewComplianceCategoryRepository(db)
	complianceCategorySvc := services.NewComplianceCategoryService(complianceCategoryRepo, sharedModule.ActivityLogService)
	complianceCategoryController := controllers.NewComplianceCategoryController(complianceCategorySvc)

	complianceStatusRepo := repositories.NewComplianceStatusRepository(db)
	complianceStatusSvc := services.NewComplianceStatusService(complianceStatusRepo)
	complianceStatusController := controllers.NewComplianceStatusController(complianceStatusSvc)

	complianceReviewStatusRepo := repositories.NewComplianceReviewStatusRepository(db)
	complianceReviewStatusSvc := services.NewComplianceReviewStatusService(complianceReviewStatusRepo)
	complianceReviewStatusController := controllers.NewComplianceReviewStatusController(complianceReviewStatusSvc)

	complianceRegisterRepo := repositories.NewComplianceRegisterRepository(db)
	complianceRegisterSvc := services.NewComplianceRegisterService(
		complianceRegisterRepo,
		sharedModule.ActivityLogService,
		masterModule.LifecycleStatusSvc,
		complianceCategorySvc,
		complianceStatusSvc,
		masterModule.ReviewFrequencySvc,
		masterModule.SeverityLevelSvc,
		organizationModule.ControlSvc,
		organizationModule.UserSettingSvc,
		organizationModule.UserSvc, commonModule.NotificationService,
	)
	complianceRegisterController := controllers.NewComplianceRegisterController(complianceRegisterSvc)

	complianceAssociatedItemRepo := repositories.NewComplianceAssociatedItemRepository(db)
	complianceAssociatedItemSvc := services.NewComplianceAssociatedItemService(complianceAssociatedItemRepo, complianceRegisterSvc)
	complianceAssociatedItemController := controllers.NewComplianceAssociatedItemController(complianceAssociatedItemSvc)

	complianceRequirementCategoryRepo := repositories.NewComplianceRequirementCategoryRepository(db)
	complianceRequirementCategorySvc := services.NewComplianceRequirementCategoryService(complianceRequirementCategoryRepo, sharedModule.ActivityLogService)
	complianceRequirementCategoryController := controllers.NewComplianceRequirementCategoryController(complianceRequirementCategorySvc)

	complianceRequirementRepo := repositories.NewComplianceRequirementRepository(db)
	complianceRequirementSvc := services.NewComplianceRequirementService(
		complianceRequirementRepo,
		sharedModule.ActivityLogService,
		complianceRegisterSvc,
		complianceRequirementCategorySvc,
		complianceStatusSvc,
		masterModule.SeverityLevelSvc,
		organizationModule.ControlSvc,
		organizationModule.StrategicObjectiveSvc,
		organizationModule.UserSvc, commonModule.NotificationService,
	)
	complianceRequirementController := controllers.NewComplianceRequirementController(complianceRequirementSvc)

	complianceActionPlanRepo := repositories.NewComplianceActionPlanRepository(db)
	complianceActionPlanSvc := services.NewComplianceActionPlanService(
		complianceActionPlanRepo,
		sharedModule.ActivityLogService,
		complianceRegisterSvc,
		complianceRequirementSvc,
		masterModule.ActionPlanStatusSvc,
		masterModule.SeverityLevelSvc,
		taskModule.TaskService,
		settingsModule.ModuleSvc, organizationModule.UserSvc, commonModule.NotificationService,
	)
	complianceActionPlanController := controllers.NewComplianceActionPlanController(complianceActionPlanSvc)

	complianceStatusUpdateRepo := repositories.NewComplianceStatusUpdateRepository(db)
	complianceStatusUpdateSvc := services.NewComplianceStatusUpdateService(complianceStatusUpdateRepo, sharedModule.ActivityLogService, complianceRegisterSvc, complianceStatusSvc)
	complianceStatusUpdateController := controllers.NewComplianceStatusUpdateController(complianceStatusUpdateSvc)

	complianceRequirementStatusUpdateRepo := repositories.NewComplianceRequirementStatusUpdateRepository(db)
	complianceRequirementStatusUpdateSvc := services.NewComplianceRequirementStatusUpdateService(complianceRequirementStatusUpdateRepo, sharedModule.ActivityLogService, complianceRequirementSvc, complianceStatusSvc)
	complianceRequirementStatusUpdateController := controllers.NewComplianceRequirementStatusUpdateController(complianceRequirementStatusUpdateSvc)

	complianceActionPlanStatusUpdateRepo := repositories.NewComplianceActionPlanStatusUpdateRepository(db)
	complianceActionPlanStatusUpdateSvc := services.NewComplianceActionPlanStatusUpdateService(complianceActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, complianceActionPlanSvc, masterModule.ActionPlanStatusSvc)
	complianceActionPlanStatusUpdateController := controllers.NewComplianceActionPlanStatusUpdateController(complianceActionPlanStatusUpdateSvc)

	complianceReviewRepo := repositories.NewComplianceReviewRepository(db)
	complianceReviewSvc := services.NewComplianceReviewService(complianceReviewRepo, sharedModule.ActivityLogService, complianceRegisterSvc, masterModule.ReviewFrequencySvc, complianceReviewStatusSvc)
	complianceReviewController := controllers.NewComplianceReviewController(complianceReviewSvc)

	slaCategoryRepo := repositories.NewSLACategoryRepository(db)
	slaCategorySvc := services.NewSLACategoryService(slaCategoryRepo, sharedModule.ActivityLogService)
	slaCategoryController := controllers.NewSLACategoryController(slaCategorySvc)

	slaRegisterRepo := repositories.NewSLARegisterRepository(db)
	slaRegisterSvc := services.NewSLARegisterService(slaRegisterRepo,
		slaCategorySvc,
		masterModule.SeverityLevelSvc,
		masterModule.LifecycleStatusSvc,
		sharedModule.ActivityLogService,
		organizationModule.ProcessSvc,
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc, commonModule.NotificationService,
	)
	slaRegisterController := controllers.NewSLARegisterController(slaRegisterSvc)

	slaCheckQuestionRepo := repositories.NewSLACheckQuestionRepository(db)
	slaCheckQuestionSvc := services.NewSLACheckQuestionService(slaCheckQuestionRepo,
		sharedModule.ActivityLogService,
		slaRegisterSvc,
		masterModule.ChecklistQuestionSvc,
		masterModule.QuestionGroupService)
	slaCheckQuestionController := controllers.NewSLACheckQuestionController(slaCheckQuestionSvc)

	slaCheckAnswerRepo := repositories.NewSLACheckAnswerRepository(db)
	slaCheckAnswerSvc := services.NewSLACheckAnswerService(slaCheckAnswerRepo,
		sharedModule.ActivityLogService,
		slaRegisterSvc,
		complianceStatusSvc,
		slaCheckQuestionSvc)
	slaCheckAnswerController := controllers.NewSLACheckAnswerController(slaCheckAnswerSvc)

	slaRenewRepo := repositories.NewSLARenewRepository(db)
	slaRenewSvc := services.NewSLARenewService(slaRenewRepo, slaRegisterRepo, masterModule.LifecycleStatusSvc,
		sharedModule.ActivityLogService)
	slaRenewController := controllers.NewSLARenewController(slaRenewSvc)

	slaActionPlanRepo := repositories.NewSLAActionPlanRepository(db)
	slaActionPlanSvc := services.NewSLAActionPlanService(
		slaActionPlanRepo,
		sharedModule.ActivityLogService,
		slaRegisterSvc,
		slaCheckQuestionSvc,
		masterModule.ActionPlanStatusSvc,
		masterModule.SeverityLevelSvc,
		taskModule.TaskService,
		organizationModule.UserSvc, commonModule.NotificationService,
	)
	slaActionPlanController := controllers.NewSLAActionPlanController(slaActionPlanSvc)

	slaActionPlanStatusUpdateRepo := repositories.NewSLAActionPlanStatusUpdateRepository(db)
	slaActionPlanStatusUpdateSvc := services.NewSLAActionPlanStatusUpdateService(slaActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, slaActionPlanSvc, masterModule.ActionPlanStatusSvc)
	slaActionPlanStatusUpdateController := controllers.NewSLAActionPlanStatusUpdateController(slaActionPlanStatusUpdateSvc)

	complianceRequirementAssociatedItemRepo := repositories.NewComplianceRequirementAssociatedItemRepository(db)
	complianceRequirementAssociatedItemSvc := services.NewComplianceRequirementAssociatedItemService(complianceRequirementAssociatedItemRepo, complianceRegisterSvc)
	complianceRequirementAssociatedItemController := controllers.NewComplianceRequirementAssociatedItemController(complianceRequirementAssociatedItemSvc)

	slaAssociatedItemRepo := repositories.NewSLAAssociatedItemRepository(db)
	slaAssociatedItemSvc := services.NewSLAAssociatedItemService(slaAssociatedItemRepo, slaRegisterSvc)
	slaAssociatedItemController := controllers.NewSLAAssociatedItemController(slaAssociatedItemSvc)

	return &ComplianceModule{
		ComplianceCategoryController:                  complianceCategoryController,
		ComplianceStatusController:                    complianceStatusController,
		ComplianceReviewStatusController:              complianceReviewStatusController,
		ComplianceRegisterController:                  complianceRegisterController,
		ComplianceAssociatedItemController:            complianceAssociatedItemController,
		ComplianceRequirementAssociatedItemController: complianceRequirementAssociatedItemController,
		ComplianceRequirementCategoryController:       complianceRequirementCategoryController,
		ComplianceRequirementController:               complianceRequirementController,
		ComplianceActionPlanController:                complianceActionPlanController,
		ComplianceStatusUpdateController:              complianceStatusUpdateController,
		ComplianceRequirementStatusUpdateController:   complianceRequirementStatusUpdateController,
		ComplianceActionPlanStatusUpdateController:    complianceActionPlanStatusUpdateController,
		ComplianceReviewController:                    complianceReviewController,
		SLACategoryController:                         slaCategoryController,
		SLARegisterController:                         slaRegisterController,
		SLACheckQuestionController:                    slaCheckQuestionController,
		SLACheckAnswerController:                      slaCheckAnswerController,
		SLARenewController:                            slaRenewController,
		SLAActionPlanController:                       slaActionPlanController,
		SLAActionPlanStatusUpdateController:           slaActionPlanStatusUpdateController,
		SLAAssociatedItemController:                   slaAssociatedItemController,

		ComplianceRequirementSvc: complianceRequirementSvc,
		ComplianceRegisterSvc:    complianceRegisterSvc,
		ComplianceActionPlanSvc:  complianceActionPlanSvc,
		SlaRegisterSvc:           slaRegisterSvc,
	}
}
