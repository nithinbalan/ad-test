package routes

import (
	"jethur/grc/auth/controllers"
	"jethur/grc/auth/repositories/implementations"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	authRepository := implementations.NewAuthRepository(db)
	authController := controllers.NewAuthController(authRepository, db)

	router.POST("/auth/login", authController.Login)
	router.POST("/auth/verify-otp", authController.Login)

	router.POST("/auth/password-recovery/get-otp", authController.Login)
	router.GET("/auth/password-recovery/resend-otp", authController.Login)
	router.GET("/auth/password-recovery/verify-otp", authController.Login)
	router.POST("/auth/password-recovery/set-password", authController.Login)
}
