package routes

import (
	"jethur/grc/auth/controllers"
	"jethur/grc/auth/repositories"
	"jethur/grc/auth/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	authRoutes := router.Group("/")

	authRoutes.POST("/auth/login", authController.Login)
	authRoutes.POST("/auth/verify-otp", authController.Login)

	authRoutes.POST("/auth/password-recovery/send-otp", authController.SendPasswordRecoveryOTP)
	authRoutes.POST("/auth/password-recovery/resend-otp", authController.ReSendPasswordRecoveryOTP)
	authRoutes.GET("/auth/password-recovery/verify-otp", authController.Login)
	authRoutes.POST("/auth/password-recovery/set-password", authController.Login)
}
