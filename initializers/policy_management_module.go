package initializers

import (
	"jethur/grc/policy-management/controllers"
	"jethur/grc/policy-management/repositories"
	"jethur/grc/policy-management/services"

	"gorm.io/gorm"
)

type PolicyManagementModule struct {
	PolicyRegisterController *controllers.PolicyRegisterController

	PolicyRegisterSvc services.PolicyRegisterServiceInterface
}

func NewPolicyManagementModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
	documentModule *DocumentModule) *PolicyManagementModule {

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

	return &PolicyManagementModule{
		PolicyRegisterController: policyRegisterController,
	}
}
