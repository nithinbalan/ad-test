package initializers

import (
	"jethur/grc/policy/controllers"
	"jethur/grc/policy/repositories"
	"jethur/grc/policy/services"

	"gorm.io/gorm"
)

type PolicyModule struct {
	PolicyRegisterController *controllers.PolicyRegisterController

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
		PolicyRegisterController: policyRegisterController,
	}
}
