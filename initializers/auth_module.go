package initializers

import (
	"jethur/grc/auth/controllers"
	"jethur/grc/auth/services"

	"gorm.io/gorm"
)

type AuthModule struct {
	AuthController *controllers.AuthController
}

func NewAuthModule(db *gorm.DB, organizationModule *OrganizationModule) *AuthModule {

	authService := services.NewAuthService(organizationModule.UserSvc, organizationModule.UserSettingSvc)
	authController := controllers.NewAuthController(db, authService)

	return &AuthModule{
		AuthController: authController,
	}

}
