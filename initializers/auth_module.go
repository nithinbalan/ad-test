package initializers

import (
	"jethur/grc/auth/controllers"
	"jethur/grc/auth/repositories"
	"jethur/grc/auth/services"

	"gorm.io/gorm"
)

type AuthModule struct {
	AuthController *controllers.AuthController
}

func NewAuthModule(db *gorm.DB, organizationModule *OrganizationModule) *AuthModule {

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(*authRepo, organizationModule.UserSvc, organizationModule.UserSettingSvc)
	authController := controllers.NewAuthController(db, authService)

	return &AuthModule{
		AuthController: authController,
	}

}
