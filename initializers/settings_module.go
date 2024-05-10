package initializers

import (
	settingControllers "jethur/grc/settings/controllers"
	settingRepositories "jethur/grc/settings/repositories"
	settingServices "jethur/grc/settings/services"

	"gorm.io/gorm"
)

type SettingsModule struct {
	OrganizationLevelController *settingControllers.OrganizationLevelController
	CurrencyController          *settingControllers.CurrencyController
	SystemPreferenceController  *settingControllers.SystemPreferenceController

	SystemPreferenceSvc settingServices.SystemPreferenceServiceInterface
	RoleSvc             settingServices.RoleServiceInterface
}

func NewSettingsModule(db *gorm.DB, sharedModule *SharedModule) *SettingsModule {
	organizationLevelRepo := settingRepositories.NewOrganizationLevelRepository(db)
	organizationLevelSvc := settingServices.NewOrganizationLevelService(organizationLevelRepo, sharedModule.ActivityLogService)
	organizationLevelController := settingControllers.NewOrganizationLevelController(organizationLevelSvc)

	currencyRepo := settingRepositories.NewCurrencyRepository(db)
	currencySvc := settingServices.NewCurrencyService(currencyRepo)
	currencyController := settingControllers.NewCurrencyController(currencySvc)

	systemPreferenceRepo := settingRepositories.NewSystemPreferenceRepository(db)
	systemPreferenceSvc := settingServices.NewSystemPreferenceService(systemPreferenceRepo, sharedModule.ActivityLogService, currencySvc)
	systemPreferenceController := settingControllers.NewSystemPreferenceController(systemPreferenceSvc)

	roleRepo := settingRepositories.NewRoleRepository(db)
	roleSvc := settingServices.NewRoleService(roleRepo, sharedModule.ActivityLogService)

	return &SettingsModule{
		OrganizationLevelController: organizationLevelController,
		CurrencyController:          currencyController,
		SystemPreferenceController:  systemPreferenceController,

		SystemPreferenceSvc: systemPreferenceSvc,
		RoleSvc:             roleSvc,
	}
}
