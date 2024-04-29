package routes

import (
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/settings/controllers"
	settingRepositories "jethur/grc/settings/repositories"
	settingServices "jethur/grc/settings/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	settingRoutes := router.Group("/")

	settingRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	organizationLevelRepo := settingRepositories.NewOrganizationLevelRepository(db)
	organizationLevelSvc := settingServices.NewOrganizationLevelService(organizationLevelRepo, activityLogSvc)
	OrganizationLevelController := controllers.NewOrganizationLevelController(organizationLevelSvc)

	settingRoutes.GET("/organization-levels", OrganizationLevelController.GetAllOrganizationLevels)
	settingRoutes.POST("/organization-levels", OrganizationLevelController.CreateOrganizationLevel)
}
