package routes

import (
	"jethur/grc/auth/controllers"
	"jethur/grc/auth/services"
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	organizationRepositories "jethur/grc/organization/repositories"
	organizationServices "jethur/grc/organization/services"
	settingRepositories "jethur/grc/settings/repositories"
	settingService "jethur/grc/settings/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	currencyRepo := settingRepositories.NewCurrencyRepository(db)
	currencySvc := settingService.NewCurrencyService(currencyRepo)

	systemPreferenceRepo := settingRepositories.NewSystemPreferenceRepository(db)
	systemPreferenceSvc := settingService.NewSystemPreferenceService(systemPreferenceRepo, activityLogSvc, currencySvc)

	designationRepo := organizationRepositories.NewDesignationRepository(db)
	designationSvc := organizationServices.NewDesignationService(designationRepo, activityLogSvc)

	userStatusRepo := organizationRepositories.NewUserStatusRepository(db)
	userStatusSvc := organizationServices.NewUserStatusService(userStatusRepo)

	roleRepo := settingRepositories.NewRoleRepository(db)
	roleSvc := settingService.NewRoleService(roleRepo, activityLogSvc)

	userRepo := organizationRepositories.NewUserRepository(db)
	userSvc := organizationServices.NewUserService(userRepo, designationSvc, userStatusSvc, activityLogSvc, roleSvc, systemPreferenceSvc)

	userSettingRepo := organizationRepositories.NewUserSettingRepository(db)
	userSettingSvc := organizationServices.NewUserSettingService(userSettingRepo, activityLogSvc)
	userSettingSvc.SetUserService(userSvc)

	authService := services.NewAuthService(userSvc, userSettingSvc)

	authController := controllers.NewAuthController(db, authService)

	authRoutes := router.Group("/")

	//Login
	authRoutes.POST("/auth/login", authController.Login)
	authRoutes.POST("/auth/token", authController.GetTokenFromRefreshToken)
	authRoutes.POST("/auth/login/resend-otp", authController.LoginResendOTP)
	authRoutes.POST("/auth/login/verify-otp", authController.LoginVerifyOTP)

	//Password Recovery
	authRoutes.POST("/auth/password-recovery/send-otp", authController.SendPasswordRecoveryOTP)
	authRoutes.POST("/auth/password-recovery/resend-otp", authController.ReSendPasswordRecoveryOTP)
	authRoutes.POST("/auth/password-recovery/verify-otp", authController.VerifyPasswordRecoveryOTP)
	authRoutes.PUT("/auth/password-recovery/set-password", authController.SetPassword)
}
