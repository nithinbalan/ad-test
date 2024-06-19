package routes

import (
	"jethur/grc/common/controllers"
	"jethur/grc/common/repositories"
	"jethur/grc/common/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	commonRoutes := router.Group("/")

	commonRoutes.Use(middlewares.AuthMiddleware())

	fileRepo := repositories.NewFileRepository()
	fileSvc := services.NewFileService(fileRepo)
	fileController := controllers.NewFileController(fileSvc, db)

	commonRoutes.POST("/upload/temp", fileController.TempFileUpload)
	commonRoutes.POST("/files/preview", fileController.FilePreview)
}
