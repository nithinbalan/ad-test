package routes

import (
	"jethur/grc/auth/controllers"
	"jethur/grc/auth/repositories"
	"jethur/grc/auth/services"
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	organizationRepositories "jethur/grc/organization/repositories"
	organizationServices "jethur/grc/organization/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	designationRepo := organizationRepositories.NewDesignationRepository(db)
	designationSvc := organizationServices.NewDesignationService(designationRepo, activityLogSvc)

	userStatusRepo := organizationRepositories.NewUserStatusRepository(db)
	userStatusSvc := organizationServices.NewUserStatusService(userStatusRepo)

	userRepo := organizationRepositories.NewUserRepository(db)
	userSvc := organizationServices.NewUserService(userRepo, designationSvc, userStatusSvc, activityLogSvc)

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, userSvc)

	authController := controllers.NewAuthController(db, authService)

	authRoutes := router.Group("/")

	authRoutes.POST("/auth/login", authController.Login)
	authRoutes.POST("/auth/token", authController.GetTokenFromRefreshToken)
	authRoutes.POST("/auth/mfa/resend-otp", authController.MFAReSendOTP)
	authRoutes.POST("/auth/mfa/verify-otp", authController.MFAVerifyOTP)

	authRoutes.POST("/auth/password-recovery/send-otp", authController.SendPasswordRecoveryOTP)
	authRoutes.POST("/auth/password-recovery/resend-otp", authController.ReSendPasswordRecoveryOTP)
	authRoutes.POST("/auth/password-recovery/verify-otp", authController.VerifyPasswordRecoveryOTP)
	authRoutes.PUT("/auth/password-recovery/set-password", authController.SetPassword)
}
