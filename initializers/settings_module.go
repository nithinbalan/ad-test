package initializers

import (
	"jethur/grc/settings/controllers"
	"jethur/grc/settings/repositories"
	"jethur/grc/settings/services"

	"gorm.io/gorm"
)

type SettingsModule struct {
	OrganizationLevelController *controllers.OrganizationLevelController
	SystemPreferenceController  *controllers.SystemPreferenceController

	SystemPreferenceSvc services.SystemPreferenceServiceInterface
	RoleSvc             services.RoleServiceInterface
}

func NewSettingsModule(db *gorm.DB, sharedModule *SharedModule, masterModule *MasterModule) *SettingsModule {
	organizationLevelRepo := repositories.NewOrganizationLevelRepository(db)
	organizationLevelSvc := services.NewOrganizationLevelService(organizationLevelRepo, sharedModule.ActivityLogService)
	organizationLevelController := controllers.NewOrganizationLevelController(organizationLevelSvc)

	systemPreferenceRepo := repositories.NewSystemPreferenceRepository(db)
	systemPreferenceSvc := services.NewSystemPreferenceService(systemPreferenceRepo, sharedModule.ActivityLogService, masterModule.CurrencySvc)
	systemPreferenceController := controllers.NewSystemPreferenceController(systemPreferenceSvc)

	roleRepo := repositories.NewRoleRepository(db)
	roleSvc := services.NewRoleService(roleRepo, sharedModule.ActivityLogService)

	return &SettingsModule{
		OrganizationLevelController: organizationLevelController,
		SystemPreferenceController:  systemPreferenceController,

		SystemPreferenceSvc: systemPreferenceSvc,
		RoleSvc:             roleSvc,
	}
}
