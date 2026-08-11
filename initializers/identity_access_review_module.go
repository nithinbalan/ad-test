package initializers

import (
	"jethur/grc/identity-access-review/controllers"
	"jethur/grc/identity-access-review/repositories"
	"jethur/grc/identity-access-review/services"

	"gorm.io/gorm"
)

type IdentityAccessReviewModule struct {
	ReviewTypeController                           *controllers.ReviewTypeController
	ReviewCampaignStatusController                 *controllers.ReviewCampaignStatusController
	IdentityReviewCampaignController               *controllers.IdentityReviewCampaignController
	IdentityReviewCampaignReviewController         *controllers.IdentityReviewCampaignReviewController
	IdentityReviewCampaignWorkflowController       *controllers.IdentityReviewCampaignWorkflowController
	IdentityReviewCampaignReviewCommentController  *controllers.IdentityReviewCampaignReviewCommentController
	IdentityReviewUserController                   *controllers.IdentityReviewUserController
	IdentityReviewUserCommentController            *controllers.IdentityReviewUserCommentController
	PrivilegedAccountController                    *controllers.PrivilegedAccountController
	AccountTypeController                          *controllers.AccountTypeController
	AccountDecisionController                      *controllers.AccountDecisionController
	IdentityReviewActionPlanController             *controllers.IdentityReviewActionPlanController
	IdentityReviewActionPlanStatusUpdateController *controllers.IdentityReviewActionPlanStatusUpdateController
	IdentityReviewActionPlanCommentController      *controllers.IdentityReviewActionPlanCommentController
}

