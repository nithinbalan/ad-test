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
	RoleController              *controllers.RoleController
	OrganizationUnitController  *controllers.OrganizationUnitController
	LogSettingsController       *controllers.LogSettingsController

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
	roleController := controllers.NewRoleController(roleSvc)

	organizationUnitRepo := repositories.NewOrganizationUnitRepository(db)
	organizationUnitSvc := services.NewOrganizationUnitService(organizationUnitRepo, sharedModule.ActivityLogService, organizationLevelSvc)
	organizationUnitController := controllers.NewOrganizationUnitController(organizationUnitSvc)

	logSettingsRepo := repositories.NewLogSettingsRepository(db)
	logSettingsSvc := services.NewLogSettingsService(logSettingsRepo, sharedModule.ActivityLogService)
	logSettingsController := controllers.NewLogSettingsController(logSettingsSvc)

	return &SettingsModule{
		OrganizationLevelController: organizationLevelController,
		SystemPreferenceController:  systemPreferenceController,
		RoleController:              roleController,
		OrganizationUnitController:  organizationUnitController,
		LogSettingsController:       logSettingsController,

		SystemPreferenceSvc: systemPreferenceSvc,
		RoleSvc:             roleSvc,
	}
}
