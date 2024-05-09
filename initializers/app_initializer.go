package initializers

import (
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	organizationControllers "jethur/grc/organization/controllers"
	organizationRepositories "jethur/grc/organization/repositories"
	organizationServices "jethur/grc/organization/services"
	settingRepositories "jethur/grc/settings/repositories"
	settingServices "jethur/grc/settings/services"

	"gorm.io/gorm"
)

type App struct {
	UserController        *organizationControllers.UserController
	UserSettingController *organizationControllers.UserSettingController
	// Add other controllers as needed
}

func NewApp(db *gorm.DB) *App {
	// Initialize Repositories
	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	currencyRepo := settingRepositories.NewCurrencyRepository(db)
	systemPreferenceRepo := settingRepositories.NewSystemPreferenceRepository(db)
	roleRepo := settingRepositories.NewRoleRepository(db)
	designationRepo := organizationRepositories.NewDesignationRepository(db)
	userStatusRepo := organizationRepositories.NewUserStatusRepository(db)
	userRepo := organizationRepositories.NewUserRepository(db)
	userSettingRepo := organizationRepositories.NewUserSettingRepository(db)

	// Initialize Services
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)
	currencySvc := settingServices.NewCurrencyService(currencyRepo)
	systemPreferenceSvc := settingServices.NewSystemPreferenceService(systemPreferenceRepo, activityLogSvc, currencySvc)
	roleSvc := settingServices.NewRoleService(roleRepo, activityLogSvc)
	designationSvc := organizationServices.NewDesignationService(designationRepo, activityLogSvc)
	userStatusSvc := organizationServices.NewUserStatusService(userStatusRepo)

	userSvc := organizationServices.NewUserService(userRepo, designationSvc, userStatusSvc, activityLogSvc, roleSvc, systemPreferenceSvc)
	userSettingSvc := organizationServices.NewUserSettingService(userSettingRepo, activityLogSvc)
	userSvc.SetUserSettingService(userSettingSvc)
	userSettingSvc.SetUserService(userSvc)

	// Initialize Controllers
	userController := organizationControllers.NewUserController(userSvc)
	userSettingController := organizationControllers.NewUserSettingController(userSettingSvc)

	return &App{
		UserController:        userController,
		UserSettingController: userSettingController,
	}
}
