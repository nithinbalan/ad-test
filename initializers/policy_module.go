package initializers

import (
	"jethur/grc/policy/controllers"
	"jethur/grc/policy/repositories"
	"jethur/grc/policy/services"

	"gorm.io/gorm"
)

type PolicyModule struct {
	PolicyFolderController   *controllers.PolicyFolderController
	PolicyRegisterController *controllers.PolicyRegisterController
	PolicyClauseController   *controllers.PolicyClauseController
	PolicyReviewController   *controllers.PolicyReviewController
	PolicyWorkflowController *controllers.PolicyWorkflowController

	PolicyFolderSvc   services.PolicyFolderServiceInterface
	PolicyRegisterSvc services.PolicyRegisterServiceInterface
	PolicyClauseSvc   services.PolicyClauseServiceInterface
}

func NewPolicyModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
	documentModule *DocumentModule) *PolicyModule {

	policyWorkflowController := controllers.NewPolicyWorkflowController(documentModule.DocumentWorkflowSvc)

	policyReviewController := controllers.NewPolicyReviewController(documentModule.DocumentReviewSvc)

	policyFolderRepo := repositories.NewPolicyFolderRepository(db)
	policyFolderSvc := services.NewPolicyFolderService(policyFolderRepo, sharedModule.ActivityLogService, documentModule.DocumentAccessTypeSvc)
	policyFolderController := controllers.NewPolicyFolderController(policyFolderSvc)

	policyRegisterRepo := repositories.NewPolicyRegisterRepository(db)
	policyRegisterSvc := services.NewPolicyRegisterService(
		policyRegisterRepo,
		documentModule.DocumentCategorySvc,
		documentModule.DocumentTypeSvc,
		documentModule.DocumentAccessTypeSvc,
		documentModule.DocumentStatusSvc,
		organizationModule.UserSvc,
		sharedModule.ActivityLogService,
		organizationModule.UserSettingSvc,
		documentModule.DocumentAccessSvc,
		taskModule.TaskService,
		commonModule.NotificationService,
		settingsModule.OrganizationLevelSvc,
		masterModule.ReviewFrequencySvc,
		settingsModule.SystemPreferenceSvc,
		documentModule.DocumentWorkflowSvc,
		documentModule.DocumentWorkflowStatusSvc,
		masterModule.WorkflowContextSvc,
		settingsModule.WorkflowSvc,
	)
	policyRegisterController := controllers.NewPolicyRegisterController(policyRegisterSvc)

	policyClauseRepo := repositories.NewPolicyClauseRepository(db)
	policyClauseSvc := services.NewPolicyClauseService(policyClauseRepo, documentModule.DocumentSvc, sharedModule.ActivityLogService)
	policyClauseController := controllers.NewPolicyClauseController(policyClauseSvc)

	return &PolicyModule{
		PolicyFolderController:   policyFolderController,
		PolicyRegisterController: policyRegisterController,
		PolicyClauseController:   policyClauseController,
		PolicyWorkflowController: policyWorkflowController,
		PolicyReviewController:   policyReviewController,
	}
}
