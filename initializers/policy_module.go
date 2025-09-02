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

	PolicyFolderSvc   services.PolicyFolderServiceInterface
	PolicyRegisterSvc services.PolicyRegisterServiceInterface
}

func NewPolicyModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
	documentModule *DocumentModule) *PolicyModule {

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
		settingsModule.SystemPreferenceSvc)
	policyRegisterController := controllers.NewPolicyRegisterController(policyRegisterSvc)

	return &PolicyModule{
		PolicyFolderController:   policyFolderController,
		PolicyRegisterController: policyRegisterController,
	}
}
