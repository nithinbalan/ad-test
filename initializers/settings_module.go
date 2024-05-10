package initializers

import (
	"jethur/grc/settings/controllers"
	"jethur/grc/settings/repositories"
	"jethur/grc/settings/services"

	"gorm.io/gorm"
)

type SettingsModule struct {
	OrganizationLevelController *controllers.OrganizationLevelController
	CurrencyController          *controllers.CurrencyController
	SystemPreferenceController  *controllers.SystemPreferenceController
	TimezoneController          *controllers.TimezoneController

	SystemPreferenceSvc services.SystemPreferenceServiceInterface
	RoleSvc             services.RoleServiceInterface
}

func NewSettingsModule(db *gorm.DB, sharedModule *SharedModule) *SettingsModule {
	organizationLevelRepo := repositories.NewOrganizationLevelRepository(db)
	organizationLevelSvc := services.NewOrganizationLevelService(organizationLevelRepo, sharedModule.ActivityLogService)
	organizationLevelController := controllers.NewOrganizationLevelController(organizationLevelSvc)

	currencyRepo := repositories.NewCurrencyRepository(db)
	currencySvc := services.NewCurrencyService(currencyRepo)
	currencyController := controllers.NewCurrencyController(currencySvc)

	systemPreferenceRepo := repositories.NewSystemPreferenceRepository(db)
	systemPreferenceSvc := services.NewSystemPreferenceService(systemPreferenceRepo, sharedModule.ActivityLogService, currencySvc)
	systemPreferenceController := controllers.NewSystemPreferenceController(systemPreferenceSvc)

	roleRepo := repositories.NewRoleRepository(db)
	roleSvc := services.NewRoleService(roleRepo, sharedModule.ActivityLogService)

	timezoneRepo := repositories.NewTimezoneRepository(db)
	timezoneSvc := services.NewTimezoneService(timezoneRepo)
	timezoneController := controllers.NewTimezoneController(timezoneSvc)

	return &SettingsModule{
		OrganizationLevelController: organizationLevelController,
		CurrencyController:          currencyController,
		SystemPreferenceController:  systemPreferenceController,
		TimezoneController:          timezoneController,

		SystemPreferenceSvc: systemPreferenceSvc,
		RoleSvc:             roleSvc,
	}
}
