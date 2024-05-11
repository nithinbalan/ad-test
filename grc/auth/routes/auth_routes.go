package routes

import (
	"jethur/initializers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	authRoutes := router.Group("/")

	authController := app.AuthModule.AuthController

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