func NewIdentityAccessReviewModule(
	db *gorm.DB,
	sharedModule *SharedModule,
	organizationModule *OrganizationModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
) *IdentityAccessReviewModule {
	reviewTypeRepo := repositories.NewReviewTypeRepository(db)
	reviewTypeSvc := services.NewReviewTypeService(reviewTypeRepo)
	reviewTypeController := controllers.NewReviewTypeController(reviewTypeSvc)

	reviewCampaignStatusRepo := repositories.NewReviewCampaignStatusRepository(db)
	reviewCampaignStatusSvc := services.NewReviewCampaignStatusService(reviewCampaignStatusRepo)
	reviewCampaignStatusController := controllers.NewReviewCampaignStatusController(reviewCampaignStatusSvc)

	identityReviewCampaignRepo := repositories.NewIdentityReviewCampaignRepository(db)
	identityReviewCampaignSvc := services.NewIdentityReviewCampaignService(
		identityReviewCampaignRepo,
		sharedModule.ActivityLogService,
		reviewTypeSvc,
		reviewCampaignStatusSvc,
		organizationModule.UserSvc,
		settingsModule.ReportTemplateSettingSvc,
		commonModule.NotificationService,
	)
	identityReviewCampaignController := controllers.NewIdentityReviewCampaignController(identityReviewCampaignSvc)

	identityReviewCampaignWorkflowRepo := repositories.NewIdentityReviewCampaignWorkflowRepository(db)
	identityReviewCampaignWorkflowSvc := services.NewIdentityReviewCampaignWorkflowService(
		identityReviewCampaignWorkflowRepo,
		settingsModule.WorkflowSvc,
		masterModule.WorkflowStatusSvc,
		identityReviewCampaignSvc,
		organizationModule.UserSvc,
		organizationModule.OrganizationUnitsSvc,
	)
	identityReviewCampaignSvc.SetIdentityReviewCampaignWorkflowService(identityReviewCampaignWorkflowSvc)
	identityReviewCampaignWorkflowController := controllers.NewIdentityReviewCampaignWorkflowController(identityReviewCampaignWorkflowSvc)

	identityReviewCampaignReviewRepo := repositories.NewIdentityReviewCampaignReviewRepository(db)
	identityReviewCampaignReviewSvc := services.NewIdentityReviewCampaignReviewService(
		identityReviewCampaignReviewRepo,
		identityReviewCampaignSvc,
		identityReviewCampaignWorkflowSvc,
		reviewCampaignStatusSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
	)
	identityReviewCampaignReviewController := controllers.NewIdentityReviewCampaignReviewController(identityReviewCampaignReviewSvc)

	identityReviewCampaignReviewCommentRepo := repositories.NewIdentityReviewCampaignReviewCommentRepository(db)
	identityReviewCampaignReviewCommentSvc := services.NewIdentityReviewCampaignReviewCommentService(identityReviewCampaignReviewCommentRepo, sharedModule.ActivityLogService, identityReviewCampaignSvc, masterModule.CommentStatusSvc)
	identityReviewCampaignReviewCommentController := controllers.NewIdentityReviewCampaignReviewCommentController(identityReviewCampaignReviewCommentSvc)

	privilegedAccountRepo := repositories.NewPrivilegedAccountRepository(db)
	privilegedAccountSvc := services.NewPrivilegedAccountService(
		privilegedAccountRepo,
		sharedModule.ActivityLogService,
	)
	privilegedAccountController := controllers.NewPrivilegedAccountController(privilegedAccountSvc)

	accountTypeRepo := repositories.NewAccountTypeRepository(db)
	accountTypeSvc := services.NewAccountTypeService(accountTypeRepo)
	accountTypeController := controllers.NewAccountTypeController(accountTypeSvc)

	accountDecisionRepo := repositories.NewAccountDecisionRepository(db)
	accountDecisionSvc := services.NewAccountDecisionService(accountDecisionRepo)
	accountDecisionController := controllers.NewAccountDecisionController(accountDecisionSvc)

	identityReviewUserRepo := repositories.NewIdentityReviewUserRepository(db)
	identityReviewUserSvc := services.NewIdentityReviewUserService(
		identityReviewUserRepo,
		sharedModule.ActivityLogService,
		settingsModule.OrganizationUnitSvc,
		privilegedAccountSvc,
		accountTypeSvc,
		accountDecisionSvc,
		identityReviewCampaignSvc,
	)
	identityReviewUserController := controllers.NewIdentityReviewUserController(identityReviewUserSvc)

	identityReviewUserCommentRepo := repositories.NewIdentityReviewUserCommentRepository(db)
	identityReviewUserCommentSvc := services.NewIdentityReviewUserCommentService(
		identityReviewUserCommentRepo,
		sharedModule.ActivityLogService,
		identityReviewUserSvc,
		accountDecisionSvc,
	)
	identityReviewUserCommentController := controllers.NewIdentityReviewUserCommentController(identityReviewUserCommentSvc)

	identityReviewActionPlanRepo := repositories.NewIdentityReviewActionPlanRepository(db)
	identityReviewActionPlanSvc := services.NewIdentityReviewActionPlanService(
		identityReviewActionPlanRepo,
		sharedModule.ActivityLogService,
		identityReviewUserSvc,
		masterModule.SeverityLevelSvc,
		masterModule.ActionPlanStatusSvc,
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc,
		taskModule.TaskService,
		commonModule.DelegateeUserService,
		accountDecisionSvc,
		privilegedAccountSvc,
		identityReviewUserCommentSvc,
		commonModule.NotificationService,
	)
	identityReviewActionPlanController := controllers.NewIdentityReviewActionPlanController(identityReviewActionPlanSvc)

	identityReviewActionPlanStatusUpdateRepo := repositories.NewIdentityReviewActionPlanStatusUpdateRepository(db)
	identityReviewActionPlanStatusUpdateSvc := services.NewIdentityReviewActionPlanStatusUpdateService(
		identityReviewActionPlanStatusUpdateRepo,
		sharedModule.ActivityLogService,
		identityReviewActionPlanSvc,
		masterModule.ActionPlanStatusSvc,
		taskModule.TaskService,
		taskModule.TaskStatusService,
	)
	identityReviewActionPlanStatusUpdateController := controllers.NewIdentityReviewActionPlanStatusUpdateController(identityReviewActionPlanStatusUpdateSvc)

	identityReviewActionPlanCommentRepo := repositories.NewIdentityReviewActionPlanCommentRepository(db)
	identityReviewActionPlanCommentSvc := services.NewIdentityReviewActionPlanCommentService(identityReviewActionPlanCommentRepo, sharedModule.ActivityLogService, identityReviewActionPlanSvc)
	identityReviewActionPlanCommentController := controllers.NewIdentityReviewActionPlanCommentController(identityReviewActionPlanCommentSvc)

	return &IdentityAccessReviewModule{
		ReviewTypeController:                           reviewTypeController,
		ReviewCampaignStatusController:                 reviewCampaignStatusController,
		IdentityReviewCampaignController:               identityReviewCampaignController,
		IdentityReviewCampaignReviewController:         identityReviewCampaignReviewController,
		IdentityReviewCampaignWorkflowController:       identityReviewCampaignWorkflowController,
		IdentityReviewCampaignReviewCommentController:  identityReviewCampaignReviewCommentController,
		IdentityReviewUserController:                   identityReviewUserController,
		IdentityReviewUserCommentController:            identityReviewUserCommentController,
		PrivilegedAccountController:                    privilegedAccountController,
		AccountTypeController:                          accountTypeController,
		AccountDecisionController:                      accountDecisionController,
		IdentityReviewActionPlanController:             identityReviewActionPlanController,
		IdentityReviewActionPlanStatusUpdateController: identityReviewActionPlanStatusUpdateController,
		IdentityReviewActionPlanCommentController:      identityReviewActionPlanCommentController,
	}
}
