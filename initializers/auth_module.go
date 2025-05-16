package initializers

import (
	"jethur/grc/auth/controllers"
	"jethur/grc/auth/services"

	"gorm.io/gorm"
)

type AuthModule struct {
	AuthController *controllers.AuthController
}

func NewAuthModule(db *gorm.DB, organizationModule *OrganizationModule, settingsModule *SettingsModule, commonModule *CommonModule) *AuthModule {

	authService := services.NewAuthService(
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc,
		commonModule.NotificationService,
		settingsModule.EmailTemplateSettingsSvc,
		settingsModule.ActiveDirectorySettingSvc,
		settingsModule.IntegrationSvc,
	)
	authController := controllers.NewAuthController(db, authService)

	return &AuthModule{
		AuthController: authController,
	}

}
