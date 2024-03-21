package routes

import (
	activityLogRepository "jethur/grc/logs/repositories"
	activityLogService "jethur/grc/logs/services"
	"jethur/grc/settings/controllers"
	organizationLevelRepository "jethur/grc/settings/repositories"
	organizationLevelService "jethur/grc/settings/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	settingRoute := router.Group("/")

	settingRoute.Use(middlewares.AuthMiddleware())

	activityLogRepo := activityLogRepository.NewActivityLogRepository(db)
	activityLogSvc := activityLogService.NewActivityLogService(activityLogRepo)

	organizationLevelRepo := organizationLevelRepository.NewOrganizationLevelRepository(db)
	organizationLevelSvc := organizationLevelService.NewOrganizationLevelService(organizationLevelRepo, activityLogSvc)
	OrganizationLevelController := controllers.NewOrganizationLevelController(organizationLevelSvc)

	settingRoute.GET("/settings/organization-level", OrganizationLevelController.GetAllOrganizationLevel)
	settingRoute.POST("/settings/organization-level", OrganizationLevelController.CreateOrganizationLevel)
}
